# API 设计规范

本文档定义了项目的 API 设计规范，确保接口设计的一致性和可维护性。

## 设计概述

本项目采用统一的 POST 请求方式处理所有 API 接口。这种设计简化了前后端交互，统一了请求处理逻辑。

## 请求规范

### 请求方法

所有 API 接口统一使用 POST 方法：

```go
apiGroup.POST("/login", ctrl.AuthController.Login)
apiGroup.POST("/rbac/user/list", ctrl.RbacUserController.List)
apiGroup.POST("/rbac/user/create", ctrl.RbacUserController.Create)
```

### URL 规范

- 使用小写字母
- 多个单词用斜杠分隔
- 体现资源层级关系
- URL 中包含操作类型（list、create、update、delete 等）

**正确**:
```
/api/rbac/user/list
/api/rbac/user/create
/api/rbac/user/update
/api/rbac/user/delete
/api/user/user/detail
```

**错误**:
```
/api/rbac/users        # 缺少操作类型
/api/rbac/user/get     # 使用 get 动词
```

### 操作类型

项目中定义了以下标准操作：

| 操作 | 说明 |
|------|------|
| list | 获取列表 |
| detail | 获取详情 |
| create | 创建 |
| update | 更新 |
| delete | 删除 |
| tree | 获取树形结构 |
| upload | 上传文件 |
| refresh-token | 刷新 Token |

### 请求体

所有请求使用 JSON 格式，请求参数使用 DTO 定义：

```json
{
    "username": "admin",
    "password": "password123"
}
```

### 分页参数

```json
{
    "page": 1,
    "limit": 20,
    "sort": "created_at",
    "order": "desc"
}
```

| 参数 | 说明 | 默认值 |
|------|------|--------|
| page | 页码 | 1 |
| limit | 每页数量 | 20 |
| sort | 排序字段 | - |
| order | 排序方向 | desc |

### 过滤参数

```json
{
    "status": 1,
    "role": "admin",
    "keyword": "search"
}
```

## 响应规范

### 响应格式

所有响应使用统一格式：

```json
{
    "errcode": 0,
    "message": "success",
    "data": {},
    "cost": "1.234ms",
    "traceId": "xxx"
}
```

### 字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| errcode | int | 错误码，0 表示成功 |
| message | string | 消息提示 |
| data | any | 响应数据 |
| cost | string | 请求耗时 |
| traceId | string | 追踪ID，用于日志定位 |

### 成功响应

**无数据**:
```json
{
    "errcode": 0,
    "message": "success",
    "data": null,
    "cost": "0.123ms",
    "traceId": "xxx"
}
```

**单个数据**:
```json
{
    "errcode": 0,
    "message": "success",
    "data": {
        "id": 1,
        "username": "admin",
        "email": "admin@example.com"
    },
    "cost": "1.234ms",
    "traceId": "xxx"
}
```

**列表数据**:
```json
{
    "errcode": 0,
    "message": "success",
    "data": {
        "list": [
            {"id": 1, "username": "user1"},
            {"id": 2, "username": "user2"}
        ],
        "total": 100
    },
    "cost": "2.345ms",
    "traceId": "xxx"
}
```

**树形数据**:
```json
{
    "errcode": 0,
    "message": "success",
    "data": [
        {
            "id": 1,
            "name": "parent",
            "children": [
                {"id": 2, "name": "child1"},
                {"id": 3, "name": "child2"}
            ]
        }
    ],
    "cost": "1.234ms",
    "traceId": "xxx"
}
```

### 错误响应

**用户级错误**:
```json
{
    "errcode": 4000,
    "message": "用户名或密码错误",
    "data": null,
    "cost": "0.234ms",
    "traceId": "xxx"
}
```

**认证错误**:
```json
{
    "errcode": 4001,
    "message": "Token已过期",
    "data": null,
    "cost": "0.123ms",
    "traceId": "xxx"
}
```

**系统错误**:
```json
{
    "errcode": 5000,
    "message": "系统错误",
    "data": null,
    "cost": "10.234ms",
    "traceId": "xxx"
}
```

## 认证规范

### 认证方式

项目使用 JWT 进行身份认证，请求头携带 Token：

```
Authorization: Bearer <access_token>
```

### Token 类型

| Token | 说明 | 时效 |
|-------|------|------|
| access_token | 访问令牌 | 1小时（admin）/ 2小时（api） |
| refresh_token | 刷新令牌 | 24小时（admin）/ 7天（api） |

### 认证流程

1. 用户登录获取 access_token 和 refresh_token
2. 每次请求携带 access_token
3. access_token 过期后使用 refresh_token 刷新
4. 刷新后返回新的 token 对

### 认证接口

```
POST /api/login           # 登录
POST /api/refresh-token   # 刷新 Token
POST /api/logout          # 退出登录
```

## 权限控制

### 权限验证

项目实现了基于 RBAC 的权限控制：

```go
// JWT 认证中间件
apiGroup.Use(adminMiddlewares.Auth.JWTAuth())

// API 权限验证中间件
apiGroup.Use(adminMiddlewares.Auth.CheckAPIPermission())
```

### 权限接口

```
POST /api/permissions      # 获取用户权限
POST /api/change-password  # 修改密码
```

## 文件上传规范

### 上传接口

```
POST /api/upload/file/upload
Content-Type: multipart/form-data
```

### 请求参数

| 参数 | 类型 | 说明 |
|------|------|------|
| file | file | 上传的文件 |
| group_id | int | 文件分组ID |

### 响应示例

```json
{
    "errcode": 0,
    "message": "success",
    "data": {
        "id": 1,
        "url": "/uploads/xxx.jpg",
        "filename": "xxx.jpg",
        "size": 102400,
        "mime_type": "image/jpeg"
    },
    "cost": "100.234ms",
    "traceId": "xxx"
}
```

## 接口分组

项目按功能模块组织接口：

### 鉴权相关

```
POST /api/login
POST /api/refresh-token
POST /api/logout
POST /api/change-password
POST /api/permissions
```

### RBAC 权限管理

```
POST /api/rbac/menu/list
POST /api/rbac/menu/create
POST /api/rbac/menu/update
POST /api/rbac/menu/delete
POST /api/rbac/api/list
POST /api/rbac/role/list
POST /api/rbac/user/list
POST /api/rbac/user/set-roles
```

### 设置管理

```
POST /api/setting/list
POST /api/setting/create
POST /api/setting/update
POST /api/setting/delete
```

### 文章管理

```
POST /api/article/category/list
POST /api/article/category/create
POST /api/article/article/list
POST /api/article/article/create
```

### 用户管理

```
POST /api/user/user/list
POST /api/user/user/create
POST /api/user/user/update
POST /api/user/user/detail
POST /api/user/points/logs
```

### 文件上传

```
POST /api/upload/group/list
POST /api/upload/file/list
POST /api/upload/file/upload
```

## 注意事项

- 所有接口统一使用 POST 方法
- 请求参数通过 JSON Body 传递
- 响应格式保持统一
- 错误码遵循项目定义的错误码规范