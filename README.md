# Distributed Job Queue in Go

A production-style distributed job processing system built in Go.

The project progressively introduces HTTP APIs, PostgreSQL persistence, concurrent workers, distributed coordination, retries, scheduling, observability, and Kubernetes deployment.

## Tech Stack

* Go
* PostgreSQL
* Podman / Podman Compose
* Kubernetes
* Prometheus
* OpenTelemetry

## Architecture

```text
                         ┌──────────────┐
                         │    Client    │
                         └──────┬───────┘
                                │ HTTP
                                ▼
                         ┌──────────────┐
                         │  API Server  │
                         └──────┬───────┘
                                │
                                ▼
                         ┌──────────────┐
                         │  PostgreSQL  │
                         │              │
                         │   Job State  │
                         │   Attempts   │
                         │  Retry Time  │
                         └──────┬───────┘
                                │
                         Claim Eligible Jobs
                                │
                                ▼
                  ┌──────────────────────────┐
                  │         Workers          │
                  │  ┌────┐ ┌────┐ ┌────┐  │
                  │  │ W1 │ │ W2 │ │ W3 │  │
                  │  └────┘ └────┘ └────┘  │
                  └───────────┬──────────────┘
                              │
                    ┌─────────┴─────────┐
                    │                   │
                 Success             Failure
                    │                   │
                    ▼                   ▼
               completed          Retry Policy
                                      │
                              ┌───────┴────────┐
                              │                │
                         Attempts < max   Attempts >= max
                              │                │
                              ▼                ▼
                           pending           failed
                              │
                         Backoff delay
                              │
                              ▼
                         PostgreSQL
```

## Job Lifecycle

```text
pending
   │
   ▼
processing
   │
   ├── success ──────────────► completed
   │
   └── failure
        │
        ├── attempts < max ──► pending
        │                         │
        │                    backoff delay
        │                         │
        │                         ▼
        │                    processing
        │
        └── attempts >= max ──► failed
```

Jobs are retried when processing fails and the maximum number of attempts has not been reached.

Retry delays use exponential backoff and are persisted through the job's retry state.

## Project Structure

```text
distributed-job-queue/

├── cmd/
│   ├── api/
│   │   └── main.go
│   └── worker/
│       └── main.go
│
├── internal/
│   ├── api/
│   ├── config/
│   ├── database/
│   ├── integration/
│   │   └── worker_test.go
│   └── job/
│       ├── default_handler.go
│       ├── handler.go
│       ├── memory_repository.go
│       ├── model.go
│       ├── postgres_repository.go
│       ├── postgres_repository_test.go
│       ├── repository.go
│       ├── retry_policy.go
│       ├── worker.go
│       └── worker_test.go
│
├── migrations/
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── go.mod
├── go.sum
├── README.md
└── .gitignore
```

## Retry Configuration

The worker supports configurable retry behavior.

Example defaults:

```text
Max Attempts:    3
Initial Backoff: 1s
Max Backoff:     30s
```

Exponential backoff:

```text
Attempt 1 → 1s
Attempt 2 → 2s
Attempt 3 → 4s
...
```

A job is permanently marked as `failed` once the maximum number of attempts has been exhausted.

## Build Project

```bash
make check

make test

make build
```

## Run Project

```bash
make run-api

make run-worker
```
