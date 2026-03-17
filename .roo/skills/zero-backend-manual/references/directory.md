# 目录结构说明

本文档详细说明项目中各目录和包的职责、设计思想以及使用注意事项。

## 整体目录结构

```
zero-backend/
├── cmd/                    # 入口程序目录
│   ├── admin/              # 管理后台服务入口
│   ├── api/                # 用户API服务入口
│   └── cli/                # CLI工具入口
├── internal/               # 内部包（不对外暴露）
│   ├── apperror/           # 错误定义和处理
│   ├── config/             # 配置管理
│   ├── constants/          # 常量定义
│   ├── ctxkeys/            # 上下文键定义
│   ├── dto/                # 数据传输对象
│   ├── logger/             # 日志组件
│   ├── middleware/         # 公共中间件
│   ├── model/              # 数据模型定义
│   ├── repository/         # 数据访问层
│   ├── request/            # 请求处理和验证
│   ├── response/           # 响应处理
│   ├── service/            # 业务逻辑层
│   └── storage/            # 存储层（数据库、缓存）
├── modules/                # 模块目录
│   ├── admin/              # 管理后台模块
│   │   ├── controller/     # 控制器
│   │   ├── middleware/     # 模块中间件
│   │   ├── server/         # 服务器配置
│   │   └── service/        # 模块服务
│   ├── api/                # 用户API模块
│   │   ├── controller/
│   │   ├── middleware/
│   │   ├── server/
│   │   └── service/
│   └── cli/                # CLI模块
├── providers/              # 依赖注入提供者
├── docs/                   # 项目文档
├── data/                   # 静态数据（如SQL脚本）
└── views/                  # 视图文件
```

## 目录详解

### cmd 目录

**作用**: 存放各服务的入口程序，每个子目录代表一个独立的服务。

**设计思想**: 采用多服务架构，将管理后台和用户 API 分离部署，实现关注点分离。这种设计允许：
- 独立部署和扩展各服务
- 独立维护和发布
- 不同的配置和权限控制

**包含内容**:
- `main.go`: 服务入口文件，调用依赖注入生成的初始化函数
- `wire.go`: Wire 依赖注入定义文件
- `wire_gen.go`: Wire 自动生成的依赖注入代码（请勿手动修改）

**注意事项**:
- 每个服务的 main.go 应该保持简洁，仅负责启动服务
- 具体的依赖组装逻辑在 wire.go 中定义

### internal 目录

**作用**: 存放项目的内部包，这些包不对外暴露，仅限本项目使用。

**设计思想**: Go 语言的 internal 包机制确保这些包只能被本项目导入，提供了良好的封装性。

#### internal/apperror

**作用**: 定义应用级别的错误类型和处理函数。

**包含内容**:
- `error.go`: 定义错误类型（UserError、SystemError、UnauthorizedError）和错误码
- `error_util.go`: 错误处理的工具函数

**设计思想**: 
- 将错误分为用户级错误、系统级错误和认证错误三类
- 使用错误码区分不同类型的错误，便于客户端处理
- 支持错误链和错误信息覆盖

**使用示例**:
```go
// 创建用户级错误
return nil, apperror.NewUserError("用户名或密码错误")

// 创建系统级错误
return nil, apperror.NewSystemError(err, "传入参数错误")
```

#### internal/config

**作用**: 管理应用配置，支持从配置文件和环境变量加载配置。

**包含内容**:
- `config.go`: 配置结构体定义和加载逻辑

**设计思想**:
- 使用 Viper 库实现配置管理，支持多种格式（YAML、JSON 等）
- 支持从 .env 文件加载环境变量
- 配置结构体使用 mapstructure 标签映射 YAML 配置

**注意事项**:
- 敏感配置（如数据库密码、JWT 密钥）应通过环境变量或 .env 文件提供
- 配置结构体应与 config.yaml 保持同步

#### internal/constants

**作用**: 存放项目中使用的常量。

**包含内容**:
- `redis_keys.go`: Redis 键常量定义
- `queue_keys.go`: 消息队列键常量定义
- `points.go`: 积分相关常量

**设计思想**: 将魔法字符串和魔法数字提取为具名常量，提高代码可读性和可维护性。

#### internal/ctxkeys

**作用**: 定义 Context 中使用的键，用于在请求生命周期中传递数据。

**包含内容**:
- `ctxkeys.go`: Context 键定义

**设计思想**: 使用类型安全的键定义，避免字符串拼写错误。

#### internal/dto

**作用**: 定义数据传输对象（Data Transfer Object），用于请求参数和响应数据的结构化。

**包含内容**:
- 各业务模块的 DTO 定义（如 user.go、article.go、auth.go 等）
- `dto.go`: 公共 DTO 定义

**设计思想**:
- DTO 与 Model 分离，DTO 用于 API 层，Model 用于数据层
- 支持请求参数验证和响应数据格式化
- 使用结构体标签定义验证规则

**注意事项**:
- DTO 应根据 API 接口需求设计，不应直接暴露 Model
- 对于复杂的验证逻辑，可以使用自定义验证器

#### internal/logger

**作用**: 日志组件，提供统一的日志记录功能。

**包含内容**:
- `logger.go`: 日志核心实现
- `writer.go`: 日志写入器
- `logger_test.go`: 单元测试

**设计思想**:
- 基于 Zerolog 库实现高性能日志
- 支持多种输出方式：控制台、文件、MongoDB
- 支持日志轮转和压缩
- 支持日志级别控制

