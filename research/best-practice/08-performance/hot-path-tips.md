---
primary_sources:
  - id: T3-PERF
    title: "go-perfbook"
    url: "https://github.com/dgryski/go-perfbook"
    author: "Damian Gryski"
    section: "Runtime and compiler; Common gotchas with the standard library"
  - id: T2-ADV
    title: "go-advices"
    url: "https://github.com/cristaloleg/go-advices"
    author: "Oleg Kiselyov"
    section: "Performance"
studied_at: "2026-08-25"
go_min: "1.26.7"
applicability: "current"
---
# Hot path tips

> **Applicability (Go >= 1.26.7):** Guidance below is current for Go >= 1.26.7.

### Source: go-perfbook — Runtime and compiler

> * cost of calls via interfaces (indirect calls on the CPU level)
> * runtime.convT2E / runtime.convT2I
> * type assertions vs. type switches
> * defer
> * special-case map implementations for ints, strings
>  * map for byte/uint16 not optimized; use a slice instead.
>  * You can fake a float64-optimized with math.Float{32,64}{from,}bits, but beware float equality issues
> * bounds check elimination
> * []byte <-> string copies, map optimizations
> * two-value range will copy an array, use the slice instead
> * use string concatenation instead of fmt.Sprintf where possible; runtime has optimized routines for it

### Source: go-perfbook — Common gotchas with the standard library

> * time.After() leaks until it fires; use t := NewTimer(); t.Stop() / t.Reset()
> * Reusing HTTP connections...; ensure the body is drained
> * rand.Int() and friends are 1) mutex protected and 2) expensive to create
>  * consider alternate random number generation (go-pcgr, xorshift)
> * binary.Read and binary.Write use reflection and are slow; do it by hand.
> * use strconv instead of fmt if possible
> * Use `strings.EqualFold(str1, str2)` instead of `strings.ToLower(str1) == strings.ToLower(str2)` or `strings.ToUpper(str1) == strings.ToUpper(str2)` to efficiently compare strings if possible.

### Source: go-advices — Performance

> do not omit `defer` — 200ns speedup is negligible in most cases
>
> always close http body aka `defer r.Body.Close()` — unless you need leaked goroutine
>
> filtering without allocating:
>
> `b := a[:0]`
> `for _, x := range a {`
> `	if f(x) {`
> `		b = append(b, x)`
> `	}`
> `}`

> To help compiler to remove bound checks see this pattern `_ = b[7]`

> `time.Time` has pointer field `time.Location` and this is bad for go GC — it's relevant only for big number of `time.Time`, use timestamp instead
>
> prefer `regexp.MustCompile` instead of `regexp.Compile` — in most cases your regex is immutable, so init it in `func init`
>
> do not overuse `fmt.Sprintf` in your hot path. It is costly due to maintaining the buffer pool and dynamic dispatches for interfaces. If you are doing `fmt.Sprintf("%s%s", var1, var2)`, consider simple string concatenation. If you are doing `fmt.Sprintf("%x", var)`, consider using `hex.EncodeToString` or `strconv.FormatInt(var, 16)`
>
> always discard body e.g. `io.Copy(io.Discard, resp.Body)` if you don't use it — HTTP client's Transport will not reuse connections unless the body is read to completion and closed
>
> don't use defer in a loop or you'll get a small memory leak — 'cause defers will grow your stack without the reason
>
> don't forget to stop ticker, unless you need a leaked channel — `ticker := time.NewTicker(1 * time.Second)` / `defer ticker.Stop()`
>
> use custom marshaler to speed up marshaling — but before using it - profile!
>
> `sync.Map` isn't a silver bullet, do not use it without a strong reasons
>
> storing non-pointer values in `sync.Pool` allocates memory
>
> use buffered I/O if you do many sequential reads or writes — to reduce number of syscalls
>
> there are 2 ways to clear a map:
>
> reuse map memory — `for k := range m { delete(m, k) }`
>
> allocate new — `m = make(map[int]int)`

### Source: go-advices — Performance (unique items)

> to hide a pointer from escape analysis you might carefully(!!!) use this func:
>
> ```go
> // noescape hides a pointer from escape analysis.  noescape is
> // the identity function but escape analysis doesn't think the
> // output depends on the input. noescape is inlined and currently
> // compiles down to zero instructions.
> func noescape(p unsafe.Pointer) unsafe.Pointer {
> 	x := uintptr(p)
> 	return unsafe.Pointer(x ^ 0)
> }
> ```
>
> for fastest atomic swap you might use this
> `m := (*map[int]int)(atomic.LoadPointer(&ptr))`
