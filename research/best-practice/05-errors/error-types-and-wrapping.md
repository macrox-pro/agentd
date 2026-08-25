---
primary_sources:
  - id: T1-CRC
    title: "Go Code Review Comments"
    url: "https://go.dev/wiki/CodeReviewComments"
    author: "Go team"
    section: "In-Band Errors"
  - id: T2-UBER
    title: "Uber Go Style Guide"
    url: "https://github.com/uber-go/guide"
    author: "Uber"
    section: "Errors (Error Types, Error Wrapping, Error Naming, Handle Errors Once)"
  - id: T3-CLEAN
    title: "Clean Go Code"
    url: "https://github.com/Pungyeon/clean-go-article"
    author: "Pungyeon"
    section: "Returning Defined Errors"
studied_at: "2026-08-25"
also_cited_in: []
go_min: "1.26.7"
applicability: "current"
---
# Error types and wrapping

> **Applicability (Go >= 1.26.7):** Guidance below is current for Go >= 1.26.7.

### Source: Go Code Review Comments — In-Band Errors

> In C and similar languages, it's common for functions to return values like -1 or null to signal errors or missing results:
>
> ```
> // Lookup returns the value for key or "" if there is no mapping for key.
> func Lookup(key string) string
>
> // Failing to check for an in-band error value can lead to bugs:
> Parse(Lookup(key))  // returns "parse failure for value" instead of "no value for key"
> ```
>
> Go's support for multiple return values provides a better solution. Instead of requiring clients to check for an in-band error value, a function should return an additional value to indicate whether its other return values are valid. This return value may be an error, or a boolean when no explanation is needed. It should be the final return value.
>
> ```
> // Lookup returns the value for key or ok=false if there is no mapping for key.
> func Lookup(key string) (value string, ok bool)
> ```
>
> This prevents the caller from using the result incorrectly:
>
> ```
> Parse(Lookup(key))  // compile-time error
> ```
>
> And encourages more robust and readable code:
>
> ```
> value, ok := Lookup(key)
> if !ok {
>     return fmt.Errorf("no value for %q", key)
> }
> return Parse(value)
> ```
>
> This rule applies to exported functions but is also useful for unexported functions.
>
> Return values like nil, "", 0, and -1 are fine when they are valid results for a function, that is, when the caller need not handle them differently from other values.
>
> Some standard library functions, like those in package "strings", return in-band error values. This greatly simplifies string-manipulation code at the cost of requiring more diligence from the programmer. In general, Go code should return additional values for errors.

### Source: Uber Go Style Guide — Error Types

> There are few options for declaring errors. Consider the following before picking the option best suited for your use case.
>
> - Does the caller need to match the error so that they can handle it? If yes, we must support the `errors.Is` or `errors.As` functions by declaring a top-level error variable or a custom type.
> - Is the error message a static string, or is it a dynamic string that requires contextual information? For the former, we can use `errors.New`, but for the latter we must use `fmt.Errorf` or a custom error type.
> - Are we propagating a new error returned by a downstream function? If so, see the section on error wrapping.
>
> | Error matching? | Error Message | Guidance |
> |-----------------|---------------|--------------------------------------------------------------------|
> | No | static | `errors.New` |
> | No | dynamic | `fmt.Errorf` |
> | Yes | static | top-level `var` with `errors.New` |
> | Yes | dynamic | custom `error` type |
>
> For example, use `errors.New` for an error with a static string. Export this error as a variable to support matching it with `errors.Is` if the caller needs to match and handle this error.
>
> For an error with a dynamic string, use `fmt.Errorf` if the caller does not need to match it, and a custom `error` if the caller does need to match it.
>
> Note that if you export error variables or types from a package, they will become part of the public API of the package.

### Source: Uber Go Style Guide — Error Wrapping

> There are three main options for propagating errors if a call fails:
>
> - return the original error as-is
> - add context with `fmt.Errorf` and the `%w` verb
> - add context with `fmt.Errorf` and the `%v` verb
>
> Return the original error as-is if there is no additional context to add. This maintains the original error type and message. This is well suited for cases when the underlying error message has sufficient information to track down where it came from.
>
> Otherwise, add context to the error message where possible so that instead of a vague error such as "connection refused", you get more useful errors such as "call service foo: connection refused".
>
> Use `fmt.Errorf` to add context to your errors, picking between the `%w` or `%v` verbs based on whether the caller should be able to match and extract the underlying cause.
>
> - Use `%w` if the caller should have access to the underlying error. This is a good default for most wrapped errors, but be aware that callers may begin to rely on this behavior. So for cases where the wrapped error is a known `var` or type, document and test it as part of your function's contract.
> - Use `%v` to obfuscate the underlying error. Callers will be unable to match it, but you can switch to `%w` in the future if needed.
>
> When adding context to returned errors, keep the context succinct by avoiding phrases like "failed to", which state the obvious and pile up as the error percolates up through the stack:
>
> Bad: `failed to create new store: %w` → `failed to x: failed to y: failed to create new store: the error`
>
> Good: `new store: %w` → `x: y: new store: the error`
>
> However once the error is sent to another system, it should be clear the message is an error (e.g. an `err` tag or "Failed" prefix in logs).

### Source: Uber Go Style Guide — Error Naming

> For error values stored as global variables, use the prefix `Err` or `err` depending on whether they're exported. This guidance supersedes the Prefix Unexported Globals with _.
>
> For custom error types, use the suffix `Error` instead.

### Source: Uber Go Style Guide — Handle Errors Once

> When a caller receives an error from a callee, it can handle it in a variety of different ways depending on what it knows about the error.
>
> These include, but not are limited to:
>
> - if the callee contract defines specific errors, matching the error with `errors.Is` or `errors.As` and handling the branches differently
> - if the error is recoverable, logging the error and degrading gracefully
> - if the error represents a domain-specific failure condition, returning a well-defined error
> - returning the error, either wrapped or verbatim
>
> Regardless of how the caller handles the error, it should typically handle each error only once. The caller should not, for example, log the error and then return it, because *its* callers may handle the error as well.
>
> **Bad**: Log the error and return it — Callers further up the stack will likely take a similar action with the error. Doing so makes a lot of noise in the application logs for little value.
>
> **Good**: Wrap the error and return it — Callers further up the stack will handle the error. Use of `%w` ensures they can match the error with `errors.Is` or `errors.As` if relevant.
>
> **Good**: Log the error and degrade gracefully — If the operation isn't strictly necessary, we can provide a degraded but unbroken experience by recovering from it.
>
> **Good**: Match the error and degrade gracefully — If the callee defines a specific error in its contract, and the failure is recoverable, match on that error case and degrade gracefully. For all other cases, wrap the error and return it. Callers further up the stack will handle other errors.

### Source: Clean Go Code — Returning Defined Errors

> By simply representing the error as a variable (`ErrItemNotFound`), we've ensured that anyone using this package can check against the variable rather than the actual string that it returns:
>
> This feels much nicer and is also much safer. Some would even say that it's easier to read as well. In the case of a more verbose error message, it certainly would be preferable for a developer to simply read `ErrItemNotFound` rather than a novel on why a certain error has been returned.
