package main

import (
	"github.com/eientei/iichan-thread-grabber/server"
	"net/http"
)

func main() {
	http.HandleFunc("/", server.Handler)
	_ = http.ListenAndServe(server.ListenAddr, nil)
}
