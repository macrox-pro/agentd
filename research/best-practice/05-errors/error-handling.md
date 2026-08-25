---
primary_sources:
  - id: T1-CRC
    title: "Go Code Review Comments"
    url: "https://go.dev/wiki/CodeReviewComments"
    author: "Go team"
    section: "Handle Errors"
  - id: T1-CRC
    title: "Go Code Review Comments"
    url: "https://go.dev/wiki/CodeReviewComments"
    author: "Go team"
    section: "Don't Panic"
  - id: T1-CRC
    title: "Go Code Review Comments"
    url: "https://go.dev/wiki/CodeReviewComments"
    author: "Go team"
    section: "Indent Error Flow"
  - id: T3-CLEAN
    title: "Clean Go Code"
    url: "https://github.com/Pungyeon/clean-go-article"
    author: "Pungyeon"
    section: "Returning Defined Errors; Returning Dynamic Errors"
  - id: T2-PG
    title: "Practical Go"
    url: "https://dave.cheney.net/practical-go/presentations/qcon-china.html"
    author: "Dave Cheney"
    section: "§7 Error handling"
studied_at: "2026-08-25"
go_min: "1.26.7"
applicability: "current"
---
# Error handling

> **Applicability (Go >= 1.26.7):** Guidance below is current for Go >= 1.26.7.

### Source: Go Code Review Comments — Handle Errors

> See https://go.dev/doc/effective_go#errors. Do not discard errors using `_` variables. If a function returns an error, check it to make sure the function succeeded. Handle the error, return it, or, in truly exceptional situations, panic.

### Source: Go Code Review Comments — Don't Panic

> See https://go.dev/doc/effective_go#errors. Don't use panic for normal error handling. Use error and multiple return values.

### Source: Go Code Review Comments — Indent Error Flow

> Try to keep the normal code path at a minimal indentation, and indent the error handling, dealing with it first. This improves the readability of the code by permitting visually scanning the normal path quickly. For instance, don't write:
>
> ```
> if err != nil {
>     // error handling
> } else {
>     // normal code
> }
> ```
>
> Instead, write:
>
> ```
> if err != nil {
>     // error handling
>     return // or continue, etc.
> }
> // normal code
> ```
>
> If the `if` statement has an initialization statement, such as:
>
> ```
> if x, err := f(); err != nil {
>     // error handling
>     return
> } else {
>     // use x
> }
> ```
>
> then this may require moving the short variable declaration to its own line:
>
> ```
> x, err := f()
> if err != nil {
>     // error handling
>     return
> }
> // use x
> ```

### Source: Clean Go Code — Returning Defined Errors

> Let's consider the normal way to return a custom error. This is a hypothetical example taken from a thread-safe map implementation that we've named `Store`:
>
> There is nothing inherently smelly about this function when we consider it in isolation. We look into the `items` map of our `Store` struct to see if we already have an item with the given `id`. If we do, we return it; otherwise, we return an error. Pretty standard. So, what is the issue with returning custom errors as string values? Well, let's look at what happens when we use this function inside another package:
>
> This is actually not too bad. However, there is one glaring problem: An error in Go is simply an `interface` that implements a function (`Error()`) returning a string; thus, we are now hardcoding the expected error code into our codebase, which isn't ideal. This hardcoded string is known as a magic string. And its main problem is flexibility: If at some point we decide to change the string value used to represent an error, our code will break (softly) unless we update it in possibly many different places. Our code is tightly coupled—it relies on that specific magic string and the assumption that it will never change as the codebase grows.
>
> An even worse situation would arise if a client were to use our package in their own code. Imagine that we decided to update our package and changed the string that represents an error—the client's software would now suddenly break. This is quite obviously something that we want to avoid. Fortunately, the fix is very simple:
>
> By simply representing the error as a variable (`ErrItemNotFound`), we've ensured that anyone using this package can check against the variable rather than the actual string that it returns:
>
> This feels much nicer and is also much safer. Some would even say that it's easier to read as well. In the case of a more verbose error message, it certainly would be preferable for a developer to simply read `ErrItemNotFound` rather than a novel on why a certain error has been returned.
>
> This approach is not limited to errors and can be used for other returned values. As an example, we are also returning a `NullItem` instead of `Item{}` as we did before. There are many different scenarios in which it might be preferable to return a defined object, rather than initialising it on return.

