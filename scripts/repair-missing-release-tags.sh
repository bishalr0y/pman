#!/usr/bin/env bash
# One-time repair when release-please manifest advanced without GitHub tags.
# Safe to re-run: skips tags that already exist on origin.
set -euo pipefail

repo="${GITHUB_REPOSITORY:-bishalr0y/pman}"
root="$(git rev-parse --show-toplevel)"
cd "$root"

create_if_missing() {
  local tag=$1
  local sha=$2
  if git ls-remote --exit-code --tags origin "refs/tags/${tag}" &>/dev/null; then
    echo "ok ${tag} already exists on origin"
    return 0
  fi
  local full_sha
  full_sha="$(git rev-parse "${sha}^{commit}")"
  echo "creating ${tag} -> ${full_sha}"
  gh api -X POST "repos/${repo}/git/refs" -f ref="refs/tags/${tag}" -f sha="${full_sha}"
}

create_if_missing "v3.0.0" 0a8dfed
create_if_missing "v4.0.0" fdf01da