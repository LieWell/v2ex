# V2EX

## 一个爬取 `v2ex.com` 数据并展示的项目。

![stars](https://img.shields.io/github/stars/LieWell/v2ex.svg)
![forks](https://img.shields.io/github/forks/LieWell/v2ex.svg)
![issues](https://img.shields.io/github/issues/LieWell/v2ex.svg)
![watchers](https://img.shields.io/github/watchers/LieWell/v2ex.svg)
![contributors](https://img.shields.io/github/contributors/LieWell/v2ex.svg)
![license](https://img.shields.io/github/license/LieWell/v2ex.svg)

## 项目说明

本项目 *不提供* 已爬取的数据，而是提供了一个爬虫自动去爬取数据。重启将从中断的地方继续爬取。

`V2EX` 限制每个 IP 每小时配额为 600，因此可能需要很长时间才可以爬取完所有数据。

也许可以通过自动更换代理的方式加快爬取速度。

## 运行

### 本地运行 ![golang](https://img.shields.io/badge/golang->=1.21.0-blue)

1. 开始前需要初始化数据库,表结构文件存在于 `mysql/init.sql`
2. 控制台运行 `go run main.go`
3. 浏览器访问：`http://localhost:12321`

### 在线预览

## 功能清单

| 功能描述     | 是否已实现 | API地址                          | 备注 |
|----------|-------|--------------------------------|----|
| 会员列表     | ✅     | /api/members/:pageNo/:pageSize |    |
| 会员地域分布   | ❌     | /api/members/location          |    |
| 会员增长趋势   | ❌     |                                |    |
| 账号字母概率分布 | ❌     |                                |    |

