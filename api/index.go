package handler

import (
	"net/http"

	"github.com/therainisme/potion/proxy"
)

func VercelHandler(w http.ResponseWriter, r *http.Request) {
	proxy.HandleRequest(w, r)
}
