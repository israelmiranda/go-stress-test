# Go Stress Test

A CLI tool developed in Go for load testing web services. The application allows you to test a service's ability to handle multiple simultaneous requests, providing a detailed report at the end of execution.

## Features

- Configurable HTTP load testing
- Concurrency control (multiple simultaneous requests)
- Detailed report with statistics
- HTTP status code distribution
- Docker execution support
- HTTPS support

## Quick Start

### Build Docker Image

To build the Docker image:

```bash
docker build -t stress-test:latest .
```

Add a version tag (optional):

```bash
docker build -t stress-test:v1.0 .
```

## Usage

```bash
docker run stress-test:latest --url http://google.com --requests 1000 --concurrency 10
```

Or using short flags:

```bash
docker run stress-test:latest -u http://google.com -r 1000 -c 10
```

### Parameters

- `--url` or `-u` (required): URL of the service to test
  - Example: `http://google.com` or `https://api.example.com`

- `--requests` or `-r` (required): Total number of requests to be made
  - Example: `1000` (makes 1000 requests)

- `--concurrency` or `-c` (optional): Number of simultaneous calls
  - Default: `1`
  - Example: `10` (10 requests in parallel)

### Execution Examples

#### Simple test (1 request at a time)

```bash
docker run stress-test:latest -u http://localhost:8080 -r 100
```

#### Test with high concurrency

```bash
docker run stress-test:latest -u https://api.github.com -r 5000 -c 50
```

## Output Report

At the end of execution, the application displays a report with the following information:

```
=== Load Test Report ===
Total Time: 5.234s
Total Requests: 1000
Successful (HTTP 200): 950
Success Rate: 95.00%
Requests per second: 191.05

Status Code Distribution:
  HTTP 200: 950
  HTTP 503: 50
========================
```

### Provided Metrics

- **Total Time**: Total time spent executing the test
- **Total Requests**: Total number of requests made
- **Successful (HTTP 200)**: Number of 200 status responses
- **Success Rate**: Success percentage
- **Requests per second**: Request rate per second
- **Status Code Distribution**: Complete distribution of HTTP status codes returned

## Project Structure

```
go-stress-test/
├── main.go           # Main application code
├── go.mod            # Go module definition
├── Dockerfile        # Docker image build configuration
├── README.md         # This file
└── cmd/
    └── root.go       # CLI command definitions
```
