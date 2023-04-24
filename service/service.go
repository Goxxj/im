package service

import (
	"im/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetIndex
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	// t, err := template.ParseFiles("index.html")
	// if err != nil {
	// 	panic(err)
	// }

	// t.Execute(c.Writer, t.Name())
	c.HTML(200, "index.html", gin.H{})
	/* c.JSON(200, gin.H{
		"message": "welcome!!",
	}) */
}
func Register(c *gin.Context) {
	c.HTML(200, "register.html", gin.H{})
}
func ToChat(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Query("userId"))
	token := c.Query("token")
	user := models.UserBasic{}
	user.ID = uint(userId)
	user.Identity = token
	c.HTML(200, "chat/index.tmpl", user)
}

func Chat(c *gin.Context) {
	models.Chat(c)
}
