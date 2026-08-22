#!/usr/bin/env bash
# Run one e2e-m*.sh with optional retries (E2E_RETRIES, default 1).
set -euo pipefail

script="${1:?usage: e2e-run.sh scripts/e2e-mN.sh}"
retries="${E2E_RETRIES:-1}"
if (( retries < 1 )); then
	retries=1
fi

attempt=1
while (( attempt <= retries )); do
	if bash "$script"; then
		exit 0
	fi
	if (( attempt < retries )); then
		echo "e2e-run: $(basename "$script") failed (attempt ${attempt}/${retries}), retrying..." >&2
		sleep "${E2E_SCRIPT_RETRY_DELAY_SECS:-2}"
	fi
	attempt=$((attempt + 1))
done

echo "e2e-run: $(basename "$script") failed after ${retries} attempts" >&2
exit 1
