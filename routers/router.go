package routers

import (
	"flapper/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/api/ai-agent/hotel-json-maker", &controllers.TestController{})
}
