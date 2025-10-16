# CloudForge

A learning path for building production-grade platform engineering automation. This repository contains all the code, infrastructure, and documentation for building and managing a cloud platform from the ground up.

## Project Status
**Phase 0: Foundation Setup (In Progress)**

## Quick Start

1.  **Prerequisites:** Ensure Go, Python, Terraform, Docker, and the AWS CLI are installed.
2.  **Start Local Environment:** `docker-compose up -d`
3.  **Run CLI:** `go run cmd/cloudforge/main.go config.yaml`
4.  **Build Binary:** `go build -o bin/cloudforge cmd/cloudforge/main.go`
5.  **Execute Go CLI Binary:** `./bin/cloudforge config.yaml`

## Repository Structure

-   `cmd/`: Main applications (CLIs) for the project.
-   `internal/`: Private Go packages used by our applications.
-   `terraform/`: All Terraform modules and environment configurations.
-   `scripts/`: Automation and helper scripts (Python, Bash).
-   `tests/`: Unit and integration tests for Go, Python, and Terraform.
-   `docs/`: All project documentation, including ADRs and runbooks.