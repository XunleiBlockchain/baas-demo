package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	sdk "github.com/XunleiBlockchain/baas-sdk-go"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/XunleiBlockchain/baas-demo/models"
	_ "github.com/XunleiBlockchain/baas-demo/routers"

	// import contracts to do `init`
	_ "github.com/XunleiBlockchain/baas-demo/contract/ERC20"
	_ "github.com/XunleiBlockchain/baas-demo/contract/KVDB"
	_ "github.com/XunleiBlockchain/baas-demo/contract/StructDataRecord"
)

var useSDK = flag.Bool("sdk", false, "Use SDK in demo or not")

func main() {
	flag.Parse()

	log := logs.NewLogger()
	log.SetLogger(logs.AdapterFile, `{"filename":"baas-demo.log"}`)

	if *useSDK {
		sdkConf, err := initBaasSDKConfig()
		if err != nil {
			panic("[initBaasSDKConfig] error: " + err.Error())
		}
		sdkInstance, err := sdk.NewSDK(sdkConf, log)
		if err != nil {
			panic("[sdk.NewSDK] error: " + err.Error())
		}
		models.InitSDK(sdkInstance)
	}

	beego.Run()
}

func initBaasSDKConfig() (*sdk.Config, error) {
	chainidConf, _ := beego.AppConfig.Int64("chainid")
	sdkConf := &sdk.Config{
		Keystore:       beego.AppConfig.String("keystore"),
		UnlockAccounts: make(map[string]string),
		RPCProtocal:    beego.AppConfig.String("rpc_protocal"),
		XHost:          beego.AppConfig.String("xhost"),
		Namespace:      beego.AppConfig.String("namespace"),
		ChainID:        chainidConf,
	}

	authInfoJSON, err := ioutil.ReadFile("./conf/auth.json")
	if err != nil {
		return nil, fmt.Errorf("can not read auth.json file, error: %v", err)
	}
	err = json.Unmarshal(authInfoJSON, &sdkConf.AuthInfo)
	if err != nil {
		return nil, fmt.Errorf("can not unmarshal auth json, error: %v", err)
	}
	passwdsJSON, err := ioutil.ReadFile("./conf/passwd.json")
	if err != nil {
		return nil, fmt.Errorf("can not read passwd.json file, error: %v", err)
	}
	err = json.Unmarshal(passwdsJSON, &sdkConf.UnlockAccounts)
	if err != nil {
		return nil, fmt.Errorf("can not unmarshal passwd json, error: %v", err)
	}

	return sdkConf, nil
}
