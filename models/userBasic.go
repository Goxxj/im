package models

import (
	"fmt"
	"im/utils"
	"time"

	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Name          string
	Password      string
	Phone         string `valid:"matches(^(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\\d{8}$)"`
	Email         string `valid:"email"`
	Identity      string
	ClientIp      string
	ClientPort    string
	Salt          string
	LoginTime     time.Time
	HeartbeatTime time.Time
	LogoutTime    time.Time
	IsLogout      bool
	DeviceInfo    string
}

//模型类
func (table *UserBasic) TableName() string {
	return "user_basic"
}

//查找用户名
func FindUserByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name =?", name).First(&user)
	return user
}
func FindUserById(id uint) UserBasic {
	user := UserBasic{}
	utils.DB.Where("id =?", id).First(&user)
	return user
}
func FindUserByNameAndPwd(name string, password string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name =? and password=?", name, password).First(&user)
	return user
}

//token
func Token(user UserBasic) UserBasic {
	s := fmt.Sprintf("%d", time.Now().Unix())
	template := utils.MD5Encode(s)
	utils.DB.Model(&user).Update("identity", template)
	return user
}

//查找用户电话
func FindUserByPhone(phone string) *gorm.DB {
	user := UserBasic{}
	return utils.DB.Where("phone =?", phone).First(&user)
}

//查找用户邮件
func FindUserByEmail(email string) *gorm.DB {
	user := UserBasic{}
	return utils.DB.Where("email =?", email).First(&user)
}

//创建用户
func CreateUser(user UserBasic) *gorm.DB {

	return utils.DB.Create(&user)
}

//获取用户列表
func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)

	for _, v := range data {
		fmt.Println(v)
	}
	return data
}

//删除用户
func DeleteUser(user UserBasic) *gorm.DB {

	return utils.DB.Delete(&user)
}

//修改用户
func UpdateUser(user UserBasic) *gorm.DB {

	return utils.DB.Model(&user).Updates(user)
}
