# SKILLS.md

This file describes common skills for working with this Go template repository.

---

## Skill: setup-new-service

**Description:**  
Initialize a new service created from this template. Adjust the module name, update imports, and verify that the project builds and tests pass.

**When to use:**

- A new repository has been created from this template.
- The user wants a first working setup with the correct module path.

**Steps:**

1. Update the `module` directive in `go.mod` to the value provided by the user.
2. Replace all occurrences of `github.com/LatiyevA/go-template` with the new module path.
3. Run:
    - `go mod tidy`
    - `golangci-lint run ./...`
    - `go test -v -race -cover ./...`
4. Optionally run `go run ./cmd/app` to verify the service starts.
5. Report any errors and propose concrete fixes.

---

## Skill: modify-business-logic

**Description:**  
Change or extend business logic while keeping the existing structure and tests in sync.

**When to use:**

- The user asks to change API behavior, add a new endpoint, or adjust business rules.
- The changes are focused on code under `internal/` and optionally `cmd/app`.

**Steps:**

1. Locate the relevant handler/service in `internal/...` based on the user’s description.
2. Implement the minimal change needed to satisfy the request.
3. Update existing tests or add new tests for the changed behavior.
4. Run `go test -v -race -cover ./...` and ensure all tests pass.
5. If the change affects database shape or queries, coordinate with the `work-with-migrations` skill.

---

## Skill: work-with-migrations

**Description:**  
Introduce or adjust database schema changes via migration files under `migrations/`.

**When to use:**

- The user wants to add or modify tables, columns, indexes, or constraints.
- Schema changes are required for new features.

**Steps:**

1. Create a new migration file in `migrations/` with a descriptive name (for example, `20260218_add_users_table.sql`).
2. Define the necessary SQL for applying the change (and optionally for rolling it back, depending on the migration tool).
3. Clearly describe what the migration does (tables/columns added, removed, or changed).
4. Suggest a command to apply migrations, for example:
    - `goose -dir migrations postgres "$(DATABASE_URL)" up`
    - or the appropriate CLI for the chosen migration tool.
5. If needed, update code in `internal/...` to use the new schema and ensure tests cover the changes.

---

## Skill: maintenance-and-tooling

**Description:**  
Work with tooling and project hygiene: linting, formatting, and keeping configuration consistent.

**When to use:**

- The user asks to fix linter issues, clean up code style, or adjust tooling config.
- There are golangci-lint, formatting, or minor refactor tasks.

**Steps:**

1. Use `gofumpt` to format and organize Go code.
2. Respect the existing `.golangci.yml` configuration; only change it when explicitly requested or clearly justified.
3. Run `golangci-lint run ./...` and fix reported issues where reasonable.
4. Avoid disabling linters globally; prefer local `//nolint` with a short comment when an exception is necessary.
5. Keep changes focused and small, so they are easy to review.

---
