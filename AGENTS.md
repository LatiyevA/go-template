# AGENTS.md

This repository is a template for backend services written in Go.

## Goals

- Provide a fast way to bootstrap new services with a consistent structure.
- Enforce unified rules for formatting, linting, and testing.
- Minimize manual setup for Docker, PostgreSQL, and CI.

## Project Structure

- `cmd/app` — application entry point.
- `internal/...` — business logic, HTTP handlers, data access.
- `configs/` — configuration files.
- `migrations/` — database migrations.
- `.github/workflows/` — CI pipelines (deploy).

## Commands

Before committing changes, prefer to run:

- `golangci-lint run ./...` — run the linter (golangci-lint).
- `go test -v -race -cover ./...` — run tests.
- `go run ./cmd/app` — start the application locally.

## Guidelines for agents

### Do

- Preserve the existing structure (`cmd/`, `internal/`, `migrations/`).
- Update or add tests when changing business logic.
- Follow Go idioms and respect linter feedback where reasonable.

### Do not

- Do not change the module name in `go.mod` unless the user explicitly requests it.
- Do not introduce heavy dependencies without a clear benefit.
- Do not rewrite the entire project structure unless explicitly requested.
- Do not disable linters globally; if needed, disable them locally and with justification.

## When in doubt

- Propose a short plan before performing large refactors.
- For database changes, propose a new migration file and clearly describe what it does.

## Human vs agent docs

- `README.md` — instructions for humans.
- `AGENTS.md` — instructions for AI agents interacting with this repository.
