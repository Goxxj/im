package service

import (
	"fmt"
	"im/utils"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	w := c.Writer
	req := c.Request
	flie, head, err := req.FormFile("file")
	if err != nil {
		utils.RespFail(w, err.Error())
	}
	suffix := ".png"
	ofliName := head.Filename
	tem := strings.Split(ofliName, ".")
	if len(tem) > 1 {
		suffix = "." + tem[len(tem)-1]
	}
	flieName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)
	dstFlie, err2 := os.Create("./asset/upload/" + flieName)
	if err2 != nil {
		utils.RespFail(w, err2.Error())
	}
	_, err3 := io.Copy(dstFlie, flie)
	if err3 != nil {
		utils.RespFail(w, err3.Error())
	}
	url := "./asset/upload/" + flieName
	utils.RespOk(w, url, "发送图片成功")

}
