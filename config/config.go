package config

import (
	"log/slog"
	"os"
)

var (
	DEFAULT_SERVER_PORT = 9090
	MAX_JOB_COUNT       = 2
	DEFAULT_MEDIA_DIR   = "media"
)

func Init() {
	info, err := os.Stat(DEFAULT_MEDIA_DIR)
	if err == nil && info.IsDir() {
		// Directory already exists, no need to create it
		return
	}
	slog.Info("Creating media directory", "path", DEFAULT_MEDIA_DIR)
	os.Mkdir(DEFAULT_MEDIA_DIR, os.ModePerm)
}
