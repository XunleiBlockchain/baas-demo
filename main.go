package main

import (
	_ "github.com/XunleiBlockchain/baas-demo/routers"

	// inport contracts to do `init`
	_ "github.com/XunleiBlockchain/baas-demo/contract/ERC20"
	_ "github.com/XunleiBlockchain/baas-demo/contract/KVDB"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func main() {
	logs.SetLogger(logs.AdapterFile, `{"filename":"baas-demo.log"}`)
	beego.Run()
}
