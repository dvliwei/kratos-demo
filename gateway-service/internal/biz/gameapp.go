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

import "context"

type GameApp struct {
	ID     uint64 `json:"id"`
	GameID int64  `json:"game_id"`
	AppID  string `json:"app_id"`
	Name   string `json:"name"`
	AppKey string `json:"app_key"`
}

type GameAppRepo interface {
	GetGameApp(ctx context.Context, gameAppID int64) (*GameApp, error)
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
