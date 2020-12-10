package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	sdk "github.com/XunleiBlockchain/baas-sdk-go"

	"github.com/XunleiBlockchain/tc-libs/common"

	"github.com/XunleiBlockchain/baas-demo/contract"
)

var (
	sdkInstance sdk.SDK
)

func InitSDK(s sdk.SDK) {
	sdkInstance = s
}

type ReqCall struct {
	Account  string        `json:"account"`
	Contract string        `json:"contract"`
	Addr     string        `json:"addr"`
	Method   string        `json:"method"`
	Params   []interface{} `json:"params"`
}

func (rc *ReqCall) Parse(rbody []byte) error {
	if err := json.Unmarshal(rbody, rc); err != nil {
		return err
	}
	if err := rc.Sanity(); err != nil {
		return err
	}
	return nil
}

func (rc *ReqCall) Sanity() error {
	if !common.IsHexAddress(rc.Account) {
		return errors.New("account address illegal")
	}
	if !common.IsHexAddress(rc.Addr) {
		return errors.New("contract address illegal")
	}
	return nil
}

func (rc *ReqCall) Call() (interface{}, error) {
	cInstance := contract.Get(rc.Contract)
	if cInstance == nil {
		return nil, errors.New("contract instance unsupported")
	}
	cdata, err := cInstance.Data(rc.Method, rc.Params)
	if err != nil {
		return nil, err
	}
	backReqParam := []interface{}{
		map[string]interface{}{
			"from": rc.Account,
			"to":   rc.Addr,
			"data": cdata,
		},
	}
	ret, sdkErr := sdkInstance.Call(backReqParam)
	if sdkErr != nil && sdkErr.Code != 0 {
		return nil, sdkErr
	}
	data, ok := ret.(string)
	if !ok {
		return nil, fmt.Errorf("ret type not match. want string got: %s", reflect.TypeOf(ret))
	}
	return cInstance.Result(rc.Method, data)
}
