package db

import (
  "github.com/gomodule/redigo/redis"
  "log"
  "fmt"
  "github.com/url_shortner/i13n"
  "sync"
)

var busPool = newPool()

func SetKey(key string, val interface{}) error {
  redisConn := busPool.Get()
  defer redisConn.Close()

  var err error

  val_str := fmt.Sprintf("%v", val)
  _, err = redis.String(redisConn.Do("set", key, val_str))

  if err != nil {
    ctx := i13n.LogFields{
      Fields: map[string]interface{}{ "method" : "SetKey", "key" : key, "val": val_str },
    }
    i13n.Error(err, ctx)
  }

  return err
}

func SetKeyAsync(key string, val interface{}, wg *sync.WaitGroup) {
  defer wg.Done()

  redisConn := busPool.Get()
  defer redisConn.Close()

  var err error

  val_str := fmt.Sprintf("%v", val)
  _, err = redis.String(redisConn.Do("set", key, val_str))

  if err != nil {
    ctx := i13n.LogFields{
      Fields: map[string]interface{}{ "method" : "SetKey", "key" : key, "val": val_str },
    }
    i13n.Error(err, ctx)
  }
}

func GetKey(key string) (string, error) {
  redisConn := busPool.Get()
  defer redisConn.Close()

  var err error
  ans, err := redis.String(redisConn.Do("get", key))

  if err != nil {
    ctx := i13n.LogFields{
      Fields: map[string]interface{}{ "method" : "GetKey", "key" : key },
    }
    i13n.Error(err, ctx)
  }

  return ans, err
}

func ExpireKeyAfter(key string, seconds int) (int, error) {
  redisConn := busPool.Get()
  defer redisConn.Close()

  log.Println("ExpireKeyAfter called..")
  var err error
  ans, err := redis.Int(redisConn.Do("expire", key, seconds))
  
  if err != nil {
    ctx := i13n.LogFields{
      Fields: map[string]interface{}{ "method" : "ExpireKeyAfter", "key" : key, "seconds": seconds },
    }
    i13n.Error(err, ctx)
  }
  log.Println("ExpireKeyAfter result..", int(ans), err)
  return int(ans), err
}
