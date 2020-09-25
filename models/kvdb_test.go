package models

import (
	"fmt"
	"testing"

	"github.com/XunleiBlockchain/baas-demo/contract"
	_ "github.com/XunleiBlockchain/baas-demo/contract/ERC20"
	_ "github.com/XunleiBlockchain/baas-demo/contract/KVDB"
)

func Test_Result(t *testing.T) {
	name := "kvdb"
	method := "get"
	result := "0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000e7468697320697320612074657374000000000000000000000000000000000000"
	instance := contract.Get(name)
	if instance == nil {
		panic("instance nil")
	}
	ret, err := instance.Result(method, result)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ret: %s\n", ret)
}
