package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/arunjeyaprasad/golive/config"
)

type Job struct {
	ID                 string           `json:"id"`
	Status             string           `json:"status"`
	CreatedAt          string           `json:"created"`
	StreamingStartedAt string           `json:"streamed_from,omitempty"`
	CompletedAt        string           `json:"completed,omitempty"`
	PlaybackURLs       []PlaybackURLs   `json:"playback_urls,omitempty"`
	Configuration      JobCreateRequest `json:"config"` // Original request that created this job
}

type JobCreateRequest struct {
	Description string       `json:"description"`
	VideoTrack  *VideoTrack  `json:"video,omitempty"`
	AudioTrack  *AudioTrack  `json:"audio,omitempty"`
	AudioConfig *AudioConfig `json:"audio_config,omitempty"`
	JobFormat
}

type JobResponse struct {
	ID string `json:"id"`
}

type VideoTrack struct {
	BitRate    string `json:"bitrate"`
	Resolution string `json:"resolution"`
	Framerate  string `json:"framerate"`
	Codec      string `json:"codec"`
}

type AudioTrack struct {
	AudioCodec      string `json:"codec"`
	AudioBitrate    string `json:"bitrate"`
	AudioSampleRate string `json:"sample_rate"`
	AudioChannels   string `json:"channels"`
}

type AudioConfig struct {
	AudioTracks          int      `json:"audio_tracks"`
	AudioLanguages       []string `json:"audio_languages,omitempty"`
	AudioDefaultLanguage string   `json:"audio_default_language,omitempty"`
}

type JobOutputFormat string

const (
	JobOutputFormatHLS  JobOutputFormat = "hls"
	JobOutputFormatDASH JobOutputFormat = "dash"
)

type JobFormat struct {
	OutputFormat  []JobOutputFormat `json:"output_format,omitempty"`
	SegmentLength int               `json:"segment_length,omitempty"` // Length of each segment in seconds
	WindowSize    int               `json:"window_size,omitempty"`    // Number of segments to keep in the playlist
}

type PlaybackURLs struct {
	Format JobOutputFormat `json:"format"`
	URL    string          `json:"url"`
}

