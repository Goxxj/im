package models

import (
	"fmt"
	"im/utils"

	"gorm.io/gorm"
)

type Community struct {
	gorm.Model
	Name    string
	OwnerId uint
	Img     string
	Desc    string
}

//创建群
func CreateCommunity(community Community) (int, string) {
	if len(community.Name) == 0 {
		return -1, "群名称不能为空"
	}
	if community.OwnerId == 0 {
		return -1, "请先登录"
	}
	if err := utils.DB.Create(&community).Error; err != nil {
		fmt.Println(err)
		return -1, "建群失败"
	}
	JoinGroup(community.OwnerId, community.ID)
	return 0, "建群成功"
}

//查询我创建的群
func LoadCommunity(id uint) ([]*Community, string) {
	data := make([]*Community, 10)
	utils.DB.Where("owner_id=?", id).Find(&data)
	return data, "查询成功"
}

//id查群
func FindCommunity(id uint) Community {
	data := Community{}
	utils.DB.Where("owner_id=?", id).Find(&data)
	return data
}
