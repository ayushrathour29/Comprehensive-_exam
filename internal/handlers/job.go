package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"final_assessment/internal/services"
	"final_assessment/internal/utils"
)

// JobHandler handles job-related HTTP requests
// Uses dependency injection for JobService
// Follows SOLID principles
type JobHandler struct {
	Service    *services.JobService
	WorkerPool *services.WorkerPool
}

// NewJobHandler creates a new JobHandler
func NewJobHandler(service *services.JobService, pool *services.WorkerPool) *JobHandler {
	return &JobHandler{Service: service, WorkerPool: pool}
}

// SubmitJob handles POST /jobs
func (h *JobHandler) SubmitJob(c *gin.Context) {
	var req struct {
		Payload string `json:"payload" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithError(err).Warn(utils.ErrInvalidPayload)
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrInvalidPayload})
		return
	}
	job, err := h.Service.CreateJob(req.Payload)
	if err != nil {
		logrus.WithError(err).Error(utils.ErrFailedToCreateJob)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.ErrFailedToCreateJob})
		return
	}
	// Enqueue the job for processing
	h.WorkerPool.Enqueue(job.ID)
	c.JSON(http.StatusCreated, job)
}

// GetJob handles GET /jobs/:id
func (h *JobHandler) GetJob(c *gin.Context) {
	id := c.Param("id")
	job, err := h.Service.GetJobByID(id)
	if err != nil {
		logrus.WithError(err).Warn(utils.LogJobNotFound)
		c.JSON(http.StatusNotFound, gin.H{"error": utils.LogJobNotFound})
		return
	}
	c.JSON(http.StatusOK, job)
}

// ListJobs handles GET /jobs
func (h *JobHandler) ListJobs(c *gin.Context) {
	limit := 10
	offset := 0
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	if o := c.Query("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}
	jobs, err := h.Service.ListJobs(limit, offset)
	if err != nil {
		logrus.WithError(err).Error(utils.LogFailedToListJobs)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.ErrFailedToListJobs})
		return
	}
	c.JSON(http.StatusOK, jobs)
}
