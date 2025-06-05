package jobs

import (
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/arunjeyaprasad/golive/config"
	"github.com/arunjeyaprasad/golive/models"
	"github.com/arunjeyaprasad/golive/streamer"
	"github.com/google/uuid"
)

var (
	jobs          = make(map[string]models.Job)
	jobProcessMap = make(map[string]*streamer.StreamingProcess)
)

type JobStatus string

const (
	JobStatusCreated   JobStatus = "created"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "error"
)

// generateJobID generates a unique ID for a job.
func generateJobID() string {
	id := uuid.New()
	return id.String()
}

// createJob creates a new job with the given description and adds it to the jobs map.
func CreateJob(description string) *models.Job {
	job := models.Job{
		ID:          generateJobID(),
		Description: description,
		Status:      string(JobStatusCreated),
		CreatedAt:   time.Now().Format(time.RFC3339),
	}
	jobs[job.ID] = job
	return &job
}

func GetJobs() []models.Job {
	var jobList []models.Job
	for _, job := range jobs {
		jobList = append(jobList, job)
	}
	return jobList
}

func GetJob(id string) (*models.Job, bool) {
	job, exists := jobs[id]
	if !exists {
		return nil, false
	}
	return &job, true
}

func DeleteJob(id string) {
	// Stop the job if it's running
	if sp, exists := jobProcessMap[id]; exists {
		sp.StopJob()
		delete(jobProcessMap, id)
	}
	// Clean up the job's output directory to reclaim space
	path := filepath.Join(config.DEFAULT_MEDIA_DIR, id)
	if err := os.RemoveAll(path); err != nil {
		slog.Error("Failed to remove job directory", "path", path, "error", err)
	}
	// Remove the job from the jobs map
	if _, exists := jobs[id]; !exists {
		return // Job not found
	}
	delete(jobs, id)
}

func StartJob(job *models.Job) error {
	sp := streamer.NewStreamingProcess(job)
	if err := sp.StartJob(); err != nil {
		job.Status = string(JobStatusFailed)
		return err
	}
	jobProcessMap[job.ID] = sp

	// Update job status to running
	job.Status = string(JobStatusRunning)
	jobs[job.ID] = *job

	return nil
}

func StopJob(jobID string) error {
	sp, exists := jobProcessMap[jobID]
	if !exists {
		return nil // Job not found
	}

	// Terminate the streaming process
	if err := sp.StopJob(); err != nil {
		return err
	}

	// Update job status to completed
	job, exists := jobs[jobID]
	if !exists {
		return nil // Job not found
	}
	job.Status = string(JobStatusCompleted)
	jobs[jobID] = job

	return nil
}
