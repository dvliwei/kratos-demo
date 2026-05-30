# Kratos Demo 微服务项目

这是一个基于 Go Kratos 的多服务示例项目，当前包含网关服务、用户服务和游戏应用服务。项目通过 `go.work` 管理多个 Go module，适合用来练习 Kratos 的 HTTP/gRPC、服务分层、网关转发、统一响应、JWT 登录和数据库访问。

## 项目结构

```text
kratos-demo/
├── gateway-service/   # 对外 HTTP 网关，聚合 user-service 与 gameapp-service
├── user-service/      # 用户服务，提供登录、用户查询、分页查询
├── gameapp-service/   # 游戏应用服务，提供游戏应用查询
├── vendor/            # 依赖缓存
├── go.work            # 多模块工作区
└── README.md          # 项目总说明
```

## 服务职责

| 服务 | 默认 HTTP | 默认 gRPC | 说明 |
| --- | --- | --- | --- |
| gateway-service | `:8080` | `:9000` | 对外统一入口，负责 HTTP 路由、统一响应包装、request_id 生成与透传 |
| user-service | `:8000` | `:9100` | 用户登录、用户详情、用户分页查询 |
| gameapp-service | `:8088` | `:9200` | 游戏应用详情查询 |

## 调用关系

```text
前端/调用方
  |
  | HTTP
  v
gateway-service
  |-- gRPC --> user-service
  |-- gRPC --> gameapp-service
```

`gateway-service` 是推荐的外部访问入口。内部微服务主要通过 gRPC 被网关调用。

## 统一响应格式

`gateway-service` 已统一包装 HTTP 响应：

```json
{
  "code": 0,
  "message": "ok",
  "request_id": "a6d79c9c245145a24f082b73768b7618",
  "server_time": 1780116826198,
  "data": {}
}
```

错误响应也会带上 `request_id` 和 `server_time`：

```json
{
  "code": 401,
  "message": "invalid email or password",
  "request_id": "a6d79c9c245145a24f082b73768b7618",
  "server_time": 1780116826198,
  "data": null
}
```

如果请求头中传入 `X-Request-Id`，网关会沿用；否则自动生成。网关调用下游 gRPC 服务时会通过 metadata 透传 `x-request-id`。

## 主要接口

通过 `gateway-service` 访问：

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| `POST` | `/v1/login` | 用户邮箱密码登录 |
| `GET` | `/v1/users/{id}` | 查询用户详情 |
| `POST` | `/v1/users` | 分页查询用户列表 |
| `GET` | `/v1/game_app/{id}` | 查询游戏应用详情 |

示例：

```bash
curl -X POST http://127.0.0.1:8080/v1/users \
  -H 'Content-Type: application/json' \
  -H 'X-Request-Id: demo-request-id' \
  -d '{"page":1,"page_size":10,"search":{"name":"test"}}'
```

## 本地启动

建议按依赖顺序启动：

```bash
cd user-service
go run ./cmd/user-service -conf ./configs
```

```bash
cd gameapp-service
go run ./cmd/gameapp-service -conf ./configs
```

```bash
cd gateway-service
go run ./cmd/gateway-service -conf ./configs
```

## 配置说明

每个服务的配置文件位于各自的 `configs/config.yaml`。

公共配置：

```yaml
server:
  http:
    addr: 0.0.0.0:8080
    timeout: 1s
  grpc:
    addr: 0.0.0.0:9000
    timeout: 1s
data:
  database:
    driver: mysql
    source: root:root123456@tcp(127.0.0.1:3306)/game_sdk_cn?parseTime=True&loc=Local
```

`user-service` 额外包含 JWT 配置：

```yaml
auth:
  jwt:
    password: change-me-user-service-jwt-secret
    expire_seconds: 86400
```

生产环境请替换数据库连接和 JWT 密钥。

## 常用命令

在单个服务目录下执行：

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

## 注意事项

- 当前项目使用 `go.work` 管理多个服务模块。
- `gateway-service` 是外部 HTTP 的统一入口，业务微服务不建议直接暴露给前端。
- 修改 `.proto` 后需要重新执行对应服务的 `make api` 或 `make config`。
- 如果 `make config` 报 `protoc` 动态库缺失，需要先修复本机 Protobuf 安装环境。
- `gameapp-service` 当前游戏应用数据是内存示例数据，不是数据库查询。
- `gateway-service` 的下游 gRPC 地址当前在代码中仍有硬编码调用点，后续可统一改为读取 `configs/config.yaml` 的 client 配置。
