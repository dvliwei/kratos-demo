# gameapp-service/internal/data

本目录是游戏应用服务的数据访问层，负责实现 `internal/biz` 定义的游戏应用仓储接口。

## 职责

- 根据游戏应用 ID 查询应用信息。
- 将数据来源转换为业务层 `biz.GameApp`。
- 屏蔽底层数据来源，方便后续从内存数据切换到数据库。

## 文件说明

| 文件 | 说明 |
| --- | --- |
| `data.go` | 数据层 ProviderSet 和共享 Data 结构 |
| `gameapp.go` | 游戏应用查询仓储实现 |
| `greeter.go` | Kratos 模板示例仓储 |

## 当前数据来源

当前 `gameapp.go` 使用内存 map 返回示例数据：

```text
id=1 -> GameApp1
id=2 -> GameApp2
```

这适合演示服务调用链路，但不是生产实现。

## 后续接入数据库建议

如果要改为数据库查询，可以按下面方式调整：

1. 在 `Data` 中保存 GORM DB 实例。
2. 定义 `gameAppModel` 映射数据库表。
3. 在 `FindByID` 中使用 `db.WithContext(ctx).First(...)` 查询。
4. 将 model 转换为 `biz.GameApp`。
5. 对未找到数据返回 Kratos NotFound 错误。

## 注意事项

- data 层不返回 proto 类型，只返回 biz 层领域对象。
- 当前未命中 ID 时返回 `nil, nil`，后续建议改为明确的 NotFound 错误。
- 不要在 data 层处理 HTTP 响应或网关统一响应格式。
