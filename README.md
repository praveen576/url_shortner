# url_shortner

Pre Requisites

1. Golang installed on your system. This service built on go version go1.11 linux/amd64
2. Redis is required.
3. This is using golang-glide for package management. so set up glide. reference https://github.com/Masterminds/glide 


Installation
1. Once repo cloned. cd to repo. Make sure repo path is $GOPATH/src/github.com/url_shortner
2. From repo foler run `glide install`
3. Then `go run server_connection.go`. This will by default run in development mode. Check config/developemt.yaml file.
