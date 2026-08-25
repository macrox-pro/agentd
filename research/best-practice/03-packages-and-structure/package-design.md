---
primary_sources:
  - id: T2-PG
    title: "Practical Go"
    url: "https://dave.cheney.net/practical-go/presentations/qcon-china.html"
    author: "Dave Cheney"
    section: "§4 Package Design"
  - id: T2-RAK
    title: "Style guideline for Go packages"
    url: "https://rakyll.org/style-packages/"
    author: "JBD (rakyll)"
    section: "Package Organization, Naming, Documentation"
studied_at: "2026-08-25"
go_min: "1.26.7"
applicability: "current"
---
# Package design

> **Applicability (Go >= 1.26.7):** Guidance below is current for Go >= 1.26.7.

### Source: Practical Go — §4 Package Design

> Write shy code - modules that don't reveal anything unnecessary to other modules and that don't rely on other modules' implementations.
>
> Each Go package is in effect it's own small Go program. Just as the implementation of a function or method is unimportant to the caller, the implementation of the functions, methods and types that comprise your package's public API—its behaviour—is unimportant for the caller.
>
> A good Go package should strive to have a low degree of source level coupling such that, as the project grows, changes to one package do not cascade across the code-base. These stop-the-world refactorings place a hard limit on the rate of change in a code base and thus the productivity of the members working in that code-base.
>
> In this section we'll talk about designing a package—including the package's name—naming types, and tips for writing methods and functions.

> #### 4.1. A good package starts with its name
>
> Writing a good Go package starts with the package's name. Think of your package's name as an elevator pitch to describe what it does using just one word.
>
> Just as I talked about names for variables in the previous section, the name of a package is very important. The rule of thumb I follow is not, "what types should I put in this package?". Instead the question I ask "what does service does package provide?" Normally the answer to that question is not "this package provides the X type", but "this package let's you speak HTTP".
>
> Name your package for what it provides, not what it contains.

> ##### 4.1.1. Good package names should be unique.
>
> Within your project, each package name should be unique. This should pretty easy to if you've followed the advice that a package's name should derive from its purpose. If you find you have two packages which need the same name, it is likely either;
>
> The name of the package is too generic.
>
> The package overlaps another package of a similar name. In this case either you should review your design, or consider merging the packages.

> #### 4.2. Avoid package names like base, common, or util
>
> A common cause of poor package names is what call utility packages. These are packages where common helpers and utility code congeals over time. As these packages contain an assortment of unrelated functions, their utility is hard to describe in terms of what the package provides. This often leads to the package's name being derived from what the package contains--utilities.
>
> Package names like `utils` or `helpers` are commonly found in larger projects which have developed deep package hierarchies and want to share helper functions without encountering import loops. By extracting utility functions to new package the import loop is broken, but because the package stems from a design problem in the project, its name doesn't reflect its purpose, only its function of breaking the import cycle.
>
> My recommendation to improve the name of `utils` or `helpers` packages is to analyse where they are called and if possible move the relevant functions into their caller's package. Even if this involves duplicating some helper code this is better than introducing an import dependency between two packages.
>
> [A little] duplication is far cheaper than the wrong abstraction.
>
> — Sandy Metz
>
> In the case where utility functions are used in many places prefer multiple packages, each focused on a single aspect, to a single monolithic package.
>
> Use plurals for naming utility packages. For example the `strings` for string handling utilities.
>
> Packages with names like `base` or `common` are often found when functionality common to two or more implementations, or common types for a client and server, has been refactored into a separate package. I believe the solution to this is to reduce the number of packages, to combine the client, server, and common code into a single package named after the function of the package.
>
> For example, the `net/http` package does not have `client` and `server` sub packages, instead it has a `client.go` and `server.go` file, each holding their respective types, and a `transport.go` file for the common message transport code.
>
> An identifier's name includes its package name. It's important to remember that the name of an identifier includes the name of its package. The `Get` function from the `net/http` package becomes `http.Get` when referenced by another package. The `Reader` type from the `strings` package becomes `strings.Reader` when imported into other packages. The `Error` interface from the `net` package is clearly related to network errors.

