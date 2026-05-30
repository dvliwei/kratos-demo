/**
 * @Title
 * @Author: liwei
 * @Description:  TODO
 * @File:  stats
 * @Version: 1.0.0
 * @Date: 2026/05/30 23:34
 * @Update liwei 2026/5/30 23:34
 */

package data

import (
	"context"
	gameappv1 "gameapp-service/api/gameapp/v1"
	"gateway-service/internal/biz"
	"gateway-service/internal/pkg/requestid"
	userv1 "user-service/api/user/v1"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type StatsRepo struct {
	client     userv1.UserClient
	gameClient gameappv1.GameAppServiceClient
	logger     *log.Helper
}

func NewStatsRepo(logger log.Logger) *StatsRepo {
	conn, err := grpc.Dial(
		"127.0.0.1:9100",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	gameConn, err := grpc.Dial(
		"127.0.0.1:9200",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return &StatsRepo{
		client:     userv1.NewUserClient(conn),
		gameClient: gameappv1.NewGameAppServiceClient(gameConn),
		logger:     log.NewHelper(logger),
	}
}

func (a *StatsRepo) GetUserGameAppStats(ctx context.Context) (*biz.GetUserGameAppStatsReply, error) {
	if rid := requestid.FromContext(ctx); rid != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, requestid.MetadataKey, rid)
	}
	result := &biz.GetUserGameAppStatsReply{}
	resp, err := a.client.GetUserTotal(ctx, &userv1.GetUserTotalRequest{})
	if err == nil {
		result.TotalUsers = int64(resp.GetTotal())
	}
	totalGameApps, err := a.gameClient.CountGameApps(ctx, &gameappv1.CountGameAppsRequest{})
	if err == nil {
		result.TotalGameApps = int64(totalGameApps.GetTotal())
	}
	return result, nil
}
