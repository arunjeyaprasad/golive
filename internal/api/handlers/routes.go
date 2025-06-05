package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterRoutes registers API routes to the provided router.
func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/jobs", createJobHandler()).Methods(http.MethodPost)
	router.HandleFunc("/jobs", getJobsHandler()).Methods(http.MethodGet)
	router.HandleFunc("/jobs/{job_id}", getJobHandler()).Methods(http.MethodGet)
	router.HandleFunc("/jobs/{job_id}", cleanUpJobHandler()).Methods(http.MethodDelete)
	router.HandleFunc("/jobs/{job_id}/start", startJobHandler()).Methods(http.MethodPut)
	router.HandleFunc("/jobs/{job_id}/stop", stopJobHandler()).Methods(http.MethodPut)

	// Media Endpoints
	router.HandleFunc("/jobs/{job_id}/{file:.+}", getMediaHandler()).Methods(http.MethodGet)
}
