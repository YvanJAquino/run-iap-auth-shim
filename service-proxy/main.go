package main

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"
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
	proxy := &Proxy{
		logger: logger,
	}
	mux := http.DefaultServeMux
	mux.Handle("/", proxy)
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

type Proxy struct {
	logger *log.Logger
	client *http.Client
}

func NewProxy(l *log.Logger) *Proxy {
	return new(Proxy).Init(l)
}

func (p *Proxy) Init(l *log.Logger) *Proxy {
	p.client = &http.Client{}
	p.logger = l
	return p
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	backendService := "https://hello-a67fdjzmma-uc.a.run.app"
	req, err := http.NewRequestWithContext(req.Context(), http.MethodGet, backendService, nil)
	if err != nil {
		p.logger.Fatal(err)
	}
	values := req.URL.Query()
	req.Header = *(*http.Header)(unsafe.Pointer(&values))
	resp, err := p.client.Do(req)
	if err != nil {
		p.logger.Fatal(err)
	}
	defer resp.Body.Close()
	for header, values := range resp.Header {
		for index := range values {
			if len(values) == 1 {
				w.Header().Set(header, values[index])
			} else {
				w.Header().Add(header, values[index])
			}
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
