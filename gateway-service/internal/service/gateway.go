/**
 * @Title
 * @Author: liwei
 * @Description:  TODO
 * @File:  gateway
 * @Version: 1.0.0
 * @Date: 2026/05/29 20:18
 * @Update liwei 2026/5/29 20:18
 */

package service

import (
	"context"
	v1 "gateway-service/api/gateway/v1"
	"gateway-service/internal/biz"
	"strconv"
)

type GatewayService struct {
	// 实现 GatewayServiceServer 接口
	v1.UnimplementedGatewayServiceServer
	uc        *biz.UserUseCase
	gameAppUC *biz.GameAppUseCase
}

func NewGatewayService(uc *biz.UserUseCase, gameAppUC *biz.GameAppUseCase) *GatewayService {
	return &GatewayService{
		uc:        uc,
		gameAppUC: gameAppUC,
	}
}

func (a *GatewayService) GetGatewayInfo(ctx context.Context, req *v1.GetUserRequest) (*v1.GetUserReply, error) {
	reply, err := a.uc.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &v1.GetUserReply{
		Id:    reply.ID,
		Name:  reply.Name,
		Email: reply.Email,
	}, nil
}

func (a *GatewayService) ListUsersWithPage(ctx context.Context, req *v1.ListUsersRequest) (*v1.ListUsersReply, error) {
	var search *biz.SearchUser
	if req.GetSearch() != nil {
		search = &biz.SearchUser{
			Name:  req.GetSearch().GetName(),
			Email: req.GetSearch().GetEmail(),
		}
	}

	users, total, err := a.uc.ListUsersWithPage(ctx, int(req.GetPage()), int(req.GetPageSize()), search)
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
		Total: int32(total),
		Users: replyUsers,
	}, nil
}

func (a *GatewayService) GetGameApp(ctx context.Context, req *v1.GetGameAppRequest) (*v1.GetGameAppReply, error) {
	gameApp, err := a.gameAppUC.GetGameApp(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &v1.GetGameAppReply{
		GameApp: &v1.GameAppInfo{
			Id:     int64(gameApp.ID),
			Name:   gameApp.Name,
			GameId: strconv.FormatInt(gameApp.GameID, 10),
			AppKey: gameApp.AppKey,
		},
	}, nil
}

func (a *GatewayService) GetUserGameAppStats(ctx context.Context, req *v1.GetUserGameAppStatsRequest) (*v1.GetUserGameAppStatsReply, error) {
	result := &v1.GetUserGameAppStatsReply{}
	totalUsers, err := a.uc.GetUserCount(ctx)
	if err == nil {
		result.TotalUsers = int64(totalUsers)
	}
	totalGameApps, err := a.gameAppUC.CountGameApps(ctx)
	if err == nil {
		result.TotalGameApps = int64(totalGameApps)
	}
	return result, nil
}

// Login 用户邮箱密码登录
func (a *GatewayService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error) {
	name, token, err := a.uc.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}
	return &v1.LoginReply{
		Name:  *name,
		Token: *token,
	}, nil
}
