---
primary_sources:
  - id: T2-ZEN
    title: "The Zen of Go"
    url: "https://the-zen-of-go.netlify.app/"
    author: "Dave Cheney"
    section: "All principles"
  - id: T2-PG
    title: "Practical Go"
    url: "https://dave.cheney.net/practical-go/presentations/qcon-china.html"
    author: "Dave Cheney"
    section: "§1 Guiding principles"
  - id: T3-T12
    title: "Twelve Go Best Practices"
    url: "https://talks.golang.org/2013/bestpractices.slide"
    author: "Francesc Campoy Flores"
    section: "Practices 3–12"
studied_at: "2026-08-25"
go_min: "1.26.7"
applicability: "current"
---
# Guiding principles

> **Applicability (Go >= 1.26.7):** Guidance below is current for Go >= 1.26.7.

### Source: The Zen of Go — All principles

> Ten engineering values for writing simple, readable, maintainable Go code. Presented at GopherCon Israel 2020.

> **Each package fulfils a single purpose**
>
> A well designed Go package provides a single idea, a set of related behaviours. A good Go package starts by choosing a good name. Think of your package's name as an elevator pitch to describe what it provides, using just one word.

> **Handle errors explicitly**
>
> Robust programs are composed from pieces that handle the failure cases before they pat themselves on the back. The verbosity of `if err != nil { return err }` is outweighed by the value of deliberately handling each failure condition at the point at which they occur. Panic and recover are not exceptions, they aren't intended to be used that way.

> **Return early rather than nesting deeply**
>
> Every time you indent you add another precondition to the programmer's stack consuming one of the 7 ±2 slots in their short term memory. Avoid control flow that requires deep indentation. Rather than nesting deeply, keep the success path to the left using guard clauses.

> **Leave concurrency to the caller**
>
> Let the caller choose if they want to run your library or function asynchronously, don't force it on them. If your library uses concurrency it should do so transparently.

> **Before you launch a goroutine, know when it will stop**
>
> Goroutines own resources; locks, variables, memory, etc. The sure fire way to free those resources is to stop the owning goroutine.

> **Avoid package level state**
>
> Seek to be explicit, reduce coupling, and spooky action at a distance by providing the dependencies a type needs as fields on that type rather than using package variables.

> **Simplicity matters**
>
> Simplicity is not a synonym for unsophisticated. Simple doesn't mean crude, it means readable and maintainable. When it is possible to choose, defer to the simpler solution.

> **Write tests to lock in the behaviour of your package's API**
>
> Test first or test later, if you shoot for 100% test coverage or are happy with less, regardless your package's API is your contract with its users. Tests are the guarantees that those contracts are written in. Make sure you test for the behaviour that users can observe and rely on.

> **If you think it's slow, first prove it with a benchmark**
>
> So many crimes against maintainability are committed in the name of performance. Optimisation tears down abstractions, exposes internals, and couples tightly. If you're choosing to shoulder that cost, ensure it is done for good reason.

> **Moderation is a virtue**
>
> Use goroutines, channels, locks, interfaces, embedding, in moderation.

> **Maintainability counts**
>
> Clarity, readability, simplicity, are all aspects of maintainability. Can the thing you worked hard to build be maintained after you're gone? What can you do today to make it easier for those that come after you?

### Source: Practical Go — §1 Guiding principles

> If I'm going to talk about best practices in any programming language I need some way to define what I mean by best. If you came to my keynote yesterday you would have seen this quote from the Go team lead, Russ Cox:
>
> Software engineering is what happens to programming when you add time and other programmers.
>
> — Russ Cox
>
> Russ is making the distinction between software programming and software engineering. The former is a program you write for yourself, the latter is a product that many people will work on over time. Engineers will come and go, teams will grow and shrink, requirements will change, features will be added and bugs fixed. This is the nature of software engineering.
>
> I'm possibly one of the earliest users of Go in this room, but to argue that my seniority gives my views more weight is false. Instead, the advice I'm going to present today is informed by what I believe to be the guiding principles underlying Go itself. They are:
>
> Simplicity
>
> Readability
>
> Productivity
>
> You'll note that I didn't say performance, or concurrency. There are languages which are a bit faster than Go, but they're certainly not as simple as Go. There are languages which make concurrency their highest goal, but they are not as readable, nor as productive. Performance and concurrency are important attributes, but not as important as simplicity, readability, and productivity.

> #### 1.1. Simplicity
>
> Simplicity is prerequisite for reliability.
>
> — Edsger W. Dijkstra
>
> Why should we strive for simplicity? Why is important that Go programs be simple?
>
> We've all been in a situation where you say "I can't understand this code", yes? We've all worked on programs where you're scared to make a change because you're worried it'll break another part of the program; a part you don't understand and don't know how to fix. This is complexity.
>
> There are two ways of constructing a software design: One way is to make it so simple that there are obviously no deficiencies, and the other way is to make it so complicated that there are no obvious deficiencies. The first method is far more difficult.
>
> — C. A. R. Hoare
>
> Complexity turns reliable software in unreliable software. Complexity is what kills software projects. Therefore simplicity is the highest goal of Go. Whatever programs we write, we should be able to agree that they are simple.

