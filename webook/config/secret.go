package config

import "os"

var (
	SecretKey = os.Getenv("SMSSecretID")
	SecretId  = os.Getenv("SMSSecretKey")
	//设备id
	SdkAPPId = "1400930134"
	SigName  = "影剧瞬间公众号"
)
