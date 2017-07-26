package main

import (
	_ "diangetai/routers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Informational(beego.AppConfig.String("num"))
}

func main() {
	beego.AddFuncMap("index", index)
	beego.Run()
}

func index(in int) (out int) {
	out = in + 1
	return
}
