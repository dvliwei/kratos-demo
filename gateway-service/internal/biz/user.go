/**
* @Title
* @Author: liwei
* @Description:  负责“业务编排”
*  比如校验 id、决定调用哪个下游服务、组装业务结
* @File:  user.go
* @Version: 1.0.0
* @Date: 2026/05/29 21:18
* @Update liwei 2026/5/29 21:18
 */

package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
)

type User struct {
	ID    int64
	Name  string
	Email string
}

type SearchUser struct {
	Name  string
	Email string
}

type UserRepo interface {
	GetUser(ctx context.Context, userID int64) (*User, error)
	ListUsersWithPage(ctx context.Context, page, pageSize int, search *SearchUser) ([]*User, int64, error)
	Login(ctx context.Context, email string, password string) (*string, *string, error)
}
type UserUseCase struct {
	repo UserRepo
}

func NewUserUsecase(repo UserRepo) *UserUseCase {
	return &UserUseCase{repo: repo}
}
func (a *UserUseCase) GetUser(ctx context.Context, userID int64) (*User, error) {
	if userID <= 0 {
		return nil, errors.BadRequest("USER_ID_INVALID", "user id is invalid")
	}
	return a.repo.GetUser(ctx, userID)
}

func (a *UserUseCase) ListUsersWithPage(ctx context.Context, page, pageSize int, search *SearchUser) ([]*User, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return a.repo.ListUsersWithPage(ctx, page, pageSize, search)
}

func (a *UserUseCase) Login(ctx context.Context, email string, password string) (*string, *string, error) {
	if email == "" || password == "" {
		return nil, nil, errors.BadRequest("USER_EMAIL_INVALID", "email or password is empty")
	}
	return a.repo.Login(ctx, email, password)
}
