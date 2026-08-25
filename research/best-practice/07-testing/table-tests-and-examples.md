---
primary_sources:
  - id: T2-UBER
    title: "Uber Go Style Guide"
    url: "https://github.com/uber-go/guide"
    author: "Uber"
    section: "Test Tables"
  - id: T3-TEST
    title: "Go testing style guide"
    url: "https://arp242.net/weblog/go-testing-style.html"
    author: "Martin Tournoij"
    section: "Table tests; subtests; want/have; aligned output"
studied_at: "2026-08-25"
go_min: "1.26.7"
applicability: "current"
---
# Table tests and examples

> **Applicability (Go >= 1.26.7):** Guidance below is current for Go >= 1.26.7.

### Source: Uber Go Style Guide — Test Tables

> Table-driven tests with subtests can be a helpful pattern for writing tests to avoid duplicating code when the core test logic is repetitive.
>
> If a system under test needs to be tested against *multiple conditions* where certain parts of the inputs and outputs change, a table-driven test should be used to reduce redundancy and improve readability.
>
> Test tables make it easier to add context to error messages, reduce duplicate logic, and add new test cases.
>
> We follow the convention that the slice of structs is referred to as `tests` and each test case `tt`. Further, we encourage explicating the input and output values for each test case with `give` and `want` prefixes.

> #### Avoid Unnecessary Complexity in Table Tests
>
> Table tests can be difficult to read and maintain if the subtests contain conditional assertions or other branching logic. Table tests should **NOT** be used whenever there needs to be complex or conditional logic inside subtests (i.e. complex logic inside the `for` loop).
>
> Large, complex table tests harm readability and maintainability because test readers may have difficulty debugging test failures that occur.
>
> Table tests like this should be split into either multiple test tables or multiple individual `Test...` functions.
>
> Some ideals to aim for are:
>
> * Focus on the narrowest unit of behavior
> * Minimize "test depth", and avoid conditional assertions (see below)
> * Ensure that all table fields are used in all tests
> * Ensure that all test logic runs for all table cases
>
> In this context, "test depth" means "within a given test, the number of successive assertions that require previous assertions to hold" (similar to cyclomatic complexity). Having "shallower" tests means that there are fewer relationships between assertions and, more importantly, that those assertions are less likely to be conditional by default.
>
> Concretely, table tests can become confusing and difficult to read if they use multiple branching pathways (e.g. `shouldError`, `expectCall`, etc.), use many `if` statements for specific mock expectations (e.g. `shouldCallFoo`), or place functions inside the table (e.g. `setupMocks func(*FooMock)`).
>
> However, when testing behavior that only changes based on changed input, it may be preferable to group similar cases together in a table test to better illustrate how behavior changes across all inputs, rather than splitting otherwise comparable units into separate tests and making them harder to compare and contrast.
>
> If the test body is short and straightforward, it's acceptable to have a single branching pathway for success versus failure cases with a table field like `shouldErr` to specify error expectations.
>
> This complexity makes it more difficult to change, understand, and prove the correctness of the test.
>
> While there are no strict guidelines, readability and maintainability should always be top-of-mind when deciding between Table Tests versus separate tests for multiple inputs/outputs to a system.

> #### Parallel Tests
>
> Parallel tests, like some specialized loops (for example, those that spawn goroutines or capture references as part of the loop body), must take care to explicitly assign loop variables within the loop's scope to ensure that they hold the expected values.
>
> In the example above, we must declare a `tt` variable scoped to the loop iteration because of the use of `t.Parallel()` below. If we do not do that, most or all tests will receive an unexpected value for `tt`, or a value that changes as they're running.

### Source: Go testing style guide — Use table-drive tests, and consistently use `tt` for a test case

> Try to use table-driven tests whenever feasible, but it's okay to just copy some code when it's not; don't force it (e.g. sometimes it's easier to write a table-driven test for all but one or two cases; be practical).
>
> Consistently using the same variable name for a test case will make it easier to work with large code bases. You don't _have_ to use `tt`, but it is the most commonly used in Go's standard library (564 times, vs. 116 for `tc`).

### Source: Go testing style guide — Use subtests

> Using subtests makes it possible to run just a single test case from a table, as well as easily see which test _exactly_ failed. Even though subtests were added back in 2016, I sometimes still see people not using them for new tests.
>
> I tend to simply use the test number if it's obvious what is being tested, and add a test name if it's not or if there are many test cases.

### Source: Go testing style guide — Use `want` and `have`

> `want` is shorter than `expected`, `have` is shorter than `actual`. Shorter names is always good, IMHO, and is especially beneficial for aligning output (see next item).

### Source: Go testing style guide — Add useful, aligned, information

> It's annoying when a test fail with a useless error message, or a noisy error message which makes it hard to see what exactly went wrong.
>
> When aligned, this is a lot easier:
>
> ```
> t.Errorf("wrong output\ngot: %q\nwant: %q", have, want)
> ```
>
> This is also who I like `want` and `have`: they're of identical length so it's easy to align.
>
> I also tend to prefer to use `%q` or `%#v`, as that will show things like trailing whitespace or unprintable characters more clearly.

