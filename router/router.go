package router

import (
	"im/docs"
	"im/service"

	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	//swagger
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	//静态资源
	r.Static("/asset", "asset/")
	r.LoadHTMLGlob("views/**/*")
	//index
	r.GET("/", service.GetIndex)
	r.GET("/index", service.GetIndex)
	r.GET("/register", service.Register)
	r.GET("/tochat", service.ToChat)
	r.GET("/chat", service.Chat)
	r.POST("/searchFriends", service.SearchFriends)
	//user
	r.POST("/user/login", service.Login)
	r.POST("/user/getUserList", service.GetUserList)
	r.POST("/user/createUser", service.CreateUser)
	r.POST("/user/deleteUser", service.DeleteUser)
	r.POST("/user/updateUser", service.UpdateUser)
	//发送消息
	r.GET("/user/sendMsg", service.SendMessage)
	//上传文件
	r.POST("/attach/upload", service.Upload)
	//加好友
	r.POST("/contact/addfriend", service.AddFriend)
	//建群
	r.POST("/contact/createCommunity", service.CreateCommunity)
	//群列表
	r.POST("/contact/loadcommunity", service.LoadCommunity)
	r.POST("/contact/joinGroup", service.JoinGroup)
	return r
}
