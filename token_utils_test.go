package gotoken

import (
	"fmt"
	"testing"
	"time"

	// "soyo/keeptalk/interface/openplatform/account/app/conf"
)
//
func getRedisConn() *GoToken {


	config := GoTokenConfig{
		RedisHost:    "10.211.55.4:6379",
		RedisPwd:     "xx",
		RedisDB:      1,
		TokenVersion: "111",
		Secret:       "vvv",
		ExpireHour:   1,
	}

	goToken := NewRedisGoToken(config)
	// goToken.TokenVersion = "v1"

	return goToken
}
func TestR(t *testing.T) {

	// query := "token=asdfasdfasdf&u=100&ct=android&cv=1.0.1&av=v1"
	// encoded := base64.StdEncoding.EncodeToString([]byte(query))
	// decoded, _ := base64.StdEncoding.DecodeString(encoded)
	//
	// fmt.Println(encoded)

	// payload := ParseFormQuery(string(decoded))
	payload, err := ParseFormBase64("YXY9djEmY3Q9YW5kcm9pZCZjdj0xLjAuMSZlYT0wJnRva2VuPWFzZGZhc2RmYXNkZiZ1PTEwMA==")

	fmt.Println(payload.Token, payload, err)
	// fmt.Println(PayloadToQuery(payload))
	// fmt.Println(PayloadToBase64(payload))

}
func TestGenerate(t *testing.T) {
	goToken := getRedisConn()

	payload := &TokenPayload{UID: 100, ClientVersion: "v1.0.1", ClientType: "ios", ExpiresAt: time.Now().Add(time.Hour * 24).Unix()}

	token, ok := goToken.GenerateAuth(payload, "aaa")
	// base64 := payload.ToBase64()
	// fmt.Println(token, ok)

	clientPayload := &ClientPayload{ClientType: "ios", Token: token}
	getPayload, ok := goToken.Parse(clientPayload)
	fmt.Println(getPayload, ok)
}
