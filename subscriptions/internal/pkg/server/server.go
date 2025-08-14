package server

import (
	"go.uber.org/fx"
	"net/http"
)

type Params struct {
	fx.In

	Config Config
	Router *Router
}

func RunServer(params Params) {
	srv := &http.Server{
		Addr:              params.Config.Address,
		Handler:           params.Router.handler,
		ReadHeaderTimeout: params.Config.ReadHeaderTimeout,
		IdleTimeout:       params.Config.IdleTimeout,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
}
