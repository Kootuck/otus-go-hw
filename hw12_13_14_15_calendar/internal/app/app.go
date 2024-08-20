package app

import (
	"context"

	//nolint:depguard
	"github.com/Kootuck/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	//nolint:depguard
	"github.com/Kootuck/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	Storage storage.EventStorage
	Logger  *logger.Logger
}

func New(logger *logger.Logger, storage storage.EventStorage) *App {
	return &App{
		Storage: storage,
		Logger:  logger,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	_ = ctx
	return a.Storage.Add(storage.Event{ID: storage.EventID(id), Title: title})
}
