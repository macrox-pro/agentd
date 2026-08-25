---
primary_sources:
  - id: T2-UBER
    title: "Uber Go Style Guide"
    url: "https://github.com/uber-go/guide"
    author: "Uber"
    section: "Test Tables (intro)"
  - id: T1-CRC
    title: "Go Code Review Comments"
    url: "https://go.dev/wiki/CodeReviewComments"
    author: "Go team"
    section: "Useful Test Failures"
  - id: T2-ADV
    title: "go-advices"
    url: "https://github.com/cristaloleg/go-advices"
    author: "Oleg Kiselyov"
    section: "Testing"
  - id: T3-TEST
    title: "Go testing style guide"
    url: "https://arp242.net/weblog/go-testing-style.html"
    author: "Martin Tournoij"
    section: "Don't ignore errors"
studied_at: "2026-08-25"
go_min: "1.26.7"
applicability: "current"
---
# Test packages and helpers

> **Applicability (Go >= 1.26.7):** Guidance below is current for Go >= 1.26.7.

### Source: Uber Go Style Guide — Test Tables (intro)

> Table-driven tests with subtests can be a helpful pattern for writing tests to avoid duplicating code when the core test logic is repetitive.
>
> If a system under test needs to be tested against *multiple conditions* where certain parts of the inputs and outputs change, a table-driven test should be used to reduce redundancy and improve readability.
>
> Test tables make it easier to add context to error messages, reduce duplicate logic, and add new test cases.
>
> We follow the convention that the slice of structs is referred to as `tests` and each test case `tt`. Further, we encourage explicating the input and output values for each test case with `give` and `want` prefixes.

### Source: Go Code Review Comments — Useful Test Failures

> Tests should fail with helpful messages saying what was wrong, with what inputs, what was actually got, and what was expected. It may be tempting to write a bunch of assertFoo helpers, but be sure your helpers produce useful error messages. Assume that the person debugging your failing test is not you, and is not your team. A typical Go test fails like:
>
> `if got != tt.want {`
> `    t.Errorf("Foo(%q) = %d; want %d", tt.in, got, tt.want) // or Fatalf, if test can't test anything more past this point`
> `}`
>
> Note that the order here is actual != expected, and the message uses that order too. Some test frameworks encourage writing these backwards: 0 != x, "expected 0, got x", and so on. Go does not.
>
> If that seems like a lot of typing, you may want to write a table-driven test.
>
> Another common technique to disambiguate failing tests when using a test helper with different input is to wrap each caller with a different TestFoo function, so the test fails with that name:
>
> `func TestSingleValue(t *testing.T) { testHelper(t, []int{80}) }`
> `func TestNoValues(t *testing.T)    { testHelper(t, []int{}) }`
>
> In any case, the onus is on you to fail with a helpful message to whoever's debugging your code in the future.

### Source: go-advices — Testing

> prefer `package_test` name for tests, rather than `package`
>
> `go test -short` allows to reduce set of tests to be runned
>
> skip test depending on architecture — `if runtime.GOARM == "arm" { t.Skip("this doesn't work under ARM") }`
>
> track your allocations with `testing.AllocsPerRun`
>
> run your benchmarks multiple times, to get rid of noise — `go test -test.bench=. -count=20`

### Source: Go testing style guide — Don't ignore errors

> I frequently see people ignore errors in tests. This is a bad idea and can make for some confusing test failures.
>
> Example:
>
> ```
> have, err := Fun()
> if err != nil {
>     t.Fatalf("unexpected error: %v", err)
> }
> ```
>
> or:
>
> ```
> have, err := Fun()
> if err != tt.wantErr {
>     t.Fatalf("wrong error\ngot:  %v\nwant: %v", err, tt.wantErr)
> }
> ```
>
> I often use ErrorContains, which is a useful helper function for testing error messages (avoids some `if err != nil && [..]`).

