# AGENTS.md

## Project Overview

This is a small Go web service that provides HTTP APIs for managing personal information.

This repository is also a hands-on practice project for deeply understanding and exercising Codex capabilities, including skills, MCP tools, GitHub workflow automation, testing, code review, and delivery discipline.

The service supports common CRUD operations:
- create personal data
- read personal data
- update personal data
- delete personal data

All APIs use the HTTP POST method for consistency.

Data is persisted in local SQLite.

Default expectation:
Implement requested changes with minimal correct modifications, add or update meaningful tests, validate with Go test commands, fix failures, and only finish when everything passes, the code has been pushed to GitHub when a Git remote is available, and the current branch has been deployed to `robot-43:/tmp/personal-manager` when SSH access is available.

## Skill and MCP Tooling

- Use available Codex skills and MCP tools when they directly support the task, especially for GitHub workflow, code review, test automation, and repository delivery.
- For deployment to `robot-43:/tmp/personal-manager`, use the repo-local skill at `.codex/skills/deploy-personal-manager-to-robot43/SKILL.md`; this skill is only for this repository.
- If a required skill or MCP tool is unavailable, continue with the best local equivalent and clearly report the missing tool or authentication blocker.
- Prefer repository-local commands and existing project tooling over adding new dependencies or services.
- When a task is useful for learning Codex behavior, briefly report which skills, MCP tools, or local fallbacks were used.

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
- use real development tasks to practice Codex capabilities end to end

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
- Keep exported names documented when they are part of the public package API.
- Use `context.Context` for request-scoped operations and database calls when practical.
- Return errors explicitly; do not panic for normal request, validation, or persistence failures.
- Wrap internal errors where useful for debugging, but keep HTTP error responses stable and free of database details.
- Keep code simple and readable.
- Keep functions small and focused.
- Avoid unnecessary dependencies.
- Keep HTTP handlers thin; move business logic to service/store layers.
- Use clear JSON request and response bodies.
- Prefer explicit error handling.
- Return consistent error responses.
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

- Every new feature, bug fix, or behavior change must include meaningful unit tests or integration-style tests.
- Do not submit code changes with failing, skipped, trivial, or assertion-free tests.
- Run the relevant focused tests while developing, then run the full suite before review:

```sh
go test ./...
```

Preferred coverage:
- handler tests for request parsing, response status, and JSON body shape
- service tests for validation and business rules
- store tests for SQLite persistence behavior
- integration-style tests when behavior crosses layers

SQLite tests should be deterministic:
- use temporary database files or in-memory databases
- isolate test data between tests
- avoid relying on test execution order

If formatting may have changed, run:

```sh
gofmt -w <changed-go-files>
```

## Branch Strategy

Use `main` as the stable trunk branch.

Do not develop features or bug fixes directly on `main` unless the user explicitly requests it.

For every new feature, bug fix, or behavior change:

1. Start from the latest `main`.
2. Create a new dedicated branch named `codex/<task-name>`.
3. Keep the branch focused on one feature, fix, or small related change set.
4. Implement the change and update meaningful tests.
5. Run `go test ./...` and fix failures before finishing.
6. Run a code review pass before opening a PR.
7. Fix any review findings, rerun the relevant tests, and repeat review if needed.
8. Commit the reviewed and tested changes.
9. Push the branch to GitHub.
10. Deploy the current branch to `robot-43:/tmp/personal-manager`.
11. Open or update a pull request back into `main` when the change is ready for review.

Branch naming examples:
- `codex/add-person-tags`
- `codex/fix-email-validation`
- `codex/sqlite-schema`

Use one branch per issue. If the worktree is already on `main`, create a new task branch before making code changes. Reuse an existing branch only when the user explicitly points to that branch for the same issue.

## GitHub Issue and PR Workflow

For every non-trivial feature, bug fix, or behavior change, create a new GitHub issue before implementation when GitHub access is available.

Issue requirements:
- Describe the problem or requested change clearly.
- Include concrete acceptance criteria.
- Link any relevant API behavior, validation rule, persistence change, or test expectation.
- Associate exactly one implementation branch with the issue.

PR requirements:
- Push the task branch to GitHub.
- Open a pull request from the task branch into `main`.
- Include `Fixes #<issue-number>` in the PR body so GitHub closes the issue automatically after the PR is merged.
- Summarize the implementation, code review result, and test command results, including `go test ./...`.

Do not manually close the linked issue immediately after pushing the branch. The issue should close automatically when the PR containing `Fixes #<issue-number>` is merged. Do not close the PR after pushing; leave it open for review and merge unless the user explicitly asks to abandon the PR.

Skip automatic issue creation only for trivial local-only changes, pure questions, exploratory analysis, or when the user explicitly asks not to create an issue. If GitHub access is unavailable, continue the local work and report the blocker clearly.

## Robot-43 Deployment

After local tests pass and the reviewed changes are committed, deploy the current branch's committed `HEAD` to `robot-43:/tmp/personal-manager` when SSH access is available.

Deployment requirements:
- Use the repo-local skill at `.codex/skills/deploy-personal-manager-to-robot43/SKILL.md`.
- Deploy only tracked, committed repository contents from the current branch.
- Do not include local untracked files, local SQLite data, `.git`, or unrelated workspace artifacts.
- Prefer the skill's bundled script: `.codex/skills/deploy-personal-manager-to-robot43/scripts/deploy_to_robot43_tmp.sh`.
- If the remote host has Go installed, run `go test ./...` in `/tmp/personal-manager` after deployment.
- If the remote host cannot run tests, verify that `/tmp/personal-manager` contains the expected branch contents and report the missing remote dependency clearly.

## Code Review Gate

Before opening a PR:
- Review the diff for correctness, API compatibility, error handling, tests, and unintended file changes.
- Prefer using an available code review skill or MCP tool when configured.
- Treat review findings as blockers until fixed or explicitly documented as non-blocking.
- Rerun `go test ./...` after fixing review findings that touch Go code.

## Final Delivery Notes

After completing a task, always report:
- The issue number and PR URL when created.
- The branch name.
- The exact validation commands that were run and whether they passed.
- The Codex skills, MCP tools, or local fallbacks used for the workflow.
- The `robot-43:/tmp/personal-manager` deployment result, including whether remote tests were run or why they could not run.
- How to start the service locally.
- How to test the changed behavior manually, including example HTTP requests when relevant.

## Change Discipline

- Inspect existing code before editing.
- Preserve existing route names, request fields, and response fields unless the task explicitly asks for an API change.
- Prefer small commits with clear messages.
- Keep documentation updates aligned with actual behavior.
- If a requirement is ambiguous, choose the smallest behavior that fits the existing project style and document any assumption in the final response.
