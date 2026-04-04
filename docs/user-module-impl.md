# KamaChat 用户管理模块实现详解

## 目录

1. [模块架构](#1-模块架构)
2. [数据模型](#2-数据模型)
3. [用户注册接口](#3-用户注册接口)
4. [用户登录接口](#4-用户登录接口)
5. [用户信息管理](#5-用户信息管理)
6. [管理员功能](#6-管理员功能)
7. [短信服务集成](#7-短信服务集成)
8. [API 路由汇总](#8-api-路由汇总)
9. [与复刻计划对比](#9-与复刻计划对比)

---

## 1. 模块架构

### 1.1 文件结构

```
KamaChat/
├── api/v1/
│   └── user_info_controller.go    # Controller 层：处理 HTTP 请求
├── internal/
│   ├── model/
│   │   └── user_info.go           # Model 层：数据模型定义
│   ├── dto/
│   │   ├── request/
│   │   │   ├── login_request.go
│   │   │   ├── register_request.go
│   │   │   └── ...
│   │   └── respond/
│   │       ├── login_respond.go
│   │       ├── register_respond.go
│   │       └── ...
│   ├── service/
│   │   ├── gorm/
│   │   │   └── user_info_service.go # Service 层：业务逻辑
│   │   ├── redis/
│   │   │   └── redis_service.go     # Redis 服务
│   │   └── sms/
│   │       └── auth_code_service.go  # 短信服务
│   ├── dao/
│   │   └── gorm.go                  # 数据库连接
│   └── https_server/
│       └── https_server.go          # 路由配置
└── pkg/
    ├── constants/
    │   └── constants.go             # 常量定义
    └── enum/
        └── user_info/
            └── user_status_enum/    # 用户状态枚举
```

### 1.2 调用流程

```
HTTP Request
    ↓
Controller (user_info_controller.go)
    ↓ 解析请求参数
Service (user_info_service.go)
    ↓ 业务逻辑处理
DAO (gorm.go) + Redis (redis_service.go)
    ↓ 数据访问
Database / Redis
    ↓
Controller
    ↓ 统一响应格式
HTTP Response
```

---

## 2. 数据模型

### 2.1 UserInfo 实体

**文件**: `internal/model/user_info.go`

```go
type UserInfo struct {
    Id            int64          `gorm:"column:id;primaryKey;comment:自增id"`
    Uuid          string         `gorm:"column:uuid;uniqueIndex;type:char(20);comment:用户唯一id"`
    Nickname      string         `gorm:"column:nickname;type:varchar(20);not null;comment:昵称"`
    Telephone     string         `gorm:"column:telephone;index;not null;type:char(11);comment:电话"`
    Email         string         `gorm:"column:email;type:char(30);comment:邮箱"`
    Avatar        string         `gorm:"column:avatar;type:char(255);default:...;comment:头像"`
    Gender        int8           `gorm:"column:gender;comment:性别，0.男，1.女"`
    Signature     string         `gorm:"column:signature;type:varchar(100);comment:个性签名"`
    Password      string         `gorm:"column:password;type:char(18);not null;comment:密码"`
    Birthday      string         `gorm:"column:birthday;type:char(8);comment:生日"`
    CreatedAt     time.Time      `gorm:"column:created_at;index;type:datetime;not null;comment:创建时间"`
    DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;comment:删除时间"`
    LastOnlineAt  sql.NullTime   `gorm:"column:last_online_at;type:datetime;comment:上次登录时间"`
    LastOfflineAt sql.NullTime   `gorm:"column:last_offline_at;type:datetime;comment:最近离线时间"`
    IsAdmin       int8           `gorm:"column:is_admin;not null;comment:是否是管理员"`
    Status        int8           `gorm:"column:status;index;not null;comment:状态，0.正常，1.禁用"`
}
```

### 2.2 用户状态枚举

**文件**: `pkg/enum/user_info/user_status_enum/`

```go
const (
    NORMAL  = 0  // 正常
    DISABLE = 1  // 禁用
)
```

### 2.3 DTO 定义

#### 请求 DTO

```go
// RegisterRequest 注册请求
type RegisterRequest struct {
    Telephone string `json:"telephone"`  // 手机号
    Password  string `json:"password"`   // 密码
    Nickname  string `json:"nickname"`   // 昵称
    SmsCode   string `json:"sms_code"`   // 短信验证码
}

// LoginRequest 登录请求
type LoginRequest struct {
    Telephone string `json:"telephone"`
    Password  string `json:"password"`
}

// SmsLoginRequest 短信登录请求
type SmsLoginRequest struct {
    Telephone string `json:"telephone"`
    SmsCode   string `json:"sms_code"`
}
```

#### 响应 DTO

```go
// LoginRespond / RegisterRespond 登录/注册响应
type RegisterRespond struct {
    Uuid      string `json:"uuid"`
    Nickname  string `json:"nickname"`
    Telephone string `json:"telephone"`
    Avatar    string `json:"avatar"`
    Email     string `json:"email"`
    Gender    int8   `json:"gender"`
    Birthday  string `json:"birthday"`
    Signature string `json:"signature"`
    CreatedAt string `json:"created_at"`
    IsAdmin   int8   `json:"is_admin"`
    Status    int8   `json:"status"`
}
```

---

## 3. 用户注册接口

### 3.1 接口信息

| 项目 | 说明 |
|------|------|
| URL | `POST /register` |
| 功能 | 用户注册 |
| 认证 | 无需认证 |

### 3.2 实现流程

```
┌─────────────────────────────────────────────────────────────────┐
│                        注册流程                                  │
├─────────────────────────────────────────────────────────────────┤
│  1. 解析请求参数 (Controller)                                    │
│     ↓                                                           │
│  2. 从 Redis 获取验证码                                          │
│     key: "auth_code_" + telephone                               │
│     ↓                                                           │
│  3. 校验验证码是否正确                                           │
│     不正确 → 返回 "验证码不正确"                                  │
│     ↓                                                           │
│  4. 删除已使用的验证码                                           │
│     ↓                                                           │
│  5. 检查手机号是否已注册                                          │
│     已存在 → 返回 "该电话已经存在"                                │
│     ↓                                                           │
│  6. 生成用户数据                                                 │
│     - UUID: "U" + 日期 + 随机数                                  │
│     - 默认头像                                                   │
│     - 创建时间                                                   │
│     - 状态: NORMAL                                               │
│     ↓                                                           │
│  7. 写入数据库                                                   │
│     ↓                                                           │
│  8. 返回用户信息                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 3.3 核心代码

**Controller 层**:

```go
func Register(c *gin.Context) {
    // 1. 解析请求
    var registerReq request.RegisterRequest
    if err := c.BindJSON(&registerReq); err != nil {
        c.JSON(http.StatusOK, gin.H{
            "code":    500,
            "message": constants.SYSTEM_ERROR,
        })
        return
    }

    // 2. 调用 Service
    message, userInfo, ret := gorm.UserInfoService.Register(registerReq)

    // 3. 返回响应
    JsonBack(c, message, ret, userInfo)
}
```

**Service 层**:

```go
func (u *userInfoService) Register(registerReq request.RegisterRequest) (string, *respond.RegisterRespond, int) {
    // 1. 验证码校验
    key := "auth_code_" + registerReq.Telephone
    code, err := myredis.GetKey(key)
    if err != nil {
        return constants.SYSTEM_ERROR, nil, -1
    }
    if code != registerReq.SmsCode {
        return "验证码不正确，请重试", nil, -2
    }

    // 2. 删除验证码
    myredis.DelKeyIfExists(key)

    // 3. 检查手机号是否存在
    message, ret := u.checkTelephoneExist(registerReq.Telephone)
    if ret != 0 {
        return message, nil, ret
    }

    // 4. 创建用户
    var newUser model.UserInfo
    newUser.Uuid = "U" + random.GetNowAndLenRandomString(11)  // 生成 UUID
    newUser.Telephone = registerReq.Telephone
    newUser.Password = registerReq.Password  // 注意：未加密！
    newUser.Nickname = registerReq.Nickname
    newUser.Avatar = "https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png"
    newUser.CreatedAt = time.Now()
    newUser.Status = user_status_enum.NORMAL

    // 5. 保存到数据库
    res := dao.GormDB.Create(&newUser)
    if res.Error != nil {
        return constants.SYSTEM_ERROR, nil, -1
    }

    // 6. 构造响应
    registerRsp := &respond.RegisterRespond{
        Uuid:      newUser.Uuid,
        Telephone: newUser.Telephone,
        Nickname:  newUser.Nickname,
        // ...
    }

    return "注册成功", registerRsp, 0
}
```

### 3.4 ⚠️ 注意事项

| 问题 | KamaChat 实现 | 建议改进 |
|------|--------------|----------|
| 密码加密 | **未加密**，明文存储 | 使用 bcrypt 加密 |
| UUID 格式 | 自定义 "U" + 日期 + 随机数 | 可用标准 UUID |
| 验证码有效期 | 1 分钟 | 建议 5 分钟 |
| 返回 Token | 无 | 应返回 JWT Token |

---

## 4. 用户登录接口

### 4.1 密码登录

**URL**: `POST /login`

**流程**:

```
1. 根据手机号查询用户
2. 检查用户是否存在
3. 比对密码
4. 返回用户信息
```

**核心代码**:

```go
func (u *userInfoService) Login(loginReq request.LoginRequest) (string, *respond.LoginRespond, int) {
    // 1. 查询用户
    var user model.UserInfo
    res := dao.GormDB.First(&user, "telephone = ?", loginReq.Telephone)

    if errors.Is(res.Error, gorm.ErrRecordNotFound) {
        return "用户不存在，请注册", nil, -2
    }

    // 2. 验证密码
    if user.Password != loginReq.Password {  // 明文比对！
        return "密码不正确，请重试", nil, -2
    }

    // 3. 构造响应
    loginRsp := &respond.LoginRespond{
        Uuid:      user.Uuid,
        Telephone: user.Telephone,
        Nickname:  user.Nickname,
        // ...
    }

    return "登陆成功", loginRsp, 0
}
```

### 4.2 短信验证码登录

**URL**: `POST /user/smsLogin`

**流程**:

```
1. 根据手机号查询用户
2. 从 Redis 获取验证码
3. 校验验证码
4. 返回用户信息
```

**核心代码**:

```go
func (u *userInfoService) SmsLogin(req request.SmsLoginRequest) (string, *respond.LoginRespond, int) {
    // 1. 查询用户
    var user model.UserInfo
    res := dao.GormDB.First(&user, "telephone = ?", req.Telephone)
    if errors.Is(res.Error, gorm.ErrRecordNotFound) {
        return "用户不存在，请注册", nil, -2
    }

    // 2. 验证码校验
    key := "auth_code_" + req.Telephone
    code, err := myredis.GetKey(key)
    if code != req.SmsCode {
        return "验证码不正确，请重试", nil, -2
    }
    myredis.DelKeyIfExists(key)

    // 3. 返回响应
    return "登陆成功", &respond.LoginRespond{...}, 0
}
```

### 4.3 ⚠️ 注意事项

| 问题 | KamaChat 实现 | 建议改进 |
|------|--------------|----------|
| JWT Token | **未生成** | 登录成功应生成 JWT |
| 密码比对 | 明文比对 | bcrypt.CompareHashAndPassword |
| 登录时间 | 未更新 | 应更新 LastOnlineAt |

---

## 5. 用户信息管理

### 5.1 获取用户信息

**URL**: `POST /user/getUserInfo`

**请求**:
```json
{
    "uuid": "U20260331123"
}
```

**流程**:
```
1. 先从 Redis 缓存读取
2. 缓存未命中则查数据库
3. 返回用户信息
```

**核心代码**:

```go
func (u *userInfoService) GetUserInfo(uuid string) (string, *respond.GetUserInfoRespond, int) {
    // 1. 尝试从 Redis 获取
    rspString, err := myredis.GetKeyNilIsErr("user_info_" + uuid)

    if errors.Is(err, redis.Nil) {
        // 2. 缓存未命中，查数据库
        var user model.UserInfo
        dao.GormDB.Where("uuid = ?", uuid).Find(&user)

        rsp := respond.GetUserInfoRespond{
            Uuid:      user.Uuid,
            Telephone: user.Telephone,
            // ...
        }
        return "获取用户信息成功", &rsp, 0
    }

    // 3. 缓存命中
    var rsp respond.GetUserInfoRespond
    json.Unmarshal([]byte(rspString), &rsp)
    return "获取用户信息成功", &rsp, 0
}
```

### 5.2 更新用户信息

**URL**: `POST /user/updateUserInfo`

**请求**:
```json
{
    "uuid": "U20260331123",
    "nickname": "新昵称",
    "email": "new@email.com",
    "birthday": "20000101",
    "signature": "个性签名",
    "avatar": "https://..."
}
```

**核心代码**:

```go
func (u *userInfoService) UpdateUserInfo(updateReq request.UpdateUserInfoRequest) (string, int) {
    // 1. 查询用户
    var user model.UserInfo
    dao.GormDB.First(&user, "uuid = ?", updateReq.Uuid)

    // 2. 更新非空字段
    if updateReq.Email != "" {
        user.Email = updateReq.Email
    }
    if updateReq.Nickname != "" {
        user.Nickname = updateReq.Nickname
    }
    // ... 其他字段

    // 3. 保存
    dao.GormDB.Save(&user)

    return "修改用户信息成功", 0
}
```

### 5.3 获取用户列表

**URL**: `POST /user/getUserInfoList`

**请求**:
```json
{
    "owner_id": "U20260331123"  // 排除当前用户
}
```

**核心代码**:

```go
func (u *userInfoService) GetUserInfoList(ownerId string) (string, []respond.GetUserListRespond, int) {
    var users []model.UserInfo
    // Unscoped() 包含已软删除的用户
    dao.GormDB.Unscoped().Where("uuid != ?", ownerId).Find(&users)

    var rsp []respond.GetUserListRespond
    for _, user := range users {
        rp := respond.GetUserListRespond{
            Uuid:      user.Uuid,
            Telephone: user.Telephone,
            Nickname:  user.Nickname,
            Status:    user.Status,
            IsAdmin:   user.IsAdmin,
            IsDeleted: user.DeletedAt.Valid,  // 是否已删除
        }
        rsp = append(rsp, rp)
    }

    return "获取用户列表成功", rsp, 0
}
```

---

## 6. 管理员功能

### 6.1 启用用户

**URL**: `POST /user/ableUsers`

**请求**:
```json
{
    "uuid_list": ["U20260331123", "U20260331124"]
}
```

**核心代码**:

```go
func (u *userInfoService) AbleUsers(uuidList []string) (string, int) {
    var users []model.UserInfo
    dao.GormDB.Model(model.UserInfo{}).Where("uuid in (?)", uuidList).Find(&users)

    for _, user := range users {
        user.Status = user_status_enum.NORMAL
        dao.GormDB.Save(&user)
    }

    return "启用用户成功", 0
}
```

### 6.2 禁用用户

**URL**: `POST /user/disableUsers`

**额外操作**: 禁用用户时会同时软删除该用户的所有会话

```go
func (u *userInfoService) DisableUsers(uuidList []string) (string, int) {
    for _, user := range users {
        user.Status = user_status_enum.DISABLE
        dao.GormDB.Save(&user)

        // 软删除相关会话
        var sessionList []model.Session
        dao.GormDB.Where("send_id = ? or receive_id = ?", user.Uuid, user.Uuid).Find(&sessionList)
        for _, session := range sessionList {
            session.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
            dao.GormDB.Save(&session)
        }
    }
    return "禁用用户成功", 0
}
```

### 6.3 删除用户

**URL**: `POST /user/deleteUsers`

**级联删除**: 删除用户时会同时软删除：
- 用户的所有会话
- 用户的所有联系人关系
- 用户的所有申请记录

```go
func (u *userInfoService) DeleteUsers(uuidList []string) (string, int) {
    for _, user := range users {
        // 1. 软删除用户
        user.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
        dao.GormDB.Save(&user)

        // 2. 软删除会话
        // 3. 软删除联系人
        // 4. 软删除申请记录
        // ... (详见源码)
    }
    return "删除用户成功", 0
}
```

### 6.4 设置管理员

**URL**: `POST /user/setAdmin`

**请求**:
```json
{
    "uuid_list": ["U20260331123"],
    "is_admin": 1  // 0=取消管理员, 1=设为管理员
}
```

```go
func (u *userInfoService) SetAdmin(uuidList []string, isAdmin int8) (string, int) {
    var users []model.UserInfo
    dao.GormDB.Where("uuid = (?)", uuidList).Find(&users)

    for _, user := range users {
        user.IsAdmin = isAdmin
        dao.GormDB.Save(&user)
    }
    return "设置管理员成功", 0
}
```

---

## 7. 短信服务集成

### 7.1 发送验证码

**URL**: `POST /user/sendSmsCode`

**请求**:
```json
{
    "telephone": "13800138000"
}
```

### 7.2 阿里云 SMS 集成

**文件**: `internal/service/sms/auth_code_service.go`

**依赖**:
```
github.com/alibabacloud-go/dysmsapi-20170525/v4/client
```

**核心代码**:

```go
func VerificationCode(telephone string) (string, int) {
    // 1. 创建阿里云 SMS 客户端
    client, err := createClient()

    // 2. 检查是否已有验证码（防止频繁发送）
    key := "auth_code_" + telephone
    code, _ := redis.GetKey(key)
    if code != "" {
        return "目前还不能发送验证码，请输入已发送的验证码", -2
    }

    // 3. 生成 6 位随机验证码
    code = strconv.Itoa(random.GetRandomInt(6))

    // 4. 存入 Redis，1 分钟有效期
    redis.SetKeyEx(key, code, time.Minute)

    // 5. 调用阿里云 API 发送短信
    sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
        SignName:      tea.String("阿里云短信测试"),
        TemplateCode:  tea.String("SMS_154950909"),
        PhoneNumbers:  tea.String(telephone),
        TemplateParam: tea.String("{\"code\":\"" + code + "\"}"),
    }

    client.SendSmsWithOptions(sendSmsRequest, runtime)

    return "验证码发送成功", 0
}
```

### 7.3 Redis 验证码存储

| Key 格式 | 值 | 有效期 |
|---------|---|--------|
| `auth_code_13800138000` | `123456` | 1 分钟 |

---

## 8. API 路由汇总

| 方法 | 路由 | 功能 | 权限 |
|------|------|------|------|
| POST | `/register` | 用户注册 | 公开 |
| POST | `/login` | 密码登录 | 公开 |
| POST | `/user/smsLogin` | 验证码登录 | 公开 |
| POST | `/user/sendSmsCode` | 发送验证码 | 公开 |
| POST | `/user/getUserInfo` | 获取用户信息 | 登录 |
| POST | `/user/updateUserInfo` | 更新用户信息 | 登录 |
| POST | `/user/getUserInfoList` | 获取用户列表 | 管理员 |
| POST | `/user/ableUsers` | 启用用户 | 管理员 |
| POST | `/user/disableUsers` | 禁用用户 | 管理员 |
| POST | `/user/deleteUsers` | 删除用户 | 管理员 |
| POST | `/user/setAdmin` | 设置管理员 | 管理员 |

---

## 9. 与复刻计划对比

### 9.1 已实现功能对比

| 复刻计划要求 | KamaChat 实现 | 状态 |
|-------------|--------------|------|
| 密码加密 (bcrypt) | ❌ 明文存储 | **不符合** |
| UUID 生成 | ✅ 自定义格式 | 符合（格式不同） |
| 手机号唯一性校验 | ✅ 有 | 符合 |
| 短信验证码校验 | ✅ Redis 存储 | 符合 |
| JWT Token 生成 | ❌ 未实现 | **不符合** |
| 更新最后登录时间 | ❌ 未实现 | **不符合** |
| 用户信息管理 | ✅ 完整实现 | 符合 |
| 管理员功能 | ✅ 完整实现 | 符合 |
| 阿里云 SMS 集成 | ✅ 完整实现 | 符合 |

### 9.2 建议改进项

| 问题 | 建议改进方案 |
|------|-------------|
| 密码明文存储 | 使用 `golang.org/x/crypto/bcrypt` 加密 |
| 无 JWT Token | 使用 `github.com/golang-jwt/jwt` 生成 Token |
| 验证码有效期短 | 从 1 分钟延长到 5 分钟 |
| 无权限校验中间件 | 添加 JWT 认证中间件 |

### 9.3 密码加密示例

```go
import "golang.org/x/crypto/bcrypt"

// 加密密码
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

// 验证密码
func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

### 9.4 JWT Token 示例

```go
import "github.com/golang-jwt/jwt/v5"

// 生成 Token
func GenerateToken(uuid string) (string, error) {
    claims := jwt.MapClaims{
        "uuid": uuid,
        "exp":  time.Now().Add(7 * 24 * time.Hour).Unix(),  // 7 天有效期
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte("your-secret-key"))
}

// 解析 Token
func ParseToken(tokenString string) (string, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return []byte("your-secret-key"), nil
    })
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return claims["uuid"].(string), nil
    }
    return "", err
}
```

---

## 10. 总结

KamaChat 用户管理模块实现了完整的用户 CRUD 功能，但存在以下问题需要改进：

1. **安全性**: 密码明文存储是严重的安全隐患
2. **认证**: 缺少 JWT Token 机制
3. **验证码**: 有效期过短（1 分钟）

复刻时建议在保持功能不变的前提下，增强安全性和认证机制。
