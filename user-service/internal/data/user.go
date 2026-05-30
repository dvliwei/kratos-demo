package data

import (
	"context"
	stderrors "errors"
	"fmt"
	"time"

	"user-service/internal/biz"
	"user-service/internal/conf"
	jwtpkg "user-service/internal/pkg/jwt"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type userRepo struct {
	log   *log.Helper
	db    *gorm.DB
	redis *redis.Client
	jwt   *jwtpkg.Manager
}

type userModel struct {
	ID              uint64         `gorm:"column:id;primaryKey;autoIncrement"`
	Name            string         `gorm:"column:name;size:191;not null"`
	Email           string         `gorm:"column:email;size:191;not null;uniqueIndex:users_email_unique"`
	EmailVerifiedAt *time.Time     `gorm:"column:email_verified_at"`
	Password        string         `gorm:"column:password;size:191;not null"`
	RememberToken   *string        `gorm:"column:remember_token;size:512"`
	CreatedAt       *time.Time     `gorm:"column:created_at"`
	UpdatedAt       *time.Time     `gorm:"column:updated_at"`
	RoleID          int16          `gorm:"column:role_id;not null;default:0;comment:角色id"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

// TableName 指定用户模型对应的数据库表名。
func (userModel) TableName() string {
	return "users"
}

// NewUserRepo 创建用户仓储实例，用于封装用户表的数据访问逻辑。
func NewUserRepo(data *Data, auth *conf.JWTAuthConfig, logger log.Logger) (biz.UserRepo, error) {
	var jwtConf *conf.JWTConfig
	if auth != nil {
		jwtConf = auth.Jwt
	}
	if jwtConf == nil {
		jwtConf = &conf.JWTConfig{}
	}
	jwtManager, err := jwtpkg.NewManager(jwtConf.Password, jwtpkg.DurationFromSeconds(jwtConf.ExpireSeconds))
	if err != nil {
		return nil, err
	}
	return &userRepo{
		log:   log.NewHelper(logger),
		db:    data.db,
		redis: data.redis,
		jwt:   jwtManager,
	}, nil
}

// FindByID 根据用户 ID 从 users 表查询用户基础信息。
func (r *userRepo) FindByID(ctx context.Context, id uint64) (*biz.User, error) {
	var user userModel
	err := r.db.WithContext(ctx).
		Select("id", "name", "email").
		First(&user, "id = ?", id).
		Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("USER_NOT_FOUND", "user not found")
		}
		r.log.Errorf("query user failed: %v", err)
		return nil, err
	}
	if user.ID == 0 {
		return nil, errors.NotFound("USER_NOT_FOUND", "user not found")
	}
	return &biz.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

// ListUsersWithPage 按页查询用户基础信息和总记录数。
func (r *userRepo) ListUsersWithPage(ctx context.Context, page, pageSize int, search *biz.SearchUser) ([]*biz.User, int64, error) {
	var users []userModel
	var total int64

	baseQuery := r.db.WithContext(ctx).
		Session(&gorm.Session{SkipDefaultTransaction: true}).
		Model(&userModel{})
	if search != nil {
		if search.Name != "" {
			baseQuery = baseQuery.Where("name LIKE ?", "%"+search.Name+"%")
		}
		if search.Email != "" {
			baseQuery = baseQuery.Where("email LIKE ?", "%"+search.Email+"%")
		}
	}
	err := baseQuery.Count(&total).Error
	if err != nil {
		r.log.Errorf("count users failed: %v", err)
		return nil, 0, err
	}

	if total == 0 {
		return []*biz.User{}, 0, nil
	}

	err = baseQuery.
		Select("id", "name", "email").
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&users).
		Error
	if err != nil {
		r.log.Errorf("find users failed: %v", err)
		return nil, total, err
	}
	bizUsers := make([]*biz.User, 0, len(users))
	for _, user := range users {
		bizUsers = append(bizUsers, &biz.User{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}
	return bizUsers, total, nil
}

// Login 用户邮箱密码登录
func (r *userRepo) Login(ctx context.Context, email string, password string) (*string, *string, error) {
	var user userModel
	query := r.db.WithContext(ctx).Model(&userModel{})
	err := query.Where("email = ?", email).First(&user).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.Unauthorized("INVALID_CREDENTIALS", "invalid email or password")
		}
		return nil, nil, err
	}
	if user.Password != password {
		return nil, nil, errors.Unauthorized("INVALID_CREDENTIALS", "invalid email or password")
	}

	token, err := r.jwt.GenerateToken(user.ID, user.Name)
	if err != nil {
		r.log.Errorf("generate token failed: %v", err)
		return nil, nil, err
	}

	err = query.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&userModel{}).
			Where("id = ?", user.ID).
			Update("remember_token", token).
			Error
		if err != nil {
			return err
		}
		redisKey := userTokenKey(user.ID)
		if err := r.redis.Set(ctx, redisKey, token, 0).Err(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		r.log.Errorf("save user token failed: %v", err)
		return nil, nil, err
	}
	return &user.Name, &token, nil
}

func userTokenKey(id uint64) string {
	return fmt.Sprintf("user:token:%d", id)
}

// GetUserCount 查询用户总数
func (r *userRepo) GetUserCount(ctx context.Context) (uint64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&userModel{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return uint64(count), nil
}
