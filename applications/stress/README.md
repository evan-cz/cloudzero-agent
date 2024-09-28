# Cirrus Stress App

The `cirrus-stress-app` is a tool designed to simulate CPU stress on Kubernetes pods. It allows you to specify the desired CPU usage, memory usage, and sleep interval, helping you test the behavior of your Kubernetes cluster under different load scenarios.

## Features

- Simulates CPU and memory stress on Kubernetes pods.
- Configurable CPU usage (in millicores), memory usage (in MB), and sleep interval.
- Helps evaluate Kubernetes cluster performance and resource allocation.

## Getting Started

### Prerequisites

To use the `cirrus-stress-app`, you need to have the following installed:

- Go (version 1.19 or higher)
- Docker with BuildX

### Building the Binary

1. Clone this repository to your local machine.
2. Navigate to the repository directory.
3. Build the binary using the following command:

  ```shell
  make build
  ```

## Running with Docker

Build the Docker image using the provided Dockerfile:

  ```shell
  make docker-build
  ```

Run the Docker container, specifying the desired CPU usage, memory usage, and sleep interval:

```shell
docker run --rm ghcr.io/cloudzero/cirrus-stress-app -cpu=0.2 -mem=100 -sleep=100
```
> This will use 20% of 1 vCPU, allocate 100MB, and sleep for 100ms inbetween CPU consumption cycles

## Parameters

The `cirrus-stress-app` allows you to specify the desired CPU usage, memory usage, and sleep interval using the following parameters:

- **CPU Usage**: The CPU usage as a fraction of vCPU milliCPUs. It is specified using the `-cpu` argument. For example, `-cpu=0.2` represents a CPU usage of 20% or 200 milliCPU (mCPU).

- **Memory Usage**: The memory usage is specified in megabytes (MB) using the `-mem` argument. For example, `-mem=100` represents a memory usage of 100 MB.

- **Sleep Interval**: The sleep interval defines the duration to sleep between CPU stress cycles. It is specified using the `-sleep` argument. Duration is in milliseconds.

Adjust these parameters according to your requirements to simulate different load scenarios.


## Kubernetes deployment

A [Kuberentes deployment template](../../deployments/stress/deployment.yaml) is available, tailor it to your needs.

## Make Targets

Use the following make targets to build and manage the project:

```shell
make help
```

Available targets:

Usage:
  make <target>

Targets:
  help             Show this help message
  docker-build     Build and push the Docker image
  build            Build the stressor locally
  fmt              Run go fmt against code
  lint             Run the linter 
  clean            Clean the build artifacts

## Contributing

Contributions are welcome! If you find a bug or have a suggestion, please open an issue or submit a pull request. For major changes, please discuss them first in the issue tracker.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
