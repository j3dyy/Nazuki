package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/j3dyy/nazuki/internal/app"
	cfg "github.com/j3dyy/nazuki/internal/config"
	"github.com/rs/zerolog/log"
)

func main() {
	defaultCfg := cfg.LoadConfigFromEnv()
	app, err := app.NewApplication(defaultCfg)
	if err != nil {
		log.Info().Msg(err.Error())
		os.Exit(1)
	}

	go listenForError(app)
	go listenForShutdown(app)

}

func listenForError(app *app.Application) {
	for {
		select {
		case err := <-app.ErrorChan:
			app.Logger.Error().Err(err).Msg(err.Error())
		case <-app.ErrorChanDone:
			return
		}
	}
}

func listenForShutdown(app *app.Application) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	shutdown(app)
	os.Exit(0)
}

func shutdown(app *app.Application) {
	app.Logger.Info().Msg("Shutting down server...")

	app.Wait.Wait()
	app.ErrorChanDone <- true

	close(app.ErrorChan)
	close(app.ErrorChanDone)
}