### Source: Clean Go Code — Returning Dynamic Errors

> There are certainly some scenarios where returning an error variable might not actually be viable. In cases where the information in custom errors is dynamic, if we want to describe error events more specifically, we can no longer define and return our static errors.
>
> So, what to do? There is no well-defined or standard method for handling and returning these kinds of dynamic errors. My personal preference is to return a new interface, with a bit of added functionality:
>
> This new data structure still works as our standard error. We can still compare it to `nil` since it's an interface implementation, and we can still call `.Error()` on it, so it won't break any existing implementations. However, the advantage is that we can now check our error type as we could previously, despite our error now containing the dynamic details:

### Source: Practical Go — §7 Error handling

> I've given several presentations about error handling and written a lot about error handling on my blog. I also spoke a lot about error handling in yesterday's session so I won't repeat what I've said.
>
> Instead I want to cover two other areas related to error handling.

> #### 7.1. Eliminate error handling by eliminating errors
>
> If you were in my presentation yesterday I talked about the draft proposals for improving error handling. But do you know what is better than an improved syntax for handling errors? Not needing to handle errors at all.
>
> I'm not saying "remove your error handling". What I am suggesting is, change your code so you do not have errors to handle.
>
> This section draws inspiration from John Ousterhout's recently book, A philosophy of Software Design. One of the chapters in that book is called "Define Errors Out of Existence". We're going to try to apply this advice to Go.

> ##### 7.1.1. Counting lines
>
> Let's write a function to count the number of lines in a file.
>
> Because we're following our advice from previous sections, `CountLines` takes an `io.Reader`, not a `*os.File`; its the job of the caller to provide the `io.Reader` who's contents we want to count.
>
> We construct a `bufio.Reader`, and then sit in a loop calling the `ReadString` method, incrementing a counter until we reach the end of the file, then we return the number of lines read.
>
> At least that's the code we want to write, but instead this function is made more complicated by error handling.
>
> But we're not done checking errors yet. `ReadString` will return `io.EOF` when it hits the end of the file. This is expected, `ReadString` needs some way of saying stop, there is nothing more to read. So before we return the error to the caller of `CountLine`, we need to check if the error was not `io.EOF`, and in that case propagate it up, otherwise we return `nil` to say that everything worked fine.
>
> I think this is a good example of Russ Cox's observation that error handling can obscure the operation of the function. Let's look at an improved version.
>
> This improved version switches from using `bufio.Reader` to `bufio.Scanner`.
>
> Under the hood `bufio.Scanner` uses `bufio.Reader`, but it adds a nice layer of abstraction which helps remove the error handling with obscured the operation of `CountLines`.
>
> `bufio.Scanner` can scan for any pattern, but by default it looks for newlines.
>
> The method, `sc.Scan()` returns `true` if the scanner has matched a line of text and has not encountered an error. So, the body of our `for` loop will be called only when there is a line of text in the scanner's buffer. This means our revised `CountLines` correctly handles the case where there is no trailing newline, and also handles the case where the file was empty.
>
> Secondly, as `sc.Scan` returns `false` once an error is encountered, our `for` loop will exit when the end-of-file is reached or an error is encountered. The `bufio.Scanner` type memoises the first error it encountered and we can recover that error once we've exited the loop using the `sc.Err()` method.
>
> Lastly, `sc.Err()` takes care of handling `io.EOF` and will convert it to a `nil` if the end of file was reached without encountering another error.
>
> When you find yourself faced with overbearing error handling, try to extract some of the operations into a helper type.

