package main

import (
	"GoViewFile/app/service"
	_ "GoViewFile/boot"
	_ "GoViewFile/router"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcron"
)

// @title       `GoViewFile`示例服务API
// @version     1.0
// @description `GoFrame`基础开发框架服务API接口文档。
// @schemes     http
func main() {
	//每天凌晨两点清理服务器文件
	gcron.Add("0 0 2 * * *", service.ClearFile)
	g.Server().Run()
}
