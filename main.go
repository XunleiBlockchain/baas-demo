package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/XunleiBlockchain/baas-demo/models"
	_ "github.com/XunleiBlockchain/baas-demo/routers"
	sdk "github.com/XunleiBlockchain/baas-sdk-go"

	// import contracts to do `init`
	_ "github.com/XunleiBlockchain/baas-demo/contract/ERC20"
	_ "github.com/XunleiBlockchain/baas-demo/contract/KVDB"
	_ "github.com/XunleiBlockchain/baas-demo/contract/StructDataRecord"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func main() {
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterFile, `{"filename":"baas-demo.log"}`)

	err := initBaasSDK()
	if err != nil {
		panic("[initBaasSDK] error: " + err.Error())
	}
	sdkInstance := sdk.GetSDK(log)
	models.InitSDK(sdkInstance)

	beego.Run()
}

func initBaasSDK() error {
	dnscacheUpdateintervalConf, _ := beego.AppConfig.Int("dnscache_updateinterval")
	chainidConf, _ := beego.AppConfig.Int64("chainid")
	sdkConf := &sdk.Config{
		Keystore:               beego.AppConfig.String("keystore"),
		UnlockAccounts:         make(map[string]string),
		DNSCacheUpdateInterval: dnscacheUpdateintervalConf,
		RPCProtocal:            beego.AppConfig.String("rpc_protocal"),
		XHost:                  beego.AppConfig.String("xhost"),
		Namespace:              beego.AppConfig.String("namespace"),
		ChainID:                chainidConf,
	}

	authInfoJSON, err := ioutil.ReadFile("./conf/auth.json")
	if err != nil {
		return fmt.Errorf("can not read auth.json file, error: %v", err)
	}
	err = json.Unmarshal(authInfoJSON, &sdkConf.AuthInfo)
	if err != nil {
		return fmt.Errorf("can not unmarshal auth json, error: %v", err)
	}
	passwdsJSON, err := ioutil.ReadFile("./conf/passwd.json")
	if err != nil {
		return fmt.Errorf("can not read passwd.json file, error: %v", err)
	}
	err = json.Unmarshal(passwdsJSON, &sdkConf.UnlockAccounts)
	if err != nil {
		return fmt.Errorf("can not unmarshal passwd json, error: %v", err)
	}

	sdk.Conf = sdkConf
	return nil
}
