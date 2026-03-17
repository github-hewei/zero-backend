# 数据库设计规范

本文档定义了项目的数据库设计规范，包括表结构设计规则、字段命名规范和设计原则。

## 设计原则

### 通用规范

- 使用 MySQL 数据库，InnoDB 存储引擎
- 使用 utf8mb4 字符集，支持 emoji 和更多字符
- 表名和字段名使用小写字母，多个单词用下划线分隔
- 时间字段使用 Unix 时间戳整数存储
- 支持软删除机制

### 表前缀

所有表使用统一前缀 `gaz_`，在配置文件中定义：

```yaml
mysql:
  prefix: "gaz_"
```

## 命名规范

### 表命名

- 使用小写字母
- 多个单词用下划线分隔
- 添加统一前缀 `gaz_`
- 使用单数形式或模块名

**正确**: `gaz_user`, `gaz_article_category`
**错误**: `gazUsers`, `gaz_article_categories`

### 字段命名

| 字段类型 | 命名规则 | 示例 |
|----------|----------|------|
| 主键 | id | `id int` |
| 创建时间 | created_at | `created_at int` |
| 更新时间 | updated_at | `updated_at int` |
| 删除时间 | deleted_at | `deleted_at int` |
| 排序 | sort | `sort int` |
| 状态 | status | `status tinyint` |
| 企业ID | store_id | `store_id int` |
| 父级ID | parent_id / pid | `parent_id int` |

### 常用字段

所有业务表建议包含以下通用字段：

```sql
`id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID'
`created_at` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间'
`updated_at` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间'
`deleted_at` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间'
`store_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '企业ID'
`sort` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '排序'
```

## 字段类型规范

### 主键

```sql
`id` int(11) unsigned NOT NULL AUTO_INCREMENT
```

- 使用无符号整型
- 从 10000 开始自增，避免与历史数据冲突

### 时间字段

```sql
`created_at` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间'
`updated_at` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间'
```

- 使用 Unix 时间戳整数存储
- 类型为 `int(11) unsigned`
- 不使用时区，存储 UTC 时间

### 软删除

```sql
`deleted_at` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间'
```

- 值为 0 表示未删除
- 值大于 0 表示删除时间（Unix 时间戳）
- GORM 使用 `gorm/plugin/soft_delete` 插件处理

### 状态字段

```sql
`status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '状态(1启用 0禁用)'
```

- 使用 tinyint 类型
- 明确注释说明各值的含义

### 字符串字段

```sql
-- 短字符串
`name` varchar(50) NOT NULL DEFAULT '' COMMENT '名称'

-- 长字符串
`description` varchar(500) NOT NULL DEFAULT '' COMMENT '描述'

-- 唯一标识
`username` varchar(32) NOT NULL DEFAULT '' COMMENT '用户名'
```

- 根据实际需求选择合适的长度
- 添加 NOT NULL 和默认值
- 添加 COMMENT 说明

### 文本字段

```sql
-- 短文本
`content` text NOT NULL COMMENT '内容'

-- 长文本
`content` longtext NOT NULL COMMENT '内容'

-- JSON 文本
`setting_values` mediumtext NOT NULL COMMENT '设置内容（json格式）'
```

## 索引规范

### 主键索引

- 每个表必须有主键
- 使用自增 ID 作为主键

### 常用查询索引

根据查询需求创建索引：

```sql
-- 用户名查询
KEY `username` (`username`)

-- 手机号查询
KEY `mobile` (`mobile`)

-- 企业ID + 状态查询
KEY `store_id_status` (`store_id`, `status`)
```

### 唯一索引

确保数据唯一性的字段应创建唯一索引：

```sql
-- 设置项唯一性
UNIQUE KEY `unique_key` (`store_id`, `setting_key`)
```

## 多租户设计

### store_id 字段

所有业务表添加 `store_id` 字段实现多租户隔离：

```sql
`store_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '企业ID'
```

- 值为 0 表示系统级数据
- 值大于 0 表示具体企业数据

### 查询规范

- 超级管理员可以访问所有企业数据（store_id = 0 或所有）
- 普通管理员只能访问所属企业数据（store_id = 当前企业ID）
- 查询时必须携带 store_id 条件

## 关联表设计

### 一对多关系

使用外键关联：

```sql
-- 文章表
`category_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '文章分类ID'

-- 索引
KEY `category_id` (`category_id`)
```

### 多对多关系

使用中间表关联：

```sql
CREATE TABLE `gaz_user_role` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
  `role_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '角色ID',
  `store_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '企业ID',
  `created_at` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `role_id` (`role_id`),
  KEY `store_id` (`store_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户角色关联记录表';
```

## 迁移规范

### 分片迁移

项目使用分片迁移方式管理数据库变更：

```sql
-- 数据库设计文件说明：
-- 这个文件应该分为多个"执行片段"，用特定的分隔符隔开。
-- 系统启动时依次执行每个"执行片段"的SQL，并记录进度，下次启动时从上次的进度开始执行。
-- 所以"执行片段"在执行之后，不应该再修改。如果需要修改数据库，应该新建一个"执行片段"。
-- 特定的分隔符为：-- [CHECK POINT] --

CREATE TABLE `gaz_example` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '名称',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='示例表';

-- [CHECK POINT] --

ALTER TABLE `gaz_example` ADD COLUMN `status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '状态';

-- [CHECK POINT] --
```

### 规范要求

- 新增表或字段时在 database.sql 中添加
- 使用 `-- [CHECK POINT] --` 分隔执行片段
- 每个片段执行后记录进度
- 不修改已执行的片段
- 添加必要的 COMMENT 说明

## Model 定义规范

### GORM 标签

使用 GORM 标签定义模型映射：

```go
type User struct {
    ID        int64          `gorm:"primaryKey" json:"id"`
    Username  string         `gorm:"size:50;uniqueIndex" json:"username"`
    Password  string         `gorm:"size:255" json:"-"`
    Status    int            `gorm:"default:1" json:"status"`
    StoreID   int            `gorm:"index" json:"storeId"`
    CreatedAt time.Time      `json:"createdAt"`
    UpdatedAt time.Time      `json:"updatedAt"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
    return "gaz_user"
}
```

### 常用标签

| 标签 | 说明 |
|------|------|
| primaryKey | 主键 |
| column:xxx | 指定列名 |
| size:xxx | 字段大小 |
| uniqueIndex | 唯一索引 |
| index | 普通索引 |
| default:xxx | 默认值 |
| not null | 非空 |

## 注意事项

- 遵循表的命名规范
- 添加必要的 COMMENT 说明字段含义
- 考虑多租户场景，添加 store_id 字段
- 合理创建索引，避免过多索引影响性能
- 使用合适的字段类型，避免浪费存储空间
- 时间字段统一使用 Unix 时间戳