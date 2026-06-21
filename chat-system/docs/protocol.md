# 通信协议

服务端默认 HTTP 地址：

```text
http://127.0.0.1:8080
```

WebSocket 地址：

```text
ws://127.0.0.1:8080/ws?token=<登录返回的 token>
```

前端开发环境通过 Vite 代理访问后端，所以客户端代码中请求的是 `/api/...` 和 `/ws`。

后端 websocket 传输适配层当前位于 [server/internal/infrastructure/websocket](../server/internal/infrastructure/websocket)。

## 通用约定

### 认证

除注册和登录外，HTTP API 都需要请求头：

```http
Authorization: Bearer <token>
```

WebSocket 使用查询参数传 token：

```text
/ws?token=<token>
```

### 常见状态码

- `200`：请求成功。
- `400`：请求参数错误，或当前业务操作不允许。
- `401`：未登录、token 无效，或登录时账号/密码错误。
- `500`：服务端依赖异常或内部错误。服务端会输出日志帮助定位。

### 错误响应

```json
{
  "error": "invalid account or password"
}
```

登录时 Redis、MySQL 等基础设施错误不会返回 401，而是返回：

```json
{
  "error": "login service unavailable"
}
```

后端日志会包含具体错误上下文。

## HTTP API

### 注册

```http
POST /api/register
Content-Type: application/json
```

请求：

```json
{
  "nickname": "张三",
  "password": "123456",
  "confirm_password": "123456"
}
```

响应：

```json
{
  "account": "10000001",
  "nickname": "张三",
  "message": "register_success"
}
```

说明：服务端自动生成 8 位账号，密码使用 bcrypt 哈希后存储。

### 登录

```http
POST /api/login
Content-Type: application/json
```

请求：

```json
{
  "account": "10000001",
  "password": "123456"
}
```

响应：

```json
{
  "type": "login_success",
  "token": "SESSION_TOKEN",
  "account": "10000001",
  "user_id": 1,
  "nickname": "张三",
  "last_login_time": "2026-06-19T10:00:00+08:00"
}
```

说明：token 会写入 Redis，默认有效期为 24 小时。

### 获取好友列表

```http
GET /api/friends
Authorization: Bearer <token>
```

响应：

```json
[
  {
    "id": 1,
    "user_id": 1,
    "friend_id": 2,
    "friend": {
      "id": 2,
      "account": "10000002",
      "nickname": "李四",
      "create_time": "2026-06-19T10:00:00+08:00",
      "last_login_time": null
    },
    "online": true
  }
]
```

说明：好友列表是独立列表，用于展示可聊天联系人和在线状态，不等同于对话列表。

### 发送好友申请

```http
POST /api/friends
Authorization: Bearer <token>
Content-Type: application/json
```

请求：

```json
{
  "account": "10000002"
}
```

响应：

```json
{
  "message": "friend request sent"
}
```

### 获取好友申请

```http
GET /api/friend-requests
Authorization: Bearer <token>
```

响应：

```json
[
  {
    "id": 3,
    "user_id": 1,
    "friend_id": 2,
    "status": "pending",
    "user": {
      "id": 1,
      "account": "10000001",
      "nickname": "张三",
      "create_time": "2026-06-19T10:00:00+08:00",
      "last_login_time": null
    }
  }
]
```

### 同意好友申请

```http
POST /api/friend-requests/:id/accept
Authorization: Bearer <token>
```

响应：

```json
{
  "message": "ok"
}
```

说明：同意后服务端会把请求记录改为 `accepted`，并创建反向好友关系。

### 获取对话列表

```http
GET /api/conversations
Authorization: Bearer <token>
```

响应：

