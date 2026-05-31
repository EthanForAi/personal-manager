#!/usr/bin/env bash
set -euo pipefail

remote="${REMOTE:-robot-43}"
remote_dir="/tmp/personal-manager"

repo_root="$(git rev-parse --show-toplevel)"
cd "$repo_root"

if [[ ! -f go.mod ]] || ! grep -qx "module personal-manager" go.mod; then
  echo "ERROR: this deployment skill must be run from the personal-manager repository." >&2
  exit 1
fi

if ! git diff --quiet --; then
  echo "ERROR: tracked working tree changes are present. Commit them before deploying HEAD." >&2
  git status --short
  exit 1
fi

if ! git diff --cached --quiet --; then
  echo "ERROR: staged changes are present. Commit them before deploying HEAD." >&2
  git status --short
  exit 1
fi

echo "Running git diff --check"
git diff --check

if [[ "${SKIP_LOCAL_TESTS:-0}" != "1" ]]; then
  echo "Running go test ./..."
  go test ./...
else
  echo "Skipping local go test ./... because SKIP_LOCAL_TESTS=1"
fi

branch="$(git branch --show-current)"
sha="$(git rev-parse --short=12 HEAD)"
printf -v quoted_branch "%q" "$branch"
printf -v quoted_sha "%q" "$sha"

remote_script="
set -e
rm -rf $remote_dir
mkdir -p $remote_dir
tar -x -C $remote_dir
cd $remote_dir
echo DEPLOYED_PATH=\$(pwd)
echo DEPLOYED_BRANCH=$quoted_branch
echo DEPLOYED_SHA=$quoted_sha
if command -v go >/dev/null 2>&1; then
  echo Running remote go test ./...
  go test ./...
else
  echo REMOTE_GO_MISSING
fi
"

echo "Deploying $branch@$sha to $remote:$remote_dir"
git archive --format=tar HEAD | ssh -o RemoteCommand=none -o RequestTTY=no "$remote" "$remote_script"
