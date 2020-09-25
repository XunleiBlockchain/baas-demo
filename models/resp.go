package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
)

//--------------------------------------------------------------------------------------

type backResp struct {
	ID      int         `json:"id"`
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Errcode int         `json:"errcode"`
	Errmsg  string      `json:"errmsg"`
}

func backCall(method string, params interface{}) (interface{}, error) {
	backReqParams := make(map[string]interface{})
	backReqParams["id"] = 0
	backReqParams["jsonrpc"] = "2.0"
	backReqParams["method"] = method
	backReqParams["params"] = params
	backReq := httplib.Post(beego.AppConfig.String("sdk-server"))
	backReq.SetTimeout(60*time.Second, 60*time.Second)
	reqBody, err := json.Marshal(backReqParams)
	if err != nil {
		return nil, err
	}
	backReq.Body(reqBody)
	logs.Debug("req params: " + string(reqBody))
	var ret backResp
	err = backReq.ToJSON(&ret)
	if err != nil {
		return nil, err
	}
	if ret.Errcode != 0 {
		return nil, fmt.Errorf("errcode: %d errmsg: %s", ret.Errcode, ret.Errmsg)
	}
	return ret.Result, nil
}

//--------------------------------------------------------------------------------------

type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func ErrResp(err *ErrWrapper) *Resp {
	return &Resp{
		Code: err.Code(),
		Msg:  err.Error(),
	}
}

func NormalResp(data interface{}) *Resp {
	return &Resp{
		Code: 0,
		Msg:  "success",
		Data: data,
	}
}
