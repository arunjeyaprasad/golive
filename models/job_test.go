package models

import "testing"

func TestJobCreateRequest_Validate(t *testing.T) {
	type fields struct {
		Description string
		VideoTrack  *VideoTrack
		AudioTrack  *AudioTrack
		AudioConfig *AudioConfig
		JobFormat   JobFormat
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Valid job create with defaults",
			fields: fields{
				Description: "Test job",
			},
			wantErr: false,
		},
		{
			name: "Valid job create with default video track",
			fields: fields{
				Description: "Test job with video",
				VideoTrack:  &VideoTrack{},
			},
			wantErr: false,
		},
		{
			name: "Valid job create with video track with only bitrate and other defaults",
			fields: fields{
				Description: "Test job with video track",
				VideoTrack: &VideoTrack{
					BitRate: "1000k",
				},
			},
			wantErr: false,
		},
		{
			name: "Valid job create with video track with all values provided",
			fields: fields{
				Description: "Test job with video track",
				VideoTrack: &VideoTrack{
					BitRate:    "1000k",
					Resolution: "1280x720",
					Framerate:  "30",
					Codec:      "h264",
				},
			},
			wantErr: false,
		},
		{
			name: "Valid job create with video track with only Resolution",
			fields: fields{
				Description: "Test job with video track",
				VideoTrack: &VideoTrack{
					Resolution: "1280x720",
				},
			},
			wantErr: false,
		},
		{
			name: "Valid job create with video track with only Framerate",
			fields: fields{
				Description: "Test job with video track",
				VideoTrack: &VideoTrack{
					Framerate: "30",
				},
			},
			wantErr: false,
		},
		{
			name: "Valid job create with video track with only Codec",
			fields: fields{
				Description: "Test job with video track",
				VideoTrack: &VideoTrack{
					Codec: "h264",
				},
			},
			wantErr: false,
		},
		{
			name: "Valid job create with audio track with all values filled",
			fields: fields{
				Description: "Test job with audio track",
				AudioTrack: &AudioTrack{
					AudioBitrate:    "128k",
					AudioCodec:      "aac",
					AudioSampleRate: "44100",
					AudioChannels:   "2",
				},
			},
			wantErr: false,
		},
		{
			name: "Valid job create with audio track with only AudioCodec",
			fields: fields{
				Description: "Test job with audio track",
				AudioTrack: &AudioTrack{
					AudioCodec: "aac",
				},
			},
			wantErr: false,
		},
		{
			name: "Valid job create with audio track with only AudioBitrate",
			fields: fields{
				Description: "Test job with audio track",
				AudioTrack: &AudioTrack{
					AudioBitrate: "128k",
				},
			},
			wantErr: false,
		},
		{
			name: "Valid job create with audio track with only AudioSampleRate",
			fields: fields{
				Description: "Test job with audio track",
				AudioTrack: &AudioTrack{
					AudioSampleRate: "44100",
				},
			},
			wantErr: false,
		},
		{
			name: "Valid job create with audio track with only AudioChannels",
			fields: fields{
				Description: "Test job with audio track",
				AudioTrack: &AudioTrack{
					AudioChannels: "2",
				},
			},
			wantErr: false,
		},
		{
			name: "Valid job with only segment length",
			fields: fields{
				Description: "Test job with segment length",
				JobFormat: JobFormat{
					SegmentLength: 10, // Valid segment length
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid job with negative segment length",
			fields: fields{
				Description: "Test job with negative segment length",
				JobFormat: JobFormat{
					SegmentLength: -5, // Invalid segment length
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid job with incorrect video params",
			fields: fields{
				Description: "Test job with invalid video track",
				VideoTrack: &VideoTrack{
					BitRate:    "1000k",
					Resolution: "1280", // Invalid resolution
					Framerate:  "30",
					Codec:      "h264",
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid job with incorrect video params",
			fields: fields{
				Description: "Test job with invalid video track",
				VideoTrack: &VideoTrack{
					BitRate:    "M", // Invalid bitrate
					Resolution: "1280x720",
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid job with incorrect video params",
			fields: fields{
				Description: "Test job with invalid video track",
				VideoTrack: &VideoTrack{
					BitRate:    "1000M", // Invalid bitrate
					Resolution: "1280x720",
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid job with incorrect video params",
			fields: fields{
				Description: "Test job with invalid video track",
				VideoTrack: &VideoTrack{
					BitRate:    "10G", // Invalid bitrate
					Resolution: "1280x720",
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid job with incorrect video params",
			fields: fields{
				Description: "Test job with invalid video track",
				VideoTrack: &VideoTrack{
					Codec: "unknown_codec", // Invalid codec
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid job with incorrect video resolution",
			fields: fields{
				Description: "Test job with invalid video track",
				VideoTrack: &VideoTrack{
					Resolution: "1280xabc", // Invalid resolution
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid job with incorrect audio params",
			fields: fields{
				Description: "Test job with invalid audio track",
				AudioTrack: &AudioTrack{
					AudioBitrate: "128M", // Invalid bitrate
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid job with incorrect audio params",
			fields: fields{
				Description: "Test job with invalid audio track",
				AudioTrack: &AudioTrack{
					AudioBitrate: "912k", // Invalid bitrate
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid job with incorrect audio params",
			fields: fields{
				Description: "Test job with invalid audio track",
				AudioTrack: &AudioTrack{
					AudioCodec: "unknown_codec", // Invalid codec
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid job with audio config with invalid sample rate",
			fields: fields{
				Description: "Test job with invalid audio config",
				AudioTrack: &AudioTrack{
					AudioSampleRate: "100", // Invalid sample rate
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid job with audio config with invalid sample rate",
			fields: fields{
				Description: "Test job with invalid audio config",
				AudioTrack: &AudioTrack{
					AudioSampleRate: "193000", // Invalid sample rate
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid job with audio config with invalid channels",
			fields: fields{
				Description: "Test job with invalid audio config",
				AudioTrack: &AudioTrack{
					AudioChannels: "20", // Invalid channels
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid job with audio config with invalid channels",
			fields: fields{
				Description: "Test job with invalid audio config",
				AudioTrack: &AudioTrack{
					AudioChannels: "0", // Invalid channels
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid job with audio config with invalid default language",
			fields: fields{
				Description: "Test job with invalid audio config",
				AudioConfig: &AudioConfig{
					AudioTracks:          2,
					AudioLanguages:       []string{"en", "fr"},
					AudioDefaultLanguage: "de", // Not in audio languages
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid job with audio config with mismatched tracks and languages",
			fields: fields{
				Description: "Test job with mismatched audio config",
				AudioConfig: &AudioConfig{
					AudioTracks:          2,
					AudioLanguages:       []string{"en"}, // Only one language provided
					AudioDefaultLanguage: "en",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jcr := JobCreateRequest{
				Description: tt.fields.Description,
				VideoTrack:  tt.fields.VideoTrack,
				AudioTrack:  tt.fields.AudioTrack,
				AudioConfig: tt.fields.AudioConfig,
				JobFormat:   tt.fields.JobFormat,
			}
			if err := jcr.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("JobCreateRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
