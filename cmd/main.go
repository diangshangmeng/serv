package main

import (
	"fmt"

	"voucher-platform/config"
	"voucher-platform/model"
	"voucher-platform/router"
	"voucher-platform/service"
	"voucher-platform/util"

	"github.com/gin-gonic/gin"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		return
	}

	err = util.InitLogger(util.LogConfig{
		Level:      config.AppConfig.LogLevel,
		Filename:   config.AppConfig.LogFilename,
		MaxSize:    config.AppConfig.LogMaxSize,
		MaxBackups: config.AppConfig.LogMaxBackups,
		MaxAge:     config.AppConfig.LogMaxAge,
		Compress:   config.AppConfig.LogCompress,
	})
	if err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
	}

	err = config.InitDB()
	if err != nil {
		fmt.Printf("初始化数据库失败: %v\n", err)
		return
	}
	defer config.DB.Close()

	config.InitRedis()

	err = model.AutoMigrate()
	if err != nil {
		fmt.Printf("数据库迁移失败: %v\n", err)
		return
	}

	service.StartOrderTimeoutChecker()

	err = util.RegisterCustomValidators()
	if err != nil {
		fmt.Printf("注册自定义验证器失败: %v\n", err)
		return
	}

	r := gin.Default()
	router.InitRouter(r)

	port := fmt.Sprintf(":%s", config.AppConfig.ServerPort)
	fmt.Printf("服务启动成功，监听端口: %s\n", port)
	r.Run(port)
}
