package memorystorage

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	//nolint:depguard
	"github.com/Kootuck/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
	//nolint:depguard
	"github.com/Kootuck/otus-go-hw/hw12_13_14_15_calendar/internal/util"
	//nolint:depguard
	"github.com/stretchr/testify/assert"
)

func TestAddEvent(t *testing.T) {
	ms := New()
	event := storage.Event{
		ID:        "ADD_1",
		Title:     "Test Event",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
	}
	err := ms.Add(event)
	assert.NoError(t, err, "should add event without error")
}

func TestAddEventWithConflict(t *testing.T) {
	ms := New()
	startTime := time.Now()
	event1 := storage.Event{
		ID:        "CONF_1",
		Title:     "Test Event 1",
		StartTime: startTime,
		EndTime:   startTime.Add(time.Hour),
	}
	event2 := storage.Event{
		ID:        "CONF_2",
		Title:     "Test Event 2",
		StartTime: startTime,
		EndTime:   startTime.Add(time.Hour),
	}
	ms.Add(event1)
	err := ms.Add(event2)
	assert.Equal(t, storage.ErrDateIsBusy, err, "should detect time conflict")
}

func TestUpdateEvent(t *testing.T) {
	ms := New()
	event := storage.Event{
		ID:        "UPD_1",
		Title:     "Test Event",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
	}
	ms.Add(event)
	event.Title = "Updated Event"
	err := ms.Update("UPD_1", event)
	assert.NoError(t, err, "should update event without error")
}

func TestDeleteEvent(t *testing.T) {
	ms := New()
	event := storage.Event{
		ID:        "DEL_1",
		Title:     "Test Event",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
	}
	ms.Add(event)
	err := ms.Delete("DEL_1")
	assert.NoError(t, err, "should delete event without error")
	events, err := ms.ListDay(time.Now())
	assert.Len(t, events, 0, "should be no events stored")
	assert.NoError(t, err, "should be no error when no events")
}

func TestListDay(t *testing.T) {
	ms := New()
	date := time.Now()

	theMainEvent := storage.Event{
		ID:        "LIST_1",
		Title:     "The Main event",
		StartTime: date,
		EndTime:   date.Add(time.Hour),
	}
	ms.Add(theMainEvent)

	for i := 1; i < 100; i++ {
		event := storage.Event{
			ID:        storage.EventID(strconv.Itoa(i)),
			Title:     fmt.Sprintf("Noise event %v", i),
			StartTime: date.Add(time.Duration(i) * time.Hour * 24),
			EndTime:   date.Add(time.Duration(i) * time.Hour * 25),
		}
		ms.Add(event)
	}
	events, err := ms.ListDay(date)
	assert.NoError(t, err, "should list events without error")
	assert.Len(t, events, 1, "should list one event")
	assert.Equal(t, theMainEvent, events[0])
}

func TestListMonth(t *testing.T) {
	first := util.FirstDayOfMonth(time.Now())
	last := util.LastDayOfMonth(time.Now())
	ms := New()

	inserted := insertEvents(first, last, ms)
	events, err := ms.ListMonth(first)

	assert.NoError(t, err, "should list events without error")
	assert.Len(t, events, inserted, "should list all inserted events")
}

func TestListWeek(t *testing.T) {
	first := util.FirstDayOfWeek(time.Now())
	last := util.LastDayOfWeek(time.Now())
	ms := New()

	// insert some events on the current week; all events last 30mins
	inserted := insertEvents(first, last, ms)
	events, err := ms.ListWeek(first)

	assert.NoError(t, err, "should list events without error")
	assert.Len(t, events, inserted, "should list all inserted events")
}

func TestConcurrency(t *testing.T) {
	ms := New()
	var wg sync.WaitGroup
	startTime := time.Now()
	for i := 0; i < 33; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			event := storage.Event{
				ID:        storage.EventID(fmt.Sprintf("%d", id)),
				Title:     fmt.Sprintf("Event %d", id),
				StartTime: startTime.Add(time.Duration(id) * time.Minute),
				EndTime:   startTime.Add(time.Duration(id+1) * time.Minute),
			}
			err := ms.Add(event)
			assert.NoError(t, err)
		}(i)
	}
	wg.Wait()
	events, err := ms.ListDay(startTime)
	assert.NoError(t, err)
	assert.Len(t, events, 33, "should have 33 events")
}

func insertEvents(first, last time.Time, memoryStorage *MemoryStorage) (inserted int) {
	runningDate := first
	for {
		if k := rand.Intn(2); k == 0 { // 0 -> add one day
			runningDate = runningDate.Add(time.Hour * 24)
		} else { // 1 -> add one hour
			runningDate = runningDate.Add(time.Hour * 1)
		}
		if runningDate.Before(last) {
			memoryStorage.addEvent(storage.Event{
				ID:          storage.EventID(runningDate.String()),
				Title:       fmt.Sprintf("Event of %v", runningDate),
				Description: fmt.Sprintf("Description for event %v", runningDate),
				StartTime:   runningDate,
				EndTime:     runningDate.Add(time.Hour * 1 / 2),
			})
			inserted++
			continue
		}
		break
	}
	return inserted
}
