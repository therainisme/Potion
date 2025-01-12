package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/therainisme/potion/proxy"
	"github.com/therainisme/potion/util"
)

func main() {
	util.LogInfo("Server starting at http://0.0.0.0:%s", util.GetPort())
	http.HandleFunc("/", proxy.HandleRequest)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", util.GetPort()), nil); err != nil {
		util.LogError("Server failed to start: %v", err)
		log.Fatal(err)
	}
}
