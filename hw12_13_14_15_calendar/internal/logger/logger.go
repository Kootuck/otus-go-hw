package logger

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
)

type Logger struct {
	Level  LogLevel
	Output io.Writer
}

func ParseLogLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warning":
		return WARNING
	case "error":
		return ERROR
	default:
		return INFO
	}
}

func New(level string, out io.Writer) *Logger {
	return &Logger{
		Level:  ParseLogLevel(level),
		Output: out,
	}
}

func (l *Logger) Log(level LogLevel, msg string) {
	if level < l.Level {
		return
	}
	logMsg := fmt.Sprintf("[%d] %s\n", level, msg)
	fmt.Fprint(l.Output, logMsg)
}

func (l *Logger) Debug(msg string) {
	l.Log(DEBUG, msg)
}

func (l *Logger) Info(msg string) {
	l.Log(INFO, msg)
}

func (l *Logger) Warning(msg string) {
	l.Log(WARNING, msg)
}

func (l *Logger) Error(err error) {
	l.ErrorS(err.Error())
}

func (l *Logger) ErrorS(msg string) {
	l.Log(ERROR, msg)
}

func (l Logger) LogHTTP(r *http.Request, startTime time.Time, statusCode int) {
	clientIP := r.RemoteAddr
	userAgent := r.Header.Get("User-Agent")
	requestTime := startTime.Format("01/Jan/2001:12:13:14 +0300")
	method := r.Method
	path := r.URL.Path
	httpVersion := r.Proto
	latency := time.Since(startTime)

	msg := fmt.Sprintf("%s [%s] %s %s %s %d %d \"%s\"",
		clientIP, requestTime, method, path, httpVersion, statusCode, latency.Milliseconds(), userAgent)

	if statusCode == http.StatusOK {
		l.Info(msg)
		return
	}

	l.ErrorS(msg)
}
