#!/usr/bin/env bash
# Forbidden terminology in user docs (docs/en + docs/ru).
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$root"

python3 <<'PY'
import re
import sys
from pathlib import Path

root = Path(".")
fail = False

INTERNAL = re.compile(
    r"CapAsk|Runner\.Decide|HookSpec|Hookedge|\bslog\b|fsnotify|e2e-m|M12|\bfive features\b",
    re.I,
)

# "daemon restart" as a command phrase (not "when the daemon process restarts")
DAEMON_RESTART_CMD = re.compile(r"daemon restart", re.I)

RU_LATIN = re.compile(
    r"(?<![`a-zA-Z])(offline|ledger|summary|checkpoint|sidecar|fallback|watcher|rollup)(?![`a-zA-Z])",
    re.I,
)

RU_DAEMON_LATIN = re.compile(
    r"(?<![`a-zA-Z/])daemon(?![`a-zA-Z])",
    re.I,
)

RU_UNITS = re.compile(r"(?<![`a-zA-Z0-9])[0-9]+(?:\.[0-9]+)?\s+(ms|s)(?![`a-zA-Z0-9])")

in_code = False


def strip_inline_code(s: str) -> str:
    return re.sub(r"`[^`]*`", "", s)


for doc in sorted(root.glob("docs/**/*.md")):
    in_code = False
    for i, line in enumerate(doc.read_text(encoding="utf-8").splitlines(), 1):
        if line.strip().startswith("```"):
            in_code = not in_code
            continue
        if in_code:
            continue

        plain = strip_inline_code(line)

        if INTERNAL.search(plain):
            print(f"docs-terms: internal ref {doc}:{i}: {line.strip()}", file=sys.stderr)
            fail = True
        if DAEMON_RESTART_CMD.search(plain):
            print(f"docs-terms: daemon restart {doc}:{i}: {line.strip()}", file=sys.stderr)
            fail = True

        if doc.parts[1] != "ru":
            continue

        # Skip CLI examples and mirrored EN section titles.
        if re.search(r"^\s*agentd\s+", line) or re.match(r"^#{1,6}\s+daemon\b", line, re.I):
            continue
        if re.match(r"^#{1,6}\s+config\b", line, re.I):
            continue
        if "policy.offline" in line:
            continue

        if RU_DAEMON_LATIN.search(plain):
            print(f"docs-terms: Latin daemon in RU {doc}:{i}: {line.strip()}", file=sys.stderr)
            fail = True
        if RU_LATIN.search(plain):
            print(f"docs-terms: Latin in RU {doc}:{i}: {line.strip()}", file=sys.stderr)
            fail = True
        if RU_UNITS.search(plain):
            print(f"docs-terms: Latin unit in RU {doc}:{i}: {line.strip()}", file=sys.stderr)
            fail = True

if fail:
    sys.exit(1)
print("docs-terms: ok")
PY
