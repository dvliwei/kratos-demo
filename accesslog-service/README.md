# accesslog-service 访问日志消费服务

`accesslog-service` 是访问日志领域微服务，当前负责消费 RabbitMQ 中的网关访问日志。现阶段消费后先打印日志并手动确认消息，后续可以在这里扩展 ClickHouse、OpenSearch、MongoDB 或 PostgreSQL 存储。

## 服务职责

- 连接 RabbitMQ。
- 声明并消费 `access_log` 队列。
- 打印访问日志消息体。
- 在 `auto_ack: false` 时消费成功后手动 `Ack`。
- 跟随 Kratos 应用生命周期启动和停止消费者。

## 默认端口

配置文件：[configs/config.yaml](./configs/config.yaml)

| 类型 | 地址 |
| --- | --- |
| HTTP | `0.0.0.0:8003` |
| gRPC | `0.0.0.0:9300` |
| RabbitMQ | `127.0.0.1:5672` |
| Queue | `access_log` |

## 消息来源

`gateway-service` 会把 HTTP 访问日志发布到 RabbitMQ 默认 exchange，并使用队列名作为 routing key：

```text
gateway-service -> RabbitMQ queue: access_log -> accesslog-service
```

消息体是 JSON，字段来自网关的访问日志结构：

```json
{
  "method": "POST",
  "path": "/v1/users",
  "query": "",
  "header": {},
  "user_token": "",
  "request_id": "demo-request-id",
  "status": 200,
  "cost_ms": 12,
  "request_body": "{}",
  "response_body": "{}",
  "created_at": "2026-06-01T01:51:34+08:00"
}
```

## RabbitMQ 配置

```yaml
rabbitmq:
  url: amqp://admin:admin123@127.0.0.1:5672/
  queue: access_log
  consumer_tag: accesslog-service
  auto_ack: false
  prefetch_count: 1
```

字段说明：

| 字段 | 说明 |
| --- | --- |
| `url` | RabbitMQ AMQP 连接地址 |
| `queue` | 消费队列名，必须与 `gateway-service.rabbitmq.topic` 一致 |
| `consumer_tag` | RabbitMQ consumer 标识 |
| `auto_ack` | 是否自动确认消息，当前建议为 `false` |
| `prefetch_count` | 单次预取消息数量，用于控制消费并发和积压 |

## 启动方式

先启动 RabbitMQ：

```bash
cd ../rabbitmq
./start-rabbitmq.sh
```

再启动本服务：

```bash
cd ../accesslog-service
go run ./cmd/accesslog-service -conf ./configs
```

或使用 Kratos：

```bash
kratos run
```

启动成功后会看到类似日志：

```text
rabbitmq consumer started: queue=access_log auto_ack=false
```

收到消息后会打印：

```text
received rabbitmq message: exchange= routing_key=access_log body=...
```

## 开发命令

运行测试：

```bash
GOWORK=off go test ./...
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
accesslog-service/
├── api/                     # 示例 proto 和生成代码
├── cmd/accesslog-service/   # 程序入口和 Wire 注入
├── configs/                 # 服务配置
├── internal/server/         # HTTP/gRPC server 与 RabbitMQ consumer
├── internal/service/        # 示例 service
├── internal/biz/            # 示例 biz
└── internal/data/           # 示例 data
```

## 注意事项

- `gateway-service.rabbitmq.topic` 和本服务 `rabbitmq.queue` 必须保持一致，当前为 `access_log`。
- 当前消费者只打印消息；接入数据库时建议保留手动 `Ack`，写入成功后再确认。
- 如果服务启动失败并提示 RabbitMQ 连接错误，先确认 `rabbitmq/start-rabbitmq.sh` 已启动成功。
- 如果队列里数据没有减少，优先检查队列名、连接账号、消费者启动日志和 RabbitMQ 管理后台的 `Ready/Unacked` 数量。
