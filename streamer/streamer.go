package streamer

import (
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/arunjeyaprasad/golive/config"
	"github.com/arunjeyaprasad/golive/models"
	"github.com/fsnotify/fsnotify"
)

type StreamingProcess struct {
	Job                  *models.Job
	Pid                  int
	OutDir               string
	monitoringChannel    chan bool
	lastSegmentCreatedAt int64 // Timestamp of the last segment created
	channelClosed        bool
}

func NewStreamingProcess(job *models.Job) *StreamingProcess {
	return &StreamingProcess{
		Job:               job,
		OutDir:            filepath.Join(config.DEFAULT_MEDIA_DIR, job.ID),
		monitoringChannel: make(chan bool),
	}
}

func (sp *StreamingProcess) StartJob() error {
	// Start the job using ffmpeg
	cmd := sp.buildCommand(sp.Job)
	// Create the output directory
	err := os.MkdirAll(sp.OutDir, os.ModePerm)
	if err != nil {
		slog.Error("Failed to create output directory", "error", err)
		return err
	}
	// Execute the command
	go func() {
		defer func() {
			if err != nil {
				slog.Error("Streaming Command failed", "error", err)
			}
		}()
		execCmd := exec.Command(cmd[0], cmd[1:]...)
		err := execCmd.Start()
		if err != nil {
			slog.Error("Failed to start command", "error", err)
			return
		}
		sp.Pid = execCmd.Process.Pid
		// Log the command and PID
		slog.Info("Streaming Command started", "command", cmd, "pid", execCmd.Process.Pid)
		// Wait for the command to finish
		err = execCmd.Wait()
		if err != nil {
			slog.Error("Encoding Command failed with error", "error", err)
		}
	}()
	if err := sp.MonitorDirectory(); err != nil {
		slog.Error("Failed to start directory monitoring", "error", err)
	}
	return nil
}

func (sp *StreamingProcess) buildCommand(job *models.Job) []string {
	// Compute Filter complex using provided description
	var text string
	if job.Description == "" {
		text = "Test Live Stream"
	} else {
		text = job.Description
	}
	if len(text) > 50 {
		text = text[:50] // Limit to 50 characters
	}
	slog.Info("Building command for job", "jobID", job.ID, "description", text)

	filterString := "[0:v]drawtext=text='REPLACE_ME':fontsize=42:fontcolor=white:x=50+500*abs(sin(t/2)):y=(h-text_h)/3:box=1:boxcolor=black@0.7,drawtext=text='Frame %{frame_num}':fontsize=28:fontcolor=cyan:x=10:y=h-40:box=1:boxcolor=black@0.7[v]; [1:a]aloop=loop=-1:size=22050[a]"
	filterString = strings.ReplaceAll(filterString, "REPLACE_ME", text)
	return []string{
		"ffmpeg",
		"-re",
		"-f", "lavfi",
		"-i", "testsrc=size=1280x720:rate=25",
		"-f", "lavfi",
		"-i", "sine=frequency=1200:duration=0.03,afade=t=out:st=0.02:d=0.01,apad=pad_dur=0.97",
		"-filter_complex", filterString,
		"-map", "[v]",
		"-map", "[a]",
		"-c:v", "libx264",
		"-g", "150",
		"-keyint_min", "150",
		"-x264-params", "scenecut=0:open_gop=0",
		"-preset", "fast",
		"-c:a", "aac",
		"-f", "dash",
		"-seg_duration", "6",
		"-window_size", "6",
		"-use_template", "1",
		"-use_timeline", "1",
		"-hls_playlist", "1",
		"-streaming", "1",
		"-write_prft", "1",
		// "-ldash", "1",
		"-y", filepath.Join(sp.OutDir, "manifest.mpd"), // Will generate HLS manifest and segments in the output directory
	}
}

func (sp *StreamingProcess) StopJob() error {
	if sp.channelClosed {
		slog.Info("Job already stopped or monitoring channel closed", "jobID", sp.Job.ID)
		return nil
	}
	slog.Info("Stopping job", "jobID", sp.Job.ID, "PID", sp.Pid)
	// Stop monitoring the directory
	if sp.monitoringChannel != nil {
		sp.monitoringChannel <- true
		sp.channelClosed = true
	}
	// Stop the process if it's running
	if sp.Pid != 0 {
		process, err := os.FindProcess(sp.Pid)
		if err != nil {
			slog.Error("Failed to find process", "error", err)
			return err
		}
		// Send SIGINT to the process
		// Send multiple SIGINT signals if ffmpeg is concurrently doing both HLS and DASH
		for i := 0; i < 3; i++ {
			slog.Info("Sending SIGINT to process", "pid", sp.Pid, "attempt", i+1)
			process.Signal(syscall.SIGINT)
			time.Sleep(1 * time.Second)
		}
		err = process.Kill()
		if err != nil {
			if err.Error() == "os: process already finished" {
				slog.Error("Process not found or already terminated", "pid", sp.Pid, "error", err)
			} else {
				slog.Error("Failed to kill process", "error", err)
				return err
			}
		}
		slog.Info("Process killed", "pid", sp.Pid)
	}
	return nil
}

func (sp *StreamingProcess) MonitorDirectory() error {
	// Implement the monitoring logic here
	// This could involve checking the output directory for new files,
	// updating the job state, etc.
	slog.Info("Monitoring directory", "directory", sp.OutDir)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					// Ignore the .tmp files
					if strings.HasSuffix(event.Name, ".m4s") || strings.HasSuffix(event.Name, ".ts") {
						slog.Info("New media file created", "file", event.Name)
						sp.lastSegmentCreatedAt = time.Now().Unix()
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				slog.Error("Watcher error", "error", err)
			case <-sp.monitoringChannel:
				slog.Info("Stopping directory monitoring", "job", sp.Job.ID)
				// Close the watcher when the monitoring channel is closed
				if err := watcher.Close(); err != nil {
					slog.Error("Failed to close watcher", "error", err)
				}
				defer close(sp.monitoringChannel)
				return
			}
		}
	}()
	err = watcher.Add(sp.OutDir)
	if err != nil {
		return err
	}
	return nil
}
