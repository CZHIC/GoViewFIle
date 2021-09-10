package main

import (
	_ "GoViewFile/boot"
	_ "GoViewFile/router"

	"github.com/gogf/gf/frame/g"
)

// @title       `GoViewFile`示例服务API
// @version     1.0
// @description `GoFrame`基础开发框架服务API接口文档。
// @schemes     http
func main() {
	g.Server().Run()
}
