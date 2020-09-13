package gotoken

import (
	"time"
	"github.com/gin-gonic/gin"
	"github.com/872409/gatom/log"
)

type Payload struct {
	Token         string `json:"-" schema:"token"`
	UID           int64  `json:"u,omitempty" schema:"u"`
	ClientType    string `json:"ct,omitempty" schema:"ct"`
	ClientVersion string `json:"cv,omitempty" schema:"cv"`
	// TokenVersion    string `json:"av,omitempty" schema:"av"`
	ExpiresAt int64 `json:"ea,omitempty" schema:"ea"`
}

type StorageHandler interface {
	SaveAuthToken(payload *Payload, authToken string) (err error)
	GetAuthToken(clientPayload *Payload) (payload *Payload, ok bool)
	Close()
}

type GoToken struct {
	TokenKeyName string
	TokenVersion string
	Secret       string
	ExpireHour   int
	storage      StorageHandler
	ginMiddleware gin.HandlerFunc
}

func (gt *GoToken) Exit() {
	gt.storage.Close()
}

func (gt *GoToken) GenerateAuth(payload *Payload, payloadSecret string) (string, error) {
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

func (gt *GoToken) ParseBase64(encoded string) (*Payload, error) {
	payload, _ := ParseFormBase64(encoded)
	return gt.Parse(payload)
}

func (gt *GoToken) Parse(clientPayload *Payload) (payload *Payload, err error) {
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
