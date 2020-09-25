package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/XunleiBlockchain/baas-demo/models"
)

type ContractController struct {
	beego.Controller
}

func (this *ContractController) Execute() {
	logs.Debug("Execute")
	object := &models.ReqExecute{}
	if err := object.Parse(this.Ctx.Input.RequestBody); err != nil {
		this.errResp(models.ErrParams.Join(err))
		return
	}
	data, err := object.Execute()
	if err != nil {
		this.errResp(models.ErrServer.Join(err))
		return
	}
	this.normalResp(data)
}

func (this *ContractController) Call() {
	logs.Debug("Call")
	object := &models.ReqCall{}
	if err := object.Parse(this.Ctx.Input.RequestBody); err != nil {
		this.errResp(models.ErrParams.Join(err))
		return
	}
	data, err := object.Call()
	if err != nil {
		this.errResp(models.ErrServer.Join(err))
		return
	}
	this.normalResp(data)
}

func (this *ContractController) GetTx() {
	logs.Debug("GetTx")
	object := &models.ReqGetTx{}
	if err := object.Parse(this.Ctx.Input.RequestBody); err != nil {
		this.errResp(models.ErrParams.Join(err))
		return
	}
	data, err := object.GetTx()
	if err != nil {
		this.errResp(models.ErrServer.Join(err))
		return
	}
	this.normalResp(data)
}

func (this *ContractController) errResp(err *models.ErrWrapper) {
	this.Data["json"] = models.ErrResp(err)
	this.ServeJSON()
}

func (this *ContractController) normalResp(data interface{}) {
	this.Data["json"] = models.NormalResp(data)
	this.ServeJSON()
}
