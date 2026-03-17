# Git 提交规范

本文档定义了项目的 Git 提交规范，旨在保持提交历史清晰，便于代码审查和问题追踪。

## 提交信息格式

项目采用 Angular 提交规范，格式如下：

```
<type>(<scope>): <subject>

[optional body]

[optional footer]
```

### 格式说明

- **type**: 提交类型
- **scope**: 影响范围（可选）
- **subject**: 简短描述
- **body**: 详细说明（可选）
- **footer**: 脚注（可选，通常用于关联 Issue）

## 提交类型（Type）

| 类型 | 说明 |
|------|------|
| feat | 新功能 |
| fix | Bug 修复 |
| docs | 文档更新 |
| style | 代码格式调整（不影响功能） |
| refactor | 重构（既不是新功能也不是修复） |
| perf | 性能优化 |
| test | 测试相关 |
| build | 构建系统或外部依赖变更 |
| ci | CI 配置文件和脚本变更 |
| chore | 其他不修改源代码或测试文件的变更 |

## 作用域（Scope）

使用以下作用域：

- **模块相关**: admin, api, cli
- **层次相关**: controller, service, repository, model, middleware
- **功能相关**: auth, rbac, upload, article, user, setting
- **基础设施**: config, logger, database, redis, mongodb
- **其他**: docs, workflow

## 提交示例

### 功能提交

```
feat(user): 添加用户头像上传功能

- 支持 JPG、PNG 格式
- 最大文件大小 2MB
- 上传后自动生成缩略图

Closes #123
```

### 修复提交

```
fix(auth): 修复 Token 过期后刷新失败的问题

- 调整 Token 验证逻辑
- 添加过期时间检查

Fixes #456
```

### 文档提交

```
docs: 更新项目架构文档

- 添加目录结构说明
- 补充依赖注入使用方式
```

### 重构提交

```
refactor(service): 优化用户服务层代码

- 提取公共方法
- 简化查询逻辑
```

### 性能优化

```
perf(repository): 优化用户列表查询性能

- 添加索引
- 使用分页缓存
```

## 提交规范

### 标题行

- 不超过 50 个字符
- 首字母小写
- 末尾不加句号
- 使用祈使句

**正确**:
```
feat(user): 添加用户注册功能
```

**错误**:
```
feat(user): 添加了用户注册功能  # 过去时
feat: 添加用户注册功能         # 缺少作用域
```

### 正文

- 标题行后空一行
- 每行不超过 72 个字符
- 说明「做什么」和「为什么这样做」

### 脚注

用于关联 Issue：

- `Closes #123`: 关闭 Issue
- `Fixes #456`: 修复 Issue
- `Refs #789`: 引用 Issue

## 分支管理

### 分支命名

- 功能分支: `feature/<feature-name>`
- 修复分支: `fix/<issue-name>`
- 热修复分支: `hotfix/<issue-name>`
- 发布分支: `release/<version>`

**示例**:
```
feature/user-avatar-upload
fix/token-refresh-issue
hotfix/security-vulnerability
```

### 分支策略

- main/master: 稳定分支，仅通过 PR 合并
- develop: 开发分支
- 功能分支: 从 develop 创建，合并回 develop

## 合并策略

### Pull Request 规范

- 标题清晰描述变更内容
- 详细说明变更原因和内容
- 关联相关 Issue
- 至少一人代码审查通过

### 合并方式

- 使用 Squash Merge 合并功能分支
- 保持提交历史整洁

## 注意事项

- 提交前检查代码格式
- 确保提交信息清晰准确
- 避免一次性提交大量不相关的变更
- 保持提交粒度适中