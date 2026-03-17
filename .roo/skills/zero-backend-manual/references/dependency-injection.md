# 依赖注入说明

本文档介绍项目中依赖注入的使用方式和原理。

## 概述

项目使用 Google Wire 实现依赖注入。Wire 是一个编译时依赖注入框架，通过代码生成实现依赖组装，具有以下优点：
- 编译时检查依赖关系，提前发现错误
- 无运行时反射，性能更好
- 代码生成，调试方便

## 目录结构

```
cmd/
├── admin/
│   ├── main.go          # 入口文件
│   ├── wire.go          # 依赖注入定义
│   └── wire_gen.go      # 自动生成的代码（请勿修改）
└── api/
    ├── main.go
    ├── wire.go
    └── wire_gen.go
```

## 使用方式

### 定义依赖

在 wire.go 文件中定义依赖关系：

```go
//go:build wireinject
// +build wireinject

package main

import (
    "zero-backend/internal/config"
    "zero-backend/internal/storage/redis"
    "zero-backend/modules/api/server"
    "zero-backend/providers"

    "github.com/google/wire"
)

func wireApp() *server.HTTPServer {
    panic(wire.Build(
        config.New,
        redis.New,
        providers.ApiControllersProviderSet,
        providers.MiddlewaresProviderSet,
        // ... 其他依赖
    ))
}
```

### 使用依赖

在 main.go 中直接使用生成的初始化函数：

```go
package main

func main() {
    app := wireApp()
    app.Run()
}
```

### 生成代码

修改 wire.go 后，运行以下命令生成代码：

```bash
cd cmd/api
wire
```

或者在项目根目录运行：

```bash
go generate ./...
```

## Provider Set 组织

项目使用 Provider Set 机制组织依赖，常见的 Provider Set 如下：

### providers/service.go

```go
var ServiceProviderSet = wire.NewSet(
    service.NewRbacMenuService,
    service.NewRbacApiService,
    service.NewRbacRoleService,
    // ... 其他 Service
)
```

### providers/repository.go

```go
var RepositoryProviderSet = wire.NewSet(
    repository.NewUserRepository,
    repository.NewArticleRepository,
    // ... 其他 Repository
)
```

### providers/controllers.go

```go
var ControllersProviderSet = wire.NewSet(
    controller.NewControllers,
    controller.NewAuthController,
    // ... 其他 Controller
)
```

## 添加新依赖

### 添加新的 Service

1. 在 internal/service/ 目录下创建 Service 文件
2. 在 providers/service.go 中添加 Provider
3. 在对应模块的 wire.go 中添加依赖
4. 运行 wire 生成代码

示例：添加新的 UserService

```go
// internal/service/user.go
type UserService struct {
    repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
    return &UserService{repo: repo}
}
```

```go
// providers/service.go
var ServiceProviderSet = wire.NewSet(
    // ... 现有 Service
    service.NewUserService,  // 添加新 Service
)
```

### 添加新的 Repository

1. 在 internal/repository/ 目录下创建 Repository 文件
2. 在 providers/repository.go 中添加 Provider
3. 在对应模块的 wire.go 中添加依赖
4. 运行 wire 生成代码

### 添加新的模块

1. 在 modules/ 目录下创建新模块目录
2. 创建 controller、service、server 等子目录
3. 在 providers/ 中添加对应的 Provider Set
4. 在 cmd/ 中创建新的入口程序
5. 在新入口的 wire.go 中添加依赖
6. 运行 wire 生成代码

## 注意事项

- wire_gen.go 是自动生成的文件，请勿手动修改
- 修改 wire.go 后需要重新运行 wire 命令
- Wire 使用 panic 方式报告错误，这是正常行为
- 确保所有依赖都能正确初始化，否则会在启动时 panic

## 最佳实践

- 保持 Provider Set 的职责清晰，按类型组织
- 公共依赖放在 providers/ 目录
- 模块特定依赖放在各模块的 service/ 目录
- 使用接口定义依赖，便于测试和替换实现