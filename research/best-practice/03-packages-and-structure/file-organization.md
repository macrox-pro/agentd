---
primary_sources:
  - id: T1-EG
    title: "Effective Go"
    url: "https://go.dev/doc/effective_go"
    author: "Go team"
    section: "Examples"
  - id: T1-EG
    title: "Effective Go"
    url: "https://go.dev/doc/effective_go"
    author: "Go team"
    section: "Package names"
  - id: T2-PG
    title: "Practical Go"
    url: "https://dave.cheney.net/practical-go/presentations/qcon-china.html"
    author: "Dave Cheney"
    section: "§5 Project Structure"
studied_at: "2026-08-25"
also_cited_in: []
go_min: "1.26.7"
applicability: "current"
---
# File organization

> **Applicability (Go >= 1.26.7):** Guidance below is current for Go >= 1.26.7.

### Source: Effective Go — Examples

> The Go package sources are intended to serve not only as the core library but also as examples of how to use the language. Moreover, many of the packages contain working, self-contained executable examples you can run directly from the go.dev web site, such as this one(if necessary, click on the word "Example" to open it up). If you have a question about how to approach a problem or how something might be implemented, the documentation, code and examples in the library can provide answers, ideas and background.

### Source: Effective Go — Package names

> Another convention is that the package name is the base name of its source directory; the package in `src/encoding/base64` is imported as `"encoding/base64"` but has name `base64`, not `encoding_base64` and not `encodingBase64`.

### Source: Practical Go — §5 Project Structure

> Let's talk about combining packages together into a project. Commonly this will be a single git repository. In the future Go developers will use the terms module and project interchangeably.
>
> Just like a package, each project should have a clear purpose. If your project is a library, it should provide one thing, say XML parsing, or logging. You should avoid combining multiple purposes into a single project, this will help avoid the dreaded `common` library.
>
> In my experience, the `common` repo ends up tightly coupled to its biggest consumer and that makes it hard to back-port fixes without upgrading both common and consumer in lock step, bringing in a lot of unrelated changes and API breakage along the way.
>
> If your project is an application, like your web application, Kubernetes controller, and so on, then you might have one or more `main` packages inside your project. For example, the Kubernetes controller I work on has a single `cmd/contour` package which serves as both the server deployed to a Kubernetes cluster, and a client for debugging purposes.

> #### 5.1. Consider fewer, larger packages
>
> One of the things I tend to pick up in code review for programmers who are transitioning from other languages to Go is they tend to overuse packages.
>
> Go does not provide elaborate ways of establishing visibility. Go lacks Java's `public`, `protected`, `private`, and implicit `default` access modifiers. There is no equivalent of C++'s notion of a `friend` classes.
>
> In Go we have only two access modifiers, public and private, indicated by the capitalisation of the first letter of the identifier. If an identifier is public, it's name starts with a capital letter, that identifier can be referenced by any other Go package.
>
> You may hear people say exported and not exported as synonyms for public and private.
>
> Given the limited controls available to control access to a package's symbols, what practices should Go programmers follow to avoid creating over-complicated package hierarchies?
>
> Every package, with the exception of `cmd/` and `internal/`, should contain some source code.
>
> The advice I find myself repeating is to prefer fewer, larger packages. Your default position should be to not create a new package. That will lead to too many types being made public creating a wide, shallow, API surface for your package.
>
> Coming from Java? If you're coming from a Java or C# background, consider this rule of thumb. - A Java package is equivalent to a single `.go` source file. - A Go package is equivalent to a whole Maven module or .NET assembly.

> ##### 5.1.1. Arrange code into files by import statements
>
> If you're arranging your packages by what they provide to callers, should you do the same for files within a Go package? How do you know when you should break up a `.go` file into multiple ones? How do you know when you've gone to far and should consider consolidating `.go` file?
>
> Here are the guidelines I use:
>
> Start each package with one `.go` file. Give that file the same name as the name of the folder. eg. `package http` should be placed in a file called `http.go` in a directory named `http`.
>
> As your package grows you may decide to split apart the various responsibilities into different files. eg, `messages.go` contains `the `Request` and `Response` types, `client.go` contains the `Client` type, `server.go` contains the `Server` type.
>
> If you find your files have similar `import` declarations, consider combining them. Alternatively, identify the differences between the import sets and move those
>
> Different files should be responsible for different areas of the package. `messages.go` may be responsible for marshalling of HTTP requests and responses on and off the network, `http.go` may contain the low level network handling logic, `client.go` and `server.go` implement the HTTP business logic of request construction or routing, and so on.
>
> Prefer nouns for source file names.
>
> The Go compiler compiles each package in parallel. Within a package the compiler compiles each function (methods are just fancy functions in Go) in parallel. Changing the layout of your code within a package should not affect compilation time.

> ##### 5.1.2. Prefer internal tests to external tests
>
> The `go` tool supports writing your `testing` package tests in two places. Assuming your package is called `http2`, you can write a `http2_test.go` file and use the `package http2` declaration. Doing so will compile the code in `http2_test.go` as if it were part of the `http2` package. This is known colloquially as an internal test.
>
> The `go` tool also supports a special package declaration, ending in `test`, ie., `package http_test`. This allows your test files to live alongside your code in the same package, however when those tests are compiled they are not part of your package's code, they live in their own package. This allows you to write your tests as if you were another package calling into your code. This is known as an _external test.
>
> I recommend using internal tests when writing unit tests for your package. This allows you to test each function or method directly, avoiding the bureaucracy of external testing.
>
> However, you should place your `Example` test functions in an external test file. This ensures that when viewed in godoc, the examples have the appropriate package prefix and can be easily copy pasted.
>
> Avoid elaborate package hierarchies, resist the desire to apply taxonomy With one exception, which we'll talk about next, the hierarchy of Go packages has no meaning to the `go` tool. For example, the `net/http` package is not a child or sub-package of the `net` package. If you find you have created intermediate directories in your project which contain no `.go` files, you may have failed to follow this advice.

> ##### 5.1.3. Use internal packages to reduce your public API surface
>
> If your project contains multiple packages you may find you have some exported functions which are intended to be used by other packages in your project, but are not intended to be part of your project's public API. If you find yourself in this situation the `go` tool recognises a special folder name—not package name--, `internal/` which can be used to place code which is public to your project, but private to other projects.
>
> To create such a package, place it in a directory named `internal/` or in a sub-directory of a directory named `internal/`. When the `go` command sees an import of a package with `internal` in its path, it verifies that the package doing the import is within the tree rooted at the parent of the `internal` directory.
>
> For example, a package `…/a/b/c/internal/d/e/f` can be imported only by code in the directory tree rooted at `…/a/b/c`. It cannot be imported by code in `…/a/b/g` or in any other repository.

> #### 5.2. Keep package main small as small as possible
>
> Your `main` function, and `main` package should do as little as possible. This is because `main.main` acts as a singleton; there can only be one `main` function in a program, including tests.
>
> Because `main.main` is a singleton there are a lot of assumptions built into the things that `main.main` will call that they will only be called during main.main or main.init, and only called once. This makes it hard to write tests for code written in `main.main`, thus you should aim to move as much of your business logic out of your main function and ideally out of your main package.
>
> `func main()` should parse flags, open connections to databases, loggers, and such, then hand off execution to a high level object.
