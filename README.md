[![Go](https://github.com/arunjeyaprasad/golive/actions/workflows/go.yml/badge.svg)](https://github.com/arunjeyaprasad/golive/actions/workflows/go.yml)

# golive
Simple Go service to generate Live HLS and DASH Streams

# Usage
The service can be run standalone or within Docker container

## Direct Run
Install golang if you don't have it already. Clone this repo and navigate to this directory and run
```go run main.go``` from your favourite terminal

## Docker Run
Build the Docker container using
```
docker build -t golive .
docker run -d -p 9090:9090 golive
```
