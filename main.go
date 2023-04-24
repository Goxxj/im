package main

import (
	"im/router"
	"im/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMysql()
	// utils.InitRedis()

	r := router.Router()

	r.Run()
}
