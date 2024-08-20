package memorystorage

import (
	"sort"
	"sync"
	"time"

	//nolint:depguard
	"github.com/Kootuck/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
	//nolint:depguard
	"github.com/Kootuck/otus-go-hw/hw12_13_14_15_calendar/internal/util"
)

type MemoryStorage struct {
	mu            sync.RWMutex
	eventsByStart []storage.Event
}

func New() *MemoryStorage {
	return &MemoryStorage{
		eventsByStart: []storage.Event{},
	}
}

// Add event to the storage.
func (m *MemoryStorage) Add(event storage.Event) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check for time conflicts.
	if err := m.checkTimeConflict(event.StartTime, event.EndTime); err != nil {
		return err
	}

	m.addEvent(event)

	return nil
}

// Check for overlapping events in the specified time range.
func (m *MemoryStorage) checkTimeConflict(start, end time.Time) error {
	for _, event := range m.eventsByStart {
		if start.Before(event.EndTime) && end.After(event.StartTime) {
			return storage.ErrDateIsBusy
		}
	}
	return nil
}

// Update an existing event.
func (m *MemoryStorage) Update(id storage.EventID, event storage.Event) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	idx, err := m.findEventByID(id)
	if err != nil {
		return err
	}
	oldEvent := m.eventsByStart[idx]

	// check that updated dates don't overlap.
	if oldEvent.StartTime != event.StartTime || oldEvent.EndTime != event.EndTime {
		if err := m.checkTimeConflict(event.StartTime, event.EndTime); err != nil {
			return err
		}
	}
	m.deleteEvent(id)
	m.addEvent(event)
	return nil
}

// Delete an event by ID.
func (m *MemoryStorage) Delete(id storage.EventID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, err := m.findEventByID(id)
	if err != nil {
		return err
	}

	m.deleteEvent(id)
	return nil
}

// List events for a specific day.
func (m *MemoryStorage) ListDay(date time.Time) ([]storage.Event, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	start := date.Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)

	return m.listRange(start, end)
}

func (m *MemoryStorage) ListWeek(date time.Time) ([]storage.Event, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	start := util.FirstDayOfWeek(date)
	end := start.AddDate(0, 0, 7)

	return m.listRange(start, end)
}

// List events for a specific month.
func (m *MemoryStorage) ListMonth(date time.Time) ([]storage.Event, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	start := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	return m.listRange(start, end)
}

// Helper function to list events within a date range.
func (m *MemoryStorage) listRange(start, end time.Time) ([]storage.Event, error) {
	var results []storage.Event
	for _, event := range m.eventsByStart {
		if event.StartTime.After(end) {
			break
		}
		if event.StartTime.Before(end) && event.EndTime.After(start) {
			results = append(results, event)
		}
	}
	return results, nil
}

// Helper function to delete an event from the sorted slice.
func (m *MemoryStorage) deleteEvent(id storage.EventID) {
	for i, event := range m.eventsByStart {
		if event.ID == id {
			m.eventsByStart = append(m.eventsByStart[:i], m.eventsByStart[i+1:]...)
			break
		}
	}
}

// Helper function to add an event and maintain sorted order (used internally).
func (m *MemoryStorage) addEvent(event storage.Event) {
	// Insert event into the sorted slice in the correct position.
	index := sort.Search(len(m.eventsByStart), func(i int) bool {
		return m.eventsByStart[i].StartTime.After(event.StartTime)
	})
	m.eventsByStart = append(m.eventsByStart, storage.Event{})
	copy(m.eventsByStart[index+1:], m.eventsByStart[index:])
	m.eventsByStart[index] = event
}

// Find the index of the event in storage slice.
func (m *MemoryStorage) findEventByID(id storage.EventID) (idx int, err error) {
	for i, e := range m.eventsByStart {
		if e.ID == id {
			return i, nil
		}
	}
	return -1, storage.ErrEventNotFound
}
