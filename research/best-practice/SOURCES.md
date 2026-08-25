# Source catalog

Target applicability: **Go >= 1.26.7**. Status `done` means studied; only excerpts that remain valid for that toolchain are retained in topic files.

| ID | Source | URL | Tier | Status | studied_at |
|----|--------|-----|------|--------|------------|
| T1-EG | Effective Go | https://go.dev/doc/effective_go | T1 | done | 2026-08-25 |
| T1-CRC | Go Code Review Comments | https://go.dev/wiki/CodeReviewComments | T1 | done | 2026-08-25 |
| T1-CM | Common Mistakes | https://go.dev/wiki/CommonMistakes | T1 | done | 2026-08-25 |
| T1-GP | Go Proverbs | https://go-proverbs.github.io | T1 | done | 2026-08-25 |
| T2-UBER | Uber Go Style Guide | https://github.com/uber-go/guide | T2 | done | 2026-08-25 |
| T2-PG | Practical Go | https://dave.cheney.net/practical-go/presentations/qcon-china.html | T2 | done | 2026-08-25 |
| T2-ZEN | The Zen of Go | https://the-zen-of-go.netlify.app/ | T2 | done | 2026-08-25 |
| T2-ADV | go-advices | https://github.com/cristaloleg/go-advices | T2 | done | 2026-08-25 |
| T2-RAK | Style guideline for Go packages | https://rakyll.org/style-packages/ | T2 | done | 2026-08-25 |
| T2-BAHL | Go Styleguide (bahlo) | https://github.com/bahlo/go-styleguide | T2 | done | 2026-08-25 |
| T3-PERF | go-perfbook | https://github.com/dgryski/go-perfbook | T3 | done | 2026-08-25 |
| T3-TEST | Go testing style guide | https://arp242.net/weblog/go-testing-style.html | T3 | done | 2026-08-25 |
| T3-CLEAN | Clean Go Code | https://github.com/Pungyeon/clean-go-article | T3 | done | 2026-08-25 |
| T3-IDIO | Idiomatic Go | https://dmitri.shuralyov.com/idiomatic-go | T3 | done | 2026-08-25 |
| T3-T12 | Twelve Go Best Practices | https://talks.golang.org/2013/bestpractices.slide | T3 | done | 2026-08-25 |
| SKIP-MEME | Evolution of a Go Programmer | https://github.com/SuperPaintman/the-evolution-of-a-go-programmer | skip | skipped | |
| SKIP-ZH | Chinese translations | — | skip | skipped | |

## Notes

- **T1-CM:** Studied; loop-iterator sections **not retained** (fixed in Go 1.22). No topic file for this source.
- **T3-T12:** Practice 7 (“go get”-able / pre-modules) **not retained**. Other practices kept.
- **T2-BAHL:** Upstream repo 404 at study time; English excerpts from [cch123/go-styleguide](https://github.com/cch123/go-styleguide). `dep` / `gometalinter` sections **not retained**.
- **T2-ZEN URL:** `the-zen-of-go.netlify.app`.
