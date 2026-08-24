#!/usr/bin/env bash

set -euo pipefail

name="${1:-}"
if [[ -z "$name" ]]; then
	echo 'Usage: ./lc new two_sum' >&2
	exit 1
fi

if [[ ! "$name" =~ ^[a-zA-Z0-9_]+$ ]]; then
	echo 'NAME must contain only letters, numbers, and underscores.' >&2
	exit 1
fi

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_dir="$(cd "$script_dir/.." && pwd)"
template_dir="$repo_dir/templates/exercise"
exercises_dir="${EXERCISES_DIR:-$repo_dir/exercises}"
target_dir="$exercises_dir/$name"

if [[ -e "$target_dir" ]]; then
	echo "Exercise already exists: $target_dir" >&2
	exit 1
fi

mkdir -p "$target_dir"
cp "$template_dir/main.go" "$target_dir/main.go"
cp "$template_dir/input.txt" "$target_dir/input.txt"

echo "Created exercise: $target_dir"
echo "Run it with: ./lc $name"
