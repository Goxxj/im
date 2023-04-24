package main

import (
	"im/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:qwe123@tcp(127.0.0.1:3306)/im?charset=utf8&parseTime=True&loc=Local" // 数据库配置 root账号  0000密码  shop数据库  utf8mb4编码
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// db.AutoMigrate(&models.UserBasic{})
	db.AutoMigrate(&models.Message{})
	// db.AutoMigrate(&models.Contact{})
	// db.AutoMigrate(&models.GroupBasic{})

	/* user := &models.UserBasic{}
	user.Name = "huang"
	db.Create(user)
	db.First(user, 1) */ // 根据整型主键查找
	// db.First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录

	// Update - 将 product 的 price 更新为 200
	// db.Model(user).Update("PassWord", "1234")
	// Update - 更新多个字段
	// db.Model(user).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
	// db.Model(user).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - 删除 product
	// db.Delete(user, 1)
}
