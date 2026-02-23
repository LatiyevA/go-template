# AGENTS.md

This file provides guidance for AI agents interacting with code in this repository.

## Commands

- `task vet` — run `go vet` for static analysis.
- `task lint` — run linters.
- `task fmt` — format the codebase.
- `task tidy` — run `go mod tidy` to clean up dependencies.
- `task update` — update all dependencies using `go get -u ./...`.
- `task test` — run tests with coverage.
- `task build` — build the application.

**Do not run:** `task migrate`.

## Architecture

### Stack Overview

- Go 1.25.4+ for backend development.
- PostgreSQL for database.
- Docker for containerization.

### Directory Structure

TODO: Add comments about the directory structure.

```bash
├───.github
│   └───workflows
├───cmd
│   └───app
├───configs
├───internal
│   ├───config
│   ├───delivery
│   │   └───http
│   │       └───v1
│   ├───dto
│   ├───entity
│   ├───metrics
│   ├───repo
│   └───usecase
├───migrations
└───pkg
    ├───gorm
    │   └───postgres
    ├───httpserver
    └───validator
```

### Key Patterns

### Code Style

### Testing

### Important Rules

- Always run `task vet`, `task lint`, and `task fmt` after completing a task to ensure no formatting or linting issues are introduced.
- After introducing new dependencies run `task tidy` and `task deps-update` to clean up dependencies.
- Follow Go idioms and best practices.

### Additional information