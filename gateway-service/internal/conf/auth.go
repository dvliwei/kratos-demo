/**
 * @Title
 * @Author: liwei
 * @Description:  TODO
 * @File:  auth
 * @Version: 1.0.0
 * @Date: 2026/05/30 19:18
 * @Update liwei 2026/5/30 19:18
 */

package conf

type RedisConfig struct {
	DB     int `json:"db"`
	UserDB int `json:"user_db"`
}

func (c *RedisConfig) Database() int {
	if c == nil {
		return 0
	}
	if c.UserDB != 0 {
		return c.UserDB
	}
	return c.DB
}

type JWTAuthConfig struct {
	Jwt *JWTConfig `json:"jwt"`
}

type JWTConfig struct {
	Password      string `json:"password"`
	ExpireSeconds int64  `json:"expire_seconds"`
}

type RabbitMQConfig struct {
	Addr  string `json:"addr"`
	Topic string `json:"topic"`
}
