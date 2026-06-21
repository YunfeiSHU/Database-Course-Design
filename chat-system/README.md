# Chat System

一个基于 Go + Vue 3 的即时聊天系统。后端采用 Gin、GORM、Redis 和 WebSocket，前端采用 Vite、Pinia 和 Element Plus。

## 项目定位

这个仓库演示的是按业务模块划分的 DDD 风格结构：`user`、`friend`、`conversation`、`message`、`presence` 各自独立，HTTP 接口放在模块内的 `interfaces/http`，WebSocket 传输层放在 `internal/infrastructure/websocket`。

## 功能

- 注册与登录，服务端自动生成 8 位账号。
- Redis 记录登录 session，HTTP 认证使用 `Authorization: Bearer <token>`。
- 好友申请、好友同意、好友列表和在线状态展示。
- 对话列表和聊天窗口分离，消息实时收发并落库。
- WebSocket 负责实时消息、心跳和在线状态广播。
- 支持历史消息查询。

## 技术栈

- 客户端：Vue 3、Vite、Pinia、Element Plus、Axios
- 服务端：Go、Gin、GORM、Gorilla WebSocket
- 数据库：MySQL
- 缓存：Redis

## 目录结构

```text
chat-system/
├── client/
│   ├── index.html
│   ├── package.json
│   ├── vite.config.js
│   └── src/
│       ├── api/
│       ├── stores/
│       ├── styles/
│       └── websocket/
├── docs/
│   └── protocol.md
├── server/
│   ├── api/
│   ├── cmd/
│   │   └── server/
│   ├── configs/
│   ├── internal/
│   │   ├── conversation/
│   │   ├── friend/
│   │   ├── infrastructure/
│   │   ├── message/
│   │   ├── presence/
│   │   └── user/
│   └── pkg/
└── sql/
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

如果已有旧表，再按需要执行迁移脚本：

```powershell
mysql -u root -p chat_system < sql/migration_friend_status.sql
```

服务端配置可以放在 `server/.env`，也可以直接用 PowerShell 环境变量：

```powershell
$env:MYSQL_DSN="root:123456@tcp(127.0.0.1:3306)/chat_system?charset=utf8mb4&parseTime=True&loc=Local"
$env:REDIS_ADDR="127.0.0.1:6379"
$env:REDIS_PASSWORD=""
$env:REDIS_DB="0"
$env:HTTP_ADDR=":8080"
```

默认值定义在 [server/configs/config.go](server/configs/config.go)。

## 启动服务端

确认 MySQL 和 Redis 已启动，然后执行：

```powershell
cd chat-system/server
go mod tidy
go run ./cmd/server
```

启动成功后会看到类似日志：

```text
chat server listening on :8080
```

## 启动客户端

```powershell
cd chat-system/client
npm install
npm run dev
```

默认访问地址：

```text
http://127.0.0.1:5173
```

如果 `5173` 被占用，Vite 会自动切换到其他端口，例如 `http://localhost:5174/`。前端代理会把 `/api` 和 `/ws` 转发到 `127.0.0.1:8080`。

## 使用流程

1. 打开客户端页面。
2. 注册两个用户，记录两个自动生成的账号。
3. 分别在两个浏览器窗口或两个浏览器中登录。
4. 在好友列表中向对方账号发送好友申请。
5. 对方在好友申请列表中点击同意。
6. 好友会出现在好友列表中，点击好友即可进入聊天。
7. 发送第一条消息后，对话会出现在对话列表中。
8. 点击对话列表项，右侧消息区会加载历史消息。

## 接口文档

HTTP API 和 WebSocket 协议说明见 [docs/protocol.md](docs/protocol.md)。

## 常见问题

### 前端提示 `ECONNREFUSED 127.0.0.1:8080`

说明 Vite 代理无法连接后端。检查服务端是否已启动、端口是否为 `8080`，以及 `HTTP_ADDR` 是否被修改。

### 登录返回 401

401 只表示账号不存在或密码错误。可以检查 `user` 表中是否存在该账号，或确认当前后端连接的是正确的 MySQL 数据库。

### 登录返回 500 `login service unavailable`

说明登录依赖的基础设施异常，例如 Redis 未启动、MySQL 更新失败或 session 写入失败。后端日志会记录具体原因。

### 注册时报 `Incorrect datetime value: '0000-00-00'`

旧版本代码可能会把 Go 零值时间写入 MySQL。当前用户模型已经使用 `gorm:"autoCreateTime"` 自动写入创建时间。更新代码后重启服务端再注册。

### 对话列表为空

对话列表只展示已经产生过消息的会话。好友刚添加成功但还没聊天时，只会显示在好友列表中；向好友发送第一条消息后，会创建 conversation 记录并出现在对话列表。

## 验证命令

服务端：

```powershell
cd chat-system/server
go test ./...
```

或者只做编译检查：

```powershell
cd chat-system/server
go build ./cmd/server
```

客户端：

```powershell
cd chat-system/client
npm run build
```
