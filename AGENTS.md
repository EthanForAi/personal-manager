# AGENTS.md

## Project Overview

This is a small Go web service that provides HTTP APIs for managing personal information.

The service supports common CRUD operations:
- create personal data
- read personal data
- update personal data
- delete personal data

All APIs use the HTTP POST method for consistency.

Data is persisted in local SQLite.

Default expectation:
Implement requested changes with minimal correct modifications, add or update meaningful tests, validate with Go test commands, fix failures, and only finish when everything passes and the code has been pushed to GitHub when a Git remote is available.

## Repository Context

Core layers:
- `internal/handler`: HTTP handlers
- `internal/service`: business logic
- `internal/store`: SQLite persistence
- `internal/model`: shared structs

Data flow:

```text
handler -> service -> store -> SQLite
```

Project goals:
- keep the code small and readable
- use SQLite for persistence
- support deterministic local tests
- maintain clear API behavior
- avoid unnecessary abstractions

## Non-Negotiable Rules

- ALWAYS write or update unit tests for any new code or behavior change.
- NEVER consider an implementation task complete until `go test ./...` passes.
- ALWAYS run `go test ./...` before finishing implementation work.
- If tests fail, fix the failures and rerun tests until all tests pass.
- NEVER claim tests passed unless they were actually executed.
- Do NOT write fake, trivial, or assertion-free tests.
- Do NOT modify unrelated files.
- Do NOT break existing APIs unless explicitly requested.
- Tasks are not complete until code is pushed to GitHub when the repository has a configured remote.
- If Git or GitHub access is unavailable, clearly report that blocker instead of implying the push happened.

## Coding Standards

- Follow idiomatic Go.
- Keep code simple and readable.
- Keep functions small and focused.
- Avoid unnecessary dependencies.
- Keep HTTP handlers thin; move business logic to service/store layers.
- Use clear JSON request and response bodies.
- Prefer explicit error handling.
- Return consistent error responses.
- Use context-aware database calls when practical.
- Keep SQLite schema changes small and covered by tests.

## API Expectations

All APIs use HTTP POST.

Personal data fields:
- `id`: server-generated integer identifier
- `name`: non-empty string
- `age`: non-negative integer
- `email`: non-empty string

Use JSON for all requests and responses.

### 1. Create

`POST /create`

Request:

```json
{
  "name": "Alice",
  "age": 18,
  "email": "alice@example.com"
}
```

Success response:

```json
{
  "id": 1,
  "name": "Alice",
  "age": 18,
  "email": "alice@example.com"
}
```

### 2. Read

`POST /read`

Request:

```json
{
  "id": 1
}
```

Success response:

```json
{
  "id": 1,
  "name": "Alice",
  "age": 18,
  "email": "alice@example.com"
}
```

### 3. Update

`POST /update`

Request:

```json
{
  "id": 1,
  "name": "Alice Smith",
  "age": 19,
  "email": "alice.smith@example.com"
}
```

Success response:

```json
{
  "id": 1,
  "name": "Alice Smith",
  "age": 19,
  "email": "alice.smith@example.com"
}
```

### 4. Delete

`POST /delete`

Request:

```json
{
  "id": 1
}
```

Success response:

```json
{
  "deleted": true
}
```

## Error Handling

Return JSON errors using a consistent shape:

```json
{
  "error": "message"
}
```

Expected HTTP behavior:
- invalid JSON: `400 Bad Request`
- validation failure: `400 Bad Request`
- missing record: `404 Not Found`
- unexpected persistence or server failure: `500 Internal Server Error`

Do not leak internal database details in API responses.

## Testing Expectations

Tests should cover behavior, not implementation details.

Preferred coverage:
- handler tests for request parsing, response status, and JSON body shape
- service tests for validation and business rules
- store tests for SQLite persistence behavior
- integration-style tests when behavior crosses layers

SQLite tests should be deterministic:
- use temporary database files or in-memory databases
- isolate test data between tests
- avoid relying on test execution order

Before finishing implementation work, run:

```sh
go test ./...
```

If formatting may have changed, run:

```sh
gofmt -w <changed-go-files>
```

## Branch Strategy

Use `main` as the stable trunk branch.

Do not develop features or bug fixes directly on `main` unless the user explicitly requests it.

For every new feature or bug fix:

1. Start from the latest `main`.
2. Create a dedicated branch named `codex/<task-name>`.
3. Keep the branch focused on one feature, fix, or small related change set.
4. Implement the change and update meaningful tests.
5. Run `go test ./...` and fix failures before finishing.
6. Push the branch to GitHub.
7. Open a pull request back into `main` when the change is ready for review.

Branch naming examples:
- `codex/add-person-tags`
- `codex/fix-email-validation`
- `codex/sqlite-schema`

If a branch already exists for the current task, continue using that branch instead of creating a duplicate. If the worktree is already on `main`, create a new task branch before making code changes.

## Completion Output

After each completed task, include the important commands the user may need next.

For this service, include relevant commands such as:
- start the server
- stop the server
- create personal data
- read personal data
- update personal data
- delete personal data
- run tests

Use concrete commands with the actual host, port, file paths, and JSON fields whenever they are known. If a command is not relevant to the completed task, omit it rather than adding noise.

## Change Discipline

- Inspect existing code before editing.
- Preserve existing route names, request fields, and response fields unless the task explicitly asks for an API change.
- Prefer small commits with clear messages.
- Keep documentation updates aligned with actual behavior.
- If a requirement is ambiguous, choose the smallest behavior that fits the existing project style and document any assumption in the final response.
