/**
 * @Title
 * @Author: liwei
 * @Description:  TODO
 * @File:  gameapp
 * @Version: 1.0.0
 * @Date: 2026/05/29 23:06
 * @Update liwei 2026/5/29 23:06
 */

package biz

import "context"

// uint64 id = 1;
// int64 game_id = 2;
// string app_id = 3;
// string name= 4;
// string app_key = 5;
type GameApp struct {
	ID     uint64 `json:"id"`
	GameID int64  `json:"game_id"`
	AppID  string `json:"app_id"`
	Name   string `json:"name"`
	AppKey string `json:"app_key"`
}

type GameAppRepo interface {
	FindByID(ctx context.Context, id uint64) (*GameApp, error)
}
type GameAppUseCase struct {
	repo GameAppRepo
}

func NewGameAppUseCase(repo GameAppRepo) *GameAppUseCase {
	return &GameAppUseCase{repo: repo}
}
func (u *GameAppUseCase) GetByIDApp(ctx context.Context, id uint64) (*GameApp, error) {
	return u.repo.FindByID(ctx, id)
}
