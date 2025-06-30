package services

import (
	"context"
	"time"

	"final_assessment/internal/models"

	"github.com/sirupsen/logrus"
)

type WorkerPool struct {
	JobService *JobService
	JobQueue   chan string
	NumWorkers int
}

func NewWorkerPool(service *JobService, numWorkers, queueSize int) *WorkerPool {
	return &WorkerPool{
		JobService: service,
		JobQueue:   make(chan string, queueSize),
		NumWorkers: numWorkers,
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.NumWorkers; i++ {
		go wp.worker(i + 1)
	}
	logrus.Infof("Started %d workers", wp.NumWorkers)
}

func (wp *WorkerPool) Enqueue(jobID string) {
	wp.JobQueue <- jobID
	logrus.WithField("job_id", jobID).Info("Job enqueued for processing")
}

func (wp *WorkerPool) worker(workerID int) {
	for jobID := range wp.JobQueue {
		logrus.WithFields(logrus.Fields{"worker": workerID, "job_id": jobID}).Info("Worker picked up job")
		// Mark job as processing
		err := wp.JobService.UpdateJobStatus(jobID, models.JobStatusProcessing, "")
		if err != nil {
			logrus.WithFields(logrus.Fields{"worker": workerID, "job_id": jobID, "error": err}).Error("Failed to mark job as processing")
			continue
		}
		// Simulate job processing
		time.Sleep(2 * time.Second) // Simulate work
		// Mark job as done
		result := "Job processed successfully"
		err = wp.JobService.UpdateJobStatus(jobID, models.JobStatusDone, result)
		if err != nil {
			logrus.WithFields(logrus.Fields{"worker": workerID, "job_id": jobID, "error": err}).Error("Failed to mark job as done")
			continue
		}
		logrus.WithFields(logrus.Fields{"worker": workerID, "job_id": jobID}).Info("Job completed")
	}
}

// UpdateJobStatus updates the status and result of a job
func (s *JobService) UpdateJobStatus(id string, status models.JobStatus, result string) error {
	ctx := context.Background()
	query := `UPDATE jobs SET status = $1, result = $2, updated_at = $3 WHERE id = $4`
	_, err := s.DB.ExecContext(ctx, query, status, result, time.Now(), id)
	return err
}
