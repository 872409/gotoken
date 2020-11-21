package gotoken

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

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
	// c, err := redis.Dial("tcp", config.RedisHost, redis.DialPassword(config.RedisPwd), redis.DialDatabase(config.RedisDB))
	// if err != nil {
	// 	log.Errorln("initGoToken: error ", err, config)
	// 	return nil, err
	// }

	expireHour := 24 * 365
	if config.ExpireHour > 0 {
		expireHour = config.ExpireHour
	}

	goToken := &GoToken{
		TokenKeyName: "go-token",
		Secret:       config.Secret,
		ExpireHour:   expireHour,
		storage:      newRedisStorage(config),
	}

	return goToken, nil
}

//
// func newRedis(secret string, dbConn redis.Conn) *GoToken {
// 	goToken := &GoToken{
// 		Secret:     secret,
// 		ExpireHour: 24 * 365,
// 		storage:    newRedisStorage(dbConn),
// 	}
// 	return goToken
// }

func newRedisStorage(config GoTokenConfig) StorageHandler {

	redisHandler := redisHandler{
		redisPool: initRedisPool(config),
	}
	return &redisHandler
}

func initRedisPool(config GoTokenConfig) *redis.Pool {
	maxIdle := 3
	timeout := time.Duration(10) * time.Second
	maxActive := 500

	return &redis.Pool{
		MaxIdle:     maxIdle, // 空闲数
		IdleTimeout: timeout,
		MaxActive:   maxActive, // 最大数
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", config.RedisHost,
				redis.DialPassword(config.RedisPwd),
				redis.DialDatabase(config.RedisDB),
				redis.DialConnectTimeout(timeout),
				redis.DialReadTimeout(timeout),
				redis.DialWriteTimeout(timeout))
			if err != nil {
				return nil, err
			}
			return con, nil
		},
	}
}

type redisHandler struct {
	StorageHandler
	redisPool *redis.Pool
}

func (gt *redisHandler) Close() {
	_ = gt.redisPool.Close()
}

func (gt *redisHandler) GetAuthToken(clientPayload *ClientPayload) (payload *TokenPayload, ok bool) {
	if gt.redisPool == nil {
		return
	}

	key := goTokenKey(clientPayload.Token)
	redisConn := gt.redisPool.Get()
	defer redisConn.Close()
	reply, err := redis.String(redisConn.Do("GET", key))

	// log.Println("GetAuthToken", reply, err)
	if err != nil {
		ok = false
		return
	}

	payload = &TokenPayload{}
	err = json.Unmarshal([]byte(reply), payload)

	if err != nil {
		ok = false
		payload = nil
		return
	}

	ok = true

	return
}

func (gt *redisHandler) SaveAuthToken(payload *TokenPayload, token string) (err error) {

	if gt.redisPool == nil {
		return errors.New("redisConn not connect")
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatalln(err)
		return
	}

	lastToken, _ := gt.getClientToken(payload.UID, string(payload.ClientType))

	redisConn := gt.redisPool.Get()
	// defer redisConn.Close()

	if err = redisConn.Send("SET", goTokenKey(token), string(bytes)); err != nil {
		return
	}

	if err = redisConn.Send("HSET", userTokenHashKey(payload.UID), payload.ClientType, token); err != nil {
		return
	}

	if err = redisConn.Flush(); err != nil {
		return
	}

	if token != lastToken {
		_, _ = redisConn.Do("DEL", goTokenKey(lastToken))
	}

	// log.Println("SaveAuthToken", token, err)
	// if err != nil {
	// 	return false
	// }

	// go gt.removeClientToken(payload.UID, payload.ClientType)

	return nil
}

func (gt *redisHandler) getClientToken(uid int64, clientType string) (string, bool) {
	redisConn := gt.redisPool.Get()
	defer redisConn.Close()
	reply, err := redis.String(redisConn.Do("HGET", userTokenHashKey(uid), clientType))
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
	redisConn := gt.redisPool.Get()
	defer redisConn.Close()
	_, err := redisConn.Do("DEL", goTokenKey(token))
	if err != nil {
		return false
	}

	return true

}
