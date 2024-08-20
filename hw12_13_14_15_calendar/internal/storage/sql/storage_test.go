//go:build manual
// +build manual

package sqlstorage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	config "github.com/Kootuck/otus-go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Kootuck/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/Kootuck/otus-go-hw/hw12_13_14_15_calendar/internal/util"
	"github.com/stretchr/testify/assert"
)

/** Before executing tests: set up test PG in docker container ->

	1. Run container
	docker run -d \
	-p 5432:5432 \
	--name my-pg \
	-e POSTGRES_PASSWORD=<password> \
	postgres:latest

	2. Connect ->
	docker exec -it my-pg psql -U postgres -d postgres

	3. Create test DB ->
	CREATE DATABASE otus_hw_db;

	4. Create test user
	CREATE USER otus_hw_user WITH ENCRYPTED PASSWORD 'password';

	5. Set-up test user and the table (repeat this step before test run)
	cat internal/storage/sql/create_table.sql | docker exec -i my-pg  psql -U postgres -d otus_hw_db

**/

func TestConnect(t *testing.T) {
	ctx := context.Background()
	sqlStorage, err := connect(ctx)
	assert.NoError(t, err, "must connect without error")
	defer sqlStorage.Close(ctx)
}

func connect(ctx context.Context) (*SqlStorage, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	cnfg, err := config.NewConfig(configPath)
	if err != nil {
		return nil, err
	}

	ss := New(cnfg.Storage.PG)
	err = ss.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return ss, nil
}

func getConfigPath() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("cannot get current file path")
	}
	dir := filepath.Dir(filename)
	configPath := filepath.Join(dir, "../../../configs/default.yaml")

	return configPath, nil
}

func TestDB(t *testing.T) {
	type tCases struct {
		start     time.Time
		end       time.Time
		increment time.Duration
	}
	now := time.Now()
	var cases []tCases = []tCases{
		{ // day
			start:     time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
			end:       time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).Add(time.Hour * 24),
			increment: time.Hour + time.Second,
		},
		{ // week
			start:     util.FirstDayOfWeek(now.AddDate(0, 0, 8)),
			end:       util.LastDayOfWeek(now.AddDate(0, 0, 8)),
			increment: time.Hour*24 + time.Second,
		},
		{ // month
			start:     util.FirstDayOfMonth(now.AddDate(0, 2, 0)),
			end:       util.LastDayOfMonth(now.AddDate(0, 2, 0)),
			increment: time.Hour*24 + time.Second,
		},
	}

	var dayCount, weekCount, monthCount int

	sql, err := connect(context.Background())
	if err != nil {
		log.Fatal("connect err:", err)
	}
	defer sql.Close(context.Background())

	var wg sync.WaitGroup

	for i, tcase := range cases {
		wg.Add(1)
		go func(c tCases, categoryIndex int) {
			defer wg.Done()
			var localCount int
			start := c.start
			for {
				start = start.Add(c.increment)
				if start.After(c.end) {
					break
				}

				event := storage.Event{
					OwnerID:     "Test_EventOwner_ID",
					Title:       fmt.Sprintf("Event of %v", start),
					Description: fmt.Sprintf("Description for event %v", start),
					StartTime:   start,
					EndTime:     start.Add(time.Hour),
				}
				err := sql.Add(context.Background(), event)
				if err != nil {
					log.Fatal("add err:", err)
				}
				localCount++
			}

			switch categoryIndex {
			case 0:
				dayCount += localCount
			case 1:
				weekCount += localCount
			case 2:
				monthCount += localCount
			}
		}(tcase, i)

		wg.Wait()

		// now let's query back added data
		events, err := sql.ListDay(context.Background(), cases[0].start)
		assert.NoError(t, err)
		assert.Len(t, events, dayCount, "should find all added for a day")

		events, err = sql.ListWeek(context.Background(), cases[1].start)
		assert.NoError(t, err)
		assert.Len(t, events, weekCount, "should find all added for a week")

		events, err = sql.ListMonth(context.Background(), cases[2].start)
		assert.NoError(t, err)
		assert.Len(t, events, monthCount, "should find all added for a month")

	}
}
