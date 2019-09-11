package gotoken

import (
	"fmt"
	"net/url"

	"github.com/872409/gatom/crypto"
	"github.com/gorilla/schema"
)

// func (payload *Payload) ToBase64() string {
// 	query := payload.ToQuery()
// 	encoded := base64.StdEncoding.EncodeToString([]byte(query))
// 	return encoded
// }

// func (payload *Payload) ToQuery() string {
// 	encoder := schema.NewEncoder()
// 	dst := url.Values{}
//
// 	if err := encoder.Encode(payload, dst); err != nil {
// 		return ""
// 	}
//
// 	return dst.Encode()
// }

func ParseFormBase64(encoded string) (payload *Payload, err error) {
	decoded, err := crypto.Base64Decode(encoded)
	if err != nil {
		return
	}

	payload, _ = ParseFormQuery(string(decoded))
	return
}

func ParseFormQuery(query string) (payload *Payload, err error) {
	queryMap, err := url.ParseQuery(query)
	if err != nil {
		return
	}

	payload = &Payload{}
	decoder := schema.NewDecoder()
	if err = decoder.Decode(payload, queryMap); err != nil {
		return
	}

	return payload, nil
}

func generateToken(payload *Payload, secret string, tokenVersion string, payloadSecret string) string {
	c := fmt.Sprintf("auth_%d%s%s%s%d%s%s", payload.UID, payload.ClientType, payload.ClientVersion, tokenVersion, payload.ExpiresAt, secret, payloadSecret)
	token := crypto.SHA256(c)
	return token
}

//
// func generateIMToken(payload Payload, secret string, payloadSecret string) string {
// 	c := fmt.Sprintf("im_%d%s%s", payload.UID, secret, payloadSecret)
// 	token := crypto.StrMD5(c)
// 	return token
// }
