# 刷新 Token — RefreshToken

## 功能概述

AccessToken 有效期仅 15 分钟，过期后客户端无法继续访问受保护接口。本接口允许客户端凭借登录时获取的 RefreshToken（有效期 7 天），在不重新输入密码的前提下换取一对全新的 AccessToken + RefreshToken，实现"无感续签"。该接口**不经过鉴权中间件**，因为调用时 AccessToken 通常已经过期。

## 路由信息

```
POST /v1/user/refresh
```

中间件链：`TreaceLoggerMiddleware` → `RefreshTokenHandler`（无 AuthHttpMiddleware）

## 完整实现流程

### 1. 中间件阶段

仅经过全局的 `TreaceLoggerMiddleware`，为请求生成唯一 `traceID`（UUID v4），绑定到 `logrus.Entry` 后注入 `context`，用于全链路日志追踪。

**不走 Auth 认证中间件**，因为用户调用此接口时 AccessToken 通常已过期，鉴权必然失败。本接口依赖 RefreshToken 自身的 JWT 签名 + Redis 存在性来保证合法性。

### 2. Handler 层 — 参数绑定

从请求体中 JSON 绑定 `RefreshTokenRequest` 结构体：

```go
type RefreshTokenRequest struct {
    RefreshToken string `json:"refresh_token" binding:"required"`
}
```

绑定失败时返回 `400001 — 请求格式错误`。

### 3. Service 层 — Token 刷新（四步操作）

`userService.RefreshTokenService` 按顺序执行以下步骤：

#### 步骤一：校验 RefreshToken 的 JWT 签名

```go
info, code, err := tools.VerifyRefreshToken(u.jwtConfig, req.RefreshToken)
```

使用 `RefreshSecret`（与 AccessToken 的 `AccessSecret` 不同）对 JWT 进行 HS256 签名校验和过期时间验证。

- JWT 解析失败（签名错误 / 已过期，errCode == 1）→ 返回 `401002 — RefreshToken无效或已过期`
- 系统性异常（errCode == 2）→ 返回 `500001 — 刷新Token时失败，请稍后重试`
- 校验通过 → 从 Claims 中解析出 `[uuid, telephone, user角色]`

#### 步骤二：Redis 存在性校验

```go
exists, err := u.redisClient.Exists(ctx, "refreshToken:" + uuid).Result()
```

登录时会将 RefreshToken 写入 Redis（key = `refreshToken:{uuid}`，TTL = 7 天）。刷新时检查该 key 是否仍然存在：

- key 不存在 → 用户可能已登出、被封禁（封禁流程会删除此 key）、或 Redis key 已过期 → 返回 `401002`
- Redis 查询异常 → 返回 `500001`
- key 存在 → 继续

#### 步骤三：角色还原

从 JWT Claims 中的角色字符串还原为数值型 `isAdmin`：

| Claims 中的 user 值 | isAdmin 数值 | 含义 |
|---------------------|-------------|------|
| `"user"` | `0` | 普通用户 |
| `"admin"` | `1` | 管理员 |
| `"superAdmin"` | `2` | 超级管理员 |

#### 步骤四：签发全新的 Token 对

```go
accessToken, refreshToken, err := tools.GenerateAccessToken(u.jwtConfig, uuid, telephone, int8(isAdmin))
```

生成一对全新的 Token：

| Token | 签名密钥 | 有效期 | 用途 |
|-------|---------|--------|------|
| AccessToken | `AccessSecret` | 15 分钟 | 请求鉴权 |
| RefreshToken | `RefreshSecret` | 168 小时（7 天） | 续签 AccessToken |

两个 Token 均使用 HS256 算法签名，Issuer 固定为 `easyChat`。

### 4. 响应返回

```go
type RefreshTokenResp struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}
```

返回 `200 — 刷新Token成功`，携带新的 Token 对。

## 双 Token 机制说明

```
登录成功
  │
  ├── 签发 AccessToken  (15min)   ←  用于所有受保护接口的鉴权
  ├── 签发 RefreshToken (7天)     ←  仅用于本接口续签
  └── Redis SET refreshToken:{uuid} = token (TTL 7天)

                    ···15 分钟后 AccessToken 过期···

客户端发现 401001 (AccessToken已过期)
  │
  ▼
调用 /v1/user/refresh 携带 RefreshToken
  │
  ├── JWT 签名校验通过
  ├── Redis 存在性校验通过
  │
  ├── 签发 新 AccessToken  (15min)
  └── 签发 新 RefreshToken (7天)
```

## 核心设计

| 设计点 | 说明 |
|--------|------|
| 双 Secret 隔离 | AccessToken 和 RefreshToken 使用不同的签名密钥，互相无法冒充 |
| Redis 存在性校验 | 仅 JWT 合法不够，还需 Redis 中的 key 存在，使得封禁/登出操作可以通过删除 key 立即使 RefreshToken 失效 |
| 无鉴权中间件 | 接口设计上跳过 AuthHttpMiddleware，避免 AccessToken 过期时的鸡生蛋问题 |
| Token 全量刷新 | 每次刷新同时签发新的 AccessToken 和 RefreshToken，RefreshToken 的有效期从刷新时刻重新开始计算（滑动窗口效果） |

## 调用链路图

```
Client
  │
  ├── Body: { "refresh_token": "eyJhbGci..." }
  │
  ▼
TreaceLoggerMiddleware   →  生成 traceID, 注入 context
  │
  ▼
UserHandler.RefreshToken →  ShouldBindJSON 绑定 refresh_token
  │
  ▼
UserService.RefreshTokenService  →  编排四步操作
  │
  ├── [1] tools.VerifyRefreshToken  →  JWT 签名 + 过期校验
  ├── [2] Redis EXISTS              →  refreshToken:{uuid} 存在性检查
  ├── [3] 角色还原                   →  "user"/"admin"/"superAdmin" → 0/1/2
  └── [4] tools.GenerateAccessToken →  签发新 AccessToken + RefreshToken
  │
  ▼
JsonBack                 →  { code: 200, msg: "刷新Token成功", data: { access_token, refresh_token } }
```

## 涉及的 Redis Key

| Key | 写入时机 | 读取时机 | TTL |
|-----|----------|----------|-----|
| `refreshToken:{uuid}` | 登录成功时 SET | 本接口 EXISTS 校验 | 168 小时（7 天） |
