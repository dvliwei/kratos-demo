/**
 * @Title
 * @Author: liwei
 * @Description:  TODO
 * @File:  gameapp
 * @Version: 1.0.0
 * @Date: 2026/05/29 23:04
 * @Update liwei 2026/5/29 23:04
 */

package data

import (
	"context"
	"gameapp-service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type gameAppRepo struct {
	data *Data
	log  *log.Helper
}

/**
 * @Description:  创建游戏应用仓库
 * @param data 数据
 * @param logger 日志记录器
 * @return biz.GameAppRepo
 * */
func NewGameAppRepo(data *Data, logger log.Logger) biz.GameAppRepo {
	return &gameAppRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

/**
 * @Description:  根据ID查询游戏应用
 */
func (r *gameAppRepo) FindByID(ctx context.Context, id uint64) (*biz.GameApp, error) {
	gameapps := make(map[uint64]*biz.GameApp)
	gameapps[1] = &biz.GameApp{ID: 1, GameID: 1001, AppID: "app_1001", Name: "GameApp1", AppKey: "key_1001"}
	gameapps[2] = &biz.GameApp{ID: 2, GameID: 1002, AppID: "app_1002", Name: "GameApp2", AppKey: "key_1002"}
	gameapp, ok := gameapps[id]
	if !ok {
		return nil, nil
	}
	return gameapp, nil
}
