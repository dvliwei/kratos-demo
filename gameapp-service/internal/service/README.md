# gameapp-service/internal/service

本目录是游戏应用服务的 gRPC 接口层，负责实现 `api/gameapp/v1/gameapp.proto` 中定义的 RPC。

## 职责

- 接收 proto request。
- 调用 `internal/biz` 业务用例。
- 将业务对象转换为 proto response。
- 不直接访问数据库或其他外部资源。

## 文件说明

| 文件 | 说明 |
| --- | --- |
| `gameapp.go` | 游戏应用 RPC 实现 |
| `greeter.go` | Kratos 模板示例接口 |
| `service.go` | Wire ProviderSet |

## 当前 RPC

| 方法 | 说明 |
| --- | --- |
| `GetGameApp` | 根据 ID 查询游戏应用信息 |

## GetGameApp 返回字段

| 字段 | 说明 |
| --- | --- |
| `id` | 游戏应用 ID |
| `game_id` | 游戏 ID |
| `app_id` | 应用标识 |
| `name` | 应用名称 |
| `app_key` | 应用密钥 |

## 新增 RPC 时的步骤

1. 在 `api/gameapp/v1/gameapp.proto` 中新增 RPC 和 message。
2. 执行 `make api` 生成 pb/grpc 代码。
3. 在 `internal/service/gameapp.go` 中实现 RPC 方法。
4. 在 `internal/biz/gameapp.go` 中新增 usecase 方法。
5. 在 `internal/data/gameapp.go` 中实现仓储方法。
6. 执行 `go test ./...` 验证。

## 注意事项

- service 层只做参数转换和响应组装。
- 复杂业务规则应放在 `internal/biz`。
- 数据查询逻辑应放在 `internal/data`。
