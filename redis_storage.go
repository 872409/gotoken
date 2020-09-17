package gotoken

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/872409/gatom/log"
	"github.com/gomodule/redigo/redis"
)

func goTokenKey(token string) string {
	return "gt:" + token
}

func userTokenHashKey(uid int64) string {
	return "ut:" + strconv.FormatInt(uid, 10)
}

func NewRedisGoToken(config GoTokenConfig) (*GoToken, error) {
	c, err := redis.Dial("tcp", config.RedisHost, redis.DialPassword(config.RedisPwd), redis.DialDatabase(config.RedisDB))
	if err != nil {
		log.Errorln("initGoToken: error ", err, config)
		return nil, err
	}

	expireHour := 24 * 365
	if config.ExpireHour > 0 {
		expireHour = config.ExpireHour
	}

	goToken := &GoToken{
		TokenKeyName: "go-token",
		Secret:       config.Secret,
		ExpireHour:   expireHour,
		storage:      newRedisStorage(c),
	}

	return goToken, nil
}

func newRedis(secret string, dbConn redis.Conn) *GoToken {
	goToken := &GoToken{
		Secret:     secret,
		ExpireHour: 24 * 365,
		storage:    newRedisStorage(dbConn),
	}
	return goToken
}

func newRedisStorage(dbConn redis.Conn) StorageHandler {
	redisHandler := redisHandler{redisConn: dbConn}
	return &redisHandler
}

type redisHandler struct {
	StorageHandler
	redisConn redis.Conn
}

func (gt *redisHandler) Close() {
	_ = gt.redisConn.Close()
}

func (gt *redisHandler) GetAuthToken(clientPayload *ClientPayload) (payload *Payload, ok bool) {
	if gt.redisConn == nil {
		return
	}

	key := goTokenKey(clientPayload.Token)
	reply, err := redis.String(gt.redisConn.Do("GET", key))

	// log.Println("GetAuthToken", reply, err)
	if err != nil {
		ok = false
		return
	}

	payload = &Payload{}
	err = json.Unmarshal([]byte(reply), payload)

	if err != nil {
		ok = false
		payload = nil
		return
	}

	ok = true

	return
}

func (gt *redisHandler) SaveAuthToken(payload *Payload, token string) (err error) {

	if gt.redisConn == nil {
		return errors.New("redisConn not connect")
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatalln(err)
		return
	}

	lastToken, _ := gt.getClientToken(payload.UID, payload.ClientType)

	if err = gt.redisConn.Send("SET", goTokenKey(token), string(bytes)); err != nil {
		return
	}

	if err = gt.redisConn.Send("HSET", userTokenHashKey(payload.UID), payload.ClientType, token); err != nil {
		return
	}

	if err = gt.redisConn.Flush(); err != nil {
		return
	}

	if token != lastToken {
		_, _ = gt.redisConn.Do("DEL", goTokenKey(lastToken))
	}

	// log.Println("SaveAuthToken", token, err)
	// if err != nil {
	// 	return false
	// }

	// go gt.removeClientToken(payload.UID, payload.ClientType)

	return nil
}

func (gt *redisHandler) getClientToken(uid int64, clientType string) (string, bool) {
	reply, err := redis.String(gt.redisConn.Do("HGET", userTokenHashKey(uid), clientType))
	// fmt.Println(reply, err)

	if err != nil {
		return "", false
	}

	return reply, true
}

func (gt *redisHandler) removeClientToken(uid int64, clientType string) bool {
	token, ok := gt.getClientToken(uid, clientType)
	if !ok {
		return false
	}

	_, err := gt.redisConn.Do("DEL", goTokenKey(token))
	if err != nil {
		return false
	}

	return true

}
