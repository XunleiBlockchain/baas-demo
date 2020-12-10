package models

import (
	sdk "github.com/XunleiBlockchain/baas-sdk-go"
)

var (
	sdkInstance sdk.SDK
	useSDK      bool
)

func InitSDK(s sdk.SDK) {
	sdkInstance = s
	useSDK = true
}
