/**
 * @Title
 * @Author: liwei
 * @Description:  负责“访问外部资源”
 * 这里不是访问数据库，而是访问 user-service 的 gRPC
 * @File:  user.go
 * @Version: 1.0.0
 * @Date: 2026/05/29 21:17
 * @Update liwei 2026/5/29 21:17
 */

package data

import (
	"context"
	"gateway-service/internal/biz"
	"gateway-service/internal/pkg/requestid"
	userv1 "user-service/api/user/v1"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	grpcmetadata "google.golang.org/grpc/metadata"
)

type userRepo struct {
	client userv1.UserClient
	logger *log.Helper
}

func NewUserRepo(logger log.Logger) biz.UserRepo {
	conn, err := grpc.Dial(
		"127.0.0.1:9100",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return &userRepo{
		client: userv1.NewUserClient(conn),
		logger: log.NewHelper(logger),
	}
}

func (r *userRepo) GetUser(ctx context.Context, userID int64) (*biz.User, error) {
	ctx = withRequestID(ctx)
	resp, err := r.client.GetUser(ctx, &userv1.GetUserRequest{Id: uint64(userID)})
	if err != nil {
		return nil, err
	}
	return &biz.User{
		ID:    int64(resp.Id),
		Name:  resp.Name,
		Email: resp.Email,
	}, nil
}

func (r *userRepo) ListUsersWithPage(ctx context.Context, page, pageSize int, search *biz.SearchUser) ([]*biz.User, int64, error) {
	ctx = withRequestID(ctx)
	var reqSearch *userv1.SearchUser
	if search != nil {
		reqSearch = &userv1.SearchUser{
			Name:  search.Name,
			Email: search.Email,
		}
	}

	resp, err := r.client.ListUsersWithPage(ctx, &userv1.ListUsersRequest{
		Page:   uint32(page),
		Size:   uint32(pageSize),
		Search: reqSearch,
	})
	if err != nil {
		return nil, 0, err
	}

	users := make([]*biz.User, 0, len(resp.Users))
	for _, u := range resp.Users {
		users = append(users, &biz.User{
			ID:    int64(u.Id),
			Name:  u.Name,
			Email: u.Email,
		})
	}
	return users, int64(resp.Total), nil
}

func withRequestID(ctx context.Context) context.Context {
	if rid := requestid.FromContext(ctx); rid != "" {
		return grpcmetadata.AppendToOutgoingContext(ctx, requestid.MetadataKey, rid)
	}
	return ctx
}

// Login 用户邮箱密码登录
func (r *userRepo) Login(ctx context.Context, email string, password string) (*string, *string, error) {
	ctx = withRequestID(ctx)
	resp, err := r.client.Login(ctx, &userv1.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, nil, err
	}
	return &resp.Name, &resp.Token, nil
}

func (r *userRepo) GetUserCount(ctx context.Context) (uint64, error) {
	ctx = withRequestID(ctx)
	resp, err := r.client.GetUserTotal(ctx, &userv1.GetUserTotalRequest{})
	if err != nil {
		return 0, err
	}
	return resp.Total, nil
}