> #### 1.2. Readability
>
> Readability is essential for maintainability.
>
> — Mark Reinhold JVM language summit 2018
>
> Why is it important that Go code be readable? Why should we strive for readability?
>
> Programs must be written for people to read, and only incidentally for machines to execute.
>
> — Hal Abelson and Gerald Sussman Structure and Interpretation of Computer Programs
>
> Readability is important because all software, not just Go programs, is written by humans to be read by other humans. The fact that software is also consumed by machines is secondary.
>
> Code is read many more times than it is written. A single piece of code will, over its lifetime, be read hundreds, maybe thousands of times.
>
> The most important skill for a programmer is the ability to effectively communicate ideas.
>
> — Gastón Jorquera
>
> Readability is key to being able to understand what the program is doing. If you can't understand what a program is doing, how can you hope to maintain it? If software cannot be maintained, then it will be rewritten; and that could be the last time your company will invest in Go.
>
> If you're writing a program for yourself, maybe it only has to run once, or you're the only person who'll ever see it, then do what ever works for you. But if this is a piece of software that more than one person will contribute to, or that will be used by people over a long enough time that requirements, features, or the environment it runs in changes, then your goal must be for your program to be maintainable.
>
> The first step towards writing maintainable code is making sure the code is readable.

> #### 1.3. Productivity
>
> Design is the art of arranging code to work today, and be changeable forever.
>
> — Sandi Metz
>
> The last underlying principle I want to highlight is productivity. Developer productivity is a sprawling topic but it boils down to this; how much time do you spend doing useful work verses waiting for your tools or hopelessly lost in a foreign code-base. Go programmers should feel that they can get a lot done with Go.
>
> The joke goes that Go was designed while waiting for a C++ program to compile. Fast compilation is a key feature of Go and a key recruiting tool to attract new developers. While compilation speed remains a constant battleground, it is fair to say that compilations which take minutes in other languages, take seconds in Go. This helps Go developers feel as productive as their counterparts working in dynamic languages without the reliability issues inherent in those languages.
>
> More fundamental to the question of developer productivity, Go programmers realise that code is written to be read and so place the act of reading code above the act of writing it. Go goes so far as to enforce, via tooling and custom, that all code be formatted in a specific style. This removes the friction of learning a project specific dialect and helps spot mistakes because they just look incorrect.
>
> Go programmers don't spend days debugging inscrutable compile errors. They don't waste days with complicated build scripts or deploying code to production. And most importantly they don't spend their time trying to understand what their coworker wrote.
>
> Productivity is what the Go team mean when they say the language must scale.

### Source: Twelve Go Best Practices — 3. Important code goes first

> Important code goes first
>
> License information, build tags, package documentation.
>
> Import statements, related groups separated by blank lines.
>
> The rest of the code starting with the most significant types, and ending with helper function and types.

### Source: Twelve Go Best Practices — 4. Document your code

> Document your code
>
> Package name, with the associated documentation before.
>
> Exported identifiers appear in `godoc`, they should be documented correctly.

### Source: Twelve Go Best Practices — 5. Shorter is better

> Shorter is better
>
> or at least longer is not always better.
>
> Try to find the shortest name that is self explanatory.
>
> Prefer `MarshalIndent` to `MarshalWithIndentation`.
>
> Don't forget that the package name will appear before the identifier you chose.
>
> In package `encoding/json` we find the type `Encoder`, not `JSONEncoder`.
>
> It is referred as `json.Encoder`.

### Source: Twelve Go Best Practices — 6. Packages with multiple files

> Packages with multiple files
>
> Should you split a package into multiple files?
>
> Avoid very long files
>
> The `net/http` package from the standard library contains 15734 lines in 47 files.
>
> Separate code and tests
>
> `net/http/cookie.go` and `net/http/cookie_test.go` are both part of the `http` package.
>
> Test code is compiled only at test time.
>
> Separated package documentation
>
> When we have more than one file in a package, it's convention to create a `doc.go` containing the package documentation.


### Source: Twelve Go Best Practices — 8. Ask for what you need

> Ask for what you need
>
> But using a concrete type makes this code difficult to test, so we use an interface.
>
> And, since we're using an interface, we should ask only for the methods we need.

### Source: Twelve Go Best Practices — 9. Keep independent packages independent

> Keep independent packages independent
>
> Avoid dependency by using an interface.

### Source: Twelve Go Best Practices — 10. Avoid concurrency in your API

> Avoid concurrency in your API
>
> What if we want to use it sequentially?
>
> Expose synchronous APIs, calling them concurrently is easy.

### Source: Twelve Go Best Practices — 11. Use goroutines to manage state

> Use goroutines to manage state
>
> Use a chan or a struct with a chan to communicate with a goroutine

### Source: Twelve Go Best Practices — 12. Avoid goroutine leaks

> Avoid goroutine leaks with buffered chans
>
> the goroutine is blocked on the chan write
> the goroutine holds a reference to the chan
> the chan will never be garbage collected
>
> Avoid goroutines leaks with quit chan

