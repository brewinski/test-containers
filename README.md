# Test Containers Demo Repo

Demo with basic run and test information

## What is test containers

Test containers is a project that makes starting and stopping containers with an application or tests simple.

See the website and documentation for more information and examples
  - Repo: https://github.com/testcontainers
  - Website: https://testcontainers.com/

## Development

### Software Requirements
  - Golang v1.21
  - Docker (recent version)
    - (Note you can use rancher or podman) 
    - https://golang.testcontainers.org/system_requirements/using_podman/
    - https://golang.testcontainers.org/system_requirements/rancher/

### Start the project
   - Clone this repo
   - Open in your IDE or text editor of choice
   - Run `go mod tidy` to download dependencies
   - Run `go run cmd/dev.go` to start the app
   - When the app starts you can see that it creates a group of containers.
   - Stop the app with `ctrl + c` and notice the containers are stopped before exit.

## Tests
The customer folder has an example of unit tests that start and stop containers.

### Running Tests
  - Run `go test ./... -v` from the root folder

## Files where example code can be found
- cmd/dev.go - Example app that starts containers
- customer/repo_test.go - Example test that uses TestContainer
- testhelpers/containers.go - Examples of helper functions that define specific containersx