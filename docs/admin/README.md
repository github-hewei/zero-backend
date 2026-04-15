# 接口文档说明

当前目录存放的是当前项目管理后台admin模块的接口文档，文档遵循 `openapi: 3.1.0` 标准。
当前目录下的 `./index.json` 是通过 `swagger-cli` 自动生成，**无需手动编辑**。

## 关于 swagger-cli

```sh

# 安装
npm install -g swagger-cli

# 打包（bundle）所有引用
swagger-cli bundle index.yaml --outfile index.json

# 验证文档
swagger-cli validate openapi.yaml

```
