# framework-gin
快速构建web项目

Quickly build http services through framework-gin

# Introduction

application.yaml 中可以配置监听地址，监听端口，日志级别，日志路径等

application-test.yaml 中可以配置jwt密钥，数据库参数等

# 运行方法
go run main.go 目前配置监听8080端口

# 接口详情
见document.pdf

# 代码结构

--common 通用枚举，函数

--controller http接口入口

--dao 数据库访问入口

--model 数据库模型定义

--middleware http中间件

--pojo http参数

--services 实际负责处理逻辑层
