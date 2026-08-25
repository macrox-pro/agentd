---
primary_sources:
  - id: T1-GP
    title: "Go Proverbs"
    url: "https://go-proverbs.github.io"
    author: "Rob Pike"
    section: "Syscall, Cgo, build tags"
  - id: T2-ADV
    title: "go-advices"
    url: "https://github.com/cristaloleg/go-advices"
    author: "Oleg Kiselyov"
    section: "Build"
studied_at: "2026-08-25"
go_min: "1.26.7"
applicability: "current"
---
# Build tags and cgo

> **Applicability (Go >= 1.26.7):** Guidance below is current for Go >= 1.26.7.

### Source: Go Proverbs — Syscall, Cgo, build tags

> Syscall must always be guarded with build tags.
>
> Cgo must always be guarded with build tags.
>
> Cgo is not Go.

### Source: go-advices — Build

> strip your binaries with this command `go build -ldflags="-s -w" ...`
>
> easy way to split test into different builds
> - use `// +build integration` and run them with `go test -v --tags integration .`
>
> tiniest Go docker image
> - `CGO_ENABLED=0 go build -ldflags="-s -w" app.go && tar C app | docker import - myimage:latest`
>
> run `go format` on CI and compare diff
> - this will ensure that everything was generated and committed
>
> check if there are mistakes in code formatting `diff -u <(echo -n) <(gofmt -d .)`
