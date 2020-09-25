package routers

import (
	"github.com/astaxie/beego"

	"github.com/XunleiBlockchain/baas-demo/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/v1/contract/execute", &controllers.ContractController{}, "post:Execute")
	beego.Router("/v1/contract/call", &controllers.ContractController{}, "post:Call")
	beego.Router("/v1/contract/getTx", &controllers.ContractController{}, "post:GetTx")
}