> ##### 7.1.2. WriteResponse
>
> My second example is inspired from the Errors are values blog post.
>
> Earlier in this presentation We've seen examples dealing with opening, writing and closing files. The error handling is present, but not overwhelming as the operations can be encapsulated in helpers like `os.ReadFile` and `os.WriteFile`. However when dealing with low level network protocols it becomes necessary to build the response directly using I/O primitives the error handling can become repetitive.
>
> First we construct the status line using `fmt.Fprintf`, and check the error. Then for each header we write the header key and value, checking the error each time. Lastly we terminate the header section with an additional `\r\n`, check the error, and copy the response body to the client. Finally, although we don't need to check the error from `io.Copy`, we need to translate it from the two return value form that `io.Copy` returns into the single return value that `WriteResponse` returns.
>
> That's a lot of repetitive work. But we can make it easier on ourselves by introducing a small wrapper type, `errWriter`.
>
> `errWriter` fulfils the `io.Writer` contract so it can be used to wrap an existing `io.Writer`. `errWriter` passes writes through to its underlying writer until an error is detected. From that point on, it discards any writes and returns the previous error.
>
> Applying `errWriter` to `WriteResponse` dramatically improves the clarity of the code. Each of the operations no longer needs to bracket itself with an error check. Reporting the error is moved to the end of the function by inspecting the `ew.err` field, avoiding the annoying translation from `io.Copy's` return values.

> #### 7.2. Only handle an error once
>
> Lastly, I want to mention that you should only handle errors once. Handling an error means inspecting the error value, and making a single decision.
>
> If you make less than one decision, you're ignoring the error. As we see here, the error from `w.WriteAll` is being discarded.
>
> But making more than one decision in response to a single error is also problematic. The following is code that I come across frequently.
>
> In this example if an error occurs during `w.Write`, a line will be written to a log file, noting the file and line that the error occurred, and the error is also returned to the caller, who possibly will log it, and return it, all the way back up to the top of the program.
>
> So you get a stack of duplicate lines in your log file,
>
> `unable to write: io.EOF`
> `could not write config: io.EOF`
>
> but at the top of the program you get the original error without any context.
>
> `err := WriteConfig(f, &conf)`
> `fmt.Println(err) // io.EOF`
>
> I want to dig into this a little further because I don't see the problems with logging and returning as just a matter of personal preference.
>
> The problem I see a lot is programmers forgetting to return from an error. As we talked about earlier, Go style is to use guard clauses, checking preconditions as the function progresses and returning early.
>
> In this example the author checked the error, logged it, but forgot to return. This has caused a subtle bug.
>
> The contract for error handling in Go says that you cannot make any assumptions about the contents of other return values in the presence of an error. As the JSON marshalling failed, the contents of `buf` are unknown, maybe it contains nothing, but worse it could contain a half written JSON fragment.
>
> Because the programmer forgot to return after checking and logging the error, the corrupt buffer will be passed to `WriteAll`, which will probably succeed and so the config file will be written incorrectly. However the function will return just fine, and the only indication that a problem happened will be a single log line complaining about marshalling JSON, not a failure to write the config.

> ##### 7.2.1. Adding context to errors
>
> The bug occurred because the author was trying to add context to the error message. They were trying to leave themselves a breadcrumb to point them back to the source of the error.
>
> Let's look at another way to do the same thing using `fmt.Errorf`.
>
> By combining the annotation of the error with returning onto one line there it is harder to forget to return an error and avoid continuing accidentally.
>
> If an I/O error occurs writing the file, the `error's `Error()` method will report something like this;
>
> `could not write config: write failed: input/output error`

> ##### 7.2.2. Wrapping errors with github.com/pkg/errors
>
> The `fmt.Errorf` pattern works well for annotating the error message, but it does so at the cost of obscuring the type of the original error. I've argued that treating errors as opaque values is important to producing software which is loosely coupled, so the face that the type of the original error should not matter if the only thing you do with an error value is
>
> Check that it is not `nil`.
>
> Print or log it.
>
> However there are some cases, I believe they are infrequent, where you do need to recover the original error. In that case you can use something like my errors package to annotate errors.
>
> Now the error reported will be the nice K&D style error,
>
> `could not read config: open failed: open /Users/dfc/.settings.xml: no such file or directory`
>
> and the error value retains a reference to the original cause.
>
> Using the `errors` package gives you the ability to add context to error values, in a way that is inspectable by both a human and a machine. If you came to my presentation yesterday you'll know that wrapping is moving into the standard library in an upcoming Go release.
