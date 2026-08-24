#!/usr/bin/env bash

set -euo pipefail

name="${1:-}"
if [[ -z "$name" ]]; then
	echo 'Usage: ./lc two_sum' >&2
	exit 1
fi
shift

if [[ ! "$name" =~ ^[a-zA-Z0-9_]+$ ]]; then
	echo 'NAME must contain only letters, numbers, and underscores.' >&2
	exit 1
fi

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_dir="$(cd "$script_dir/.." && pwd)"
exercise_dir="$repo_dir/exercises/$name"
binary_path="$repo_dir/.bin/$name"

if [[ ! -f "$exercise_dir/main.go" ]]; then
	echo "Exercise does not exist: exercises/$name" >&2
	exit 1
fi

bash "$script_dir/build-exercise.sh" "$name"

"$binary_path" "$@"
