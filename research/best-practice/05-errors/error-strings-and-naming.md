---
primary_sources:
  - id: T1-CRC
    title: "Go Code Review Comments"
    url: "https://go.dev/wiki/CodeReviewComments"
    author: "Go team"
    section: "Error Strings"
  - id: T3-IDIO
    title: "Idiomatic Go"
    url: "https://dmitri.shuralyov.com/idiomatic-go"
    author: "Dmitri Shuralyov"
    section: "Error variable naming"
studied_at: "2026-08-25"
go_min: "1.26.7"
applicability: "current"
---
# Error strings and naming

> **Applicability (Go >= 1.26.7):** Guidance below is current for Go >= 1.26.7.

### Source: Go Code Review Comments — Error Strings

> Error strings should not be capitalized (unless beginning with proper nouns or acronyms) or end with punctuation, since they are usually printed following other context. That is, use `fmt.Errorf("something bad")` not `fmt.Errorf("Something bad")`, so that `log.Printf("Reading %s: %v", filename, err)` formats without a spurious capital letter mid-message. This does not apply to logging, which is implicitly line-oriented and not combined inside other messages.

### Source: Idiomatic Go — Error variable naming

> Do this:
>
> ```
> // Package level exported error.
> var ErrSomething = errors.New("something went wrong")
>
> func main() {
> 	// Normally you call it just "err",
> 	result, err := doSomething()
> 	// and use err right away.
>
> 	// But if you want to give it a longer name, use "somethingError".
> 	var specificError error
> 	result, specificError = doSpecificThing()
>
> 	// ... use specificError later.
> }
> ```
>
> Don't do this:
>
> ```
> var ErrorSomething = errors.New("something went wrong")
> var SomethingErr = errors.New("something went wrong")
> var SomethingError = errors.New("something went wrong")
>
> func main() {
> 	var specificErr error
> 	result, specificErr = doSpecificThing()
>
> 	var errSpecific error
> 	result, errSpecific = doSpecificThing()
>
> 	var errorSpecific error
> 	result, errorSpecific = doSpecificThing()
> }
> ```
>
> For consistency. See https://go.dev/talks/2014/names.slide#14.

