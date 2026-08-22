#!/usr/bin/env bash
# Fail if any internal package lacks a package comment with >= 3 non-empty lines.
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
min_lines=3
failed=0

count_comment_lines() {
	local file=$1
	local count=0
	local line
	while IFS= read -r line || [[ -n $line ]]; do
		if [[ $line =~ ^//[[:space:]]*$ ]]; then
			continue
		fi
		if [[ $line =~ ^// ]]; then
			count=$((count + 1))
			continue
		fi
		break
	done <"$file"
	echo "$count"
}

find_comment_file() {
	local dir=$1
	local name=$2
	if [[ -f $dir/$name.go ]]; then
		echo "$dir/$name.go"
		return
	fi
	local f
	for f in "$dir"/*.go; do
		[[ -f $f ]] || continue
		[[ $f == *_test.go ]] && continue
		if head -1 "$f" | grep -q '^// Package '; then
			echo "$f"
			return
		fi
	done
}

check_pkg() {
	local dir=$1
	local name
	name=$(basename "$dir")
	local file
	file=$(find_comment_file "$dir" "$name") || true
	if [[ -z ${file:-} ]]; then
		echo "intent-check: no package comment in $dir" >&2
		failed=1
		return
	fi
	local n
	n=$(count_comment_lines "$file")
	if [[ $n -lt $min_lines ]]; then
		echo "intent-check: $file has $n comment line(s), need >= $min_lines" >&2
		failed=1
	fi
}

while IFS= read -r dir; do
	check_pkg "$dir"
done < <(find "$root/internal" -mindepth 1 -maxdepth 1 -type d | sort)

check_pkg "$root/internal/dispatch/targets"

if [[ $failed -ne 0 ]]; then
	exit 1
fi
echo "intent-check: ok"
