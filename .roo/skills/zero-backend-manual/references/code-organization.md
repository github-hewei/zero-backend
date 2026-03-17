# 代码组织规范

本文档定义了项目的代码组织规范，包括目录结构、分层架构等内容。

## 分层架构

项目采用经典的分层架构，各层职责清晰：

```
请求 → Controller → Service → Repository → Storage
                ↓
              Model
```

### 各层职责

#### Controller 层（控制器）

**位置**: `modules/xxx/controller/`

**职责**:
- 处理 HTTP 请求和响应
- 参数绑定和验证
- 调用 Service 层处理业务逻辑
- 返回统一格式的响应

**规范**:
- 控制器应该保持简洁，不包含业务逻辑
- 每个控制器方法对应一个 API 端点
- 使用 DTO 接收请求参数
- 统一使用 response 包返回响应

**示例**:
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

#### Service 层（服务）

**位置**: `internal/service/`

**职责**:
- 处理核心业务逻辑
- 协调多个 Repository
- 事务管理
- 业务规则校验

**规范**:
- Service 应该专注于业务逻辑，不处理 HTTP 相关内容
- 使用依赖注入获取 Repository
- 返回 Model 或 DTO，不直接返回数据库对象

**示例**:
```go
type UserService struct {
    db   *gorm.DB
    repo *repository.UserRepository
}

func NewUserService(db *gorm.DB, repo *repository.UserRepository) *UserService {
    return &UserService{
        db:   db,
        repo: repo,
    }
}
```

#### Repository 层（仓储）

**位置**: `internal/repository/`

**职责**:
- 数据访问封装
- 数据库操作
- 查询条件构建

**规范**:
- Repository 只负责数据访问，不包含业务逻辑
- 使用 Filter 模式构建动态查询
- 支持分页、排序等通用功能

**示例**:
```go
type UserRepository struct {
    db *gorm.DB
}

func (r *UserRepository) FindByFilter(ctx context.Context, filter *UserFilter, pagination *repository.Pagination) ([]*model.User, error) {
    query := r.db.WithContext(ctx)
    
    // 应用过滤条件
    if filter.Username != "" {
        query = query.Where("username = ?", filter.Username)
    }
    
    // 分页
    query = query.Offset(pagination.Offset()).Limit(pagination.Limit)
    
    var users []*model.User
    err := query.Find(&users).Error
    
    return users, err
}
```

#### Model 层（模型）

**位置**: `internal/model/`

**职责**:
- 数据结构定义
- 数据库表映射

**规范**:
- Model 应该与数据库表结构一一对应
- 不包含业务逻辑
- 使用 GORM 标签定义表关系

**示例**:
```go
type User struct {
    ID        int64          `gorm:"primaryKey" json:"id"`
    Username  string         `gorm:"size:50;uniqueIndex" json:"username"`
    Password  string         `gorm:"size:255" json:"-"`
    Status    int            `gorm:"default:1" json:"status"`
    CreatedAt time.Time      `json:"createdAt"`
    UpdatedAt time.Time      `json:"updatedAt"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
    return "user"
}
```

## 包组织

### 内部包（internal）

internal 目录下的包为内部包，只能被本项目导入：

- `internal/apperror`: 错误定义
- `internal/config`: 配置管理
- `internal/constants`: 常量定义
- `internal/ctxkeys`: Context 键
- `internal/dto`: 数据传输对象
- `internal/logger`: 日志组件
- `internal/middleware`: 公共中间件
- `internal/model`: 数据模型
- `internal/repository`: 数据访问层
- `internal/request`: 请求处理
- `internal/response`: 响应处理
- `internal/service`: 业务逻辑层
- `internal/storage`: 存储层

### 模块包（modules）

modules 目录按功能模块组织：

- `modules/admin`: 管理后台模块
- `modules/api`: 用户 API 模块

每个模块内部包含：
- `controller/`: 控制器
- `middleware/`: 模块中间件
- `server/`: 服务器配置
- `service/`: 模块服务

### 提供者包（providers）

providers 目录存放依赖注入相关代码：

- `providers/service.go`: Service 依赖集合
- `providers/repository.go`: Repository 依赖集合
- `providers/controllers.go`: Controller 依赖集合
- 等等

## 文件组织

### 单文件原则

- 每个类型定义应该放在单独的文件中
- 文件名应该与主要类型名匹配
- 避免在一个文件中定义多个不相关的类型

### 导入顺序

Go 文件的导入应该按以下顺序组织：

1. 标准库
2. 第三方库
3. 项目内部包

组内按字母顺序排列：

```go
import (
    "context"
    "fmt"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/google/wire"
    "gorm.io/gorm"

    "zero-backend/internal/config"
    "zero-backend/internal/model"
    "zero-backend/internal/repository"
)
```

## 依赖管理

### 依赖注入

项目使用 Google Wire 进行依赖注入：

- 在 `cmd/xxx/wire.go` 中定义依赖关系
- 运行 `wire` 命令生成代码
- 在 `providers/` 目录组织 Provider Set

### 依赖原则

- 依赖应该明确注入，不使用全局变量
- 依赖应该是接口，而不是具体实现
- 避免循环依赖

## 目录创建规范

新增功能时，按以下结构组织代码：

```
internal/
    ├── repository/
    │   └── new_feature.go    # 新功能的 Repository
    ├── service/
    │   └── new_feature.go    # 新功能的 Service
    └── model/
        └── new_feature.go    # 新功能的 Model

modules/
    └── admin/
        ├── controller/
        │   └── new_feature.go  # 新功能的 Controller
        └── middleware/
            └── new_feature.go  # 新功能的中间件
```

## 注意事项

- 保持各层之间的依赖方向一致：Controller → Service → Repository
- 不要在 Model 中包含业务逻辑
- 不要在 Controller 中直接操作数据库
- 使用依赖注入管理依赖关系