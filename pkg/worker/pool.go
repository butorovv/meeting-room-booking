package worker

import (
	"log"
	"sync"
	"time"
)

// NotificationJob представляет задачу для отправки уведомления
type NotificationJob struct {
	UserID    string
	UserEmail string
	StartTime time.Time
	BookingID string
}

// WorkerPool управляет пулом воркеров для обработки задач
type WorkerPool struct {
	jobs   chan NotificationJob
	mu     sync.RWMutex
	closed bool
	wg     sync.WaitGroup
	once   sync.Once
}

// NewWorkerPool создает и запускает новый пул воркеров
func NewWorkerPool(workers, queueSize int) *WorkerPool {
	if workers <= 0 {
		workers = 1
	}
	if queueSize <= 0 {
		queueSize = 100
	}

	wp := &WorkerPool{
		jobs: make(chan NotificationJob, queueSize),
	}

	wp.wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wp.wg.Done()
			for job := range wp.jobs {
				log.Printf("[worker] sending notification for booking %s to %s", job.BookingID, job.UserEmail)
			}
		}()
	}

	return wp
}

// Add добавляет новую задачу в очередь
// возвращает false, если пул закрыт или очередь заполнена (неблокирующая операция)
func (wp *WorkerPool) Add(job NotificationJob) bool {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	if wp.closed {
		return false
	}

	select {
	case wp.jobs <- job:
		return true
	default:
		// Очередь заполнена, задача не добавлена
		return false
	}
}

// Shutdown инициирует корректное завершение работы пула воркеров
// метод идемпотентен
func (wp *WorkerPool) Shutdown() {
	wp.once.Do(func() {
		wp.mu.Lock()
		wp.closed = true
		wp.mu.Unlock()

		close(wp.jobs)
		wp.wg.Wait()
	})
}
