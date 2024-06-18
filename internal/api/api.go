package api

import (
	"log"
)

type Migrator interface {
	Migrate() error
}

type Server interface {
	Serve() error
	Addr() string
}

type Api struct {
	loggingEnabled bool
	migrator       Migrator
	servers        []Server
}

func New(opts ...ApiOption) *Api {
	api := Api{}
	for _, opt := range opts {
		opt.apply(&api)
	}
	return &api
}

func (api *Api) Start() error {
	if api.migrator != nil {
		if err := api.migrator.Migrate(); err != nil {
			return err
		}
	}

	errs := make(chan error)
	for _, server := range api.servers {
		if api.loggingEnabled {
			log.Printf("start serving at %s", server.Addr())
		}

		go func(server Server) {
			if err := server.Serve(); err != nil {
				errs <- err
			}
		}(server)
	}

	return <-errs
}
