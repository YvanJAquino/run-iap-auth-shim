package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	parent := context.Background()
	logger := log.Default()
	router := &Router{
		logger: logger,
	}
	mux := http.DefaultServeMux
	mux.Handle("/", router)
	server := &http.Server{
		Addr:        addr,
		Handler:     mux,
		BaseContext: func(net.Listener) context.Context { return parent },
	}
	err := server.ListenAndServe()
	if err != nil {
		logger.Fatal(err)
	}
}

type Router struct {
	logger *log.Logger
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	targetURL := "https://hello-a67fdjzmma-uc.a.run.app"
	http.Redirect(w, req, targetURL, http.StatusMovedPermanently)
	for key, values := range req.Header {
		for idx := range values {
			if idx == 0 {
				w.Header().Set(key, values[idx])
			} else {
				w.Header().Add(key, values[idx])
			}
		}
	}
}
