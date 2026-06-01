---
name: deploy-personal-manager-to-robot43
description: Deploy this personal-manager repository's current committed branch HEAD to robot-43:/tmp/personal-manager over SSH after validation passes. Use only in this repository when the user asks to deploy personal-manager, deploy the current branch to robot-43, deploy to /tmp, or follow the repository's final delivery flow.
---

# Deploy Personal Manager To Robot43

## Scope

Use this skill only for the `personal-manager` repository. Do not use it for `transferhk` or any other project.

The deployment target is fixed:

```text
robot-43:/tmp/personal-manager
```

## Workflow

1. Work from the repository root.
2. Ensure the intended changes are committed. Untracked files may exist, but tracked changes must be clean because deployment uses committed `HEAD`.
3. Run the bundled deployment script:

```sh
.codex/skills/deploy-personal-manager-to-robot43/scripts/deploy_to_robot43_tmp.sh
```

The script runs:
- `git diff --check`
- `go test ./...`
- `git archive HEAD | ssh ... tar -x -C /tmp/personal-manager`
- remote `go test ./...` with `/usr/local/go/bin` added to `PATH`

It deploys only tracked, committed files from `HEAD`; it does not copy `.git`, local SQLite data, untracked files, or workspace artifacts.

## Reporting

Report:
- local validation results
- branch and commit SHA deployed
- deployment path
- whether remote `go test ./...` ran
- if remote tests did not run, the missing remote dependency

If SSH access is unavailable, stop before implying deployment succeeded and report the blocker.
