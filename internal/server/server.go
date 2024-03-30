package server

import (
	"context"
	"net/http"
)

type WebServer struct {
	http.Server
}

func NewWebServer(host string, handler http.Handler) *WebServer {
	return &WebServer{
		Server: http.Server{
			Addr:    host,
			Handler: handler,
		},
	}
}

func (srv *WebServer) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()

		if err := srv.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()
	return srv.Server.ListenAndServe()
}
