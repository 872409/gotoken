package gotoken

import (
	"errors"

	"github.com/872409/gatom/gc"
	"github.com/872409/gatom/log"
	"github.com/gin-gonic/gin"
)

const (
	HeaderToken    = "go-token"
	queryTokenName = "go-token"
	ginPayloadKey  = "go-token"
)

func GetAuthPayload(c *gin.Context) (payload *Payload, exists bool) {
	t, exists := c.Get(ginPayloadKey)
	if !exists {
		return nil, false
	}

	payload, ok := t.(*Payload)

	if !ok {
		return nil, false
	}

	return
}

func decodePayload(encoded string) (payload *Payload, err error) {

	if encoded == "" {
		err = ErrorGoTokenHeaderParams
		return
	}

	payload, er := ParseFormBase64(encoded)

	if er != nil {
		err = ErrorGoTokenHeaderEncoded
		return
	}

	return
}

func GetClientPayload(c *gin.Context) (payload *Payload, err error) {
	encoded := c.GetHeader(HeaderToken)
	if len(encoded) > 0 {
		return decodePayload(encoded)
	}

	if encoded, ok := c.GetQuery(queryTokenName); ok {
		return decodePayload(encoded)
	}

	return nil, errors.New("client payload not exists")
}

//
// func GetClientQueryPayload(c *gin.Context, queryName ...string) (payload *Payload, err error) {
// 	_queryName := queryTokenName
//
// 	if len(queryName) > 0 {
// 		_queryName = queryName[0]
// 	}
//
// 	if encoded, ok := c.GetQuery(_queryName); ok {
// 		return decodePayload(encoded)
// 	}
//
// 	return nil, errors.New("query not exists")
//
// }

//
// func GetClientHeaderPayload(c *gin.Context) (payload *Payload, err error) {
// 	encoded := c.GetHeader(HeaderToken)
// 	return decodePayload(encoded)
// }

func (gt *GoToken) Middleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		g := gc.New(c)
		clientPayload, err := GetClientPayload(c)

		if err == nil {
			var serverPayload *Payload
			serverPayload, err = gt.Parse(clientPayload)
			// log.Println("serverPayload", serverPayload, err)

			if err == nil && serverPayload.UID > 0 {
				g.SetAuthID(serverPayload.UID)
				g.Set(ginPayloadKey, serverPayload)
			}
		}

		if err != nil || g.AuthID() == 0 {
			log.Println("GoToken Middleware", clientPayload, err)
			g.JSONCodeError(err)
			return
		}

		// log.Println("clientPayload", clientPayload)

		g.Next()
	}
}
