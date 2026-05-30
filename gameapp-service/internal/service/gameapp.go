/**
 * @Title
 * @Author: liwei
 * @Description:  TODO
 * @File:  gameapp
 * @Version: 1.0.0
 * @Date: 2026/05/29 23:11
 * @Update liwei 2026/5/29 23:11
 */

package service

import (
	"context"
	v1 "gameapp-service/api/gameapp/v1"
	"gameapp-service/internal/biz"
)

type GameAppService struct {
	v1.UnimplementedGameAppServiceServer
	uc *biz.GameAppUseCase
}

func NewGameAppService(uc *biz.GameAppUseCase) *GameAppService {
	return &GameAppService{uc: uc}
}

func (s *GameAppService) GetGameApp(ctx context.Context, req *v1.GetGameAppRequest) (*v1.GetGameAppResponse, error) {
	gameApp, err := s.uc.GetByIDApp(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	info := &v1.GameAppInfo{
		Id:     gameApp.ID,
		GameId: gameApp.GameID,
		AppId:  gameApp.AppID,
		Name:   gameApp.Name,
		AppKey: gameApp.AppKey,
	}
	return &v1.GetGameAppResponse{
		Info: info,
	}, nil
}

func (s *GameAppService) CountGameApps(ctx context.Context, req *v1.CountGameAppsRequest) (*v1.CountGameAppsResponse, error) {
	total, err := s.uc.CountGameApps(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.CountGameAppsResponse{
		Total: total,
	}, nil
}
