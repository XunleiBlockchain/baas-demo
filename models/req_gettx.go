package models

import (
	"encoding/json"
	"errors"

	"github.com/XunleiBlockchain/tc-libs/common"
)

type ReqGetTx struct {
	Account string `json:"account"`
	Hash    string `json:"hash"`
}

func (rgt *ReqGetTx) Parse(rbody []byte) error {
	if err := json.Unmarshal(rbody, rgt); err != nil {
		return err
	}
	if !common.HasHexPrefix(rgt.Hash) {
		rgt.Hash = "0x" + rgt.Hash
	}
	if err := rgt.Sanity(); err != nil {
		return err
	}
	return nil
}

func (rgt *ReqGetTx) Sanity() error {
	if !common.IsHexAddress(rgt.Account) {
		return errors.New("account address illegal")
	}
	if len(rgt.Hash) == 2+2*common.HashLength && common.IsHex(rgt.Hash[2:]) {
		return nil
	}
	return errors.New("tx hash illegal")
}

func (rgt *ReqGetTx) GetTx() (interface{}, error) {
	backReqParam := []interface{}{
		rgt.Account,
		rgt.Hash,
	}
	var ret interface{}
	if useSDK {
		res, sdkErr := sdkInstance.GetTransactionReceipt(backReqParam)
		if sdkErr != nil && sdkErr.Code != 0 {
			return nil, sdkErr
		}
		ret = res
	} else {
		res, err := backCall("call", backReqParam)
		if err != nil {
			return nil, err
		}
		ret = res
	}
	return ret, nil
}
