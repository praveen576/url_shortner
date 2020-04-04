package lib

import (
	"fmt"
	"crypto/md5"
  "github.com/url_shortner/i13n"
  "github.com/url_shortner/db"
  "sync"
)

func ShortenUrl(url_path string) string {
	val, err :=  db.GetKey(url_path)
	if err != nil {
		ctx := i13n.LogFields{
		  Fields: map[string]interface{}{ "method" : "ShortenUrl", "url_path" : url_path  },
		}
		i13n.Error(err, ctx)
	} else {
		return val
	}


	h := md5.New()
	h.Write([]byte(url_path))
	bs := h.Sum(nil)
	new_url_path := fmt.Sprintf("%x", bs)

	var wg sync.WaitGroup
	wg.Add(2)
	go db.SetKeyAsync(url_path, new_url_path, &wg)
	go db.SetKeyAsync(new_url_path, url_path, &wg)
  wg.Wait()
	return  new_url_path
}

func ExtendUrl(url_path string) (string, error) {
	val, err :=  db.GetKey(url_path)
	if err != nil {
		ctx := i13n.LogFields{
		  Fields: map[string]interface{}{ "method" : "ExtendUrl", "url_path" : url_path  },
		}
		i13n.Error(err, ctx)
	}
	return val, err
}
