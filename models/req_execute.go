package models

import (
	"encoding/json"
	"errors"

	"github.com/XunleiBlockchain/tc-libs/common"

	"github.com/XunleiBlockchain/baas-demo/contract"
)

type ReqExecute struct {
	Account  string        `json:"account"`
	Contract string        `json:"contract"`
	Addr     string        `json:"addr"`
	Method   string        `json:"method"`
	Params   []interface{} `json:"params"`
}

func (re *ReqExecute) Parse(rbody []byte) error {
	if err := json.Unmarshal(rbody, re); err != nil {
		return err
	}
	if err := re.Sanity(); err != nil {
		return err
	}
	return nil
}

func (re *ReqExecute) Sanity() error {
	if !common.IsHexAddress(re.Account) {
		return errors.New("account address illegal")
	}
	if !common.IsHexAddress(re.Addr) {
		return errors.New("contract address illegal")
	}
	return nil
}

func (re *ReqExecute) Execute() (interface{}, error) {
	cInstance := contract.Get(re.Contract)
	if cInstance == nil {
		return nil, errors.New("contract instance unsupported")
	}
	cdata, err := cInstance.Data(re.Method, re.Params)
	if err != nil {
		return nil, err
	}
	//from address should be unlocked in sdk-server
	backReqParam := []interface{}{
		map[string]interface{}{
			"from": re.Account,
			"to":   re.Addr,
			"data": cdata,
		},
	}
	var ret interface{}
	if useSDK {
		res, sdkErr := sdkInstance.SendContractTransaction(backReqParam)
		if sdkErr != nil && sdkErr.Code != 0 {
			return nil, sdkErr
		}
		ret = res
	} else {
		res, err := backCall("sendContractTransaction", backReqParam)
		if err != nil {
			return nil, err
		}
		ret = res
	}
	return ret, nil
}
