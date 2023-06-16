package main

var (
	secretId  = "AKID1111111111"
	secretKey = "12312312312312"
)

func GetAuthInfo(secretId string) (uint64, string) {
	//select appId,secretKey from app where secretId=? limit 1;
	return 1258888, secretKey
}
