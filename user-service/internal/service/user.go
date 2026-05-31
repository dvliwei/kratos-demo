package service

import (
	"context"

	v1 "user-service/api/user/v1"
	"user-service/internal/biz"
)

// UserService 实现用户服务的 gRPC 接口。
type UserService struct {
	v1.UnimplementedUserServer

	uc *biz.UserUsecase
}

// NewUserService 创建用户服务实例。
func NewUserService(uc *biz.UserUsecase) *UserService {
	return &UserService{uc: uc}
}

// GetUser 根据用户 ID 查询用户基础信息。
func (s *UserService) GetUser(ctx context.Context, in *v1.GetUserRequest) (*v1.GetUserReply, error) {
	user, err := s.uc.GetUser(ctx, in.GetId())
	if err != nil {
		return nil, err
	}
	return &v1.GetUserReply{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

// ListUsersWithPage 分页查询用户列表。
func (s *UserService) ListUsersWithPage(ctx context.Context, in *v1.ListUsersRequest) (*v1.ListUsersReply, error) {
	var search *biz.SearchUser
	if in.GetSearch() != nil {
		search = &biz.SearchUser{
			Name:  in.GetSearch().GetName(),
			Email: in.GetSearch().GetEmail(),
		}
	}

	users, total, err := s.uc.ListUsersWithPage(ctx, int(in.GetPage()), int(in.GetSize()), search)
	if err != nil {
		return nil, err
	}

	replyUsers := make([]*v1.SearchUser, 0, len(users))
	for _, user := range users {
		replyUsers = append(replyUsers, &v1.SearchUser{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}

	return &v1.ListUsersReply{
		Total: uint32(total),
		Users: replyUsers,
	}, nil
}

// Login 用户邮箱密码登录
func (s *UserService) Login(ctx context.Context, in *v1.LoginRequest) (*v1.LoginReply, error) {
	email, password := in.GetEmail(), in.GetPassword()
	name, token, err := s.uc.Login(ctx, email, password)
	if err != nil {
		return nil, err
	}
	return &v1.LoginReply{
		Name:  *name,
		Token: *token,
	}, nil
}

func (s *UserService) GetUserTotal(ctx context.Context, in *v1.GetUserTotalRequest) (*v1.GetUserTotalReply, error) {
	total, err := s.uc.GetUserCount(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserTotalReply{
		Total: uint64(total),
	}, nil
}

func (s *UserService) UpdateUserName(ctx context.Context, in *v1.UpdateUserNameRequest) (*v1.UpdateUserNameReply, error) {
	id, name, err := s.uc.UpdateUserName(ctx, in.GetId(), in.GetName())
	if err != nil {
		return nil, err
	}
	return &v1.UpdateUserNameReply{
		Id:   id,
		Name: name,
	}, nil
}
