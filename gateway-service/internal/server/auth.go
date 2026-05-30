package server

import (
	"context"
	"fmt"
	"net/http"

	"gateway-service/internal/pkg/jwt"

	kratoserrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/redis/go-redis/v9"
)

const userTokenHeader = "user-token"

type authUserContextKey struct{}

type authUser struct {
	ID   uint64
	Name string
}

func authMiddleware(jwtManager *jwt.Manager, redisClient *redis.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if skipAuth(r) {
				next.ServeHTTP(w, r)
				return
			}

			token := r.Header.Get(userTokenHeader)
			if token == "" {
				unifiedErrorEncoder(w, r, kratoserrors.Unauthorized("UNAUTHORIZED", "missing user token"))
				return
			}

			claims, err := jwtManager.ValidateToken(token)
			if err != nil {
				unifiedErrorEncoder(w, r, kratoserrors.Unauthorized("UNAUTHORIZED", "invalid user token"))
				return
			}
			fmt.Println(userTokenKey(claims.UserID))
			storedToken, err := redisClient.Get(r.Context(), userTokenKey(claims.UserID)).Result()
			if err != nil {
				if err == redis.Nil {
					unifiedErrorEncoder(w, r, kratoserrors.Unauthorized("UNAUTHORIZED", "user token expired"))
					return
				}
				unifiedErrorEncoder(w, r, kratoserrors.InternalServer("REDIS_ERROR", "check user token failed"))
				return
			}
			if storedToken != token {
				unifiedErrorEncoder(w, r, kratoserrors.Unauthorized("UNAUTHORIZED", "user token has been replaced"))
				return
			}

			ctx := context.WithValue(r.Context(), authUserContextKey{}, authUser{
				ID:   claims.UserID,
				Name: claims.Name,
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func skipAuth(r *http.Request) bool {
	return r.Method == http.MethodPost && r.URL.Path == "/v1/login"
}

func userTokenKey(userID uint64) string {
	return fmt.Sprintf("user:token:%d", userID)
}
