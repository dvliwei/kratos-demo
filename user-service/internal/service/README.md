# user-service/internal/service

本目录是用户服务的 gRPC 接口层，负责实现 `api/user/v1/user.proto` 中定义的 RPC。

## 职责

- 接收 proto request。
- 调用 `internal/biz` 业务用例。
- 将业务对象转换为 proto response。
- 不直接访问数据库。
- 不实现复杂业务规则。

## 文件说明

| 文件 | 说明 |
| --- | --- |
| `user.go` | 用户服务 RPC 实现 |
| `greeter.go` | Kratos 模板示例接口 |
| `service.go` | Wire ProviderSet |

## 当前 RPC

| 方法 | 说明 |
| --- | --- |
| `Login` | 用户邮箱密码登录 |
| `GetUser` | 根据用户 ID 查询用户 |
| `ListUsersWithPage` | 分页查询用户列表 |

## 类型转换

service 层负责 proto 类型和 biz 类型之间的转换。例如分页查询中：

- proto 入参：`*v1.SearchUser`
- biz 入参：`*biz.SearchUser`

这两个类型名字相同，但包不同，不能直接传递，需要显式转换。

## 新增 RPC 时的步骤

1. 在 `api/user/v1/user.proto` 中新增 RPC 和 message。
2. 执行 `make api` 生成 pb/grpc 代码。
3. 在 `internal/service/user.go` 中实现 RPC 方法。
4. 在 `internal/biz/user.go` 中新增 usecase 方法。
5. 在 `internal/data/user.go` 中实现仓储方法。
6. 执行 `go test ./...` 验证。

## 注意事项

- 如果 proto 已生成 RPC，但 service 没实现，运行时会返回 `method Xxx not implemented`。
- service 层不要写 SQL，也不要直接操作 GORM。
- service 层应尽量保持薄，只负责参数转换和响应组装。
