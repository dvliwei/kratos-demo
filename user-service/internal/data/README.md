# user-service/internal/data

本目录是用户服务的数据访问层，负责和 MySQL 等外部资源交互，并实现 `internal/biz` 定义的仓储接口。

## 职责

- 初始化和使用 GORM 访问 `users` 表。
- 实现用户登录、用户详情查询、分页查询。
- 登录成功后生成 JWT token 并更新 `remember_token`。
- 将数据库模型转换为业务层领域模型。

## 文件说明

| 文件 | 说明 |
| --- | --- |
| `data.go` | 初始化 GORM 数据库连接，注册 data 层 ProviderSet |
| `user.go` | 用户数据访问逻辑，包括登录、查询、分页 |
| `greeter.go` | Kratos 模板示例仓储 |

## user.go 主要方法

| 方法 | 说明 |
| --- | --- |
| `Login` | 邮箱密码登录，生成 JWT，事务更新 `remember_token` |
| `FindByID` | 根据用户 ID 查询用户基础信息 |
| `ListUsersWithPage` | 分页查询用户列表，支持姓名和邮箱模糊搜索 |

## 登录事务

登录成功后使用 GORM 闭包事务：

```go
db.Transaction(func(tx *gorm.DB) error {
    return tx.Model(&userModel{}).
        Where("id = ?", user.ID).
        Update("remember_token", token).
        Error
})
```

闭包事务的好处是：

- 返回 error 时自动回滚。
- 返回 nil 时自动提交。
- 避免忘记 `Commit` 或 `Rollback`。

## JWT 配置

JWT 密钥和过期时间来自 `configs/config.yaml`：

```yaml
auth:
  jwt:
    password: change-me-user-service-jwt-secret
    expire_seconds: 86400
```

`password` 是 JWT 签名密钥，不是用户登录密码。

## 注意事项

- 当前用户密码是明文比较，后续建议改为密码哈希校验。
- `remember_token` 字段长度建议至少 `varchar(512)` 或使用 `text`。
- data 层不应该返回 proto 类型，只返回 biz 层定义的领域对象。
- 查询错误要转换为 Kratos errors，例如 NotFound、Unauthorized。
