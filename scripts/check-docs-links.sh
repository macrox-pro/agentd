#!/usr/bin/env bash
# Verify relative .md links and GitHub-style anchors in docs/.
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$root"

python3 <<'PY'
import re
import sys
from pathlib import Path

root = Path(".")


def github_slug(text: str) -> str:
    s = text.strip().lower()
    s = re.sub(r"[^\w\s\u0400-\u04ff-]", "", s, flags=re.UNICODE)
    s = re.sub(r"[\s_]+", "-", s, flags=re.UNICODE)
    return s.strip("-")


def headings(path: Path) -> set[str]:
    slugs: set[str] = set()
    for line in path.read_text(encoding="utf-8").splitlines():
        m = re.match(r"^#{1,6}\s+(.+?)(?:\s+#.*)?$", line)
        if m:
            slugs.add(github_slug(m.group(1)))
    return slugs


link_re = re.compile(r"\[[^\]]*\]\(([^)]+)\)")
fail = False

for doc in sorted(root.glob("docs/**/*.md")):
    text = doc.read_text(encoding="utf-8")
    for raw in link_re.findall(text):
        if ".md" not in raw:
            continue
        target_part, _, anchor = raw.partition("#")
        if target_part:
            target = (doc.parent / target_part).resolve()
        else:
            target = doc.resolve()
        if not target.is_file():
            print(f"docs-links: missing file in {doc} -> {raw}", file=sys.stderr)
            fail = True
            continue
        if anchor:
            slugs = headings(target)
            if anchor not in slugs:
                print(
                    f"docs-links: missing anchor #{anchor} in {target} (from {doc})",
                    file=sys.stderr,
                )
                fail = True

if fail:
    sys.exit(1)
print("docs-links: ok")
PY
