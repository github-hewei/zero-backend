---
name: zero-backend-manual
description: zero-backend 项目的手册，当需要熟悉/开发 zero-backend 项目时，请必须阅读此手册。
modeSlugs:
  - orchestrator
---

# 项目文档

本文档是项目的完整技术文档，包含架构设计和开发规范两大核心内容。

## 文档目的

本项目文档实现两个主要目标：
- 帮助阅读者了解当前项目的架构设计
- 让阅读者继续开发时遵循约定的规范

## 文档结构

### 架构设计

- [目录结构说明](references/directory.md) - 详细描述各目录/包的职责和设计（优先阅读）
- [依赖注入说明](references/dependency-injection.md) - 依赖注入的使用方式和原理
- [数据库设计规范](references/database.md) - 数据库设计规范和注意事项
- [错误处理规范](references/error-handling.md) - 错误处理机制和规范

### 开发规范

- [代码组织规范](references/code-organization.md) - 代码结构和分层规范
- [命名规范](references/naming.md) - 变量、函数、文件等命名规范
- [注释规范](references/comments.md) - 代码注释规范
- [Git 提交规范](references/git-commit.md) - Git 提交信息规范
- [API 设计规范](references/api-design.md) - API 设计规范

## 快速开始

### 技术栈

- **Web 框架**: Gin
- **ORM**: GORM
- **依赖注入**: Google Wire
- **日志**: Zerolog
- **缓存**: Redis
- **数据库**: MySQL / MongoDB
- **对象存储**: 本地存储 / 七牛云

### 项目架构

项目采用经典的分层架构，将系统划分为清晰的层次：

- **入口层（cmd）**: 负责服务的启动和依赖组装
- **模块层（modules）**: 负责 HTTP 请求的处理和路由管理
- **服务层（internal/service）**: 负责业务逻辑的处理
- **仓储层（internal/repository）**: 负责数据访问和持久化
- **基础设施层（internal/storage）**: 负责数据库、缓存等基础组件
- **内部包（internal/*）**: 提供公共的工具和基础组件

### 多服务架构

项目支持多服务部署：
- **Admin 服务**: 管理后台服务（端口 8081）
- **API 服务**: 用户 API 服务（端口 8082）

## 相关资源

- 项目配置: `config.yaml`
- 数据库脚本: `data/database.sql`
- 依赖管理: `go.mod`
