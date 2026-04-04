# 通用 CRUD 工具封装设计

## 1. 设计思路

KamaChat 直接使用 GORM 方法，每个 Service 重复写相似的代码。封装通用 CRUD 工具可以：

- 减少代码重复
- 统一错误处理
- 统一返回格式
- 便于维护和扩展

---

## 2. 架构层次

```
Controller 层 (api/v1/)
      ↓
Service 层 (internal/service/)
      ↓
Repository 层 (internal/repository/)  ← 新增通用 CRUD 封装
      ↓
DAO 层 (internal/dao/)
      ↓
Database (MySQL)
```

---

## 3. 接口设计

### 3.1 通用 Repository 接口

```go
// BaseRepository 通用 CRUD 接口
type BaseRepository[T any] interface {
    // Create 创建记录
    Create(entity *T) error

    // CreateBatch 批量创建
    CreateBatch(entities []*T) error

    // FindByID 根据 ID 查询
    FindByID(id int64) (*T, error)

    // FindByUUID 根据 UUID 查询
    FindByUUID(uuid string) (*T, error)

    // FindOne 根据条件查询单条
    FindOne(conditions map[string]interface{}) (*T, error)

    // FindList 根据条件查询列表
    FindList(conditions map[string]interface{}) ([]*T, error)

    // FindAll 查询所有
    FindAll() ([]*T, error)

    // Update 更新记录
    Update(entity *T) error

    // UpdateFields 更新指定字段
    UpdateFields(id int64, fields map[string]interface{}) error

    // Delete 软删除
    Delete(id int64) error

    // DeleteByUUID 软删除（根据 UUID）
    DeleteByUUID(uuid string) error

    // HardDelete 硬删除
    HardDelete(id int64) error

    // Count 统计数量
    Count(conditions map[string]interface{}) (int64, error)

    // Exists 判断是否存在
    Exists(conditions map[string]interface{}) (bool, error)

    // Paginate 分页查询
    Paginate(page, pageSize int, conditions map[string]interface{}) ([]*T, int64, error)
}
```

### 3.2 基础实现

```go
// BaseRepositoryImpl 通用 CRUD 实现
type BaseRepositoryImpl[T any] struct {
    db    *gorm.DB
    model *T // 用于推断表名
}

func NewBaseRepository[T any](db *gorm.DB) BaseRepository[T] {
    return &BaseRepositoryImpl[T]{db: db}
}
```

---

## 4. 使用示例

### 4.1 定义实体

```go
// internal/model/user_info.go
type UserInfo struct {
    Id       int64          `gorm:"primaryKey"`
    Uuid     string         `gorm:"uniqueIndex;type:char(20)"`
    Nickname string         `gorm:"type:varchar(20)"`
    // ...
}
```

### 4.2 创建具体 Repository

```go
// internal/repository/user_repository.go
type UserRepository interface {
    BaseRepository[model.UserInfo]

    // 扩展方法：用户特有操作
    FindByTelephone(telephone string) (*model.UserInfo, error)
    FindByNickname(nickname string) ([]*model.UserInfo, error)
    UpdatePassword(uuid, password string) error
}

type userRepositoryImpl struct {
    *BaseRepositoryImpl[model.UserInfo]
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepositoryImpl{
        BaseRepositoryImpl: NewBaseRepository[model.UserInfo](db).(*BaseRepositoryImpl[model.UserInfo]),
    }
}

// 实现扩展方法
func (r *userRepositoryImpl) FindByTelephone(telephone string) (*model.UserInfo, error) {
    var user model.UserInfo
    err := r.db.Where("telephone = ?", telephone).First(&user).Error
    return &user, err
}
```

### 4.3 Service 层使用

```go
// internal/service/user_service.go
type UserService struct {
    userRepo repository.UserRepository
}

func (s *UserService) Login(telephone, password string) (*model.UserInfo, error) {
    user, err := s.userRepo.FindByTelephone(telephone)
    if err != nil {
        return nil, errors.New("用户不存在")
    }
    if user.Password != password {
        return nil, errors.New("密码错误")
    }
    return user, nil
}

func (s *UserService) Register(user *model.UserInfo) error {
    return s.userRepo.Create(user)
}
```

---

## 5. 目录结构

