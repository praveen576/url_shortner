package db

import (
  "errors"
  "github.com/gomodule/redigo/redis"
  "github.com/url_shortner/config"
  "strings"
  "time"
  "log"
)

func newPool() redis.Pool{
  log.Println("in newPool check config", config.Config.Redis.Pool)
  return redis.Pool{
    IdleTimeout: 60 * time.Second,
    MaxActive: config.Config.Redis.Pool,
    Wait:true,
    Dial: newConn,
    TestOnBorrow: testConn,
  }
}

func newConn() (redis.Conn, error) {
  return redis.Dial("tcp", config.Config.Redis.Addr, redis.DialDatabase(config.Config.Redis.Db))
}

func testConn(c redis.Conn, _ time.Time) error {
  reply, err := redis.String(c.Do("PING"))
  if err != nil {
    return err
  } else if strings.ToLower(reply) != "pong" {
    return errors.New("unknown redis ping reply: " + reply)
  } else {
    return nil
  }
}
