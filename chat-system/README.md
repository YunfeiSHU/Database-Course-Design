# Chat System

Go + Vue 3 即时聊天系统，包含 Web 客户端、Gin HTTP API、Gorilla WebSocket、MySQL 持久化和 Redis session/在线状态。

## 功能

- 账号注册与登录，服务端自动生成 8 位账号。
- Redis 保存登录 session，HTTP 鉴权使用 `Authorization: Bearer <token>`。
- 好友申请、好友同意、好友列表与在线状态。
- 独立的好友列表、对话列表和消息区。
- WebSocket 实时收发消息，消息落库后更新对话列表。
- 历史消息查询，聊天窗口展示最近多条消息。
- 服务端关键错误带日志，Redis/MySQL 等基础设施异常会返回 500，账号密码错误返回 401。

## 技术栈

- 客户端：Vue 3、Element Plus、Pinia、Axios、Vite、WebSocket
- 服务端：Go、Gin、Gorilla WebSocket、GORM
- 数据库：MySQL
- 缓存：Redis

## 目录结构

```text
chat-system
├── client
│   ├── index.html
│   ├── package.json
│   ├── vite.config.js
│   └── src
│       ├── api
│       ├── stores
│       ├── styles
│       └── websocket
├── docs
│   └── protocol.md
├── server
│   ├── api
│   │   └── websocket
│   ├── cmd
│   │   └── server
│   ├── configs
│   ├── internal
│   │   ├── conversation
│   │   ├── friend
│   │   ├── infrastructure
│   │   ├── message
│   │   ├── notification
│   │   └── user
│   └── pkg
└── sql
	├── schema.sql
	└── migration_friend_status.sql
```

## 环境准备

需要本地已有：

- Go
- Node.js 和 npm
- MySQL
- Redis

先创建数据库和表：

```powershell
mysql -u root -p < sql/schema.sql
```

如果已有旧表，再按需要执行：

```powershell
mysql -u root -p chat_system < sql/migration_friend_status.sql
```

服务端配置可放在 `server/.env`，也可以使用 PowerShell 环境变量：

```powershell
$env:MYSQL_DSN="root:123456@tcp(127.0.0.1:3306)/chat_system?charset=utf8mb4&parseTime=True&loc=Local"
$env:REDIS_ADDR="127.0.0.1:6379"
$env:REDIS_PASSWORD=""
$env:REDIS_DB="0"
$env:HTTP_ADDR=":8080"
```

默认值定义在 [server/configs/config.go](server/configs/config.go)。

## 启动服务端

确认 MySQL 和 Redis 已启动，然后运行：

```powershell
cd chat-system/server
go mod tidy
go run ./cmd/server
```

看到类似日志表示 HTTP/WebSocket 服务已启动：

```text
chat server listening on :8080
```

## 启动客户端

```powershell
cd chat-system/client
npm install
npm run dev
```

默认访问：

```text
http://127.0.0.1:5173
```

如果 `5173` 被占用，Vite 会自动切换到下一个端口，例如 `http://localhost:5174/`。前端代理会把 `/api` 和 `/ws` 转发到 `127.0.0.1:8080`。

## 使用流程

1. 打开客户端页面。
2. 注册两个用户，记录两个自动生成的账号。
3. 分别在两个浏览器窗口或不同浏览器中登录。
4. 在好友列表中向对方账号发送好友申请。
5. 对方在“好友申请”中点击同意。
6. 好友会出现在好友列表中，点击好友可进入聊天。
7. 发送第一条消息后，对话会出现在独立的对话列表中。
8. 点击对话列表项，右侧消息区会加载最近多条历史消息。

## 接口文档

HTTP API 和 WebSocket 协议见 [docs/protocol.md](docs/protocol.md)。

## 常见问题

### 前端提示 `ECONNREFUSED 127.0.0.1:8080`

说明 Vite 代理无法连接后端。检查服务端是否已启动，端口是否为 `8080`，以及 `HTTP_ADDR` 是否被改过。

### 登录返回 401

401 只表示账号不存在或密码错误。可以检查 `user` 表中是否存在该账号，或确认当前后端连接的是正确的 MySQL 数据库。

### 登录返回 500 `login service unavailable`

表示登录依赖的基础设施异常，例如 Redis 未启动、MySQL 更新失败、session 写入失败。后端日志会记录具体原因，例如 `store session for account ... connect refused`。

### 注册时报 `Incorrect datetime value: '0000-00-00'`

旧版本代码可能会把 Go 零值时间写入 MySQL。当前用户模型已使用 `gorm:"autoCreateTime"` 自动写入创建时间。更新代码后重启服务端再注册。

### 对话列表为空

对话列表只展示已经产生过消息的会话。好友刚添加成功但还没聊过时，只会显示在好友列表中；向好友发送第一条消息后，会创建 conversation 记录并出现在对话列表。

## 验证命令

服务端：

```powershell
cd chat-system/server
go test ./...
```

客户端：

```powershell
cd chat-system/client
npm run build
```
