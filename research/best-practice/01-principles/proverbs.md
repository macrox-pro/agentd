---
primary_sources:
  - id: T1-GP
    title: "Go Proverbs"
    url: "https://go-proverbs.github.io"
    author: "Rob Pike"
studied_at: "2026-08-25"
also_cited_in: []
go_min: "1.26.7"
applicability: "current"
---
# Go Proverbs

> **Applicability (Go >= 1.26.7):** Guidance below is current for Go >= 1.26.7.

### Source: Go Proverbs — Simple, Poetic, Pithy

> Don't communicate by sharing memory, share memory by communicating.
>
> Concurrency is not parallelism.
>
> Channels orchestrate; mutexes serialize.
>
> The bigger the interface, the weaker the abstraction.
>
> Make the zero value useful.
>
> interface{} says nothing.
>
> > **Applicability (Go >= 1.26.7):** Prefer `any` in new code; the proverb still means "empty interfaces convey no contract."
>
> Gofmt's style is no one's favorite, yet gofmt is everyone's favorite.
>
> A little copying is better than a little dependency.
>
> Syscall must always be guarded with build tags.
>
> Cgo must always be guarded with build tags.
>
> Cgo is not Go.
>
> With the unsafe package there are no guarantees.
>
> Clear is better than clever.
>
> Reflection is never clear.
>
> Errors are values.
>
> Don't just check errors, handle them gracefully.
>
> Design the architecture, name the components, document the details.
>
> Documentation is for users.
>
> Don't panic.
