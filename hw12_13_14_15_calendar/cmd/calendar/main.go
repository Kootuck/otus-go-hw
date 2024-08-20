package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	//nolint:depguard
	"github.com/Kootuck/otus-go-hw/hw12_13_14_15_calendar/internal/app"
	//nolint:depguard
	config "github.com/Kootuck/otus-go-hw/hw12_13_14_15_calendar/internal/config"
	//nolint:depguard
	"github.com/Kootuck/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	//nolint:depguard
	internalhttp "github.com/Kootuck/otus-go-hw/hw12_13_14_15_calendar/internal/server/http"
	//nolint:depguard
	memorystorage "github.com/Kootuck/otus-go-hw/hw12_13_14_15_calendar/internal/storage/memory"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/default.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := config.NewConfig(configFile)
	if err != nil {
		log.Fatalf(err.Error())
	}

	logg := logger.New(config.Logger.Level, os.Stdout)

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, config.HTTP, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.HTTP.ShutdownTimeout))
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error(fmt.Errorf("failed to stop http server: %w", err))
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(); err != nil {
		logg.Error(fmt.Errorf("failed to start http server: %w", err))
		cancel()
	}
}
