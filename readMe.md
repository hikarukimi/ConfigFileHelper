# ConfigReader

## 概述

`ConfigReader` 是一个用于读取配置文件并将配置内容赋值给结构体或字符串的 Go 语言工具。它可以方便地从配置文件中提取配置信息，并将其映射到应用程序中使用的结构体中。

## 安装

### 前提条件

- Go 1.14 或更高版本

### 安装步骤

克隆仓库：

```
git clone https://github.com/hikarukimi/ConfigReader.git
```

## 使用方法
### 读取配置文件并映射到结构体

1. 创建一个配置文件（例如 config.yaml）：

   ```
   Database:
     Host: "localhost"
     Port: 5432
     User: "admin"
     Password: "password"
   ```

2. 定义一个结构体来接收配置：

   ```
   type DatabaseConfig struct {
    Host     string
    Port     int
    User     string
    Password string
   }
   ```

3. 使用 ConfigReader 读取配置文件并映射到结构体：

   ```
   package main
   
   import (
       "fmt"
       "./configreader"
   )
   
   type DatabaseConfig struct {
       Host     string
       Port     int
       User     string
       Password string
   }
   
   func main() {
       cr := &configreader.ConfigReader{configFilePath: "config.yaml"}
       var dbConfig DatabaseConfig
       cr.AssignMapConfigToStruct("Database", &dbConfig)
   	fmt.Printf("Database Configuration: %+v\n", dbConfig)
   }
   ```

### 读取单个配置项

1. 使用 ConfigReader 读取单个配置项：

   ```
   package main
   
   import (
       "fmt"
       "./configreader"
   )
   
   func main() {
       cr := &configreader.ConfigReader{configFilePath: "config.yaml"}
       var host string
       cr.AssignSingleConfigToString("Host", &host)fmt.Printf("Database Host: %s\n", host)
   }
   ```

## 常见问题

### Q: 配置文件中的值不同类型的值都需要添加双引号吗？

A: 目前的版本不同类型的变量需要不同类型的值不能统一加上双引号，例如int类型的属性配置项的值必须没有双引号。

### Q: 支持哪些配置文件格式？

A: 当前支持 YAML 格式的配置文件。如果你需要支持其他格式的配置文件，可以扩展 getConfigFileContext 和 getSingleConfig 函数。或者等待后续更新

## 贡献

欢迎贡献代码和提出建议！请遵循以下步骤：

1. Fork 本仓库
2. 创建一个新的分支 
3. 提交你的更改 
4. 推送到新的分支
5. 发起 Pull Request