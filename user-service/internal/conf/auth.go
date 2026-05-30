package conf

type JWTAuthConfig struct {
	Jwt *JWTConfig `json:"jwt"`
}

type JWTConfig struct {
	Password      string `json:"password"`
	ExpireSeconds int64  `json:"expire_seconds"`
}

type RedisConfig struct {
	DB int `json:"db"`
}
