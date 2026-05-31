# RabbitMQ 本地环境

本目录提供本地联调用的 RabbitMQ Docker Compose 配置。`gateway-service` 会把访问日志发布到 RabbitMQ，`accesslog-service` 会消费 `access_log` 队列。

## 默认配置

| 项目 | 值 |
| --- | --- |
| AMQP 地址 | `amqp://admin:admin123@127.0.0.1:5672` |
| 管理后台 | `http://localhost:15672` |
| 用户名 | `admin` |
| 密码 | `admin123` |
| 访问日志队列 | `access_log` |

## 启动

```bash
./start-rabbitmq.sh
```

脚本会创建本地数据目录并执行：

```bash
docker compose up -d
```

## 停止

```bash
docker compose down
```

## 联调检查

- 先启动 RabbitMQ，再启动 `gateway-service` 和 `accesslog-service`。
- 在管理后台确认 `access_log` 队列存在。
- 如果消息一直停留在 `Ready`，检查 `accesslog-service` 是否启动，以及队列名是否为 `access_log`。
- 如果消息进入 `Unacked`，检查消费者是否处理后成功 `Ack`。
