package config

import (
	"log/slog"
	"os"
)

var (
	MAX_JOB_COUNT              = 2
	DEFAULT_SERVER_PORT        = 9090
	DEFAULT_MEDIA_DIR          = "media"
	DEFAULT_SEGMENT_LENGTH     = 6      // 6 seconds
	DEFAULT_WINDOW_SIZE        = 6      // 6 segments
	DEFAULT_VIDEO_CODEC        = "h264" // Default video codec
	DEFAULT_AUDIO_CODEC        = "aac"  // Default audio codec
	DEFAULT_VIDEO_BITRATE_MBPS = 1      // 1 Mbps
	DEFAULT_AUDIO_BITRATE_KBPS = 128    // 128 Kbps
	DEFAULT_VIDEO_FPS          = 30     // 30 FPS
	DEFAULT_VIDEO_WIDTH        = 1280   // 1280 pixels (HD)
	DEFAULT_VIDEO_HEIGHT       = 720    // 720 pixels (HD)
	MAX_VIDEO_BITRATE_MBPS     = 35     // 35 Mbps
	MAX_AUDIO_BITRATE_KBPS     = 512    // 512 Kbps
	MAX_VIDEO_FPS              = 60     // 60 FPS
	MAX_VIDEO_WIDTH            = 3840   // 3840 pixels (4K)
	MAX_VIDEO_HEIGHT           = 2160   // 2160 pixels (4K)
	MAX_AUDIO_LANGUAGES        = 16     // Maximum number of audio languages supported
	VALID_VIDEO_CODECS         = []string{"h264", "hevc", "vp9", "av1"}
	VALID_AUDIO_CODECS         = []string{"aac", "mp3"}
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
