package queue

import (
	"errors"
	"fmt"
	"log"
	"time"
)

// Worker processes jobs from the queue.
func (q *Queue) Worker() {
	for {
		// Get the next job
		job, err := q.GetJob()
		if err != nil {
			log.Printf("Error getting job: %v\n", err)
			time.Sleep(5 * time.Second)
			continue
		}
		fmt.Printf("Processing job: %v\n", job.ID)
		// Process the job using the Process function with the payload
		if err := q.ProcessJob(job.JobType, job.Payload); err != nil {
			log.Printf("Error processing job: %v\n", err)
			time.Sleep(5 * time.Second)
			continue
		}
		fmt.Printf("Job processed: %v\n", job.ID)
		// Mark the job as processed
		if err := q.MarkJobAsProcessed(job); err != nil {
			log.Printf("Error marking job as processed: %v\n", err)
			time.Sleep(5 * time.Second)
			continue
		}
		fmt.Printf("Job marked as processed: %v\n", job.ID)
	}
}

func (q *Queue) ProcessJob(jobType, payload string) error {
	switch jobType {
	case "email":
		return q.ProcessEmailJob(payload)
	default:
		return errors.New("unknown job type")
	}
}
