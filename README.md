# V2EX

## 一个爬取 `v2ex.com` 数据并展示的项目。

![stars](https://img.shields.io/github/stars/LieWell/v2ex.svg)
![forks](https://img.shields.io/github/forks/LieWell/v2ex.svg)
![issues](https://img.shields.io/github/issues/LieWell/v2ex.svg)
![watchers](https://img.shields.io/github/watchers/LieWell/v2ex.svg)
![contributors](https://img.shields.io/github/contributors/LieWell/v2ex.svg)
![license](https://img.shields.io/github/license/LieWell/v2ex.svg)

## 项目说明

本项目提供了一个爬虫自动去爬取 `V2EX` 的数据，包括会员信息、帖子信息、头像等。

本项目 **不提供** 已爬取的数据， `V2EX API` 每个 IP 每小时配额为 600，因此可能需要很长时间才可以爬取完所有数据。

## 在线预览

[v2ex.liewell.fun](https://v2ex.liewell.fun)

## 如何使用

### 本地运行 ![golang](https://img.shields.io/badge/golang->=1.21.0-blue)

```shell
# 初始化数据库
执行脚本 sql/init.sql

# 启动服务
go run main.go

# 访问
http://localhost:12321
```

### 线上部署

```shell
# 进入项目根目录
# 打包镜像
docker build -t v2ex:nightly .
# 启动服务
docker compose up -d
```

## 功能清单

| 功能描述    | 是否已实现 | 备注 |
|---------|-------|----|
| 注册统计    | ✅     |    |
| 增长趋势    | ✅     |    |
| 头像马赛克拼图 | ❌     |    |
| 地域分布    | ❌     |    |
| 帖子列表    | ❌     |    |