```
easyChat/
├── internal/
│   ├── dao/
│   │   └── gorm.go              # 数据库连接
│   ├── model/
│   │   ├── user_info.go
│   │   ├── group_info.go
│   │   └── ...
│   ├── repository/              # 新增：Repository 层
│   │   ├── base_repository.go   # 通用 CRUD 实现
│   │   ├── user_repository.go   # 用户 Repository
│   │   ├── group_repository.go  # 群组 Repository
│   │   └── ...
│   └── service/
│       ├── user_service.go
│       └── ...
```

---

## 6. 核心方法实现说明

### 6.1 Create

```go
func (r *BaseRepositoryImpl[T]) Create(entity *T) error {
    return r.db.Create(entity).Error
}
```

### 6.2 FindByCondition

```go
func (r *BaseRepositoryImpl[T]) FindOne(conditions map[string]interface{}) (*T, error) {
    var entity T
    err := r.db.Where(conditions).First(&entity).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    return &entity, err
}
```

### 6.3 Update

```go
func (r *BaseRepositoryImpl[T]) Update(entity *T) error {
    return r.db.Save(entity).Error
}

func (r *BaseRepositoryImpl[T]) UpdateFields(id int64, fields map[string]interface{}) error {
    var entity T
    return r.db.Model(&entity).Where("id = ?", id).Updates(fields).Error
}
```

### 6.4 Delete (软删除)

```go
func (r *BaseRepositoryImpl[T]) Delete(id int64) error {
    var entity T
    return r.db.Where("id = ?", id).Delete(&entity).Error
}
```

### 6.5 Paginate

```go
func (r *BaseRepositoryImpl[T]) Paginate(page, pageSize int, conditions map[string]interface{}) ([]*T, int64, error) {
    var entities []*T
    var total int64

    offset := (page - 1) * pageSize

    if err := r.db.Model(new(T)).Where(conditions).Count(&total).Error; err != nil {
        return nil, 0, err
    }

    if err := r.db.Where(conditions).Offset(offset).Limit(pageSize).Find(&entities).Error; err != nil {
        return nil, 0, err
    }

    return entities, total, nil
}
```

---

## 7. 可选增强功能

### 7.1 事务支持

```go
type TransactionRepository interface {
    Begin() *gorm.DB
    Commit() error
    Rollback() error
}

// 使用
tx := repo.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

repo.CreateWithTx(tx, user)
repo.CreateWithTx(tx, contact)

tx.Commit()
```

### 7.2 乐观锁

```go
type BaseEntity struct {
    Id      int64 `gorm:"primaryKey"`
    Version int   `gorm:"column:version;default:0"`
}

// 更新时检查版本
func (r *BaseRepositoryImpl[T]) UpdateWithVersion(entity *T) error {
    return r.db.Where("version = ?", entity.Version).
        Updates(entity).Error
}
```

### 7.3 缓存集成

```go
type CachedRepository[T any] struct {
    repo  BaseRepository[T]
    cache redis.Client
}

func (r *CachedRepository[T]) FindByID(id int64) (*T, error) {
    key := fmt.Sprintf("entity:%d", id)

    // 先查缓存
    if cached, err := r.cache.Get(key); err == nil {
        return cached, nil
    }

    // 再查数据库
    entity, err := r.repo.FindByID(id)
    if err != nil {
        return nil, err
    }

    // 写入缓存
    r.cache.Set(key, entity, time.Minute*10)

    return entity, nil
}
```

---

## 8. 与 KamaChat 对比

| 对比项 | KamaChat 方式 | 封装后 |
|--------|--------------|--------|
| 代码量 | 每个 Service 重复写 | 统一封装，减少 60%+ |
| 错误处理 | 分散在各处 | 统一处理 |
| 可测试性 | 需要 mock GORM | 可 mock Repository 接口 |
| 扩展性 | 直接改 Service | 实现 Repository 接口即可 |
| 复杂度 | 低 | 稍高（需要理解接口和泛型） |

---

## 9. 注意事项

1. **Go 版本要求**: 泛型需要 Go 1.18+
2. **复杂查询**: 复杂的业务查询仍需在具体 Repository 中实现
3. **性能**: 简单封装不会影响性能，GORM 底层一样
4. **学习成本**: 团队需要理解 Repository 模式和泛型

---

## 10. 是否需要封装？

### 适合封装的场景
- 项目较大，实体多（10+ 个表）
- 团队协作，需要统一规范
- 需要 mock 测试

### 不需要封装的场景
- 项目小，实体少
- 快速开发原型
- 团队不熟悉 Repository 模式

**建议**: 复刻阶段可以先不封装，等基础功能完成后，再考虑重构封装。