**配置说明**:
```yaml
logger:
  level: "debug"  # 日志级别
  writers:        # 输出方式
    - "console"
    - "file"
    - "mongodb"
  file:
    path: "runtime/logs"
    filename: "app.log"
    max_size: 100    # MB
    max_age: 30      # 天
    max_backups: 3
    compress: true
```

#### internal/middleware

**作用**: 存放公共的 HTTP 中间件。

**包含内容**:
- `middleware.go`: 中间件集合定义
- `cors.go`: CORS 跨域中间件
- `before.go`: 请求预处理中间件

**设计思想**:
- 中间件应该保持简洁，专注于单一职责
- 公共中间件放在 internal/middleware，模块特定中间件放在 modules/xxx/middleware

#### internal/model

**作用**: 定义数据模型，对应数据库表结构。

**包含内容**:
- 各业务模块的 Model 定义（如 user.go、article.go、rbac.go 等）
- `model.go`: 公共 Model 定义

**设计思想**:
- 使用 GORM 的模型定义规范
- 支持软删除（通过 gorm/plugin/soft_delete）
- 表名前缀通过配置统一管理

**注意事项**:
- Model 应该与数据库表结构一一对应
- 避免在 Model 中添加业务逻辑

#### internal/repository

**作用**: 数据访问层，负责与数据库交互。

**包含内容**:
- 各业务模块的 Repository 实现
- `repository.go`: 公共 Repository 接口和工具

**设计思想**:
- Repository 封装数据访问逻辑，提供清晰的数据操作接口
- 支持分页、排序、筛选等通用功能
- 使用 Filter 模式实现动态查询条件

**核心接口**:
```go
// Filter 接口用于构建查询条件
type Filter interface {
    Apply(db *gorm.DB) *gorm.DB
}

// Pagination 分页参数
type Pagination struct {
    Page  int
    Limit int
}
```

**注意事项**:
- Repository 应该只关注数据访问，不包含业务逻辑
- 复杂的业务逻辑应该在 Service 层处理

#### internal/request

**作用**: 请求处理和参数验证。

**包含内容**:
- `request.go`: 请求处理工具
- `validate.go`: 参数验证工具

**设计思想**:
- 使用 go-playground/validator 进行参数验证
- 支持自定义验证规则和翻译

#### internal/response

**作用**: 统一响应格式处理。

**包含内容**:
- `response.go`: 响应结构体和输出函数

**设计思想**:
- 统一 API 响应格式，包含错误码、消息、数据、耗时、追踪ID
- 支持错误自动识别和格式化

**响应格式**:
```json
{
    "errcode": 0,
    "message": "success",
    "data": {},
    "cost": "1.234ms",
    "traceId": "xxx"
}
```

#### internal/service

**作用**: 业务逻辑层，处理核心业务逻辑。

**包含内容**:
- 各业务模块的 Service 实现
- `uploader/`: 文件上传服务

**设计思想**:
- Service 层是业务逻辑的核心，负责协调 Repository 和其他服务
- 使用依赖注入获取所需的 Repository 和其他 Service
- 每个 Service 应该对应一个领域概念

**文件上传设计**:
- 使用工厂模式支持多种存储后端（本地、七牛云）
- 通过配置切换存储类型
- 接口化设计便于扩展新的存储方式

```go
type Uploader interface {
    Upload(ctx context.Context, file *multipart.FileHeader, savePath string) (domain string, err error)
    Delete(ctx context.Context, filePath string) error
}
```

#### internal/storage

**作用**: 存储层，管理数据库和缓存连接。

**包含内容**:
- `mysql/`: MySQL 数据库连接管理
- `redis/`: Redis 缓存连接管理
- `mongodb/`: MongoDB 数据库连接管理

**设计思想**:
- 封装底层存储细节，提供统一的初始化接口
- 支持懒加载和连接池管理

### modules 目录

**作用**: 存放各模块的 HTTP 处理逻辑，按模块组织。

**设计思想**: 
- 每个模块（admin、api）都有独立的 controller、middleware、server、service
- 模块之间保持独立，便于单独维护和扩展

#### modules/admin

管理后台模块，包含：
- `controller/`: HTTP 控制器，处理请求并调用 Service
- `middleware/`: 模块特定的中间件（如权限验证）
- `server/`: 服务器配置和路由定义
- `service/`: 模块特定的服务（如认证服务）

#### modules/api

用户 API 模块，结构与 admin 模块类似。

### providers 目录

**作用**: 依赖注入提供者，定义 Wire 依赖集合。

**包含内容**:
- `service.go`: Service 层依赖集合
- `repository.go`: Repository 层依赖集合
- `controllers.go`: Controller 依赖集合
- `mysql.go`: MySQL 依赖
- `mongodb.go`: MongoDB 依赖
- `redis.go`: Redis 依赖
- `logger.go`: 日志依赖
- `middleware.go`: 中间件依赖
- `server.go`: 服务器依赖

**设计思想**:
- 使用 Wire 的 Provider Set 机制组织依赖
- 公共依赖放在 providers/，模块特定依赖放在各模块的 service/

**注意事项**:
- Wire 生成的代码（wire_gen.go）请勿手动修改
- 修改依赖关系后需要重新运行 wire 命令

### data 目录

**作用**: 存放静态数据文件。

**包含内容**:
- `database.sql`: 数据库初始化脚本

### docs 目录

**作用**: 项目文档。

### views 目录

**作用**: 存放视图模板文件（如 HTML）。

## 包命名规范

项目采用以下包命名规范：
- 使用简洁、描述性的包名
- 避免使用复数形式（如使用 user 而不是 users）
- 使用通用术语（如 repository、service、controller）
- 模块特定包放在 modules/ 目录下