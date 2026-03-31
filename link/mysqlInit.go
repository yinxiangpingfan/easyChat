package link

import (
	"easyChat/global"
	"easyChat/model"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//gorm链接mysql

func InitGorm() *gorm.DB {
	//连接数据库
	// data source name
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", global.Config.Database.User, global.Config.Database.Password, global.Config.Database.Host, global.Config.Database.Port, global.Config.Database.Name)
	db, err := gorm.Open(mysql.Open(dsn), nil)
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&model.UserInfo{}, &model.GroupInfo{}, &model.Message{}, &model.Session{}, &model.UserContact{}, &model.ContactApply{})
	if err != nil {
		panic(err)
	}

	//连接池控制参数
	sqlDB, _ := db.DB()
	//池子里空闲连接的数量上限（超出此上限就把相应的连接关闭掉）
	sqlDB.SetMaxIdleConns(20)
	//最多开这么多连接
	sqlDB.SetMaxOpenConns(100)
	//一个连接最多可使用这么长时间，超时后连接会自动关闭（因为数据库本身可能也对NoActive连接设置了超时时间，我们的应对办法：定期ping，或者SetConnMaxLifetime）
	sqlDB.SetConnMaxLifetime(time.Hour)
	// 执行 Ping 检查
	if err := sqlDB.Ping(); err != nil {
		panic("数据库无法响应（Ping失败）: " + err.Error())
	}
	return db
}
