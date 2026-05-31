package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
)

// User 是业务层使用的用户领域模型。
type User struct {
	ID    uint64
	Name  string
	Email string
}

// UserRepo 定义业务层依赖的数据访问能力。
type UserRepo interface {
	Login(context.Context, string, string) (*string, *string, error)
	FindByID(context.Context, uint64) (*User, error)
	ListUsersWithPage(context.Context, int, int, *SearchUser) ([]*User, int64, error)
	GetUserCount(context.Context) (uint64, error)
	UpdateUserName(context.Context, uint64, string) (uint64, string, error)
}

// UserUsecase 包含用户业务规则，定义了用户相关操作的业务逻辑。
type UserUsecase struct {
	repo UserRepo
}

// NewUserUsecase 创建用户业务用例实例。
func NewUserUsecase(repo UserRepo) *UserUsecase {
	return &UserUsecase{repo: repo}
}

// GetUser 根据用户 ID 查询用户基础信息。
func (uc *UserUsecase) GetUser(ctx context.Context, id uint64) (*User, error) {
	if id == 0 {
		return nil, errors.BadRequest("USER_ID_INVALID", "user id is required")
	}
	return uc.repo.FindByID(ctx, id)
}

// Login 用户邮箱密码登录
// biz
func (uc *UserUsecase) Login(ctx context.Context, email string, password string) (*string, *string, error) {
	if email == "" || password == "" {
		return nil, nil, errors.BadRequest("USER_EMAIL_INVALID", "email or password is empty")
	}
	return uc.repo.Login(ctx, email, password)
}

type SearchUser struct {
	Name  string
	Email string
}

// ListUsersWithPage 分页查询用户列表，并返回总记录数。
func (uc *UserUsecase) ListUsersWithPage(ctx context.Context, page, pageSize int, search *SearchUser) ([]*User, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return uc.repo.ListUsersWithPage(ctx, page, pageSize, search)
}

func (uc *UserUsecase) GetUserCount(ctx context.Context) (uint64, error) {
	return uc.repo.GetUserCount(ctx)
}

func (uc *UserUsecase) UpdateUserName(ctx context.Context, id uint64, name string) (uint64, string, error) {
	if id == 0 || name == "" {
		return 0, "", errors.BadRequest("USER_ID_INVALID", "user id or name is required")
	}
	return uc.repo.UpdateUserName(ctx, id, name)
}
