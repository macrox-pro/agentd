---
primary_sources:
  - id: T3-PERF
    title: "go-perfbook"
    url: "https://github.com/dgryski/go-perfbook"
    author: "Damian Gryski"
    section: "Garbage Collection"
studied_at: "2026-08-25"
go_min: "1.26.7"
applicability: "current"
---
# Memory and garbage collection

> **Applicability (Go >= 1.26.7):** Guidance below is current for Go >= 1.26.7.

### Source: go-perfbook — Garbage Collection

> You pay for memory allocation more than once. The first is obviously when you
> allocate it. But you also pay every time the garbage collection runs.
>
> Reduce/Reuse/Recycle.
> -- @bboreham
>
> * Stack vs. heap allocations
> * What causes heap allocations?
> * Understanding escape analysis (and the current limitation)
> * /debug/pprof/heap , and -base
> * API design to limit allocations:
>  * allow passing in buffers so caller can reuse rather than forcing an allocation
>  * you can even modify a slice in place carefully while you scan over it
>  * passing in a struct could allow caller to stack allocate it
> * reducing pointers to reduce gc scan times
>  * pointer-free slices
>  * maps with both pointer-free keys and values
> * GOGC
> * buffer reuse (sync.Pool vs or custom via go-slab, etc)
> * slicing vs. offset: pointer writes while GC is running need writebarrier
>  * no writebarrier if writing to stack
> * use error variables instead of errors.New() / fmt.Errorf() at call site (performance or style? interface requires pointer, so it escapes to heap anyway)
> * use structured errors to reduce allocation (pass struct value), create string at error printing time
> * size classes
> * beware pinning larger allocation with smaller substrings or slices
