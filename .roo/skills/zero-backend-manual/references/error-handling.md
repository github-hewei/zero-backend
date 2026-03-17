# 错误处理规范

本文档详细介绍项目中的错误处理机制和规范。

## 错误类型定义

项目定义了三种基础错误类型，分别用于不同的场景：

### 1. 用户级错误（UserError）

用于客户端可以理解并展示给用户的错误，如参数验证失败、业务逻辑错误等。

```go
// 创建用户级错误
return nil, apperror.NewUserError("用户名不能为空")
```

**使用场景**：
- 参数验证失败
- 业务规则校验失败
- 资源不存在
- 权限不足

### 2. 系统级错误（SystemError）

用于系统内部错误，如数据库连接失败、第三方服务调用失败等。

```go
// 创建系统级错误
return nil, apperror.NewSystemError(err, "传入参数错误")
```

**使用场景**：
- 数据库操作失败
- 外部服务调用失败
- 文件操作失败
- 未知异常

### 3. 认证错误（UnauthorizedError）

用于认证和授权相关的错误。

```go
// 创建认证错误
return nil, apperror.NewUnauthorizedError("Token已过期")
```

**使用场景**：
- Token 无效
- Token 已过期
- 未登录
- 登录已失效

## 错误码定义

错误码用于精确区分不同类型的错误：

```go
const (
    ErrorCodeNone         ErrorCode = 0     // 无错误
    ErrorCodeUser         ErrorCode = 4000  // 用户级错误起始码
    ErrorCodeUnauthorized ErrorCode = 4001  // 认证错误
    ErrorCodeSystem       ErrorCode = 5000  // 系统级错误起始码
)
```

**错误码规则**：
- 0: 成功
- 4000-4999: 用户级错误
- 5000+: 系统级错误

## 错误响应格式

项目统一了 API 响应格式，错误响应示例：

```json
{
    "errcode": 4000,
    "message": "用户名或密码错误",
    "data": null,
    "cost": "1.234ms",
    "traceId": "xxx"
}
```

## 错误处理流程

### 控制器层

控制器层应该直接返回错误，由响应处理中间件统一处理：

```go
func (c *UserController) Create(ctx *gin.Context) {
    var req dto.UserCreateRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        response.Error(ctx, apperror.NewUserError("请求参数错误"))
        return
    }

    user, err := c.userService.Create(ctx, &req)
    if err != nil {
        response.Error(ctx, err)
        return
    }

    response.Success(ctx, user)
}
```

### 服务层

服务层应该返回有意义的错误信息：

```go
func (s *UserService) Create(ctx context.Context, req *dto.UserCreateRequest) (*model.User, error) {
    // 业务逻辑
    existingUser, err := s.repo.FindByUsername(ctx, req.Username)
    if err != nil {
        return nil, apperror.NewSystemError(err, "查询用户失败")
    }
    if existingUser != nil {
        return nil, apperror.NewUserError("用户名已存在")
    }

    // 创建用户
    user := &model.User{
        Username: req.Username,
        // ...
    }
    if err := s.repo.Create(ctx, user); err != nil {
        return nil, apperror.NewSystemError(err, "创建用户失败")
    }

    return user, nil
}
```

### 仓储层

仓储层应该返回底层错误，由服务层包装：

```go
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
    var user model.User
    err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &user, nil
}
```

## 错误链

项目支持错误链，便于追踪错误根源：

```go
// 使用 WithErr 添加错误链
return nil, apperror.NewSystemError(err, "操作失败")
```

错误链允许在最外层获取完整的错误信息：

```go
if errors.As(err, &systemError) {
    fmt.Println(systemError.Message)   // 用户可见的错误消息
    fmt.Println(systemError.Err)       // 原始错误（用于日志记录）
}
```

## 全局错误处理

项目使用中间件实现全局错误处理：

```go
// internal/response/response.go
func Error(c *gin.Context, err error) {
    var userError *apperror.UserError
    if errors.As(err, &userError) {
        output(c, Response{
            ErrCode: int(userError.Code),
            Message: userError.Message,
        })
        return
    }

    var systemError *apperror.SystemError
    if errors.As(err, &systemError) {
        output(c, Response{
            ErrCode: int(systemError.Code),
            Message: systemError.Message,
        })
        return
    }

    // ... 其他错误类型
}
```

## 最佳实践

### 错误消息规范

- 用户级错误消息应该是面向用户的，清晰描述问题
- 系统级错误消息可以包含技术细节，但不应暴露敏感信息
- 错误消息应保持简洁，避免过长的描述

### 错误处理原则

- 优先返回有意义的错误，而不是直接返回底层错误
- 在最接近数据源的地方处理特定错误（如 GORM 的 RecordNotFound）
- 在服务层进行业务规则校验
- 在控制器层进行参数验证

### 日志记录

系统级错误应该被记录，以便排查问题：

```go
if err != nil {
    logger.Error("操作失败").Err(err).Str("operation", "create_user").Send()
    return nil, apperror.NewSystemError(err, "操作失败")
}
```

### 避免的错误

- 避免返回空错误（nil error），应该明确返回成功或失败
- 避免在错误消息中包含敏感信息（如密码、密钥）
- 避免过度包装错误，导致错误信息丢失