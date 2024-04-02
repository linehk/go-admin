package router

import "net/http"

func Setup() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /users", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	return mux
}
