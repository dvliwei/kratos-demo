# gateway-service 网关服务

`gateway-service` 是项目的对外 HTTP 网关，负责接收外部请求、统一响应格式、生成和透传 `request_id`，并通过 gRPC 调用内部的 `user-service` 和 `gameapp-service`。
同时，网关会采集 HTTP 访问日志并异步发布到 RabbitMQ，供 `accesslog-service` 消费。

## 服务职责

- 对外提供 HTTP API。
- 聚合用户服务和游戏应用服务。
- 提供用户总数与游戏应用总数的聚合统计接口。
- 提供游戏应用分页查询接口，并支持 `type_os` 可选筛选。
- 统一包装 HTTP 响应，返回 `code`、`message`、`request_id`、`server_time`、`data`。
- 从请求头读取或生成 `X-Request-Id`。
- 调用下游 gRPC 服务时透传 `x-request-id`。
- 采集请求方法、路径、请求体、响应体、状态码、耗时和 `request_id`。
- 将访问日志发布到 RabbitMQ `access_log` 队列。

## 默认端口

配置文件：[configs/config.yaml](./configs/config.yaml)

| 类型 | 地址 |
| --- | --- |
| HTTP | `0.0.0.0:8080` |
| gRPC | `0.0.0.0:9000` |
| user-service | `127.0.0.1:9100` |
| gameapp-service | `127.0.0.1:9200` |
| RabbitMQ | `127.0.0.1:5672` |

下游服务地址从 [configs/config.yaml](./configs/config.yaml) 的 `clients` 节读取：

```yaml
clients:
  user:
    endpoint: 127.0.0.1:9100
  game_app:
    endpoint: 127.0.0.1:9200
```

如果配置缺失，代码会回退到本地默认地址，便于本地开发。

访问日志 RabbitMQ 配置：

```yaml
rabbitmq:
  addr: amqp://admin:admin123@127.0.0.1:5672
  topic: access_log
```

`topic` 当前作为队列名使用。需要和 `accesslog-service` 的 `rabbitmq.queue` 保持一致。

## 访问日志

网关通过 HTTP filter 记录每次请求，最多保留请求体和响应体前 `4096` 字节，然后异步写入 RabbitMQ。

消息体是 JSON，核心字段包括：

```json
{
  "method": "POST",
  "path": "/v1/users",
  "query": "",
  "request_id": "demo-request-id",
  "status": 200,
  "cost_ms": 12,
  "request_body": "{}",
  "response_body": "{}",
  "created_at": "2026-06-01T01:51:34+08:00"
}
```

如果 RabbitMQ 不可用，网关启动时会返回错误；本地联调请先启动 `rabbitmq` 目录下的 Docker Compose。

## 接口列表

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| `POST` | `/v1/login` | 用户邮箱密码登录 |
| `GET` | `/v1/users/{id}` | 查询用户详情 |
| `POST` | `/v1/users` | 分页查询用户列表 |
| `GET` | `/v1/game_app/{id}` | 查询游戏应用详情 |
| `POST` | `/v1/game_apps` | 分页查询游戏应用列表 |
| `GET` | `/v1/user_game_app_stats` | 查询用户总数和游戏应用总数 |

## 响应格式

成功响应：

```json
{
  "code": 0,
  "message": "ok",
  "request_id": "demo-request-id",
  "server_time": 1780116826198,
  "data": {}
}
```

错误响应：

```json
{
  "code": 401,
  "message": "invalid email or password",
  "request_id": "demo-request-id",
  "server_time": 1780116826198,
  "data": null
}
```

## 请求示例

登录：

```bash
curl -X POST http://127.0.0.1:8080/v1/login \
  -H 'Content-Type: application/json' \
  -H 'X-Request-Id: demo-login-request' \
  -d '{"email":"test@example.com","password":"123456"}'
```

分页查询用户：

```bash
curl -X POST http://127.0.0.1:8080/v1/users \
  -H 'Content-Type: application/json' \
  -d '{"page":1,"page_size":10,"search":{"name":"test","email":""}}'
```

查询用户详情：

```bash
curl http://127.0.0.1:8080/v1/users/1
```

查询游戏应用：

```bash
curl http://127.0.0.1:8080/v1/game_app/1
```

分页查询游戏应用：

```bash
curl -X POST http://127.0.0.1:8080/v1/game_apps \
  -H 'Content-Type: application/json' \
  -d '{"page":1,"page_size":10,"search":{"name":"demo","type_os":0}}'
```

`type_os` 是可选字段：

- 不传：不按平台筛选。
- 传 `0`：查询 `type_os = 0`。

查询用户和游戏应用统计：

```bash
curl http://127.0.0.1:8080/v1/user_game_app_stats
```

响应中的 `total_users` 来自 `user-service.GetUserTotal`，`total_game_apps` 来自 `gameapp-service.CountGameApps`。

## 启动方式

请先启动下游服务：

```bash
cd ../rabbitmq
./start-rabbitmq.sh
```

```bash
cd ../user-service
go run ./cmd/user-service -conf ./configs
```

```bash
cd ../gameapp-service
go run ./cmd/gameapp-service -conf ./configs
```

再启动网关：

```bash
cd ../gateway-service
go run ./cmd/gateway-service -conf ./configs
```

## 开发命令

```bash
go test ./...
```

```bash
make api
make config
go generate ./...
```

`gateway.proto` 中的 `GameAppsSearch.type_os` 使用 proto3 optional。执行 `make api` 时需要 `Makefile` 包含：

```bash
--experimental_allow_proto3_optional
```

## 目录说明

```text
gateway-service/
├── api/gateway/v1/       # 网关 proto、HTTP/gRPC 生成代码
├── cmd/gateway-service/  # 程序入口和 Wire 注入
├── configs/              # 配置文件
├── internal/service/     # HTTP/gRPC handler
├── internal/biz/         # 业务编排
├── internal/data/        # 下游 gRPC client 调用
└── internal/server/      # HTTP/gRPC server 初始化和统一响应
```

## 注意事项

- 修改 `api/gateway/v1/gateway.proto` 后，需要执行 `make api`。
- 如果新增网关接口，需要同时补齐 proto、service、biz、data 调用链。
- 游戏应用分页查询由网关转发到 `gameapp-service.ListGameAppsWithPage`。
- 聚合统计接口会分别调用 `user-service.GetUserTotal` 和 `gameapp-service.CountGameApps`。
- 下游服务地址从 `configs/config.yaml` 的 `clients` 节读取，修改端口或远程地址时优先改配置。
- RabbitMQ 队列名必须和 `accesslog-service` 保持一致，当前为 `access_log`。
- 如果 HTTP 返回 `method xxx not implemented`，通常是 proto 已生成路由，但 `internal/service` 中没有实现对应方法。
