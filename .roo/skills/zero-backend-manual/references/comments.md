# 注释规范

本文档定义了项目的注释规范，包括代码注释和文档注释的要求。

## 注释原则

- 注释应该解释「为什么」，而不是「是什么」
- 保持注释与代码同步更新
- 避免显而易见的注释
- 与项目现有风格保持一致

## 包注释

每个包应该在包声明前添加包注释：

```go
// Package user 提供用户相关的业务逻辑处理
// 包括用户的增删改查、权限验证等功能
package user
```

如果包包含多个文件，可以在单独的文件中添加包注释（通常命名为 `doc.go`）：

```go
// Package user 提供用户管理功能。
//
// 本包处理用户 CRUD 操作、认证和权限管理。
package user
```

## 导出函数/类型注释

导出（公开）的函数、类型、结构体、变量等应该添加注释：

### 函数注释

```go
// NewUserService 创建用户服务实例
// 参数 db 为数据库连接，repo 为用户仓储
func NewUserService(db *gorm.DB, repo *UserRepository) *UserService {
    return &UserService{
        db:   db,
        repo: repo,
    }
}
```

### 结构体注释

```go
// UserService 处理用户业务逻辑
type UserService struct {
    db   *gorm.DB
    repo *UserRepository
}
```

### 结构体字段注释

```go
type User struct {
    ID        int64          `gorm:"primaryKey" json:"id"`         // 用户ID
    Username  string         `gorm:"size:50" json:"username"`      // 用户名
    Password  string         `gorm:"size:255" json:"-"`            // 密码（不暴露）
    Status    int            `gorm:"default:1" json:"status"`      // 用户状态：1=启用，0=禁用
    CreatedAt time.Time      `json:"createdAt"`                    // 创建时间
    UpdatedAt time.Time      `json:"updatedAt"`                    // 更新时间
}
```

### 方法注释

```go
// List 根据过滤条件获取分页用户列表
// 返回用户列表和总数
func (s *UserService) List(ctx context.Context, req *ListRequest) ([]*User, int64, error) {
    // implementation
}
```

### 接口注释

```go
// Uploader 定义文件上传接口
type Uploader interface {
    // Upload 上传文件到存储并返回域名路径
    Upload(ctx context.Context, file *multipart.FileHeader, savePath string) (string, error)
    
    // Delete 从存储中删除已上传文件
    Delete(ctx context.Context, filePath string) error
}
```

### 常量注释

```go
const (
    // StatusActive 表示启用状态
    StatusActive = 1
    // StatusInactive 表示禁用状态
    StatusInactive = 0
)
```

## 行内注释

对于复杂的逻辑或关键步骤，可以添加行内注释：

```go
// 转换为 UTC 时间，避免时区问题
createdAt := time.Now().UTC()
```

## 注释风格

### 句子风格

- 使用完整的句子
- 句首大写
- 句末句号

**正确**:
```go
// GetUserByID 根据ID获取用户
func GetUserByID(id int) (*User, error) {}
```

**错误**:
```go
// get user by id
func GetUserByID(id int) (*User, error) {}
```

### 动词使用

- 函数/方法注释使用第三人称单数
- 描述性注释使用现在时

**正确**:
```go
// 创建新用户
// 处理存储前的密码哈希
```

## 特殊注释

### TODO 注释

```go
// TODO(username): 实现缓存失效机制
// TODO: 添加错误处理单元测试
```

### FIXME 注释

```go
// FIXME: 处理并发更新时的竞态条件
```

### NOTE 注释

```go
// NOTE: 本方法非线程安全
```

### 弃用注释

```go
// Deprecated: 请使用 NewUserServiceV2 代替
func OldUserService() {}
```

## 注释位置

### 包注释位置

包注释应该放在包声明之前的注释块中：

```go
// Package user 提供用户管理功能。
package user
```

### 变量声明注释

变量注释应该放在变量声明之前：

```go
// DefaultPageSize 是每页默认显示数量
const DefaultPageSize = 20
```

## 注意事项

- 导出函数必须添加注释
- 复杂逻辑应该添加注释解释
- 保持注释简洁明了
- 定期检查和更新注释
- 删除不再需要的注释

## 文档生成

项目可以使用 Go 自带的文档工具生成 API 文档：

```bash
# 生成包文档
go doc -all ./internal/service

# 启动文档服务器
godoc -http=:6060
```

生成的文档可以通过 `http://localhost:6060` 访问。