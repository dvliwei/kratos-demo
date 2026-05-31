/**
 * @Title
 * @Author: liwei
 * @Description:  TODO
 * @File:  gameapp
 * @Version: 1.0.0
 * @Date: 2026/05/29 23:43
 * @Update liwei 2026/5/29 23:43
 */

package data

import (
	"context"
	gameappv1 "gameapp-service/api/gameapp/v1"
	"gateway-service/internal/biz"
	"gateway-service/internal/conf"
	"gateway-service/internal/pkg/requestid"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	grpcmetadata "google.golang.org/grpc/metadata"
)

type GameAppRepo interface {
	GetGameApp(ctx context.Context, gameAppID int64) (*biz.GameApp, error)
}

type gameAppRepo struct {
	client gameappv1.GameAppServiceClient
	logger *log.Helper
}

func NewGameAppRepo(clients *conf.Clients, logger log.Logger) biz.GameAppRepo {
	endpoint := "127.0.0.1:9200"
	if clients != nil && clients.GetGameApp().GetEndpoint() != "" {
		endpoint = clients.GetGameApp().GetEndpoint()
	}
	conn, err := grpc.Dial(
		endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return &gameAppRepo{
		client: gameappv1.NewGameAppServiceClient(conn),
		logger: log.NewHelper(logger),
	}
}

func (r *gameAppRepo) GetGameApp(ctx context.Context, gameAppID int64) (*biz.GameApp, error) {
	if rid := requestid.FromContext(ctx); rid != "" {
		ctx = grpcmetadata.AppendToOutgoingContext(ctx, requestid.MetadataKey, rid)
	}
	resp, err := r.client.GetGameApp(ctx, &gameappv1.GetGameAppRequest{Id: uint64(gameAppID)})
	if err != nil {
		return nil, err
	}
	return &biz.GameApp{
		ID:     uint64(resp.Info.Id),
		Name:   resp.Info.Name,
		AppID:  resp.Info.AppId,
		AppKey: resp.Info.AppKey,
		GameID: resp.Info.GameId,
	}, nil
}

func (r *gameAppRepo) CountGameApps(ctx context.Context) (uint64, error) {
	resp, err := r.client.CountGameApps(ctx, &gameappv1.CountGameAppsRequest{})
	if err != nil {
		return 0, err
	}
	return resp.Total, nil
}

func (r *gameAppRepo) ListGameAppsWithPage(ctx context.Context, page, pageSize int, search *biz.GameAppsSearch) ([]*biz.GameAppInfoV2, int64, error) {
	if rid := requestid.FromContext(ctx); rid != "" {
		ctx = grpcmetadata.AppendToOutgoingContext(ctx, requestid.MetadataKey, rid)
	}
	appSearch := &gameappv1.GameAppsSearch{
		Name:   search.Name,
		TypeOs: search.TypeOs,
	}
	resp, err := r.client.ListGameAppsWithPage(ctx, &gameappv1.ListGameAppsRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   appSearch,
	})
	if err != nil {
		return nil, 0, err
	}
	infos := make([]*biz.GameAppInfoV2, 0, len(resp.Infos))
	for _, info := range resp.Infos {
		infos = append(infos, &biz.GameAppInfoV2{
			Id:     info.Id,
			Name:   info.Name,
			AppId:  info.AppId,
			AppKey: info.AppKey,
			GameId: info.GameId,
		})
	}
	return infos, int64(resp.Total), nil
}
