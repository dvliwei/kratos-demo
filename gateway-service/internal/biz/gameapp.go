/**
 * @Title
 * @Author: liwei
 * @Description:  TODO
 * @File:  gameapp
 * @Version: 1.0.0
 * @Date: 2026/05/29 23:44
 * @Update liwei 2026/5/29 23:44
 */

package biz

import (
	"context"
)

type GameApp struct {
	ID     uint64 `json:"id"`
	GameID int64  `json:"game_id"`
	AppID  string `json:"app_id"`
	Name   string `json:"name"`
	AppKey string `json:"app_key"`
}

type GameAppRepo interface {
	GetGameApp(ctx context.Context, gameAppID int64) (*GameApp, error)
	CountGameApps(ctx context.Context) (uint64, error)
	ListGameAppsWithPage(ctx context.Context, page, pageSize int, search *GameAppsSearch) ([]*GameAppInfoV2, int64, error)
}
type GameAppUseCase struct {
	repo GameAppRepo
}

func NewGameAppUseCase(repo GameAppRepo) *GameAppUseCase {
	return &GameAppUseCase{repo: repo}
}
func (a *GameAppUseCase) GetGameApp(ctx context.Context, gameAppID int64) (*GameApp, error) {
	return a.repo.GetGameApp(ctx, gameAppID)
}

func (a *GameAppUseCase) CountGameApps(ctx context.Context) (uint64, error) {
	return a.repo.CountGameApps(ctx)
}

type GameAppsSearch struct {
	Name   string `json:"name"`
	TypeOs *int32 `json:"type_os"`
}
type ListGameAppsRequest struct {
	PageNum  uint64         `json:"page_num"`
	PageSize uint64         `json:"page_size"`
	Search   GameAppsSearch `json:"search"`
}
type GameAppInfoV2 struct {
	Id     uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	GameId int64  `protobuf:"varint,2,opt,name=game_id,json=gameId,proto3" json:"game_id,omitempty"`
	AppId  string `protobuf:"bytes,3,opt,name=app_id,json=appId,proto3" json:"app_id,omitempty"`
	Name   string `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	AppKey string `protobuf:"bytes,5,opt,name=app_key,json=appKey,proto3" json:"app_key,omitempty"`
}

func (a *GameAppUseCase) ListGameAppsWithPage(ctx context.Context, page, pageSize int, search *GameAppsSearch) ([]*GameAppInfoV2, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 30
	}
	return a.repo.ListGameAppsWithPage(ctx, page, pageSize, search)
}
