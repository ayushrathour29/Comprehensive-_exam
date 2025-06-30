package services

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"final_assessment/internal/db"
	"final_assessment/internal/models"
	"final_assessment/internal/utils"

	"github.com/sirupsen/logrus"
)

// JobService provides job-related operations
// Follows SOLID principles by depending on the db.Database interface
type JobService struct {
	DB db.Database
}

// NewJobService creates a new JobService
func NewJobService(database db.Database) *JobService {
	return &JobService{DB: database}
}

// CreateJob inserts a new job into the database
func (s *JobService) CreateJob(payload string) (*models.Job, error) {
	job := &models.Job{
		Payload:   payload,
		Status:    models.JobStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	ctx := context.Background()
	query := `INSERT INTO jobs (payload, status, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id`
	err := s.DB.QueryRowContext(ctx, query, job.Payload, job.Status, job.CreatedAt, job.UpdatedAt).Scan(&job.ID)
	if err != nil {
		logrus.WithError(err).Error(utils.LogFailedToInsertJob)
		return nil, err
	}
	logrus.WithField("job_id", job.ID).Info(utils.LogJobCreated)
	return job, nil
}

// GetJobByID retrieves a job by its ID
func (s *JobService) GetJobByID(id string) (*models.Job, error) {
	var job models.Job
	ctx := context.Background()
	query := `SELECT id, payload, status, result, created_at, updated_at FROM jobs WHERE id = $1`
	row := s.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&job.ID, &job.Payload, &job.Status, &job.Result, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		logrus.WithError(err).Error(utils.LogFailedToRetrieveJob)
		return nil, err
	}
	return &job, nil
}

// ListJobs returns a paginated list of jobs
func (s *JobService) ListJobs(limit, offset int) ([]*models.Job, error) {
	ctx := context.Background()
	query := `SELECT id, payload, status, result, created_at, updated_at FROM jobs ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := s.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		logrus.WithError(err).Error(utils.LogFailedToListJobs)
		return nil, err
	}
	defer rows.Close()

	var jobs []*models.Job
	for rows.Next() {
		var job models.Job
		if err := rows.Scan(&job.ID, &job.Payload, &job.Status, &job.Result, &job.CreatedAt, &job.UpdatedAt); err != nil {
			logrus.WithError(err).Error(utils.LogFailedToScanJobRow)
			continue
		}
		jobs = append(jobs, &job)
	}
	if err := rows.Err(); err != nil {
		logrus.WithError(err).Error(utils.LogRowIterationError)
		return nil, err
	}
	return jobs, nil
}
