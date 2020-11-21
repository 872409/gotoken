package gotoken

import (
	"testing"
	"time"
)

func TestGoToken_Exit(t *testing.T) {
	config := GoTokenConfig{
		RedisHost:    "172.16.123.100:6379",
		RedisPwd:     "xx",
		RedisDB:      1,
		TokenVersion: "111",
		Secret:       "vvv",
		ExpireHour:   1,
	}

	goToken, _ := NewRedisGoToken(config)

	tokenPayload := &TokenPayload{
		UID:           1000,
		ClientType:    "ios",
		ClientVersion: "1",
		ExpiresAt:     time.Now().Add(time.Hour).Unix(),
	}

	token, err := goToken.GenerateAuth(tokenPayload, "asdfasdfasd")
	t.Log("token", err, token)
	clientPayload := &ClientPayload{
		Token:         token,
		ClientType:    "ios",
		ClientVersion: "1",
	}
	payload, err2 := goToken.Parse(clientPayload)
	t.Log("payload", err2, payload)
}
