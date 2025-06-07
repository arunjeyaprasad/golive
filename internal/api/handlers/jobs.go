package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/arunjeyaprasad/golive/internal/api/middleware"
	"github.com/arunjeyaprasad/golive/internal/api/postprocessor"
	"github.com/arunjeyaprasad/golive/jobs"
	"github.com/arunjeyaprasad/golive/models"
)

func createJobHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var (
			request models.JobCreateRequest
			job     *models.Job
		)
		// Decode the request body into the job struct
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request payload; Missing Description in Body", http.StatusBadRequest)
			return
		}
		// Validate the request
		if err := request.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		job = jobs.CreateJob(request)
		postprocessor.FormatResponse(w, job, http.StatusCreated)
	}
}

func getJobsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobsList := jobs.GetJobs()
		postprocessor.FormatResponse(w, jobsList, http.StatusOK)
	}
}

func getJobHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the job id from the context route params
		jobid := r.Context().Value(middleware.RouteParamsKey).(map[string]string)["job_id"]
		if job, ok := jobs.GetJob(jobid); !ok {
			http.Error(w, "Job not found", http.StatusNotFound)
			return
		} else {
			postprocessor.FormatResponse(w, job, http.StatusOK)
		}
	}
}

func startJobHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobid := r.Context().Value(middleware.RouteParamsKey).(map[string]string)["job_id"]
		var (
			job *models.Job
			ok  bool
		)
		if job, ok = jobs.GetJob(jobid); !ok {
			http.Error(w, "Job not found", http.StatusNotFound)
			return
		}
		if err := jobs.StartJob(job); err != nil {
			slog.Error("Failed to start job", "job_id", jobid, "error", err)
			http.Error(w, "Failed to start job", http.StatusInternalServerError)
			return
		}

		postprocessor.FormatResponse(w, models.JobResponse{ID: jobid}, http.StatusOK)
	}
}

func stopJobHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobid := r.Context().Value(middleware.RouteParamsKey).(map[string]string)["job_id"]
		// Check if the job exists
		if _, ok := jobs.GetJob(jobid); !ok {
			http.Error(w, "Job not found", http.StatusNotFound)
			return
		}
		if err := jobs.StopJob(jobid); err != nil {
			slog.Error("Failed to stop job", "job_id", jobid, "error", err)
			http.Error(w, "Failed to stop job", http.StatusInternalServerError)
			return
		}
		postprocessor.FormatResponse(w, models.JobResponse{ID: jobid}, http.StatusOK)
	}
}

func cleanUpJobHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobid := r.Context().Value(middleware.RouteParamsKey).(map[string]string)["job_id"]
		var (
			job *models.Job
			ok  bool
		)
		if job, ok = jobs.GetJob(jobid); !ok {
			http.Error(w, "Job not found", http.StatusNotFound)
			return
		}
		if job.Status != string(jobs.JobStatusCompleted) {
			http.Error(w, "Job is not completed", http.StatusBadRequest)
			return
		}
		jobs.DeleteJob(jobid)
		postprocessor.FormatResponse(w, models.JobResponse{ID: jobid}, http.StatusOK)
	}
}

func FileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		slog.Error("File does not exist", "file", fileName)
		return false
	}
	if err != nil {
		slog.Error("Error checking file", "file", fileName, "error", err)
		return false
	}
	return true
}
