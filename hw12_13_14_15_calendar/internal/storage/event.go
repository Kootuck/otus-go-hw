package storage

import (
	"errors"
	"time"
)

// Common business errors for event storage.
var (
	ErrEventNotFound = errors.New("event not found")
	ErrInvalidEvent  = errors.New("invalid event data")
	ErrDateIsBusy    = errors.New("that time is occupied by another event")
)

type EventID string

type Event struct {
	ID          EventID
	Title       string
	StartTime   time.Time
	EndTime     time.Time
	Description string
	OwnerID     string
	NotifyOn    time.Time
}

type EventStorage interface {
	Add(event Event) error
	Update(id EventID, event Event) error
	Delete(id EventID) error
	ListDay(date time.Time) ([]Event, error)
	ListWeek(date time.Time) ([]Event, error)
	ListMonth(date time.Time) ([]Event, error)
}
