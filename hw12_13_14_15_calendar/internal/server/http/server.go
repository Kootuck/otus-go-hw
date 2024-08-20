package internalhttp

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"time"

	//nolint:depguard
	app "github.com/Kootuck/otus-go-hw/hw12_13_14_15_calendar/internal/app"
	//nolint:depguard
	conf "github.com/Kootuck/otus-go-hw/hw12_13_14_15_calendar/internal/config"
	//nolint:depguard
	logger "github.com/Kootuck/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
)

type Server struct {
	HTTPServer *http.Server
	Logger     *logger.Logger
	Config     conf.HTTPConf
}

type MyHandler struct {
	Logger *logger.Logger
}

func NewServer(logger *logger.Logger, config conf.HTTPConf, calendar *app.App) *Server {
	_ = calendar
	handler := &MyHandler{
		Logger: logger,
	}

	mux := http.NewServeMux()
	mux.Handle("/", handler)
	loggedHandler := LoggingMiddleware(logger, mux)

	server := &http.Server{
		Addr:         net.JoinHostPort(config.Host, config.Port),
		Handler:      loggedHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &Server{
		HTTPServer: server,
		Logger:     logger,
		Config:     config,
	}
}

func (s *Server) Start() error {
	if err := s.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.HTTPServer.Shutdown(ctx); err != nil {
		return err
	}
	s.Logger.Info("server stopped gracefully")
	return nil
}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode("hello")
}
