//go:generate go generate ./model
package main

import (
	"github.com/eientei/iichan-thread-grabber/server"
	"net/http"
)

func main() {
	http.HandleFunc(server.PublicPrefix+"/", server.Handler)
	_ = http.ListenAndServe(server.ListenAddr, nil)
}
