# 命名规范

本文档定义了项目的命名规范，包括变量、函数、文件、包等的命名规则。

## 通用原则

- 清晰表达意图：名称应该能够表达其用途
- 简洁但不晦涩：在清晰的前提下尽量简洁
- 使用英文：所有命名使用英文单词
- 避免缩写：除非是广泛认可的缩写（如 ID、API、URL）

## 包命名

### 目录命名

- 使用小写字母
- 使用单数形式
- 多个单词用连字符分隔（如 `upload-file`）
- 避免使用下划线

**正确**: `internal/user`
**错误**: `internal/users` 或 `internal/user_service`

### 包内命名

- 在同一个包内，类型名不需要包含包名
- 避免使用与标准库冲突的名称

**正确**:
```go
package user

type Service struct {}
```

**错误**:
```go
package user

type UserService struct {}  // 冗余
```

## 变量命名

### 变量名

- 使用驼峰命名法（CamelCase）
- 简短且有描述性
- 布尔变量使用 `is`、`has`、`can`、`should` 等前缀

**正确**:
```go
var userName string
var isActive bool
var canEdit bool
```

**错误**:
```go
var user_name string  // 下划线风格
var active bool      // 缺乏描述性
var flag int          // 无意义
```

### 常量名

- 使用大写字母
- 多个单词用下划线分隔
- 相关常量使用相同前缀分组

**正确**:
```go
const (
    StatusActive   = 1
    StatusInactive = 0
)
```

**错误**:
```go
const (
    status_active = 1  // 小写
    ACTIVE         = 1 // 缺乏分组
)
```

### 结构体命名

- 使用驼峰命名法
- 名词或名词短语
- 避免使用 `Info`、`Data` 等泛化名称

**正确**:
```go
type UserService struct {}
type Config struct {}
```

**错误**:
```go
type UserServiceStruct struct {}  // 冗余后缀
type UserInfo struct {}            // 泛化名称
```

### 接口命名

- 使用驼峰命名法
- 以 `er` 结尾
- 名词或名词短语

**正确**:
```go
type Reader interface {}
type Uploader interface {}
```

**错误**:
```go
type IReader interface {}     // 不需要 I 前缀
type UploadInterface interface {}  // 冗余后缀
```

## 函数命名

### 函数名

- 使用驼峰命名法
- 动词或动词短语
- 描述其行为

**正确**:
```go
func GetUserByID(id int) (*User, error)
func CreateUser(req *CreateUserRequest) (*User, error)
func DeleteUser(id int) error
```

**错误**:
```go
func get_user_by_id(id int) (*User, error)  // 下划线风格
func UserGet(id int) (*User, error)         // 语序不正确
```

### 方法命名

- 使用驼峰命名法
- 描述其行为

**正确**:
```go
func (s *UserService) List(ctx context.Context, req *ListRequest) ([]*User, error)
func (s *UserService) Create(ctx context.Context, user *User) error
```

### 构造函数

- 使用 `New` 或 `New` + 类型名
- 如果需要特定前缀以区分，使用有意义的名称

**正确**:
```go
func NewUserService(db *gorm.DB) *UserService {}
func NewLocalUploader() Uploader {}
```

## 文件命名

### 源文件

- 使用小写字母
- 多个单词用下划线分隔
- 与主要类型名匹配

**正确**:
```go
user_service.go
user_repository.go
user_controller.go
```

**错误**:
```go
UserService.go      // 大写
userService.go     // 驼峰
user.service.go    // 点分隔
```

### 测试文件

- 源文件名 + `_test.go` 后缀

**正确**:
```go
user_service_test.go
```

### Wire 文件

- `wire.go`: 依赖注入定义
- `wire_gen.go`: 自动生成（请勿修改）

## 目录命名

- 使用小写字母
- 多个单词用连字符分隔
- 避免使用下划线

**正确**:
```go
modules/admin/controller
internal/service/uploader
```

**错误**:
```go
modules/admin/Controller  // 大写
modules/admin/controller/  // 斜杠结尾
```

## 数据库命名

### 表名

- 使用小写字母
- 多个单词用下划线分隔
- 添加统一前缀（在配置中定义）

**正确**:
```go
// 配置前缀为 "gaz_"
gaz_user
gaz_article_category
```

### 字段名

- 使用小写字母
- 多个单词用下划线分隔
- 避免与 SQL 关键字冲突

**正确**:
```go
created_at
updated_at
is_active
```

**错误**:
```go
CreatedAt        // 大写
createdAt        // 驼峰
```

## Redis 键命名

- 使用冒号分隔层级
- 使用大写字母
- 有意义的命名

**正确**:
```go
const RedisAdminLoginKey = "ZAG:ADMIN:LOGIN"
const RedisUserLoginKey = "ZAG:USER:LOGIN"
```

## 路由命名

### URL 路径

- 使用小写字母
- 多个单词用连字符分隔
- 名词使用复数形式
- 体现资源层级

**正确**:
```go
POST /api/users
GET /api/users/:id
PUT /api/users/:id
DELETE /api/users/:id
GET /api/users/:id/articles
```

**错误**:
```go
POST /api/User          // 大写
POST /api/create_user   // 动词
GET /api/user/:id      // 单数
```

### 路由组

- 使用有意义的组名
- 体现功能模块

**正确**:
```go
apiGroup := r.Group("/api")
rbacGroup := apiGroup.Group("/rbac")
```

## 注意事项

- 保持命名的一致性
- 避免使用中文拼音
- 避免使用无意义的名称（如 `temp`、`tmp`、`data`）
- 避免在同一个作用域内使用相似的名称