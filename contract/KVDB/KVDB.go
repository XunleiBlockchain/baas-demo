package KVDB

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/XunleiBlockchain/tc-libs/accounts/abi"
	"github.com/XunleiBlockchain/tc-libs/common"

	"github.com/XunleiBlockchain/baas-demo/contract"
)

func init() {
	contract.Register(newKVDB())
}

type KVDB struct {
	name string
	def  string
	abi  abi.ABI
}

func newKVDB() *KVDB {
	kvdb := &KVDB{
		name: "kvdb",
		def: `[{"constant":true,"inputs":[{"internalType":"bytes32","name":"key","type":"bytes32"}],"name":"get","outputs":[{"internalType":"bytes","name":"","type":"bytes"}],"payable":false,"stateMutability":"view","type":"function"},
			   {"constant":false,"inputs":[{"internalType":"bytes32","name":"key","type":"bytes32"},{"internalType":"bytes","name":"value","type":"bytes"}],"name":"set","outputs":[{"internalType":"bool","name":"","type":"bool"}],"payable":true,"stateMutability":"payable","type":"function"}]`,
	}
	abi, err := abi.JSON(strings.NewReader(kvdb.def))
	if err != nil {
		panic(fmt.Errorf("newKVDB abi.JSON err: %v", err))
	}
	kvdb.abi = abi
	return kvdb
}

func (self *KVDB) Name() string {
	return self.name
}

func (self *KVDB) Def() string {
	return self.def
}

func (self *KVDB) Data(method string, params []interface{}) (string, error) {
	switch method {
	case "set":
		return self.set(params)
	case "get":
		return self.get(params)
		//TODO other method
	}
	return "", errors.New("unsupported method")
}

func (self *KVDB) Result(method string, ret string) (interface{}, error) {
	var val interface{}
	switch method {
	case "get":
		val = []byte{}
	default:
		return nil, errors.New("unsupported method")
	}
	return self.unpack(&val, method, common.FromHex(ret))
}

// -----------------------------------------------------------------

func (self *KVDB) set(params []interface{}) (string, error) {
	if len(params) != 2 {
		return "", fmt.Errorf("set param number not match. want: 2 got: %d", len(params))
	}
	key, ok := params[0].(string)
	if !ok {
		return "", fmt.Errorf("set param-0 type not match. want: string got: %s", reflect.TypeOf(params[0]))
	}
	val, ok := params[1].(string)
	if !ok {
		return "", fmt.Errorf("set param-1 type not match. want: string got: %s", reflect.TypeOf(params[1]))
	}
	return self.pack("set", common.FromHex(key), []byte(val))
}

func (self *KVDB) get(params []interface{}) (string, error) {
	if len(params) != 1 {
		return "", fmt.Errorf("get param number not match. want: 1 got: %d", len(params))
	}
	key, ok := params[0].(string)
	if !ok {
		return "", fmt.Errorf("get param-0 type not match. want: string got: %s", reflect.TypeOf(params[0]))
	}
	return self.pack("get", common.FromHex(key))
}

func (self *KVDB) pack(name string, args ...interface{}) (string, error) {
	data, err := self.abi.Pack(name, args...)
	if err != nil {
		return "", fmt.Errorf("%s abi.Pack err: %v", name, err)
	}
	return fmt.Sprintf("0x%x", data), nil
}

func (self *KVDB) unpack(v interface{}, name string, output []byte) (interface{}, error) {
	err := self.abi.Unpack(&v, name, output)
	if err != nil {
		return nil, fmt.Errorf("%s abi.Unpack err: %v", name, err)
	}
	return string(v.([]byte)), nil
}
