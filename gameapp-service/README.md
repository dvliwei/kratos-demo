# gameapp-service 游戏应用服务

`gameapp-service` 是游戏应用领域微服务，提供游戏应用详情、分页查询和总数统计的 gRPC 能力。网关 `gateway-service` 会调用本服务并对外暴露 HTTP 接口。

## 服务职责

- 根据游戏应用 ID 查询应用信息。
- 分页查询游戏应用列表。
- 支持按应用名称和 `type_os` 筛选。
- 统计游戏应用总数。
- 返回应用 ID、游戏 ID、应用标识、名称和 AppKey。
- 作为内部服务被 `gateway-service` 聚合调用。

## 默认端口

配置文件：[configs/config.yaml](./configs/config.yaml)

| 类型 | 地址 |
| --- | --- |
| HTTP | `0.0.0.0:8088` |
| gRPC | `0.0.0.0:9200` |
| MySQL | `127.0.0.1:3306/game_sdk_cn` |
| Redis | `127.0.0.1:6379` |

## gRPC 接口

定义文件：[api/gameapp/v1/gameapp.proto](./api/gameapp/v1/gameapp.proto)

| 方法 | 说明 |
| --- | --- |
| `GetGameApp` | 根据 ID 查询游戏应用信息 |
| `ListGameAppsWithPage` | 分页查询游戏应用列表 |
| `CountGameApps` | 查询游戏应用总数 |

返回字段：

| 字段 | 说明 |
| --- | --- |
| `id` | 游戏应用 ID |
| `game_id` | 游戏 ID |
| `app_id` | 应用标识 |
| `name` | 应用名称 |
| `app_key` | 应用密钥 |

## 分页查询

`ListGameAppsWithPage` 请求参数：

```json
{
  "page": 1,
  "page_size": 10,
  "search": {
    "name": "demo",
    "type_os": 0
  }
}
```

`type_os` 是 proto3 `optional int32` 字段：

- 不传 `type_os`：不按平台筛选。
- 传 `type_os: 0`：明确查询 `type_os = 0`。
- 传其他值：按对应平台类型筛选。

查询实现中会先 `Count` 获取总数，再按 `id DESC`、`Offset`、`Limit` 查询当前页。

## 数据来源

当前 `internal/data/gameapp.go` 通过 GORM 查询 MySQL 表 `tab_game_app`。

主要字段：

| 字段 | 说明 |
| --- | --- |
| `id` | 游戏应用 ID |
| `game_id` | 游戏 ID |
| `app_id` | 应用标识，唯一 |
| `name` | 应用名称 |
| `app_key` | 应用 Key |
| `type_os` | 平台类型 |
| `pay_status` | 支付状态 |
| `created_at` | 创建时间 |
| `updated_at` | 更新时间 |

## 启动方式

```bash
cd gameapp-service
go run ./cmd/gameapp-service -conf ./configs
```

通过网关访问：

```bash
curl http://127.0.0.1:8080/v1/game_app/1
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

因为 `GameAppsSearch.type_os` 使用了 proto3 optional，`Makefile` 的 `make api` 已包含：

```bash
--experimental_allow_proto3_optional
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
gameapp-service/
├── api/gameapp/v1/       # 游戏应用服务 proto 和 gRPC 生成代码
├── cmd/gameapp-service/  # 程序入口和 Wire 注入
├── configs/              # 服务配置
├── internal/service/     # gRPC handler
├── internal/biz/         # 游戏应用业务逻辑
├── internal/data/        # GORM 数据访问
└── internal/server/      # HTTP/gRPC server 初始化
```

## 注意事项

- 修改 `api/gameapp/v1/gameapp.proto` 后，需要执行 `make api`。
- 当前服务查询 MySQL `tab_game_app` 表，需要确认数据库连接和表结构存在。
- 如果 `gateway-service` 查询游戏应用失败，先确认本服务 gRPC 端口 `9200` 已启动。
