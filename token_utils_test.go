package gotoken

import (
	"fmt"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	// "soyo/keeptalk/interface/openplatform/account/app/conf"
)

func getRedisConn() *GoToken {

	c, err := redis.Dial("tcp", "10.211.55.4:6379", redis.DialPassword("xx"), redis.DialDatabase(1))
	if err != nil {
		fmt.Println("Connect to storage error", err)
		return nil
	}
	goToken := NewRedis("aaa", c)
	goToken.TokenVersion = "v1"

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

	payload := &Payload{UID: 100, ClientVersion: "v1.0.1", ClientType: "ios", ExpiresAt: time.Now().Add(time.Hour * 24).Unix()}

	token, ok := goToken.GenerateAuth(payload, "aaa")
	// base64 := payload.ToBase64()
	// fmt.Println(token, ok)

	clientPayload := &Payload{ClientType: "ios", Token: token}
	getPayload, ok := goToken.Parse(clientPayload)
	fmt.Println(getPayload, ok)
}
