## Distributed Job Queue in Go

A production-style distributed job processing system built in Go.

The project progressively introduces HTTP APIs, PostgreSQL persistence, concurrent workers, distributed coordination, retries, scheduling, observability, and Kubernetes deployment.

Tech Stack
Go
PostgreSQL
Podman / Podman Compose
Kubernetes
Prometheus
OpenTelemetry
Architecture

                         ┌──────────────┐
                         │    Client    │
                         └──────┬───────┘
                                │ HTTP
                                ▼
                         ┌──────────────┐
                         │   API Server │
                         └──────┬───────┘
                                │
                                ▼
                         ┌──────────────┐
                         │  PostgreSQL  │
                         └──────┬───────┘
                                │
                                ▼
                    ┌──────────────────────┐
                    │       Workers        │
                    │  ┌────┐ ┌────┐ ┌────┐│
                    │  │ W1 │ │ W2 │ │ W3 ││
                    │  └────┘ └────┘ └────┘│
                    └──────────────────────┘

## Project Structure

```text
distributed-job-queue/
├── cmd/
│   ├── api/
│   │   └── main.go
│   └── worker/
│       └── main.go
├── internal/
│   ├── api/
│   ├── config/
│   ├── job/
│   └── worker/
├── migrations/
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── README.md
└── .gitignore
