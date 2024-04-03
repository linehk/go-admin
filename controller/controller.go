package controller

import (
	"net/http"
)

func Setup() *http.ServeMux {
	mux := http.NewServeMux()
	var userImpl UserImpl
	options := StdHTTPServerOptions{
		BaseRouter: mux,
	}
	HandlerWithOptions(&userImpl, options)
	return mux
}
