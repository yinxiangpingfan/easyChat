# KamaChat 数据库表结构

> **注意**: KamaChat 没有在数据库层面设置外键约束，所有关联关系通过应用层的 uuid 逻辑关联。

---

## 1. user_info (用户表)

| 字段 | 类型 | 约束 | 索引 | 说明 |
|------|------|------|------|------|
| id | int64 | primaryKey | PK | 自增id |
| uuid | char(20) | not null | **uniqueIndex** | 用户唯一id |
| nickname | varchar(20) | not null | - | 昵称 |
| telephone | char(11) | not null | **index** | 电话 |
| email | char(30) | - | - | 邮箱 |
| avatar | char(255) | default, not null | - | 头像 |
| gender | int8 | - | - | 性别 (0.男, 1.女) |
| signature | varchar(100) | - | - | 个性签名 |
| password | char(18) | not null | - | 密码 |
| birthday | char(8) | - | - | 生日 |
| created_at | datetime | not null | **index** | 创建时间 |
| deleted_at | datetime | - | - | 删除时间 (软删除) |
| last_online_at | datetime | - | - | 上次登录时间 |
| last_offline_at | datetime | - | - | 最近离线时间 |
| is_admin | int8 | not null | - | 是否管理员 (0.否, 1.是) |
| status | int8 | not null | **index** | 状态 (0.正常, 1.禁用) |

**索引汇总:**
- `uuid` - 唯一索引
- `telephone` - 普通索引
- `created_at` - 普通索引
- `status` - 普通索引

---

## 2. group_info (群组表)

| 字段 | 类型 | 约束 | 索引 | 说明 |
|------|------|------|------|------|
| id | int64 | primaryKey | PK | 自增id |
| uuid | char(20) | not null | **uniqueIndex** | 群组唯一id |
| name | varchar(20) | not null | - | 群名称 |
| notice | varchar(500) | - | - | 群公告 |
| members | json | - | - | 群组成员 |
| member_cnt | int | default: 1 | - | 群人数 |
| owner_id | char(20) | not null | - | 群主uuid |
| add_mode | int8 | default: 0 | - | 加群方式 (0.直接, 1.审核) |
| avatar | char(255) | default, not null | - | 头像 |
| status | int8 | default: 0 | - | 状态 (0.正常, 1.禁用, 2.解散) |
| created_at | datetime | not null | **index** | 创建时间 |
| updated_at | datetime | not null | - | 更新时间 |
| deleted_at | datetime | - | **index** | 删除时间 (软删除) |

**索引汇总:**
- `uuid` - 唯一索引
- `created_at` - 普通索引
- `deleted_at` - 普通索引

---

## 3. session (会话表)

| 字段 | 类型 | 约束 | 索引 | 说明 |
|------|------|------|------|------|
| id | int64 | primaryKey | PK | 自增id |
| uuid | char(20) | - | **uniqueIndex** | 会话uuid |
| send_id | char(20) | not null | **index** | 创建会话人id |
| receive_id | char(20) | not null | **index** | 接受会话人id |
| receive_name | varchar(20) | not null | - | 名称 |
| avatar | char(255) | default, not null | - | 头像 |
| last_message | TEXT | - | - | 最新的消息 |
| last_message_at | datetime | - | - | 最近接收时间 |
| created_at | datetime | - | **index** | 创建时间 |
| deleted_at | datetime | - | **index** | 删除时间 (软删除) |

**索引汇总:**
- `uuid` - 唯一索引
- `send_id` - 普通索引
- `receive_id` - 普通索引
- `created_at` - 普通索引
- `deleted_at` - 普通索引

---

## 4. message (消息表)

| 字段 | 类型 | 约束 | 索引 | 说明 |
|------|------|------|------|------|
| id | int64 | primaryKey | PK | 自增id |
| uuid | char(20) | not null | **uniqueIndex** | 消息uuid |
| session_id | char(20) | not null | **index** | 会话uuid |
| type | int8 | not null | - | 消息类型 (0.文本, 1.语音, 2.文件, 3.通话) |
| content | TEXT | - | - | 消息内容 |
| url | char(255) | - | - | 消息url |
| send_id | char(20) | not null | **index** | 发送者uuid |
| send_name | varchar(20) | not null | - | 发送者昵称 |
| send_avatar | varchar(255) | not null | - | 发送者头像 |
| receive_id | char(20) | not null | **index** | 接受者uuid |
| file_type | char(10) | - | - | 文件类型 |
| file_name | varchar(50) | - | - | 文件名 |
| file_size | char(20) | - | - | 文件大小 |
| status | int8 | not null | - | 状态 (0.未发送, 1.已发送) |
| created_at | datetime | not null | - | 创建时间 |
| send_at | datetime | - | - | 发送时间 |
| av_data | - | - | - | 通话传递数据 |

**索引汇总:**
- `uuid` - 唯一索引
- `session_id` - 普通索引
- `send_id` - 普通索引
- `receive_id` - 普通索引

---

## 5. user_contact (用户联系人表)

