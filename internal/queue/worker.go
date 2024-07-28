package queue

import (
	"log"
	"time"
)

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
