#!/usr/bin/env bash

set -euo pipefail

name="${1:-}"
if [[ -z "$name" ]]; then
	echo 'Usage: ./lc build two_sum' >&2
	exit 1
fi

if [[ ! "$name" =~ ^[a-zA-Z0-9_]+$ ]]; then
	echo 'NAME must contain only letters, numbers, and underscores.' >&2
	exit 1
fi

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_dir="$(cd "$script_dir/.." && pwd)"
exercise_dir="$repo_dir/exercises/$name"
binary_dir="$repo_dir/.bin"
binary_path="$binary_dir/$name"

if [[ ! -f "$exercise_dir/main.go" ]]; then
	echo "Exercise does not exist: exercises/$name" >&2
	exit 1
fi

needs_build=false
if [[ ! -x "$binary_path" ]]; then
	needs_build=true
else
	while IFS= read -r source_file; do
		if [[ "$source_file" -nt "$binary_path" ]]; then
			needs_build=true
			break
		fi
	done < <(find "$exercise_dir" "$repo_dir/internal/runner" -type f -name '*.go')

	for module_file in "$repo_dir/go.mod" "$repo_dir/go.sum"; do
		if [[ -f "$module_file" && "$module_file" -nt "$binary_path" ]]; then
			needs_build=true
		fi
	done
fi

if [[ "$needs_build" == true ]]; then
	mkdir -p "$binary_dir"
	(
		cd "$repo_dir"
		go build -o "$binary_path" "./exercises/$name"
	)
	echo "Built: .bin/$name"
else
	echo "Up to date: .bin/$name"
fi