| 字段 | 类型 | 约束 | 索引 | 说明 |
|------|------|------|------|------|
| id | int64 | primaryKey | PK | 自增id |
| user_id | char(20) | not null | **index** | 用户唯一id |
| contact_id | char(20) | not null | **index** | 对应联系id |
| contact_type | int8 | not null | - | 联系类型 (0.用户, 1.群聊) |
| status | int8 | not null | - | 联系状态 |
| created_at | datetime | not null | - | 创建时间 |
| update_at | datetime | not null | - | 更新时间 |
| deleted_at | datetime | - | **index** | 删除时间 (软删除) |

**索引汇总:**
- `user_id` - 普通索引
- `contact_id` - 普通索引
- `deleted_at` - 普通索引

**status 状态枚举:**
| 值 | 说明 |
|----|------|
| 0 | 正常 |
| 1 | 拉黑 |
| 2 | 被拉黑 |
| 3 | 删除好友 |
| 4 | 被删除好友 |
| 5 | 被禁言 |
| 6 | 退出群聊 |
| 7 | 被踢出群聊 |

---

## 6. contact_apply (联系人申请表)

| 字段 | 类型 | 约束 | 索引 | 说明 |
|------|------|------|------|------|
| id | int64 | primaryKey | PK | 自增id |
| uuid | char(20) | - | **uniqueIndex** | 申请id |
| user_id | char(20) | not null | **index** | 申请人id |
| contact_id | char(20) | not null | **index** | 被申请id |
| contact_type | int8 | not null | - | 被申请类型 (0.用户, 1.群聊) |
| status | int8 | not null | - | 申请状态 |
| message | varchar(100) | - | - | 申请信息 |
| last_apply_at | datetime | not null | - | 最后申请时间 |
| deleted_at | datetime | - | **index** | 删除时间 (软删除) |

**索引汇总:**
- `uuid` - 唯一索引
- `user_id` - 普通索引
- `contact_id` - 普通索引
- `deleted_at` - 普通索引

**status 状态枚举:**
| 值 | 说明 |
|----|------|
| 0 | 申请中 |
| 1 | 通过 |
| 2 | 拒绝 |
| 3 | 拉黑 |

---

## 索引总览

| 表名 | 唯一索引 | 普通索引 |
|------|---------|---------|
| user_info | uuid | telephone, created_at, status |
| group_info | uuid | created_at, deleted_at |
| session | uuid | send_id, receive_id, created_at, deleted_at |
| message | uuid | session_id, send_id, receive_id |
| user_contact | - | user_id, contact_id, deleted_at |
| contact_apply | uuid | user_id, contact_id, deleted_at |

---

## 外键说明

**KamaChat 没有使用数据库外键约束**，原因可能是：

1. **性能考虑**: 外键约束会影响写入性能
2. **灵活性**: 应用层控制关联逻辑更灵活
3. **分布式场景**: 便于后续分库分表
4. **GORM 习惯**: Go 项目通常在应用层处理关联关系

关联关系通过 `uuid` 字段在应用层维护：

```
user_info.uuid  ──(逻辑关联)──►  user_contact.user_id
user_info.uuid  ──(逻辑关联)──►  session.send_id
group_info.uuid ──(逻辑关联)──►  user_contact.contact_id (当 contact_type=1)
session.uuid    ──(逻辑关联)──►  message.session_id
```

---

## 表关系图

```
┌─────────────┐       ┌─────────────────┐       ┌─────────────┐
│  user_info  │       │  user_contact   │       │ group_info  │
├─────────────┤       ├─────────────────┤       ├─────────────┤
│ uuid (UQ)   │◄──────│ user_id (IDX)   │       │ uuid (UQ)   │
│ nickname    │       │ contact_id(IDX) │──────►│ name        │
│ telephone   │       │ contact_type    │       │ owner_id    │
│ ...         │       │ status          │       │ members     │
└─────────────┘       └─────────────────┘       └─────────────┘
       │                                               │
       │                                               │
       ▼                                               ▼
┌─────────────┐       ┌─────────────────┐       ┌─────────────┐
│   session   │       │    message      │       │contact_apply│
├─────────────┤       ├─────────────────┤       ├─────────────┤
│ uuid (UQ)   │◄──────│ session_id (IDX)│       │ user_id(IDX)│
│ send_id(IDX)│       │ send_id (IDX)   │       │contact_id(IDX)
│receive_id(IDX)      │ receive_id(IDX) │       │ status      │
│ last_message│       │ content         │       │ ...         │
└─────────────┘       └─────────────────┘       └─────────────┘

图例: UQ = 唯一索引, IDX = 普通索引
```

---

## 注意事项

1. **软删除**: 所有表都使用 GORM 的 `DeletedAt` 字段实现软删除
2. **UUID**: 使用 char(20) 存储唯一标识符
3. **时间字段**: 使用 `datetime` 类型，配合 GORM 自动管理
4. **无外键**: 关联关系在应用层维护，不使用数据库外键约束
