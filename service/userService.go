package service

import (
	"fmt"
	"im/models"
	"im/utils"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// GetUserList
// @Summary 所有用户
// @Tags 用户模块
// @Success 200 {string} json{"code","message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := models.GetUserList()
	c.JSON(200, gin.H{
		"code":    0, //0成功  -1失败
		"message": "获取用户列表",
		"data":    data,
	})
}

// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [post]
func CreateUser(c *gin.Context) {
	data := models.UserBasic{}
	name := c.PostForm("name")
	password := c.PostForm("password")
	repassword := c.PostForm("repassword")

	if password == "" || name == "" {
		c.JSON(-1, gin.H{
			"code":    -1,
			"message": "用户名和密码不能为空",
			"data":    data,
		})
		return
	}
	if password != repassword {
		c.JSON(-1, gin.H{
			"code":    -1,
			"message": "密码前后不一致",
			"data":    data,
		})
		return
	}

	salt := fmt.Sprintf("%06d", rand.Int31())
	user := models.FindUserByName(name)
	if user.Name != "" {
		c.JSON(-1, gin.H{
			"code":    -1,
			"message": "用户名已注册",
			"data":    user,
		})
		return
	}
	user = models.UserBasic{
		Name:     name,
		Password: password,
	}
	user.Password = utils.MakePassword(password, salt)
	user.Salt = salt
	models.CreateUser(user)
	c.JSON(200, gin.H{
		"code":    0,
		"message": "新增用户成功",
		"data":    user,
	})

}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id query string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /user/deleteUser [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(200, gin.H{
		"code":    0,
		"message": "删除用户成功",
		"data":    user,
	})

}

// UpdateUser
// @Summary 修改用户
// @Tags 用户模块
// @param id formData string false "id"
// @param name formData string false "name"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @param password formData string false "password"
// @Success 200 {string} json{"code","message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.Password = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "修改用户失败",
			"data":    user,
		})
	} else {
		models.UpdateUser(user)
		c.JSON(200, gin.H{
			"code":    0,
			"message": "修改用户成功",
			"data":    user,
		})
	}

}

// Login
// @Summary 用户登录
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/login [post]
func Login(c *gin.Context) {
	data := models.UserBasic{}
	name := c.PostForm("name")
	pwd := c.PostForm("password")
	fmt.Println(pwd)
	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(-1, gin.H{
			"code":    -1,
			"message": "用户名不存在",
			"data":    user,
		})
		return
	}
	b := utils.ValidPassword(pwd, user.Salt, user.Password)
	if !b {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "用户密码校验失败",
			"data":    data,
		})
		return
	}
	User := models.Token(user)
	c.JSON(200, gin.H{
		"code":    0, //0成功  -1失败
		"message": "登录成功",
		"data":    User,
	})
}

//防止跨域站点伪造请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMessage(c *gin.Context) {
	conn, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 启动一个协程，每隔 1s 向客户端发送一次心跳消息

	go func() {

		var (
			err error
		)

		for {

			if err = conn.WriteMessage(websocket.TextMessage, []byte("heartbeat")); err != nil {

				return

			}

			time.Sleep(1 * time.Second)

		}

	}()

	// 得到 websocket 的长链接之后，就可以对客户端传递的数据进行操作了

	for {

		// 通过 websocket 长链接读到的数据可以是 text 文本数据，也可以是二进制 Binary

		_, data, err := conn.ReadMessage()
		if err != nil {

			fmt.Println(err)
			break

		}

		err = conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {

			fmt.Println(err)
			break

		}

	}

	conn.Close()

}

func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	for {
		msg := utils.Subscribe(c, utils.PublishKey)
		fmt.Println("订阅了消息。。。")
		tm := time.Now().Format("2006-01-02 03:04:05")
		m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
		err := ws.WriteMessage(1, []byte(m))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func SearchFriends(c *gin.Context) {
	userid, _ := strconv.Atoi(c.PostForm("userId"))
	users := models.SearchFriend(uint(userid))
	/* c.JSON(200, gin.H{
		"code":    0,
		"message": "查询好友成功",
		"data":    users,
	}) */
	utils.RespOkList(c.Writer, users, len(users))
}

//添加好友
func AddFriend(c *gin.Context) {
	userId, _ := strconv.Atoi(c.PostForm("userId"))
	targetName := c.PostForm("targetName")
	user := models.FindUserById(uint(userId))
	if user.ID == uint(userId) {
		utils.RespFail(c.Writer, "不能添加自己")
		return
	}
	ub := models.SearchFriend(uint(userId))
	for _, v := range ub {
		if v.Name == targetName {
			utils.RespFail(c.Writer, "你们已经是好友了")
			return
		}
	}

	i := models.AddFriend(uint(userId), targetName)
	if i == 0 {
		utils.RespOk(c.Writer, i, "添加好友成功")
	} else {
		utils.RespFail(c.Writer, "添加失败")
	}

}

//建群
func CreateCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.Request.FormValue("ownerId"))
	name := c.Request.FormValue("name")
	community := models.Community{
		OwnerId: uint(ownerId),
		Name:    name,
	}
	i, msg := models.CreateCommunity(community)
	if i == 0 {
		utils.RespOk(c.Writer, i, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

//群列表
func LoadCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.Request.FormValue("ownerId"))
	_, c2 := models.JoinGroups(uint(ownerId))
	data, msg := models.LoadCommunity(uint(ownerId))
	if len(data) != 0 {
		utils.RespList(c.Writer, 0, c2, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

//加群
func JoinGroup(c *gin.Context) {
	comId, _ := strconv.Atoi(c.Request.FormValue("comId"))
	u := models.GroupPeople(uint(comId))

	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	for _, v := range u {
		if v == uint(userId) {
			utils.RespFail(c.Writer, "您已经在群聊中")
			return
		}
	}
	i := models.JoinGroup(uint(userId), uint(comId))
	utils.RespOk(c.Writer, i, "加入群聊成功")
}
