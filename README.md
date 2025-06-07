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

# ScreenShot
<img src="./assets/output.gif" width="400" alt="Demo"/>

