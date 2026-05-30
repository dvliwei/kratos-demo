# user-service 用户服务

`user-service` 是用户领域微服务，主要通过 gRPC 提供用户登录、用户详情查询和分页查询能力。网关 `gateway-service` 会调用本服务完成用户相关操作。

## 服务职责

- 用户邮箱密码登录。
- 登录成功后生成 JWT token。
- 将 token 写入 `users.remember_token` 字段。
- 根据用户 ID 查询用户基础信息。
- 分页查询用户列表，支持按姓名、邮箱搜索。
- 查询用户总数，供网关统计接口聚合使用。

## 默认端口

配置文件：[configs/config.yaml](./configs/config.yaml)

| 类型 | 地址 |
| --- | --- |
| HTTP | `0.0.0.0:8000` |
| gRPC | `0.0.0.0:9100` |
| MySQL | `127.0.0.1:3306/game_sdk_cn` |
| Redis | `127.0.0.1:6379` |

## gRPC 接口

定义文件：[api/user/v1/user.proto](./api/user/v1/user.proto)

| 方法 | 说明 |
| --- | --- |
| `Login` | 用户邮箱密码登录 |
| `GetUser` | 根据 ID 查询用户 |
| `ListUsersWithPage` | 分页查询用户列表 |
| `GetUserTotal` | 查询用户总数 |

### GetUserTotal

`GetUserTotal` 没有请求参数，返回：

```json
{
  "total": 100
}
```

该接口主要供 `gateway-service` 的 `/v1/user_game_app_stats` 聚合接口通过 gRPC 调用。

## 登录说明

登录流程：

```text
LoginRequest(email, password)
  -> 查询 users 表
  -> 校验密码
  -> 生成 JWT token
  -> 事务更新 remember_token
  -> 返回 name 和 token
```

JWT 配置位于 `configs/config.yaml`：

```yaml
auth:
  jwt:
    password: change-me-user-service-jwt-secret
    expire_seconds: 86400
```

`password` 是 JWT 签名密钥，不是用户登录密码。生产环境必须替换。

## 数据表字段

当前 `userModel` 对应表名为 `users`，主要字段：

| 字段 | 说明 |
| --- | --- |
| `id` | 用户 ID |
| `name` | 用户姓名 |
| `email` | 用户邮箱 |
| `password` | 用户密码 |
| `remember_token` | 登录后保存的 token |
| `deleted_at` | GORM 软删除字段 |

注意：当前密码校验是明文比较，后续建议改为 bcrypt/argon2 等安全哈希方式。

## 启动方式

```bash
cd user-service
go run ./cmd/user-service -conf ./configs
```

## 开发命令

运行测试：

```bash
go test ./...
```

生成 API 代码：

```bash
make api
```

生成配置代码：

```bash
make config
```

生成 Wire 注入代码：

```bash
go generate ./...
```

## 目录说明

```text
user-service/
├── api/user/v1/          # 用户服务 proto 和 gRPC 生成代码
├── cmd/user-service/     # 程序入口和 Wire 注入
├── configs/              # 服务配置
├── internal/service/     # gRPC handler
├── internal/biz/         # 用户业务逻辑
├── internal/data/        # GORM 数据访问
├── internal/pkg/jwt/     # JWT 生成与校验
└── internal/server/      # HTTP/gRPC server 初始化
```

## 注意事项

- 修改 `api/user/v1/user.proto` 后，需要执行 `make api`。
- 修改 `internal/conf/conf.proto` 后，需要执行 `make config`。
- 如果本机 `protoc` 报动态库缺失，需要先修复 Protobuf 安装环境。
- `make api` 中已支持 proto3 optional 兼容参数，便于与其他服务的 proto 生成方式保持一致。
- `remember_token` 字段建议在数据库中使用足够长度，例如 `varchar(512)` 或 `text`。
- 当前服务主要作为内部 gRPC 服务被网关调用。
