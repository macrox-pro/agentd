# Go Best Practices Research

Structured verbatim excerpts from Go best-practice sources listed in [smallnest/go-best-practices](https://github.com/smallnest/go-best-practices).

**Target runtime: Go >= 1.26.7** (matches this repository). Content that only applies to older toolchains is **omitted**, not archived as “historical tips.”

## How to read

- Each topic file under `01-principles/` … `09-build-modules-tools/` starts with YAML front matter (`primary_sources`, `studied_at`, `go_min`, `applicability`).
- Body text is verbatim from primary sources; see `### Source:` headings.
- Where a source used a deprecated API name in an otherwise valid tip, the retained text uses the current stdlib name (e.g. `io.Discard`, `os.ReadFile`).

## Omitted as inapplicable to Go >= 1.26.7

| Topic | Why omitted |
|-------|-------------|
| Loop-iterator capture bugs | Fixed in Go 1.22 — wiki “Common Mistakes” sections not retained |
| `ioutil` APIs | Deprecated since Go 1.16 — not retained as advice |
| `dep` / `Gopkg.toml` / GOPATH-mode workflows | Replaced by Go modules |
| `gometalinter` / `golint` | Archived / deprecated — use `golangci-lint`, `staticcheck`, `go vet` |
| Pre-modules “go get”-able packaging (Twelve BP §7) | Replaced by `go.mod` publishing |
| `math/rand.NewSource` concurrency caveats as primary guidance | Prefer `math/rand/v2` for new code — old tip removed |

## Topic index

| Topic | File | Primary sources |
|-------|------|-----------------|
| Proverbs | [01-principles/proverbs.md](./01-principles/proverbs.md) | T1-GP |
| Guiding principles | [01-principles/guiding-principles.md](./01-principles/guiding-principles.md) | T2-ZEN, T2-PG, T3-T12 |
| Formatting | [02-style-and-formatting/formatting.md](./02-style-and-formatting/formatting.md) | T1-CRC, T1-EG, T2-BAHL |
| Naming | [02-style-and-formatting/naming.md](./02-style-and-formatting/naming.md) | T1-CRC, T2-PG, T3-CLEAN, T3-IDIO, T2-BAHL |
| Comments / godoc | [02-style-and-formatting/comments-and-godoc.md](./02-style-and-formatting/comments-and-godoc.md) | T1-CRC, T2-PG |
| Package design | [03-packages-and-structure/package-design.md](./03-packages-and-structure/package-design.md) | T1-EG, T2-PG, T2-RAK, T2-BAHL |
| Package naming | [03-packages-and-structure/package-naming.md](./03-packages-and-structure/package-naming.md) | T1-CRC, T2-RAK, T2-BAHL |
| File organization | [03-packages-and-structure/file-organization.md](./03-packages-and-structure/file-organization.md) | T2-RAK, T2-PG |
| Functions | [04-api-and-types/functions.md](./04-api-and-types/functions.md) | T1-CRC, T2-PG, T2-UBER, T3-CLEAN |
| Interfaces | [04-api-and-types/interfaces.md](./04-api-and-types/interfaces.md) | T1-CRC, T2-UBER, T3-CLEAN, T3-IDIO |
| Structs / zero values | [04-api-and-types/structs-and-zero-values.md](./04-api-and-types/structs-and-zero-values.md) | T1-CRC, T2-UBER |
| Context | [04-api-and-types/context.md](./04-api-and-types/context.md) | T1-CRC |
| Error handling | [05-errors/error-handling.md](./05-errors/error-handling.md) | T1-CRC, T2-UBER, T2-PG, T3-CLEAN |
| Error types / wrapping | [05-errors/error-types-and-wrapping.md](./05-errors/error-types-and-wrapping.md) | T1-CRC, T2-UBER |
| Error strings | [05-errors/error-strings-and-naming.md](./05-errors/error-strings-and-naming.md) | T1-CRC |
| Goroutine lifecycle | [06-concurrency/goroutine-lifecycle.md](./06-concurrency/goroutine-lifecycle.md) | T1-CRC, T2-UBER, T2-PG |
| Channels / sync | [06-concurrency/channels-and-sync.md](./06-concurrency/channels-and-sync.md) | T2-UBER, T2-ADV |
| Test packages / helpers | [07-testing/test-packages-and-helpers.md](./07-testing/test-packages-and-helpers.md) | T3-TEST, T2-UBER, T1-CRC |
| Table tests | [07-testing/table-tests-and-examples.md](./07-testing/table-tests-and-examples.md) | T3-TEST, T2-UBER |
| Benchmarks / `-short` | [07-testing/benchmarks-and-short.md](./07-testing/benchmarks-and-short.md) | T3-TEST, T2-ADV |
| When to optimize | [08-performance/when-and-how-to-optimize.md](./08-performance/when-and-how-to-optimize.md) | T3-PERF |
| Memory / GC | [08-performance/memory-and-gc.md](./08-performance/memory-and-gc.md) | T3-PERF |
| Hot-path tips | [08-performance/hot-path-tips.md](./08-performance/hot-path-tips.md) | T3-PERF, T2-ADV |
| Modules | [09-build-modules-tools/modules.md](./09-build-modules-tools/modules.md) | T2-ADV |
| Build tags / cgo | [09-build-modules-tools/build-tags-and-cgo.md](./09-build-modules-tools/build-tags-and-cgo.md) | T1-GP, T2-ADV |
| Linting / tools | [09-build-modules-tools/linting-and-tools.md](./09-build-modules-tools/linting-and-tools.md) | T2-UBER |

## Meta

- Source registry: [SOURCES.md](./SOURCES.md)
- Agent continuation: [AGENT_HANDOFF.md](./AGENT_HANDOFF.md)
- File checklist: [MANIFEST.md](./MANIFEST.md)
- Topic file pattern: [TEMPLATE.md](./TEMPLATE.md)

## Invariants

1. YAML front matter on every topic `.md` file (`go_min: "1.26.7"`).
2. Verbatim excerpts only (minimal glue).
3. No project-specific mapping sections.
4. Do not retain advice that is false or harmful on Go >= 1.26.7.
5. Update [SOURCES.md](./SOURCES.md) when a source is fully processed.
