# gateway-service/internal/data

本目录是网关服务的数据访问层。这里的“数据访问”不是直接访问数据库，而是封装对下游微服务的 gRPC 调用。

## 职责

- 创建并保存下游服务 gRPC client。
- 将网关业务层的领域对象转换为下游 proto request。
- 调用下游服务并把 proto response 转换为网关业务对象。
- 透传 `request_id` 到下游服务 metadata。
- 屏蔽下游服务的调用细节，让 `internal/biz` 只依赖接口。

## 文件说明

| 文件 | 说明 |
| --- | --- |
| `data.go` | 数据层 ProviderSet 和共享 Data 结构 |
| `user.go` | 调用 `user-service`，包括用户详情、分页列表、登录 |
| `gameapp.go` | 调用 `gameapp-service` 查询游戏应用 |
| `greeter.go` | Kratos 模板示例仓储 |

## 当前下游调用

`user.go` 中封装了：

- `GetUser`
- `ListUsersWithPage`
- `Login`

`gameapp.go` 中封装了：

- `GetGameApp`

## request_id 透传

网关 HTTP 层会生成或读取 `X-Request-Id`，data 层调用下游 gRPC 时会将它写入 metadata：

```text
x-request-id: <request_id>
```

这样下游服务可以在日志和链路追踪中关联同一次请求。

## 新增下游接口时的步骤

1. 在对应下游服务 proto 中定义 RPC。
2. 重新生成下游服务 API 代码。
3. 在网关 proto 中定义对外接口。
4. 重新生成网关 API 代码。
5. 在 `internal/biz` 的 repo 接口中增加方法。
6. 在本目录对应文件中实现 gRPC 调用。
7. 在 `internal/service` 中实现 HTTP/gRPC handler。

## 注意事项

- 当前下游地址存在硬编码调用点，例如 `127.0.0.1:9100` 和 `127.0.0.1:9200`，后续建议统一改为读取配置。
- 这里不写业务规则，只负责调用下游和对象转换。
- 不要在 data 层直接拼 HTTP 响应结构，统一响应由 `internal/server` 处理。
