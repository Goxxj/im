package models

import (
	"fmt"
	"im/utils"

	"gorm.io/gorm"
)

//人员关系
type Contact struct {
	gorm.Model
	OwnerId  uint //谁的关系
	TargetId uint //对应的谁
	Type     int  //对应的类型 1好友 2群 3
	Desc     string
}

func (table *Contact) TableName() string {
	return "contact"
}

func SearchFriend(userId uint) []UserBasic {
	contacts := make([]Contact, 0)
	objIds := make([]uint64, 0)
	utils.DB.Where("owner_id=? and type=1", userId).Find(&contacts)
	for _, v := range contacts {
		fmt.Println(v)
		objIds = append(objIds, uint64(v.TargetId))
	}
	users := make([]UserBasic, 0)
	utils.DB.Where("id in ?", objIds).Find(&users)
	return users
}

func AddFriend(userId uint, targetName string) int {
	user := UserBasic{}

	if len(targetName) != 0 {
		user = FindUserByName(targetName)

		if user.Name != "" {
			contact := Contact{
				OwnerId:  userId,
				TargetId: user.ID,
				Type:     1,
			}
			utils.DB.Create(&contact)
			contact = Contact{
				OwnerId:  user.ID,
				TargetId: userId,
				Type:     1,
			}
			utils.DB.Create(&contact)
			return 0
		}
		return -1
	}
	return -1
}

//加入的群
func JoinGroups(userId uint) (int, []*Community) {
	contacts := make([]*Contact, 0)
	Communitys := make([]*Community, 0)
	utils.DB.Where("owner_id=? and type =2", userId).Find(&contacts)
	if len(contacts) != 0 {
		for _, v := range contacts {
			community := Community{}
			utils.DB.Where("id=?", v.TargetId).First(&community)
			Communitys = append(Communitys, &community)
		}
		return 0, Communitys
	} else {
		return -1, Communitys
	}

}

//加群
func JoinGroup(userId uint, comID uint) int {
	contact := Contact{
		OwnerId:  userId,
		TargetId: comID,
		Type:     2,
	}
	utils.DB.Create(&contact)
	return 0
}

//群里的人
func GroupPeople(comID uint) []uint {
	contacts := []Contact{}
	utils.DB.Where("target_id =? and type =2", comID).Find(&contacts)
	people := []uint{}
	for _, v := range contacts {
		people = append(people, v.OwnerId)
	}
	return people
}
