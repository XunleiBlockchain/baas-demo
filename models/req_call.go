package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/XunleiBlockchain/tc-libs/common"

	"github.com/XunleiBlockchain/baas-demo/contract"
)

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
		map[string]string{
			"from": rc.Account,
			"to":   rc.Addr,
			"data": cdata,
		},
	}
	ret, err := backCall("call", backReqParam)
	if err != nil {
		return nil, err
	}
	data, ok := ret.(string)
	if !ok {
		return nil, fmt.Errorf("ret type not match. want string got: %s", reflect.TypeOf(ret))
	}
	return cInstance.Result(rc.Method, data)
}
