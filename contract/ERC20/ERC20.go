package ERC20

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
	contract.Register(newERC20())
}

type ERC20 struct {
	name string
	def  string
	abi  abi.ABI
}

func newERC20() *ERC20 {
	erc20 := &ERC20{
		name: "erc20",
		def: `[{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"value","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
				{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},
				{"constant":false,"inputs":[{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"value","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
				{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"addedValue","type":"uint256"}],"name":"increaseAllowance","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
				{"constant":true,"inputs":[{"name":"owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},
				{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"subtractedValue","type":"uint256"}],"name":"decreaseAllowance","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
				{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},
				{"constant":true,"inputs":[{"name":"owner","type":"address"},{"name":"spender","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},
				{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"},
				{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"spender","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Approval","type":"event"}]`,
	}
	abi, err := abi.JSON(strings.NewReader(erc20.def))
	if err != nil {
		panic(fmt.Errorf("newERC20 abi.JSON err: %v", err))
	}
	erc20.abi = abi
	return erc20
}

func (self *ERC20) Name() string {
	return self.name
}

func (self *ERC20) Def() string {
	return self.def
}

func (self *ERC20) Data(method string, params []interface{}) (string, error) {
	switch method {
	case "transfer":
		return self.transfer(params)
	case "balanceOf":
		return self.balanceOf(params)
		//TODO other method
	}
	return "", errors.New("unsupported method")
}

func (self *ERC20) Result(method string, ret string) (interface{}, error) {
	var val interface{}
	switch method {
	case "balanceOf":
		val = big.NewInt(0)
	default:
		return nil, errors.New("unsupported method")
	}
	return self.unpack(&val, method, common.FromHex(ret))
}

// -----------------------------------------------------------------

func (self *ERC20) transfer(params []interface{}) (string, error) {
	if len(params) != 2 {
		return "", fmt.Errorf("transfer param number not match. want: 2 got: %d", len(params))
	}
	to, ok := params[0].(string)
	if !ok {
		return "", fmt.Errorf("transfer param-0 type not match. want: string got: %s", reflect.TypeOf(params[0]))
	}
	amount, ok := params[1].(string)
	if !ok {
		return "", fmt.Errorf("transfer param-1 type not match. want: string got: %s", reflect.TypeOf(params[1]))
	}
	val, err := strconv.ParseInt(amount, 0, 64)
	if err != nil {
		return "", fmt.Errorf("transfer param-1 value illegal. err: %v", err)
	}
	return self.pack("transfer", common.HexToAddress(to), big.NewInt(val))
}

func (self *ERC20) balanceOf(params []interface{}) (string, error) {
	if len(params) != 1 {
		return "", fmt.Errorf("BalanceOf param number not match. want: 1 got: %d", len(params))
	}
	addr, ok := params[0].(string)
	if !ok {
		return "", fmt.Errorf("transfer param-0 type not match. want: string got: %s", reflect.TypeOf(params[0]))
	}
	return self.pack("balanceOf", common.HexToAddress(addr))
}

func (self *ERC20) pack(name string, args ...interface{}) (string, error) {
	data, err := self.abi.Pack(name, args...)
	if err != nil {
		return "", fmt.Errorf("%s abi.Pack err: %v", name, err)
	}
	return fmt.Sprintf("0x%x", data), nil
}

func (self *ERC20) unpack(v interface{}, name string, output []byte) (interface{}, error) {
	err := self.abi.Unpack(&v, name, output)
	if err != nil {
		return nil, fmt.Errorf("%s abi.Unpack err: %v", name, err)
	}
	return v, nil
}
