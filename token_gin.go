package gotoken

import (
	"sync"

	"github.com/872409/gatom/gc"
	"github.com/872409/gatom/log"
	"github.com/gin-gonic/gin"
)

const (
	HeaderToken    = "go-token"
	queryTokenName = "go-token"
	ginPayloadKey  = "go-token"
)

func GetAuthPayload(c *gin.Context) (payload *TokenPayload, exists bool) {
	t, exists := c.Get(ginPayloadKey)
	if !exists {
		return nil, false
	}

	payload, ok := t.(*TokenPayload)

	if !ok {
		return nil, false
	}

	return
}

func decodeClientPayload(encoded string) (payload *ClientPayload, err error) {

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

func GetClientPayload(c *gin.Context) (payload *ClientPayload, err error) {
	encoded := c.GetHeader(HeaderToken)
	if len(encoded) > 0 {
		payload, err = decodeClientPayload(encoded)
	} else if encoded, ok := c.GetQuery(queryTokenName); ok {
		payload, err = decodeClientPayload(encoded)
	} else {
		err = ErrorClientPayloadNotExists
	}

	if err != nil {
		return nil, err
	}

	return wrapClientPayload(c, payload), nil
}

func GetClientPayloadDefault(c *gin.Context) (payload *ClientPayload) {
	payload, err := GetClientPayload(c)
	if err != nil {
		payload = wrapClientPayload(c, &ClientPayload{})
	}

	return payload
}

func wrapClientPayload(c *gin.Context, payload *ClientPayload) *ClientPayload {
	payload.IP = c.ClientIP()
	payload.UserAgent = c.GetHeader("User-Agent")
	return payload
}

//
// func GetClientQueryPayload(c *gin.Context, queryName ...string) (payload *TokenPayload, err error) {
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
// func GetClientHeaderPayload(c *gin.Context) (payload *TokenPayload, err error) {
// 	encoded := c.GetHeader(HeaderToken)
// 	return decodePayload(encoded)
// }

var ginMiddleWareOnce = sync.Once{}

func (gt *GoToken) GetClientPayload(g *gc.GContext) (*ClientPayload, error) {
	return GetClientPayload(g.Context)
}

func (gt *GoToken) SetMiddlewarePayloadProvider(provider PayloadProvider) {
	gt.middlewarePayloadProvider = provider
}

func (gt *GoToken) GinMiddleware() gin.HandlerFunc {

	ginMiddleWareOnce.Do(func() {
		gt.ginMiddleware = func(c *gin.Context) {
			g := gc.New(c)

			var clientPayload *ClientPayload
			var err error

			if gt.middlewarePayloadProvider != nil {
				clientPayload, err = gt.middlewarePayloadProvider.GetClientPayload(g)
			} else {
				clientPayload, err = gt.GetClientPayload(g)
			}

			if err == nil {
				var serverPayload *TokenPayload
				if gt.middlewarePayloadProvider != nil {
					serverPayload, err = gt.middlewarePayloadProvider.Parse(clientPayload)
				} else {
					serverPayload, err = gt.Parse(clientPayload)
				}

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
	})

	return gt.ginMiddleware
}
