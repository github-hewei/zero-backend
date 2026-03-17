---
name: golang-project-engineering-guide
description: A collection of concise, practical, and highly consistent engineering best practices for Golang projects, helping teams build maintainable, testable, secure, and modern codebases. This specification must be referenced and followed when developing Golang projects.
---

# Golang 项目工程化指南

**版本**：v0.1  
**目标**：为 Golang 项目提供一套简洁、实用、高一致性的工程化最佳实践集合，帮助团队构建可维护、可测试、安全且现代化的代码库。开发golang项目时必须参考并遵循它。
**适用场景**：中大型业务项目、微服务、开源库、CLI 工具等。  

## 概述原则
- 优先显式依赖注入，避免全局状态  
- 错误处理要显式、一致、可断言  
- 上下文规范传递，严禁滥用顶层 context  
- 日志统一接口 + 依赖注入 + 结构化 + TraceID 支持  
- 项目结构职责清晰、模块化  
- 高单元测试覆盖率，优先外部包（black-box）测试  
- 注释简洁规范，遵循“代码即注释”  
- 外部依赖谨慎引入，需确认  
- 持续采用现代 Go 特性（泛型、range-over-func 等）  
- 函数选项模式作为复杂对象构造首选  
- 强制 linter 门禁、接口小而美、依赖最小化、并发安全优先无锁  

## 全局变量使用约束
全局变量会引入隐式依赖、状态共享、测试困难等问题。  
**规则**：非必要严禁使用全局变量。  

**最佳实践**：  
- 所有依赖通过构造函数或接口注入（Dependency Injection）  
- 复杂依赖关系推荐使用 Google Wire 自动生成注入代码  
- 常量（const）可全局定义，但仅限配置值、枚举、魔法数字等  
- 测试中易于 mock 替换  

**示例**：
```go
type Service struct {
    db     *sql.DB
    logger Logger
}

func NewService(db *sql.DB, logger Logger) *Service {
    return &Service{db: db, logger: logger}
}
```

## 错误处理规范
**规则**：错误必须显式检查与处理，优先返回错误而非 panic。  

**最佳实践**：  
- 始终使用 `if err != nil { ... }`  
- 特定场景定义并导出 sentinel 错误，便于 `errors.Is` / `errors.As` 断言  
- 使用 `fmt.Errorf` 或 `errors.Wrap` 添加上下文（视项目确认是否引入 pkg/errors）  
- API 项目中建议统一错误码 + 结构化错误类型  
- 禁止忽略错误（使用 `_` 接收除非明确有意）  

**示例**：
```go
var ErrUserNotFound = errors.New("user not found")

func GetUser(ctx context.Context, id string) (*User, error) {
    // ...
    return nil, ErrUserNotFound
}

// 调用方
if errors.Is(err, ErrUserNotFound) {
    // 友好处理
}
```

## 上下文（Context）规范
**规则**：上下文是请求生命周期的唯一合法传递载体。  

**最佳实践**：  
- 从 HTTP/gRPC 请求或主流程上下文派生所有子上下文  
- 使用 `context.WithTimeout`、`WithCancel`、`WithValue`（值仅限请求级元数据）  
- 严禁在库/业务代码中使用 `context.Background()` 或 `context.TODO()` 作为默认  
- 所有 goroutine 必须监听 `ctx.Done()` 并及时清理  
- defer cancel() 必须成对出现  

**示例**：
```go
ctx, cancel := context.WithTimeout(parentCtx, 10*time.Second)
defer cancel()

// 向下游传递 ctx
data, err := repo.Fetch(ctx, id)
```

## 日志处理规范
**规则**：日志必须统一、可追踪、可替换、结构化。  

**最佳实践**：  
- 定义统一的 Logger 接口（含 ctx 参数）  
- 通过依赖注入传入 Logger 实例，**禁止使用全局日志对象**  
- 采用结构化日志（键值对），自动携带 TraceID（从 ctx 提取）  
- 级别清晰：Debug（开发）、Info（关键事件）、Warn、Error  
- 支持链式 With() 添加字段  

**示例**：
```go
type Logger interface {
    Info(ctx context.Context, msg string, args ...any)
    Error(ctx context.Context, msg string, args ...any)
    With(args ...any) Logger
}

func NewUserService(logger Logger, ...) *UserService {
    return &UserService{logger: logger}
}
```

