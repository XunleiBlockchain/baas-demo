package StructDataRecord

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/XunleiBlockchain/tc-libs/accounts/abi"
	"github.com/XunleiBlockchain/tc-libs/common"

	"github.com/XunleiBlockchain/baas-demo/contract"
)

func init() {
	contract.Register(newStructDataRecord())
}

type StructDataRecord struct {
	name string
	def  string
	abi  abi.ABI
}

func newStructDataRecord() *StructDataRecord {
	structDataRecord := &StructDataRecord{
		name: "structdatarecord",
		def:  `[{"inputs":[{"internalType":"address","name":"newOwner","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"id","type":"uint256"},{"indexed":true,"internalType":"string","name":"optKey","type":"string"},{"indexed":true,"internalType":"string","name":"optVal","type":"string"},{"indexed":false,"internalType":"string","name":"desc","type":"string"},{"indexed":false,"internalType":"string","name":"data","type":"string"}],"name":"NewRecord","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previousOwner","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"inputs":[{"internalType":"uint256","name":"_id","type":"uint256"},{"internalType":"string","name":"_optKey","type":"string"},{"internalType":"string","name":"_optVal","type":"string"},{"internalType":"string","name":"_desc","type":"string"},{"internalType":"string","name":"_data","type":"string"}],"name":"addRecord","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"_id","type":"uint256"}],"name":"getRecord","outputs":[{"internalType":"string","name":"_optKey","type":"string"},{"internalType":"string","name":"_optVal","type":"string"},{"internalType":"string","name":"_desc","type":"string"},{"internalType":"string","name":"_data","type":"string"},{"internalType":"uint256","name":"_createTime","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"}]`,
	}
	abi, err := abi.JSON(strings.NewReader(structDataRecord.def))
	if err != nil {
		panic(fmt.Errorf("new structdatarecord abi.JSON err: %v", err))
	}
	structDataRecord.abi = abi
	return structDataRecord
}

func (self *StructDataRecord) Name() string {
	return self.name
}

func (self *StructDataRecord) Def() string {
	return self.def
}

func (self *StructDataRecord) Data(method string, params []interface{}) (string, error) {
	switch method {
	case "addRecord":
		return self.addRecord(params)
	case "getRecord":
		return self.getRecord(params)
	}
	return "", errors.New("unsupported method")
}

type GetRecordRes struct {
	OptKey     string
	OptVal     string
	Desc       string
	Data       string
	CreateTime *big.Int
}

func (self *StructDataRecord) Result(method string, ret string) (interface{}, error) {
	var val interface{}
	switch method {
	case "getRecord":
		val = &GetRecordRes{}
	default:
		return nil, errors.New("unsupported method")
	}
	err := self.abi.Unpack(val, method, common.FromHex(ret))
	if err != nil {
		return nil, err
	}
	return val, nil
}

// -----------------------------------------------------------------

func (self *StructDataRecord) addRecord(params []interface{}) (string, error) {
	if len(params) != 5 {
		return "", fmt.Errorf("addRecord param number not match. want: 5 got: %d", len(params))
	}
	id, ok := params[0].(string)
	if !ok {
		return "", fmt.Errorf("addRecord param-0 type not match. want: string got: %s", reflect.TypeOf(params[0]))
	}
	idValue, err := strconv.ParseInt(id, 0, 64)
	if err != nil {
		return "", fmt.Errorf("addRecord param-0 can not parse to int, error: %v", err)
	}
	optKey, ok := params[1].(string)
	if !ok {
		return "", fmt.Errorf("addRecord param-1 type not match. want: string got: %s", reflect.TypeOf(params[1]))
	}
	optVal, ok := params[2].(string)
	if !ok {
		return "", fmt.Errorf("addRecord param-2 type not match. want: string got: %s", reflect.TypeOf(params[2]))
	}
	desc, ok := params[3].(string)
	if !ok {
		return "", fmt.Errorf("addRecord param-3 type not match. want: string got: %s", reflect.TypeOf(params[3]))
	}
	data, ok := params[4].(string)
	if !ok {
		return "", fmt.Errorf("addRecord param-4 type not match. want: string got: %s", reflect.TypeOf(params[4]))
	}
	return self.pack("addRecord", big.NewInt(idValue), optKey, optVal, desc, data)
}

func (self *StructDataRecord) getRecord(params []interface{}) (string, error) {
	if len(params) != 1 {
		return "", fmt.Errorf("getRecord param number not match. want: 1 got: %d", len(params))
	}
	id, ok := params[0].(string)
	if !ok {
		return "", fmt.Errorf("getRecord param-0 type not match. want: string got: %s", reflect.TypeOf(params[0]))
	}
	idValue, err := strconv.ParseInt(id, 0, 64)
	if err != nil {
		return "", fmt.Errorf("getRecord param-0 can not parse to int, error: %v", err)
	}
	return self.pack("getRecord", big.NewInt(idValue))
}

func (self *StructDataRecord) pack(name string, args ...interface{}) (string, error) {
	data, err := self.abi.Pack(name, args...)
	if err != nil {
		return "", fmt.Errorf("%s abi.Pack err: %v", name, err)
	}
	return fmt.Sprintf("0x%x", data), nil
}
