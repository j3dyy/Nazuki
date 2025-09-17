package main

import (
	"errors"
	"fmt"

	cfg "github.com/j3dyy/nazuki/internal/config"
	"github.com/j3dyy/nazuki/internal/service"
	"github.com/j3dyy/nazuki/internal/store"
)

var app cfg.Application

func main() {

	defaultCfg := cfg.LoadConfigFromEnv()

	app = cfg.NewApplication(service.Service{}, store.Store{})
	app.Logger.Err(errors.New("error occured"))

	fmt.Printf("Default Config: %+v\n", app)
	fmt.Printf("Default Config: %+v\n", defaultCfg)

}
