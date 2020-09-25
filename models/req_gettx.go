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
	if err := rgt.Sanity(); err != nil {
		return err
	}
	return nil
}

func (rgt *ReqGetTx) Sanity() error {
	if !common.IsHexAddress(rgt.Account) {
		return errors.New("account address illegal")
	}
	if len(rgt.Hash) == 2+2*common.HashLength && common.IsHex(rgt.Hash) {
		return nil
	}
	if len(rgt.Hash) == 2*common.HashLength && common.IsHex("0x"+rgt.Hash) {
		return nil
	}
	return errors.New("tx hash illegal")
}

func (rgt *ReqGetTx) GetTx() (interface{}, error) {
	backReqParam := []string{
		rgt.Account,
		rgt.Hash,
	}
	return backCall("getTransactionReceipt", backReqParam)
}