```json
[
  {
    "id": 1,
    "user_id": 1,
    "peer_id": 2,
    "last_message_id": 10,
    "status": "normal",
    "update_time": "2026-06-19T10:30:00+08:00",
    "peer": {
      "id": 2,
      "account": "10000002",
      "nickname": "李四",
      "create_time": "2026-06-19T10:00:00+08:00",
      "last_login_time": null
    },
    "last_message": {
      "id": 10,
      "sender_id": 1,
      "receiver_id": 2,
      "content": "hello",
      "status": "delivered",
      "send_time": "2026-06-19T10:30:00+08:00"
    }
  }
]
```

说明：对话列表来自 `conversation` 表，只展示已经产生过消息的会话。刚添加但未聊天的好友只会出现在好友列表中。

### 获取历史消息

```http
GET /api/history?friend_id=2
Authorization: Bearer <token>
```

也可以按账号查询：

```http
GET /api/history?account=10000002
Authorization: Bearer <token>
```

响应：

```json
[
  {
    "id": 10,
    "sender_id": 1,
    "receiver_id": 2,
    "content": "hello",
    "status": "delivered",
    "send_time": "2026-06-19T10:30:00+08:00"
  },
  {
    "id": 11,
    "sender_id": 2,
    "receiver_id": 1,
    "content": "hi",
    "status": "delivered",
    "send_time": "2026-06-19T10:31:00+08:00"
  }
]
```

说明：当前返回最近 50 条消息，并按发送时间正序展示。

`GET /api/messages` 当前等同于 `GET /api/history`。

## WebSocket

### 连接

```text
ws://127.0.0.1:8080/ws?token=<token>
```

token 缺失或无效时，服务端返回 401，不升级为 WebSocket 连接。

WebSocket 连接只是传输入口，不直接承载业务逻辑；在线状态、消息投递和协议封装由后端的 websocket 适配层处理。

### 消息信封

所有 WebSocket 消息统一使用：

```json
{
  "type": "chat",
  "data": {}
}
```

### 发送聊天消息

客户端发送：

```json
{
  "type": "chat",
  "data": {
    "to": "10000002",
    "content": "hello"
  }
}
```

服务端行为：

1. 根据 `to` 查找接收方。
2. 校验双方是否为好友。
3. 保存消息到 MySQL。
4. 更新双方 conversation 记录。
5. 向发送方回推消息；如果接收方在线，也向接收方推送。

服务端推送：

```json
{
  "type": "chat",
  "data": {
    "from": "10000001",
    "to": "10000002",
    "content": "hello",
    "send_time": "2026-06-19T10:30:00+08:00",
    "status": "delivered"
  }
}
```

如果不是好友，服务端会返回系统消息：

```json
{
  "type": "system",
  "data": {
    "content": "receiver is not your friend"
  }
}
```

### 系统消息

客户端可以发送：

```json
{
  "type": "system",
  "data": {
    "content": "今晚22点系统维护"
  }
}
```

服务端会广播该系统消息。

当前实现中，客户端也可以发送 `system` 消息，服务端会原样广播，用于调试或全局通知场景。

### 在线状态

用户连接成功后广播：

```json
{
  "type": "online",
  "data": {
    "account": "10000001",
    "user_id": 1
  }
}
```

用户断开后广播：

```json
{
  "type": "offline",
  "data": {
    "account": "10000001"
  }
}
```

前端收到 `online` / `offline` 后会更新好友列表的在线状态展示。

### 心跳

客户端发送：

```json
{
  "type": "heartbeat",
  "data": {}
}
```

服务端响应：

```json
{
  "type": "heartbeat",
  "data": {
    "time": "2026-06-19T10:30:00+08:00"
  }
}
```

心跳用于保持连接存活，前端当前每 30 秒发送一次。

## 调试建议

- 前端代理错误 `ECONNREFUSED 127.0.0.1:8080`：检查后端是否启动。
- 登录 401：检查账号和密码。
- 登录 500：检查后端日志，重点看 Redis/MySQL 是否正常。
- WebSocket 401：检查连接 URL 中的 token 是否存在且未过期。
- 对话列表为空：先和好友发送一条消息，产生 conversation 记录后再刷新对话列表。
