package main

import (
  "context"
  "github.com/gorilla/handlers"
  "github.com/traveltriangle/url_shortner/lib"
  "github.com/traveltriangle/url_shortner/config"
  "github.com/traveltriangle/url_shortner/i13n"
  "github.com/traveltriangle/url_shortner/helpers/url_mapper_helpers"

  "net"
  "net/http"
  "net/url"
  "os"
  "os/signal"
  "syscall"
  "time"
  "fmt"
  "strings"
  "encoding/json"
)

const (
  ShutdownTimeout = 5 * time.Second
)

func dsl() map[string]func(http.ResponseWriter, *http.Request) {
  return map[string]func(http.ResponseWriter, *http.Request){
    "/map_url": url_mapper,
  }
}


func url_mapper(w http.ResponseWriter, r *http.Request) {
  fmt.Println("check req header", r.Header.Get("X-CLIENT-API-TOKEN") )
  api_token := "93f67dbdacf3906955b5c529bb692fb11ac13dc3015c87f25c22fae51a5d79290841127492e312f2"

  if api_token != r.Header.Get("X-CLIENT-API-TOKEN") {
    w.WriteHeader(401)
    message := map[string]string{ "message": "API Access key not matched" }
    json.NewEncoder(w).Encode(message)
    return
  }

  r.ParseForm()
  url_str := r.Form.Get("url")
  op_type := strings.ToLower(r.Form.Get("op_type"))

  ok := url_mapper_helpers.ValidateOpType(op_type)
  if !ok {
    w.WriteHeader(400)
    message := map[string]string{ "message": "bad param op_type: should be one of shorten or extend" }
    json.NewEncoder(w).Encode(message)
    return        
  }

  url_obj, err := url.Parse(url_str)

  if err != nil {
    ctx := i13n.LogFields{
      Fields: map[string]interface{}{ "method" : "url_mapper", "request" : r.URL  },
    }
    i13n.Error(err, ctx)
    w.WriteHeader(400)
    message := map[string]string{ "message": "bad param url" }
    json.NewEncoder(w).Encode(message)
    return
  }

  var new_url_path string
  if op_type == "shorten" {
    new_url_path = lib.ShortenUrl(strings.TrimLeft(url_obj.RequestURI(), "/"))
  } else {
    new_url_path, err = lib.ExtendUrl(strings.TrimLeft(url_obj.RequestURI(), "/"))
    if err != nil {
      w.WriteHeader(500)
      message := map[string]string{ "message": "Some Server Error Occured..!" }
      json.NewEncoder(w).Encode(message)
      return   
    }
  }

  var new_url string
  if url_obj.Scheme != "" {
    new_url += (url_obj.Scheme + "://")
  }
  new_url += (url_obj.Hostname() + "/" +new_url_path)

  response_data := map[string]interface{}{ "success": true, "url": new_url }

  json.NewEncoder(w).Encode(response_data)
}

func setupServer() *http.ServeMux {
  m := http.NewServeMux()
  fs := http.FileServer(http.Dir("public"))
  m.Handle("/", fs)

  for r, h := range dsl() {
    m.HandleFunc(r, h)
  }
  return m
}

func exit() {
  // if any exit pre-operation.
}

func InitSystem() {
  config.ReadConfig()
  i13n.Init()
  i13n.Info("Started", i13n.LogFields{
    Fields: map[string]interface{}{
      "hostname": config.Config.Hostname,
      "pid":      config.Config.Pid,
    },
  })
}

func main() {
  InitSystem()

  stop := make(chan os.Signal, 0)
  signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

  m := setupServer()
  s := http.Server{Handler: handlers.ProxyHeaders(m)}

  go func() {
    l, err := net.Listen(config.Config.Server.Network, config.Config.Server.Addr)

    i13n.Fatal(err)
    i13n.Info("Listening", i13n.LogFields{
      Fields: map[string]interface{}{
        "addr": l.Addr(),
      },
    })
    i13n.Error(s.Serve(l))
  }()

  sig := <-stop
  i13n.Info("Shutting down", i13n.LogFields{
    Fields: map[string]interface{}{
      "signal": sig.String(),
    },
  })

  ctx, _ := context.WithTimeout(context.Background(), ShutdownTimeout)
  s.Shutdown(ctx)
  exit()
  i13n.Info("Shut down complete")
}