> #### 4.3. Return early rather than nesting deeply
>
> As Go does not use exceptions for control flow there is no requirement to deeply indent your code just to provide a top level structure for the `try` and `catch` blocks. Rather than the successful path nesting deeper and deeper to the right, Go code is written in a style where the success path continues down the screen as the function progresses. My friend Mat Ryer calls this practice 'line of sight' coding.
>
> This is achieved by using guard clauses; conditional blocks with assert preconditions upon entering a function.
>
> Upon entering `UnreadRune` the state of `b.lastRead` is checked and if the previous operation was not `ReadRune` an error is returned immediately. From there the rest of the function proceeds with the assertion that `b.lastRead` is greater that `opInvalid`.
>
> The body of the successful case, the most common, is nested inside the first `if` condition and the successful exit condition, `return nil`, has to be discovered by careful matching of closing braces. The final line of the function now returns an error, and the called must trace the execution of the function back to the matching opening brace to know when control will reach this point.
>
> This is more error prone for the reader, and the maintenance programmer, hence why Go prefer to use guard clauses and returning early on errors.

> #### 4.4. Make the zero value useful
>
> Every variable declaration, assuming no explicit initialiser is provided, will be automatically initialised to a value that matches the contents of zeroed memory. This is the values zero value. The type of the value determines the value's zero value; for numeric types it is zero, for pointer types nil, the same for slices, maps, and channels.
>
> This property of always setting a value to a known default is important for safety and correctness of your program and can make your Go programs simpler and more compact. This is what Go programmers talk about when they say "give your structs a useful zero value".
>
> Consider the `sync.Mutex` type. `sync.Mutex` contains two unexported integer fields, representing the mutex's internal state. Thanks to the zero value those fields will be set to will be set to 0 whenever a `sync.Mutex` is declared. `sync.Mutex` has been deliberately coded to take advantage of this property, making the type usable without explicit initialisation.
>
> Another example of a type with a useful zero value is `bytes.Buffer`. You can declare a `bytes.Buffer` and start writing to it without explicit initialisation.
>
> A useful property of slices is their zero value is `nil`. This means you don't need to explicitly `make` a slice, you can just declare it.
>
> A useful, albeit surprising, property of uninitialised pointer variables—nil pointers—is you can call methods on types that have a nil value. This can be used to provide default values simply.

> #### 4.5. Avoid package level state
>
> The key to writing maintainable programs is that they should be loosely coupled—a change to one package should have a low probability of affecting another package that does not directly depend on the first.
>
> There are two excellent ways to achieve loose coupling in Go:
>
> Use interfaces to describe the behaviour your functions or methods require.
>
> Avoid the use of global state.
>
> In Go we can declare variables at the function or method scope, and also at the package scope. When the variable is public, given a identifier starting with a capital letter, then its scope is effectively global to the entire program—any package may observe the type and contents of that variable at any time.
>
> Mutable global state introduces tight coupling between independent parts of your program as global variables become an invisible parameter to every function in your program! Any function that relies on a global variable can be broken if that variable's type changes. Any function that relies on the state of a global variable can be broken if another part of the program changes that variable.
>
> If you want to reduce the coupling a global variable creates,
>
> Move the relevant variables as fields on structs that need them.
>
> Use interfaces to reduce the coupling between the behaviour and the implementation of that behaviour.

### Source: Style guideline for Go packages — Package Organization

> Go is about naming and organization as much as everything else in the language. Well-organized Go code is easy to discover, use and read. Well-organized code is as critical as well designed APIs. The location, name, and the structure of your packages are the first elements your users see and interact with.
>
> This document's goal is to guide you with common good practices not to set rules. You will always need to use your own judgement to pick the most elegant solution for your specific case.
>
> All Go code is organized into packages. A package in Go is simply a directory/folder with one or more `.go` files inside of it. Go packages provide isolation and organization of code similar to how directories/folders organize files on a computer.
>
> All Go code lives in a package and a package is the entry point to access Go code. Understanding and establishing good practices around packages is important to write effective Go code.

> #### Use multiple files
>
> A package is a directory with one or more Go files. Feel free to separate your code into as many files as logically make sense for optimal readability.
>
> For example, an HTTP package might have been separated into different files according to the HTTP aspect the file handles. In the following example, an HTTP package is broken down into a few files: header types and code, cookie types and code, the actual HTTP implementation, and documentation of the package.
>
> `- doc.go       // package documentation`
> `- headers.go   // HTTP headers types and code`
> `- cookies.go   // HTTP cookies types and code`
> `- http.go      // HTTP client implementation, request and response types, etc.`

