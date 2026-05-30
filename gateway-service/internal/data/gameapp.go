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

func NewGameAppRepo(logger log.Logger) biz.GameAppRepo {
	conn, err := grpc.Dial(
		"127.0.0.1:9200",
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
