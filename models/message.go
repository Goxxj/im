package models

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

//消息
type Message struct {
	gorm.Model
	FromId   int64  //发送者
	TargetId int64  //接收者
	Type     int    //发送类型 群聊 私聊 广播
	Media    int    //消息类型 文字 图片 表情
	Context  string //消息内容
	Pic      string
	Url      string
	Desc     string
	Amount   int //其他数字统计
}

func (table *Message) TableName() string {
	return "message"
}

// 发送消息的类型
type SendMsg struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
}

// 回复的消息
type ReplyMsg struct {
	From    string `json:"from"`
	Code    int    `json:"code"`
	Content string `json:"content"`
}

// 用户类
type Client struct {
	userId string
	// SendID string
	Socket *websocket.Conn
	Send   chan []byte
}

// 广播类，包括广播内容和源用户
type Broadcast struct {
	Client  *Client
	Message []byte
	Type    int
}

// 用户管理
type ClientManager struct {
	Clients    map[string]*Client
	Broadcast  chan *Broadcast
	Reply      chan *Client
	Register   chan *Client
	Unregister chan *Client
}

var Manager = ClientManager{
	Clients:    make(map[string]*Client), // 参与连接的用户，出于性能的考虑，需要设置最大连接数
	Broadcast:  make(chan *Broadcast),
	Register:   make(chan *Client, 50),
	Reply:      make(chan *Client),
	Unregister: make(chan *Client, 50),
}
var (
	up = websocket.Upgrader{
		// 检查区域 可以自行设置是POST 或者GET请求 还有URL等信息 这里直接设置表示都接受
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

var clientMap map[int64]*Client = make(map[int64]*Client)
var rwLocker sync.RWMutex

//升级器函数
func Chat(c *gin.Context) {
	Id := c.Query("userId") // 自己的id
	// toUid := c.Query("toUid") // 对方的id
	userId, _ := strconv.ParseInt(Id, 10, 64)
	conn, err := up.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("upgrede failed")
		return
	} else {
		fmt.Println("now is websocket")
	}

	client := &Client{
		userId: string(userId),
		// SendID: createId(toUid, uid),
		Socket: conn,
		Send:   make(chan []byte),
	}
	rwLocker.Lock()
	clientMap[userId] = client
	rwLocker.Unlock()
	Manager.Register <- client

	//开启协程
	go client.Read()
	go client.Write()

	// defer conn.Close()
}

//读取消息
func (c *Client) Read() {
	defer func() { // 避免忘记关闭，所以要加上close
		Manager.Unregister <- c
		_ = c.Socket.Close()
	}()
	for {
		c.Socket.PongHandler()
		_, p, err := c.Socket.ReadMessage()
		if err != nil {
			log.Println("数据格式不正确1", err)
			Manager.Unregister <- c
			_ = c.Socket.Close()
			break
		}
		c.Send <- p
		msg := Message{}
		err2 := json.Unmarshal(p, &msg)
		if err2 != nil {
			log.Println("数据格式不正确2", err2)
		}
		switch msg.Type {
		case 1:
			sendMsg(msg.TargetId, p)
			/* case 2:
				sendGroup()
			case 3:
				sendAllMsg() */
		}
	}
}

//接收消息
func (c *Client) Write() {
	defer func() {
		_ = c.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			log.Println(c.userId, "接受消息:", string(message))
			_ = c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

//发送私聊消息给TargetId
func sendMsg(TargetId int64, msg []byte) {
	rwLocker.RLock()
	client, ok := clientMap[TargetId]
	rwLocker.RUnlock()
	if ok {
		client.Send <- msg
	}
}