## 项目模块拆分与职责划分
**规则**：遵循单一职责原则，模块边界清晰。  

**最佳实践**：  
- 采用标准 Go 项目布局：  
  - cmd/          → 可执行文件入口  
  - internal/     → 私有业务逻辑、domain、service、repository  
  - pkg/          → 可复用导出包（慎用）  
  - api/          → OpenAPI、proto 定义  
- 典型分层：domain → repository → service → handler/transport  
- 避免 God package，单个包职责单一  
- 使用 Go Modules 管理依赖，禁止循环依赖  

## 单元测试规范
**规则**：追求高覆盖率（业务逻辑建议 ≥80%），优先外部包测试。  

**最佳实践**：  
- 测试文件命名为 xxx_test.go  
- 优先表驱动测试（table-driven tests）  
- 使用 testify 或原生 testing + 子测试  
- mock 依赖使用 gomock / testify/mock  
- 避免测试内部实现，聚焦接口行为  
- CI 强制最低覆盖率门限  

**示例**：
```go
func TestGetUser(t *testing.T) {
    tests := []struct {
        name    string
        id      string
        wantErr bool
    }{...}

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // setup mock
            got, err := svc.GetUser(context.Background(), tt.id)
            // assert
        })
    }
}
```

## 注释规范
**规则**：注释要少而精，代码本身应自解释。  

**最佳实践**：  
- 包、导出函数、类型必须写 Godoc 注释（第一句概括用途）  
- 关键逻辑、性能敏感处、容易误解处写为什么而非做什么  
- 禁止无意义注释（如 `i++ // 递增`）  
- 使用有意义命名 + 短函数 + 清晰结构减少注释需求  

**示例**：
```go
// GetUser retrieves a user by ID.
// It returns ErrUserNotFound when the user does not exist.
func GetUser(ctx context.Context, id string) (*User, error)
```

## 外部依赖管理
**规则**：外部依赖是最大风险源，必须谨慎。  

**最佳实践**：  
- 新依赖引入前与团队/负责人确认  
- 优先标准库 → 成熟社区包  
- 评估指标：活跃度、star、最近提交、许可证、CVE 历史  
- 版本固定，使用 go.mod 明确声明  
- 定期运行 govulncheck / osv-scanner  

## Go 版本与现代语法采用
**规则**：保持与社区同步，适度拥抱新特性。  

**最佳实践**：  
- go.mod 中显式声明 `go 1.26`（或更高）  
- 优先使用：泛型、slices/maps 包、range-over-func 迭代器、slog 等  
- 新特性引入时考虑向后兼容（build tag 或渐进替换）  

## 函数选项模式（Functional Options）
**规则**：复杂对象构造首选函数选项模式。  

**最佳实践**：  
- 支持默认值、可选参数、链式调用、向后兼容  
- 广泛用于 Server、Client、Config 等构造  

**示例**：
```go
type ServerOption func(*Server)

func WithTimeout(d time.Duration) ServerOption { ... }

srv := NewServer(":8080",
    WithTimeout(5*time.Second),
    WithLogger(logger),
)
```

## 代码质量门禁 - golangci-lint
**规则**：所有代码提交必须通过统一 linter 检查。  

**最佳实践**：  
- 使用 golangci-lint（最新版）  
- 开启核心检查器：staticcheck、revive、gosec、nilaway、errcheck、contextcheck 等  
- 配置 .golangci.yml，设置严格模式  
- 本地 + CI 双重门禁  

## 接口设计原则
**规则**：小接口、单一职责、使用方定义接口。  

**最佳实践**：  
- 接口粒度小（1~3 方法最佳）  
- 接口定义放在消费方包中  
- 函数/结构体接受接口而非具体类型  

## 依赖管理与安全
**规则**：最小化 + 持续扫描。  

**最佳实践**：  
- 每次变更后 `go mod tidy && go mod verify`  
- CI 集成 govulncheck  
- 定期 dependabot / renovate 自动更新  

## 并发安全规范
**规则**：优先无锁设计，强制 race 检测。  

**最佳实践**：  
- 优先 channel、errgroup、context 取消  
- sync.Once / sync.Pool / atomic 优先于 mutex  
- 所有并发测试带 `-race`  
- 锁使用：defer unlock、最小范围、不锁内调用外部
