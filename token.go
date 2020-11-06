package gotoken

import (
	"time"

	"github.com/872409/gatom/gc"
	"github.com/872409/gatom/log"
	"github.com/gin-gonic/gin"
)

type ClientType string

const (
	ClientType_iOS     = "ios"
	ClientType_Android = "android"
)

type PayloadProvider interface {
	GetClientPayload(g *gc.GContext) (*ClientPayload, error)
	Parse(client *ClientPayload) (*TokenPayload, error)
}

type TokenPayload struct {
	Token         string     `json:"-" schema:"token"`
	UID           int64      `json:"u,omitempty" schema:"u"`
	ClientType    ClientType `json:"ct,omitempty" schema:"ct"`
	ClientVersion string     `json:"cv,omitempty" schema:"cv"`
	// TokenVersion    string `json:"av,omitempty" schema:"av"`
	ExpiresAt int64 `json:"ea,omitempty" schema:"ea"`
}
type ClientPayload struct {
	Token         string     `json:"-" schema:"token"`
	ClientType    ClientType `json:"ct,omitempty" schema:"ct"`
	ClientVersion string     `json:"cv,omitempty" schema:"cv"`
	IP            string
	UserAgent     string
}

type StorageHandler interface {
	SaveAuthToken(payload *TokenPayload, authToken string) (err error)
	GetAuthToken(clientPayload *ClientPayload) (payload *TokenPayload, ok bool)
	Close()
}

type GoToken struct {
	TokenKeyName              string
	TokenVersion              string
	Secret                    string
	ExpireHour                int
	storage                   StorageHandler
	ginMiddleware             gin.HandlerFunc
	middlewarePayloadProvider PayloadProvider
}

func (gt *GoToken) Exit() {
	gt.storage.Close()
}

func (gt *GoToken) GenerateAuth(payload *TokenPayload, payloadSecret string) (string, error) {
	if payload.ExpiresAt == 0 {
		payload.ExpiresAt = time.Now().Add(time.Hour * time.Duration(gt.ExpireHour)).Unix()
	}

	token := generateToken(payload, gt.Secret, gt.TokenVersion, payloadSecret)

	if err := gt.storage.SaveAuthToken(payload, token); err != nil {
		log.Infoln(err)
		return "", ErrorGoTokenGen
	}

	return token, nil
}

func (gt *GoToken) ParseBase64(encoded string) (*TokenPayload, error) {
	payload, _ := ParseFormBase64(encoded)
	return gt.Parse(payload)
}

func (gt *GoToken) Parse(clientPayload *ClientPayload) (payload *TokenPayload, err error) {
	storedPayload, ok := gt.storage.GetAuthToken(clientPayload)
	// log.Infoln("Parse", storedPayload.ClientType, ok, clientPayload.ClientType, storedPayload.ExpiresAt < time.Now().Unix())
	if !ok || storedPayload == nil {
		err = ErrorGoTokenInvalid
		return
	}
	if storedPayload.ExpiresAt < time.Now().Unix() {
		err = ErrorGoTokenExpired
		return
	}

	if storedPayload.ClientType != clientPayload.ClientType {
		err = ErrorGoTokenInvalid
		return
	}

	payload = storedPayload

	return
}
