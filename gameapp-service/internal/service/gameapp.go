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

	"google.golang.org/protobuf/reflect/protoreflect"
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

func (s *GameAppService) ListGameAppsWithPage(ctx context.Context, req *v1.ListGameAppsRequest) (*v1.ListGameAppsResponse, error) {
	if req.GetPage() < 1 {
		req.Page = 1
	}
	if req.GetPageSize() < 1 {
		req.PageSize = 10
	}
	if req.GetPageSize() > 20 {
		req.PageSize = 20
	}
	var search *biz.GameAppsSearch
	if req.GetSearch() != nil {
		search = &biz.GameAppsSearch{
			Name:   req.GetSearch().GetName(),
			TypeOs: gameAppsSearchTypeOS(req.GetSearch()),
		}
	}
	gameApps, total, err := s.uc.ListGameAppsWithPage(ctx, int(req.GetPage()), int(req.GetPageSize()), search)
	if err != nil {
		return nil, err
	}
	gameAppInfos := make([]*v1.GameAppInfo, 0, len(gameApps))
	for _, gameApp := range gameApps {
		gameAppInfos = append(gameAppInfos, &v1.GameAppInfo{
			Id:     gameApp.ID,
			AppId:  gameApp.AppID,
			Name:   gameApp.Name,
			AppKey: gameApp.AppKey,
		})
	}
	return &v1.ListGameAppsResponse{
		Infos: gameAppInfos,
		Total: uint64(total),
	}, nil
}

func gameAppsSearchTypeOS(search *v1.GameAppsSearch) *int32 {
	if search == nil {
		return nil
	}
	field := search.ProtoReflect().Descriptor().Fields().ByName(protoreflect.Name("type_os"))
	if field == nil || !search.ProtoReflect().Has(field) {
		return nil
	}
	typeOS := search.GetTypeOs()
	return &typeOS
}