> #### Keep types close
>
> As a rule of thumb, keep types closer to where they are used. This makes it easy for any maintainer (not just the original author) to find a type. A good place for a Header struct type might be in `headers.go`.
>
> Even though, the Go language doesn't restrict where you define types, it is often a good practice to keep the core types grouped at the top of a file.

> #### Organize by responsibility
>
> A common practise from other languages is to organize types together in a package called models or types. In Go, we organize code by their functional responsibilities.
>
> `package models // DON'T DO IT!!!`
>
> Rather than creating a models package and declare all entity types there, a User type should live in a service-layer package.

> #### Optimize for godoc
>
> It is a great exercise to use godoc in the early phases of your package's API design to see how your concepts will be rendered on doc. Sometimes, the visualization also has an impact on the design. Godoc is the way your users will consume a package, so it is ok to tweak things to make them more accessible. Run `godoc -http=` to start a godoc server locally.

> #### Provide examples to fill the gaps
>
> In some cases, you may not be able to provide all related types from a single package. It might be noisy to do so, or you might want to publish concrete implementations of a common interface from a separate package, or those types could be owned by a third-party package. Give examples to help the user to discover and understand how they are used together.
>
> If your API requires many non-standard packages to be imported, it is often useful to add a Go example to give your users some working code.
>
> Examples are a good way to increase visibility of a less discoverable package.

> #### Don't export from main
>
> An identifier may be exported to permit access to it from another package.
>
> Main packages are not importable, so exporting identifiers from main packages is unnecessary. Don't export identifiers from a main package if you are building the package to a binary.

### Source: Style guideline for Go packages — Package Naming

> A package name and import path are both significant identifiers of your package and represent everything your package contains. Naming your packages canonically not just improves your code quality but also your users'.

> #### Lowercase only
>
> Package names should be lowercase. Don't use snake_case or camelCase in package names. The Go blog has a comprehensive guide about naming packages with a good variety of examples.

> #### Short, but representative names
>
> Package names should be short, but should be unique and representative. Users of the package should be able to grasp its purpose from just the package's name.
>
> Avoid overly broad package names like "common" and "util".
>
> `import "pkgs.org/common" // DON'T!!!`
>
> Avoid duplicate names in cases where user may need to import the same package.
>
> If you cannot avoid a bad name, it is very likely that there is a problem with your overall structure and code organization.

> #### Clean import paths
>
> Avoid exposing your custom repository structure to your users. Avoid having src/, pkg/ sections in your import paths.
>
> `github.com/user/repo/src/httputil   // DON'T DO IT, AVOID SRC!!`
>
> `github.com/user/repo/gosrc/httputil // DON'T DO IT, AVOID GOSRC!!`

> #### No plurals
>
> In go, package names are not plural. This is surprising to programmers who came from other languages and are retaining an old habit of pluralizing names. Don't name a package httputils, but httputil!
>
> `package httputils  // DON'T DO IT, USE SINGULAR FORM!!`

> #### Renames should follow the same rules
>
> If you are importing more than one packages with the same name, you can locally rename the package names. The renames should follow the same rules mentioned on this article. There is no rule which package you should rename. If you are renaming the standard package library, it is nice to add a go prefix to make the name self document that it is "Go standard library's" package, e.g. `gourl`, `goos`.

> #### Enforce vanity URLs
>
> `go get` supports getting packages by a URL that is different than the URL of the package's repo. These URLs are called vanity URLs and require you to serve a page with specific meta tags the Go tools recognize.
>
> To do that, add an import statement to the package. The go tool will reject any import of this package from any other path and will display a friendly error to the user. If you don't enforce your vanity URLs, there will be two copies of your package that cannot work together due to the different namespace.
>
> `package datastore // import "cloud.google.com/go/datastore"`

### Source: Style guideline for Go packages — Package Documentation

> Always document the package. Package documentation is a top-level comment immediately preceding the package clause. For non-main packages, godoc always starts with "Package {pkgname}" and follows with a description. For main packages, documentation should explain the binary.
>
> `// Package path implements utility routines for manipulating slash-separated paths.`
> `package path`
>
> `// Command gops lists all the processes running on your system.`
> `package main`

> #### Use doc.go
>
> Sometimes, package docs can get very lengthy, especially when they provide details of usage and guidelines. Move the package godoc to a `doc.go` file.
