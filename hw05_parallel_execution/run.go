package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	// 0. Инит
	var errCount int32
	var wg sync.WaitGroup
	var numWorkers int

	if n > len(tasks) {
		numWorkers = len(tasks)
	} else {
		numWorkers = n
	}

	// work := make(chan func() error, numWorkers)
	work := make(chan Task)
	errs := make(chan error, len(tasks))
	done := make(chan struct{})

	defer func() {
		close(errs)
		if !isDone(done) {
			close(done)
		}
	}()

	// 1. Создать n обработчиков заданий.
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(work, errs, &wg)
	}

	// 2. Горутина отслеживающаяя превышение лимита ошибок.
	go errorCountThresholdObserver(errs, done, &errCount, m)

	// 3. Отправляем работу воркерам, а если накопился критический порог ошибок - выходим.
loop:
	for _, t := range tasks {
		select {
		case <-done:
		default:
			select {
			case <-done:
				break loop
			case work <- t:
			}
		}
	}

	// 4. Закрываем канал задач, ждем завершения и выходим.
	close(work)
	wg.Wait()

	if m > 0 && errCount >= int32(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}

// Done канал закрывается если превышен лимит ошибок.
// Если нет - то потребуется закрыть его при выходе из основной горутины, чтобы observer мог завершиться.
func isDone(done <-chan struct{}) bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

// Выполняет задачу, если произошла ошибка - пишет в канал ошибок.
func worker(work <-chan Task, errs chan<- error, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()

	for {
		t, ok := <-work
		if ok {
			err := t()      // Получили задачу, выполняем.
			if err != nil { // Ошибка - передать в канал ошибок.
				errs <- err
			}
		}
		if !ok { // канал был закрыт или больше нет работы.
			return
		}
	}
}

// Отслеживает канал ошибок и прерывает выполнение через done канал в случае превышения порога.
func errorCountThresholdObserver(errs <-chan error, done chan<- struct{}, errCount *int32, maxErrors int) {
	for err := range errs {
		if err != nil {
			atomic.AddInt32(errCount, 1)
		}
		if maxErrors > 0 {
			if *errCount >= int32(maxErrors) {
				close(done)
				return
			}
		}
	}
}
