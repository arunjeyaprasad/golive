[![Go](https://github.com/arunjeyaprasad/golive/actions/workflows/go.yml/badge.svg)](https://github.com/arunjeyaprasad/golive/actions/workflows/go.yml)

# golive
A comprehensive testing tool for generating customizable HLS and DASH live streams to validate enterprise streaming applications, media players, and CDN infrastructure under realistic conditions written using Go leveraging ffmpeg.

# Overview
This tool provides a robust, configurable live streaming test environment that generates both HLS (HTTP Live Streaming) and DASH (Dynamic Adaptive Streaming over HTTP) streams with customizable parameters. It's designed specifically for enterprise teams who need to test their streaming infrastructure, applications, and players against various streaming scenarios without relying on external content sources.

# Key Features
<ul>
<li><b>Live Stream Simulation:</b> Real-time stream generation with configurable segment durations and playlist updates
<li><b>Multi-Protocol Support:</b> Generate simultaneous HLS and DASH streams from the same source
<li><b>RESTful API:</b> Comprehensive API for programmatic stream creation, modification, and monitoring
<li><b>Integration Ready:</b> Docker containerization, Kubernetes helm charts, and CI/CD pipeline integration
</ul>


# Quick Start
### Prerequisites
<ul>
<li>Docker (Only, If you want to run as a container)
<li>FFmpeg
<li>Go
</ul>

The service can be run standalone or within Docker container

### Direct Run
Install golang if you don't have it already. Clone this repo and navigate to this directory and run
```go run main.go``` or ```make run```
from your favourite terminal

### Running from Docker
Build the Docker container using
```
docker build -t golive .
docker run -d -p 9090:9090 golive
```
or use the Makefile
```
make dockerbuild
make dockerrun
```

# REST API
### Create Stream
```
http

POST http://localhost:9090/jobs
```
Request Body:
```json
{
        "description": "Ajz Test Live Stream",
        "video": {
            "bitrate": "1.2M",
            "resolution": "1280x720",
            "framerate": "30",
            "codec": "h264"
        },
        "audio": {
            "codec": "aac",
            "bitrate": "192k",
            "sample_rate": "44100",
            "channels": "2"
        },
        "segment_length": 6,
        "window_size": 6
}
```

Response
```json
{
    "id": "41877717-01cc-47a7-a960-efd73fdb0d2f",
    "status": "created",
    "created": "2025-06-07T20:11:05+05:30",
    "config": {
        "description": "Ajz Test Live Stream",
        "video": {
            "bitrate": "1.2M",
            "resolution": "1280x720",
            "framerate": "30",
            "codec": "h264"
        },
        "audio": {
            "codec": "aac",
            "bitrate": "192k",
            "sample_rate": "44100",
            "channels": "2"
        },
        "segment_length": 6,
        "window_size": 6
    }
}
```

## Get Stream
```
http
GET http://localhost:9090/jobs/{{job_id}}
```

Response
```json
{
    "id": "f6deb708-eb18-4c0a-8a75-b414bb41f63a",
    "status": "running",
    "created": "2025-06-07T20:33:00+05:30",
    "streamed_from": "2025-06-07T20:33:02+05:30",
    "playback_urls": [
        {
            "format": "dash",
            "url": "http://localhost:9090/jobs/f6deb708-eb18-4c0a-8a75-b414bb41f63a/manifest.mpd"
        },
        {
            "format": "hls",
            "url": "http://localhost:9090/jobs/f6deb708-eb18-4c0a-8a75-b414bb41f63a/master.m3u8"
        }
    ],
    "config": {
        "description": "Ajz Test Live Stream",
        "video": {
            "bitrate": "1.2M",
            "resolution": "1280x720",
            "framerate": "30",
            "codec": "h264"
        },
        "audio": {
            "codec": "aac",
            "bitrate": "192k",
            "sample_rate": "44100",
            "channels": "2"
        },
        "segment_length": 6,
        "window_size": 6
    }
}
```
Note: When the Job is in `running` state the playback URLs are also returned.

You can use Safari browser to natively play the HLS streams. Alternatively use ffplay or VLC app to play the HLS/DASH URLs

# ScreenShot
<img src="./assets/output.gif" width="400" alt="Demo"/>

# License
This project is licensed under the MIT License
