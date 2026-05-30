# gateway-service/internal/service

本目录是网关服务的接口适配层，负责实现由 proto 生成的 HTTP/gRPC handler 接口。

## 职责

- 接收网关 proto request。
- 做简单参数提取和对象转换。
- 调用 `internal/biz` 完成业务编排。
- 将业务对象转换为网关 proto response。
- 不直接访问数据库，也不直接写下游 gRPC client 调用逻辑。

## 文件说明

| 文件 | 说明 |
| --- | --- |
| `gateway.go` | 网关核心接口实现 |
| `greeter.go` | Kratos 模板示例接口 |
| `service.go` | Wire ProviderSet |

## 当前接口

`gateway.go` 中主要实现：

- `GetGatewayInfo`：查询用户详情
- `ListUsersWithPage`：分页查询用户列表
- `GetGameApp`：查询游戏应用

如果 proto 中已经声明了接口，但这里没有实现对应方法，请求会落到 `UnimplementedGatewayServiceServer`，返回类似：

```json
{
  "code": 501,
  "message": "method Xxx not implemented"
}
```

## 新增接口时的检查清单

1. `api/gateway/v1/gateway.proto` 是否定义 RPC 和 HTTP 路由。
2. 是否执行过 `make api` 生成代码。
3. `gateway.go` 是否实现了对应方法。
4. `internal/biz` 是否有业务方法。
5. `internal/data` 是否有下游调用方法。
6. `internal/server/http.go` 是否注册了 `GatewayServiceHTTPServer`。

## 注意事项

- 这一层只做接口转换，不放复杂业务逻辑。
- 外部 HTTP 的统一响应包装不在这里处理，而是在 `internal/server/response.go` 中处理。
- 如果方法已在 proto 中生成，但 service 未实现，因为结构体嵌入了 `UnimplementedGatewayServiceServer`，编译可能不报错，但运行时会返回 501。