func (jcr *JobCreateRequest) Validate() error {
	// Validate the Video and Audio Params
	var (
		errs []error
	)
	// Step 1: Assign Defaults if not provided
	if jcr.VideoTrack != nil {
		if jcr.VideoTrack.BitRate == "" {
			// Assign a Default Bitrate
			jcr.VideoTrack.BitRate = "1M" // Default bitrate
		}
		if jcr.VideoTrack.Resolution == "" {
			jcr.VideoTrack.Resolution = "1280x720" // Default resolution
		}
		if jcr.VideoTrack.Framerate == "" {
			jcr.VideoTrack.Framerate = "30" // Default framerate
		}
		if jcr.VideoTrack.Codec == "" {
			jcr.VideoTrack.Codec = "h264" // Default codec
		}
	} else {
		jcr.VideoTrack = &VideoTrack{
			BitRate:    "1M",       // Default bitrate
			Resolution: "1280x720", // Default resolution
			Framerate:  "30",       // Default framerate
			Codec:      "h264",     // Default codec
		}
	}
	if jcr.AudioTrack != nil {
		if jcr.AudioTrack.AudioCodec == "" {
			jcr.AudioTrack.AudioCodec = "aac" // Default audio codec
		}
		if jcr.AudioTrack.AudioBitrate == "" {
			jcr.AudioTrack.AudioBitrate = "128k" // Default audio bitrate
		}
		if jcr.AudioTrack.AudioSampleRate == "" {
			jcr.AudioTrack.AudioSampleRate = "44100" // Default sample rate
		}
		if jcr.AudioTrack.AudioChannels == "" {
			jcr.AudioTrack.AudioChannels = "2" // Default channels
		}
	} else {
		jcr.AudioTrack = &AudioTrack{
			AudioCodec:      "aac",   // Default audio codec
			AudioBitrate:    "128k",  // Default audio bitrate
			AudioSampleRate: "44100", // Default sample rate
			AudioChannels:   "2",     // Default channels
		}
	}
	if jcr.AudioConfig != nil {
		if jcr.AudioConfig.AudioTracks <= 0 || jcr.AudioConfig.AudioTracks > config.MAX_AUDIO_LANGUAGES {
			errs = append(errs, fmt.Errorf("audio_tracks must be greater than 0"))
		}
		if len(jcr.AudioConfig.AudioLanguages) != jcr.AudioConfig.AudioTracks {
			errs = append(errs, fmt.Errorf("audio_languages must match the number of audio_tracks"))
		}
		if jcr.AudioConfig.AudioDefaultLanguage != "" {
			found := false
			for _, lang := range jcr.AudioConfig.AudioLanguages {
				if strings.EqualFold(lang, jcr.AudioConfig.AudioDefaultLanguage) {
					found = true
					break
				}
			}
			if !found {
				errs = append(errs, fmt.Errorf("audio_default_language must be one of the audio_languages"))
			}
		}
	}
	if jcr.JobFormat.SegmentLength == 0 {
		jcr.JobFormat.SegmentLength = config.DEFAULT_SEGMENT_LENGTH // Default segment length in seconds
	} else if jcr.JobFormat.SegmentLength < 0 {
		errs = append(errs, fmt.Errorf("segment_length must be greater than 0"))
	}
	if jcr.JobFormat.WindowSize == 0 {
		jcr.JobFormat.WindowSize = config.DEFAULT_WINDOW_SIZE // Default window size
	}
	// Step 2: Now validate the Video and Audio Params
	if jcr.VideoTrack != nil {
		// Validate the bitrate
		if len(jcr.VideoTrack.BitRate) < 2 {
			errs = append(errs, fmt.Errorf("video bitrate must be at least 2 characters long"))
		}
		// Check if ends with k or M
		unit := jcr.VideoTrack.BitRate[len(jcr.VideoTrack.BitRate)-1]
		if !(len(jcr.VideoTrack.BitRate) > 1 && (unit == 'k' || unit == 'M')) {
			errs = append(errs, fmt.Errorf("video bitrate must end with k or M"))
		}
		// Extract the numeric part
		numericPart := jcr.VideoTrack.BitRate[:len(jcr.VideoTrack.BitRate)-1]
		videoBitRate, nerr := strconv.ParseFloat(numericPart, 64)
		if nerr != nil {
			errs = append(errs, fmt.Errorf("video bitrate must be a valid number"))
		}
		if unit == 'k' {
			if videoBitRate < 10 || videoBitRate > float64(config.MAX_VIDEO_BITRATE_MBPS*1000) {
				errs = append(errs, fmt.Errorf("video bitrate must be between 10k and %dk", config.MAX_VIDEO_BITRATE_MBPS*1000))
			}
		}
		if unit == 'M' {
			if videoBitRate < 0.01 || videoBitRate > float64(config.MAX_VIDEO_BITRATE_MBPS) {
				errs = append(errs, fmt.Errorf("video bitrate must be between 0.01M and %dM", config.MAX_VIDEO_BITRATE_MBPS))
			}
		}

		// Validate resolution
		resParts := strings.Split(jcr.VideoTrack.Resolution, "x")
		if len(resParts) != 2 {
			errs = append(errs, fmt.Errorf("video resolution must be in the format WxH (e.g., 1280x720)"))
		} else {
			width, werr := strconv.Atoi(resParts[0])
			height, herr := strconv.Atoi(resParts[1])
			if werr != nil || herr != nil {
				errs = append(errs, fmt.Errorf("video resolution must be valid integers"))
			} else {
				if width <= 0 || height <= 0 {
					errs = append(errs, fmt.Errorf("video resolution must be greater than 0"))
				}
				if width > config.MAX_VIDEO_WIDTH || height > config.MAX_VIDEO_HEIGHT {
					errs = append(errs, fmt.Errorf("video resolution must not exceed 3840x2160 (4K)"))
				}
			}
		}

		// Validate framerate
		framerate, ferr := strconv.Atoi(jcr.VideoTrack.Framerate)
		if ferr != nil {
			errs = append(errs, fmt.Errorf("video framerate must be a valid integer"))
		} else {
			if framerate <= 0 || framerate > config.MAX_VIDEO_FPS {
				errs = append(errs, fmt.Errorf("video framerate must be between 1 and %d", config.MAX_VIDEO_FPS))
			}
		}

		// Validate codec
		validCodec := false
		for _, codec := range config.VALID_VIDEO_CODECS {
			if strings.EqualFold(jcr.VideoTrack.Codec, codec) {
				validCodec = true
				break
			}
		}
		if !validCodec {
			errs = append(errs, fmt.Errorf("video codec must be one of: %v", config.VALID_VIDEO_CODECS))
		}
	}

	if jcr.AudioTrack != nil {
		// Validate audio codec
		validAudioCodec := false
		for _, codec := range config.VALID_AUDIO_CODECS {
			if strings.EqualFold(jcr.AudioTrack.AudioCodec, codec) {
				validAudioCodec = true
				break
			}
		}
		if !validAudioCodec {
			errs = append(errs, fmt.Errorf("audio codec must be one of: %v", config.VALID_AUDIO_CODECS))
		}
		// Validate audio bitrate
		if len(jcr.AudioTrack.AudioBitrate) < 3 {
			errs = append(errs, fmt.Errorf("audio bitrate must be at least 3 characters long"))
		}
		// Check if ends with k
		unit := jcr.AudioTrack.AudioBitrate[len(jcr.AudioTrack.AudioBitrate)-1]
		numericPart := jcr.AudioTrack.AudioBitrate[:len(jcr.AudioTrack.AudioBitrate)-1]

		if !(unit == 'k') {
			errs = append(errs, fmt.Errorf("audio bitrate must end with k"))
		}
		audioBitRate, nerr := strconv.ParseInt(numericPart, 10, 16)
		if nerr != nil {
			errs = append(errs, fmt.Errorf("audio bitrate must be a valid number"))
		}
		if audioBitRate < 32 || audioBitRate > int64(config.MAX_AUDIO_BITRATE_KBPS) {
			errs = append(errs, fmt.Errorf("audio bitrate must be between 32k and %dk", config.MAX_AUDIO_BITRATE_KBPS))
		}
		// Validate audio sample rate
		audioSampleRate, aerr := strconv.Atoi(jcr.AudioTrack.AudioSampleRate)
		if aerr != nil {
			errs = append(errs, fmt.Errorf("audio sample rate must be a valid integer"))
		} else {
			if audioSampleRate <= 0 {
				errs = append(errs, fmt.Errorf("audio sample rate must be greater than 0"))
			}
			if audioSampleRate < 8000 || audioSampleRate > 192000 {
				errs = append(errs, fmt.Errorf("audio sample rate must be between 8000 and 192000 Hz"))
			}
		}
		// Validate audio channels
		audioChannels, cerr := strconv.Atoi(jcr.AudioTrack.AudioChannels)
		if cerr != nil {
			errs = append(errs, fmt.Errorf("audio channels must be a valid integer"))
		} else {
			if audioChannels <= 0 {
				errs = append(errs, fmt.Errorf("audio channels must be greater than 0"))
			}
			if audioChannels > 8 {
				errs = append(errs, fmt.Errorf("audio channels must not exceed 8"))
			}
		}
	}

	return errors.Join(errs...)
}
