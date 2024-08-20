package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	//nolint:depguard
	conf "github.com/Kootuck/otus-go-hw/hw12_13_14_15_calendar/internal/config"
	//nolint:depguard
	"github.com/Kootuck/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
	//nolint:depguard
	"github.com/Kootuck/otus-go-hw/hw12_13_14_15_calendar/internal/util"
	//nolint:depguard
	"github.com/jmoiron/sqlx"
)

// _ "github.com/jackc/pgx/v4/stdlib".
type SQLStorage struct {
	// internal data structure to hold events
	db     *sqlx.DB
	config conf.PGConnection
}

func (s *SQLStorage) Connect() error {
	var err error

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		s.config.DBUser, s.config.DBPassword, s.config.DBName, s.config.DBHost, s.config.DBPort)

	s.db, err = sqlx.Open("pgx", dsn)
	if err != nil {
		return err
	}
	err = s.db.Ping()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	return nil
}

func (s *SQLStorage) Close() error {
	e := s.db.Close()
	if e != nil {
		return e
	}
	return nil
}

func New(c conf.PGConnection) *SQLStorage {
	return &SQLStorage{config: c}
}

func (s *SQLStorage) Add(ctx context.Context, e storage.Event) error {
	query := `INSERT INTO events (owner_id, title, descr, start_time, end_time) VALUES ($1, $2, $3, $4, $5)`
	dura := time.Second * time.Duration(s.config.Timeout)
	ctx, cancel := context.WithTimeout(ctx, dura)
	defer cancel()
	_, err := s.db.ExecContext(ctx, query, e.OwnerID, e.Title, e.Description, e.StartTime, e.EndTime)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLStorage) Update(ctx context.Context, id storage.EventID, event storage.Event) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(s.config.Timeout))
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var oldStartTime, oldEndTime time.Time
	query := `SELECT start_time, 
					 end_time 
			  FROM events 
			  WHERE id = $1`
	err = tx.QueryRowContext(ctx, query, id).Scan(&oldStartTime, &oldEndTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("event not found")
		}
		return err
	}

	// Check if the times are being changed and if they conflict with existing events.
	if oldStartTime != event.StartTime || oldEndTime != event.EndTime {
		if err := s.checkTimeConflict(ctx, tx, event.StartTime, event.EndTime); err != nil {
			return err
		}
	}

	updateQuery := `UPDATE events 
				    SET owner_id = $2, 
						title = $3, 
						descr = $4, 
						start_time = $5, 
						end_time = $6 
					WHERE id = $1`
	_, err = tx.ExecContext(ctx, updateQuery, id, event.OwnerID,
		event.Title, event.Description, event.StartTime, event.EndTime)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *SQLStorage) Delete(ctx context.Context, id storage.EventID) error {
	query := `DELETE 
			  FROM events 
			  WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(s.config.Timeout))
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLStorage) ListDay(ctx context.Context, date time.Time) (events []storage.Event, err error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	events, err = s.listRange(ctx, startOfDay, endOfDay)
	return events, err
}

func (s *SQLStorage) ListMonth(ctx context.Context, date time.Time) ([]storage.Event, error) {
	start := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	return s.listRange(ctx, start, end)
}

func (s *SQLStorage) ListWeek(ctx context.Context, date time.Time) ([]storage.Event, error) {
	start := util.FirstDayOfWeek(date)
	end := start.AddDate(0, 0, 7)

	return s.listRange(ctx, start, end)
}

// checkTimeConflict checks for time conflicts with other events.
func (s *SQLStorage) checkTimeConflict(ctx context.Context, tx *sql.Tx, startTime, endTime time.Time) error {
	var count int
	query := `SELECT COUNT(*) 
			  FROM events 
			  WHERE (start_time <= $2 AND end_time >= $1)`
	err := tx.QueryRowContext(ctx, query, startTime, endTime).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return storage.ErrDateIsBusy
	}
	return nil
}

func (s *SQLStorage) listRange(ctx context.Context, start, end time.Time) (events []storage.Event, err error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(s.config.Timeout))
	defer cancel()

	query := `SELECT id, owner_id, title, descr, start_time, end_time 
			  FROM events 
			  WHERE 
			  start_time <= $2 AND 
			  end_time >= $1`
	rows, err := s.db.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e storage.Event
		var id int64
		err = rows.Scan(&id, &e.OwnerID, &e.Title, &e.Description, &e.StartTime, &e.EndTime)
		if err != nil {
			return nil, err
		}
		e.ID = storage.EventID(strconv.FormatInt(id, 10))
		events = append(events, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
