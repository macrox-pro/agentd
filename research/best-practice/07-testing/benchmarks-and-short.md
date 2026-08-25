---
primary_sources:
  - id: T2-ADV
    title: "go-advices"
    url: "https://github.com/cristaloleg/go-advices"
    author: "Oleg Kiselyov"
    section: "Testing; Tools (benchstat)"
  - id: T3-PERF
    title: "go-perfbook"
    url: "https://github.com/dgryski/go-perfbook"
    author: "Damian Gryski"
    section: "Optimization Workflow (benchmarks / benchstat)"
studied_at: "2026-08-25"
go_min: "1.26.7"
applicability: "current"
---
# Benchmarks and short mode

> **Applicability (Go >= 1.26.7):** Guidance below is current for Go >= 1.26.7.

### Source: go-advices — Testing

> prefer `package_test` name for tests, rather than `package`
>
> `go test -short` allows to reduce set of tests to be runned
>
> ```go
> func TestSomething(t *testing.T) {
> 	if testing.Short() {
> 		t.Skip("skipping test in short mode.")
> 	}
> }
> ```
>
> skip test depending on architecture
>
> ```go
> if runtime.GOARM == "arm" {
> 	t.Skip("this doesn't work under ARM")
> }
> ```
>
> track your allocations with `testing.AllocsPerRun`
> — https://godoc.org/testing#AllocsPerRun
>
> run your benchmarks multiple times, to get rid of noise
> — `go test -test.bench=. -count=20`

### Source: go-advices — Tools (benchstat)

> for fast benchmark comparison we've a `benchstat` tool
> — https://godoc.org/golang.org/x/perf/cmd/benchstat

### Source: go-perfbook — Optimization Workflow (benchmarks)

> The benchmarks you are using must be correct and provide reproducible numbers
> on representative workloads. If individual runs have too high a variance, it
> will make small improvements more difficult to spot. You will need to use
> benchstat or equivalent statistical tests
> and won't be able just to eyeball it.
> (Note that using statistical tests is a good idea anyway.) The steps to run
> the benchmarks should be documented, and any custom scripts and tooling
> should be committed to the repository with instructions for how to run them.
> Be mindful of large benchmark suites that take a long time to run: it will
> make the development iterations slower.
