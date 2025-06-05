package handlers

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/arunjeyaprasad/golive/config"
	"github.com/arunjeyaprasad/golive/internal/api/middleware"
	"github.com/arunjeyaprasad/golive/jobs"
)

func getMediaHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobid := r.Context().Value(middleware.RouteParamsKey).(map[string]string)["job_id"]
		file := r.Context().Value(middleware.RouteParamsKey).(map[string]string)["file"]

		if _, ok := jobs.GetJob(jobid); !ok {
			http.Error(w, "Job not found", http.StatusNotFound)
			return
		}

		// Check if file exists
		fileName := filepath.Join(config.DEFAULT_MEDIA_DIR, jobid, file)
		if !FileExists(fileName) {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		if strings.HasSuffix(fileName, ".mpd") {
			// For DASH, we need to return the MPD file
			w.Header().Set("Content-Type", "application/dash+xml")
		} else if strings.HasSuffix(fileName, ".m3u8") {
			// For HLS, we need to return the M3U8 file
			w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
		} else {
			// For other media files, set the appropriate content type
			w.Header().Set("Content-Type", "application/octet-stream")
		}

		w.WriteHeader(http.StatusOK)

		// Simulate sending the file content
		http.ServeFile(w, r, fileName)
	}
}
