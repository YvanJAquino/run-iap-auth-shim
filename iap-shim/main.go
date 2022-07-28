package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"unsafe"
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
	proxyURL, err := url.Parse("https://hello-a67fdjzmma-uc.a.run.app")
	if err != nil {
		r.logger.Fatal(err)
	}
	values := *(*url.Values)(unsafe.Pointer(&req.Header))
	proxyURL.RawQuery = values.Encode()
	http.Redirect(w, req, proxyURL.String(), http.StatusMovedPermanently)
}
