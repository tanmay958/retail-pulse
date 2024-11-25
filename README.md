# Retail Pulse

This backend service is built in Go language to efficiently handle image processing requests for calculating the perimeter of images collected from retail stores. The service provides robust error handling and real-time job status tracking, ensuring reliability.

## Installation

Installation without Docker (No database,only - Map,struct etc)

```bash
  git clone https://github.com/tanmay958/retail-pulse.git
  cd retail-pulse
  cd '.\app\api gateway\'
  go mod tidy
  go run main.go
  #start testing after the server starts running
```

Installation with Docker (database-Mysql)

```bash
  git clone https://github.com/tanmay958/retail-pulse.git
  cd retail-pulse
  cd app-docker
  docker compose up --build
  #start testing when container is up and running
```

## Running Apis

To run tests, hit the following endpoint

```bash
 http://localhost:8080/api/submit #post request
 http://localhost:8080/api/status?jobid=123 #get request
```
