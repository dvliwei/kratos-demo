/**
 * @Title
 * @Author: liwei
 * @Description:  TODO
 * @File:  stats
 * @Version: 1.0.0
 * @Date: 2026/05/30 23:35
 * @Update liwei 2026/5/30 23:35
 */

package biz

import "context"

type GetUserGameAppStatsReply struct {
	TotalGameApps int64
	TotalUsers    int64
}
type StatsRepo interface {
	GetUserGameAppStats(ctx context.Context) (*GetUserGameAppStatsReply, error)
}

type StatsUseCase struct {
	repo StatsRepo
}

func NewStatsRepo(repo StatsRepo) *StatsUseCase {
	return &StatsUseCase{repo: repo}
}
func (a *StatsUseCase) GetUserGameAppStats(ctx context.Context) (*GetUserGameAppStatsReply, error) {
	return a.repo.GetUserGameAppStats(ctx)
}
