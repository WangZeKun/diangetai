package routers

import (
	"diangetai/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.Router("/add",&controllers.MainController{},"post:Post")
	beego.Router("/after",&controllers.AfterController{})
	beego.Router("/end",&controllers.AfterController{},"post:End")
	beego.Router("/plain",&controllers.PlainController{})
}
