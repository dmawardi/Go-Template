package queue

import (
	"log"
	"sync"
	"time"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/email"
	"gorm.io/gorm"
)

// Queue represents a job queue backed by a SQL database.
type Queue struct {
	db          *gorm.DB   // Database connection
	mu          sync.Mutex // Mutex for synchronizing access
	cond        *sync.Cond // Condition variable for signaling
	mailService email.Email
}

// Class method for creating a new job queue
// Backed by the given database (uses the job table).
func NewQueue(db *gorm.DB, mailService email.Email) *Queue {
	// Create the queue
	q := &Queue{
		db:          db,
		mailService: mailService,
	}
	// Initialize the mutex and condition variable
	q.cond = sync.NewCond(&q.mu)
	return q
}

// AddJob adds a new job to the queue.
func (q *Queue) AddJob(jobType, payload string, process func(string) error) error {
	// Lock the queue
	q.mu.Lock()
	// Unlock the queue when the function returns
	defer q.mu.Unlock()

	// Create a new job
	job := db.Job{
		JobType: jobType,
		Payload: payload,
		Process: process,
	}

	// Store the job in the database
	if err := q.db.Create(&job).Error; err != nil {
		return err
	}
	q.cond.Signal() // Signal any waiting workers that a job is available
	return nil
}

// GetJob retrieves the next unprocessed job from the queue.
func (q *Queue) GetJob() (*db.Job, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	var job db.Job
	// Fetch the first unprocessed job from the database
	if err := q.db.Where("processed = ?", false).First(&job).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

// MarkJobAsProcessed marks a job as processed in the database.
func (q *Queue) MarkJobAsProcessed(job *db.Job) error {
	// Lock the queue
	q.mu.Lock()
	defer q.mu.Unlock()

	// Mark the job as processed
	job.Processed = true
	// Update the job in the database, returning any error
	return q.db.Save(job).Error
}

// Worker processes jobs from the queue.
func (q *Queue) Worker() {
	// Loop indefinitely
	for {
		// Get the next job
		job, err := q.GetJob()

		// If there are no jobs available, wait for a signal
		if err != nil {
			log.Println("No jobs available, waiting...")
			time.Sleep(1 * time.Second)
			continue
		}

		// Process the job
		err = job.Process(job.Payload)
		if err != nil {
			log.Printf("Failed to process job %d: %v", job.ID, err)
		} else {
			// Mark the job as processed
			err := q.MarkJobAsProcessed(job)
			if err != nil {
				log.Printf("Failed to mark job %d as processed: %v", job.ID, err)
			} else {
				log.Printf("Successfully processed job %d", job.ID)
			}
		}
	}
}
