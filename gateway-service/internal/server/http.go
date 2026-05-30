package server

import (
	gatewayv1 "gateway-service/api/gateway/v1"
	v1 "gateway-service/api/helloworld/v1"
	"gateway-service/internal/conf"
	"gateway-service/internal/data"
	"gateway-service/internal/pkg/jwt"
	"gateway-service/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, auth *conf.JWTAuthConfig, data *data.Data, greeter *service.GreeterService,
	gateway *service.GatewayService,
	logger log.Logger) *http.Server {
	jwtManager, err := newJWTManager(auth)
	if err != nil {
		panic(err)
	}
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
		http.Filter(
			requestIDFilter,
			authMiddleware(jwtManager, data.Redis()),
		),
		http.ResponseEncoder(unifiedResponseEncoder),
		http.ErrorEncoder(unifiedErrorEncoder),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	// 注册网关服务
	gatewayv1.RegisterGatewayServiceHTTPServer(srv, gateway)
	return srv
}

func newJWTManager(auth *conf.JWTAuthConfig) (*jwt.Manager, error) {
	var jwtConf *conf.JWTConfig
	if auth != nil {
		jwtConf = auth.Jwt
	}
	if jwtConf == nil {
		jwtConf = &conf.JWTConfig{}
	}
	return jwt.NewManager(jwtConf.Password, jwt.DurationFromSeconds(jwtConf.ExpireSeconds))
}
