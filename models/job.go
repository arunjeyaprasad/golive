package models

type Job struct {
	ID                 string `json:"id"`
	Description        string `json:"description"`
	Status             string `json:"status"`
	PID                int    `json:"-"`
	CreatedAt          string `json:"created"`
	StreamingStartedAt string `json:"streamed_from,omitempty"`
	CompletedAt        string `json:"completed,omitempty"`
	JobFormat
	PlaybackURLs []PlaybackURLs `json:"playback_urls,omitempty"`
}

type JobCreateRequest struct {
	Description string `json:"description"`
}

type JobResponse struct {
	ID string `json:"id"`
}

type JobConfig struct {
	BitRate         string `json:"bitrate"`
	Resolution      string `json:"resolution"`
	Framerate       string `json:"framerate"`
	Codec           string `json:"codec"`
	AudioCodec      string `json:"audio_codec"`
	AudioBitrate    string `json:"audio_bitrate"`
	AudioSampleRate string `json:"audio_sample_rate"`
	AudioChannels   string `json:"audio_channels"`
}

type JobOutputFormat string

const (
	JobOutputFormatHLS  JobOutputFormat = "hls"
	JobOutputFormatDASH JobOutputFormat = "dash"
)

type JobFormat struct {
	OutputFormat []JobOutputFormat `json:"output_format,omitempty"`
}

type PlaybackURLs struct {
	Format JobOutputFormat `json:"format"`
	URL    string          `json:"url"`
}
