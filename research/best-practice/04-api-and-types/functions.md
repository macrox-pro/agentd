---
primary_sources:
  - id: T1-CRC
    title: "Go Code Review Comments"
    url: "https://go.dev/wiki/CodeReviewComments"
    author: "Go team"
    section: "Named Result Parameters"
  - id: T1-CRC
    title: "Go Code Review Comments"
    url: "https://go.dev/wiki/CodeReviewComments"
    author: "Go team"
    section: "Naked Returns"
  - id: T1-CRC
    title: "Go Code Review Comments"
    url: "https://go.dev/wiki/CodeReviewComments"
    author: "Go team"
    section: "Pass Values"
  - id: T1-CRC
    title: "Go Code Review Comments"
    url: "https://go.dev/wiki/CodeReviewComments"
    author: "Go team"
    section: "Synchronous Functions"
  - id: T3-CLEAN
    title: "Clean Go Code"
    url: "https://github.com/Pungyeon/clean-go-article"
    author: "Pungyeon"
    section: "Cleaning Functions"
  - id: T2-PG
    title: "Practical Go"
    url: "https://dave.cheney.net/practical-go/presentations/qcon-china.html"
    author: "Dave Cheney"
    section: "§6 API Design (functions)"
  - id: T2-UBER
    title: "Uber Go Style Guide"
    url: "https://github.com/uber-go/guide"
    author: "Uber"
    section: "Defer, Time"
studied_at: "2026-08-25"
go_min: "1.26.7"
applicability: "current"
---
# Functions

> **Applicability (Go >= 1.26.7):** Guidance below is current for Go >= 1.26.7.

### Source: Go Code Review Comments — Named Result Parameters

> Consider what it will look like in godoc. Named result parameters like:
>
> ```
> func (n *Node) Parent1() (node *Node) {}
> func (n *Node) Parent2() (node *Node, err error) {}
> ```
>
> will be repetitive in godoc; better to use:
>
> ```
> func (n *Node) Parent1() *Node {}
> func (n *Node) Parent2() (*Node, error) {}
> ```
>
> On the other hand, if a function returns two or three parameters of the same type, or if the meaning of a result isn't clear from context, adding names may be useful in some contexts. Don't name result parameters just to avoid declaring a var inside the function; that trades off a minor implementation brevity at the cost of unnecessary API verbosity.
>
> ```
> func (f *Foo) Location() (float64, float64, error)
> ```
>
> is less clear than:
>
> ```
> // Location returns f's latitude and longitude.
> // Negative values mean south and west, respectively.
> func (f *Foo) Location() (lat, long float64, err error)
> ```
>
> Naked returns are okay if the function is a handful of lines. Once it's a medium sized function, be explicit with your return values. Corollary: it's not worth it to name result parameters just because it enables you to use naked returns. Clarity of docs is always more important than saving a line or two in your function.
>
> Finally, in some cases you need to name a result parameter in order to change it in a deferred closure. That is always OK.

### Source: Go Code Review Comments — Naked Returns

> A `return` statement without arguments returns the named return values. This is known as a "naked" return.
>
> ```
> func split(sum int) (x, y int) {
>     x = sum * 4 / 9
>     y = sum - x
>     return
> }
> ```
>
> See Named Result Parameters.

### Source: Go Code Review Comments — Pass Values

> Don't pass pointers as function arguments just to save a few bytes. If a function refers to its argument `x` only as `*x` throughout, then the argument shouldn't be a pointer. Common instances of this include passing a pointer to a string (`*string`) or a pointer to an interface value (`*io.Reader`). In both cases the value itself is a fixed size and can be passed directly. This advice does not apply to large structs, or even small structs that might grow.

### Source: Go Code Review Comments — Synchronous Functions

> Prefer synchronous functions - functions which return their results directly or finish any callbacks or channel ops before returning - over asynchronous ones.
>
> Synchronous functions keep goroutines localized within a call, making it easier to reason about their lifetimes and avoid leaks and data races. They're also easier to test: the caller can pass an input and check the output without the need for polling or synchronization.
>
> If callers need more concurrency, they can add it easily by calling the function from a separate goroutine. But it is quite difficult - sometimes impossible - to remove unnecessary concurrency at the caller side.

### Source: Clean Go Code — Function Length

> How small should a function be? Smaller than that! – Robert C. Martin
>
> When writing clean code, our primary goal is to make our code easily digestible. The most effective way to do this is to make our functions as short as possible. It's important to understand that we don't necessarily do this to avoid code duplication. The more important reason is to improve code comprehension.
>
> By writing short functions (which are typically 5–8 lines in Go), we can create code that reads almost as naturally as our description above:
>
> Using smaller functions also eliminates another horrible habit of writing code: indentation hell. Indentation hell typically occurs when a chain of `if` statements are carelessly nested in a function. This makes it very difficult for human beings to parse the code and should be eliminated whenever spotted.
>
> So, how do we clean this function? Fortunately, it's actually quite simple. On our first iteration, we will try to ensure that we are returning an error as soon as possible. Instead of nesting the `if` and `else` statements, we want to "push our code to the left," so to speak.
>
> If the `value, err :=` pattern is repeated more than once in a function, this is an indication that we can split the logic of our code into smaller pieces:
>
> Notice that cleaning the `GetItem` function resulted in more lines of code overall. However, the code itself is now much easier to read. It's layered in an onion-style fashion, where we can ignore "layers" that we aren't interested in and simply peel back the ones that we do want to examine.

### Source: Clean Go Code — Function Signatures

> Creating a good function naming structure makes it easier to read and understand the intent of the code. As we saw above, making our functions shorter helps us understand the function's logic. The last part of cleaning our functions involves understanding the context of the function input. With this comes another easy-to-follow rule: Function signatures should only contain one or two input parameters. In certain exceptional cases, three can be acceptable, but this is where we should start considering a refactor. Much like the rule that our functions should only be 5–8 lines long, this can seem quite extreme at first. However, I feel that this rule is much easier to justify.
>
> Therefore, it is recommended to replace these input parameters with an 'Options' `struct` instead:
>
> This solves two problems: misusing comments, and accidentally labeling the variables incorrectly. Of course, we can still confuse properties with the wrong value, but in these cases, it will be much easier to determine where our mistake lies within the code. The ordering of the properties also doesn't matter anymore, so incorrectly ordering the input values is no longer a concern. The last added bonus of this technique is that we can use our `QueueOptions` struct to infer the default values of our function's input parameters.

### Source: Practical Go — §6 API Design

> The last piece of design advice I'm going to give today I feel is the most important.
>
> All of the suggestions I've made so far are just that, suggestions. These are the way I try to write my Go, but I'm not going to push them hard in code review.
>
> However when it comes to reviewing APIs during code review, I am less forgiving. This is because everything I've talked about so far can be fixed without breaking backward compatibility; they are, for the most part, implementation details.
>
> When it comes to the public API of a package, it pays to put considerable thought into the initial design, because changing that design later is going to be disruptive for people who are already using your API.

> #### 6.1. Design APIs that are hard to misuse.
>
> APIs should be easy to use and hard to misuse.
>
> — Josh Bloch
>
> If you take anything away from this presentation, it should be this advice from Josh Bloch. If an API is hard to use for simple things, then every invocation of the API will look complicated. When the actual invocation of the API is complicated it will be less obvious and more likely to be overlooked.

> ##### 6.1.1. Be wary of functions which take several parameters of the same type
>
> A good example of a simple looking, but hard to use correctly API is one which takes two or more parameters of the same type. Let's compare two function signatures:
>
> `func Max(a, b int) int`
> `func CopyFile(to, from string) error`
>
> Max is commutative; the order of its parameters does not matter. The maximum of eight and ten is ten regardless of if I compare eight and ten or ten and eight.
>
> However, this property does not hold true for `CopyFile`.
>
> Which one of these statements made a backup of your presentation and which one overwrite your presentation with last week's version? You can't tell without consulting the documentation. A code reviewer cannot know if you've got the order correct without consulting the documentation.
>
> APIs with multiple parameters of the same type are hard to use correctly.

> #### 6.2. Design APIs for their default use case
>
> A few years ago I gave a talk about using functional options to make APIs easier to use for their default case.
>
> The gist of this talk was you should design your APIs for the common use case. Said another way, your API should not require the caller to provide parameters which they don't care about.

> ##### 6.2.1. Discourage the use of nil as a parameter
>
> I opened this chapter with the suggestion that you shouldn't force the caller of your API into providing you parameters when they don't really care what those parameters mean. This is what I mean when I say design APIs for their default use case.
>
> `ListenAndServe` takes two parameters, a TCP address to listen for incoming connections, and `http.Handler` to handle the incoming HTTP request. `Serve` allows the second parameter to be `nil`, and notes that usually the caller will pass `nil` indicating that they want to use `http.DefaultServeMux` as the implicit parameter.
>
> Now the caller of `Serve` has two ways to do the same thing.
>
> This `nil` behaviour is viral. The `http` package also has a `http.Serve` helper, which you can reasonably imagine that `ListenAndServe` builds upon.
>
> Because `ListenAndServe` permits the caller to pass `nil` for the second parameter, `http.Serve` also supports this behaviour. In fact, `http.Serve` is the one that implements the "if `handler` is `nil`, use `DefaultServeMux`" logic.
>
> Accepting `nil` for one parameter may lead the caller into thinking they can pass `nil` for both parameters. However calling `Serve` like `http.Serve(nil, nil)` results in an ugly panic.
>
> Don't mix `nil` and non `nil`-able parameters in the same function signature.
>
> The author of `http.ListenAndServe` was trying to make the API user's life easier in the common case, but possibly made the package harder to use safely.
>
> There is no difference in line count between using `DefaultServeMux` explicitly, or implicitly via `nil`.
>
> Give serious consideration to how much time helper functions will save the programmer. Clear is better than concise.
>
> Avoid public APIs with test only parameters Avoid exposing APIs with values who only differ in test scope. Instead, use Public wrappers to hide those parameters, use test scoped helpers to set the property in test scope.

> ##### 6.2.2. Prefer var args to []T parameters
>
> It's very common to write a function or method that takes a slice of values.
>
> `func ShutdownVMs(ids []string) error`
>
> This is just an example I made up, but its common to a lot of code I've worked on. The problem with signatures like these is they presume that they will be called with more than one entry. However, what I have found is many times these type of functions are called with only one argument, which has to be "boxed" inside a slice just to meet the requirements of the functions signature.
>
> Additionally, because the `ids` parameter is a slice, you can pass an empty slice or `nil` to the function and the compiler will be happy. This adds extra testing load because you should cover these cases in your testing.
>
> However there is a problem with `anyPositive`, someone could accidentally invoke it like `if anyPositive() { ... }`
>
> In this case `anyPositive` would return `false` because it would execute zero iterations and immediately return `false`. This isn't the worst thing in the world — that would be if `anyPositive` returned `true` when passed no arguments.
>
> Nevertheless it would be be better if we could change the signature of `anyPositive` to enforce that the caller should pass at least one argument. We can do that by combining normal and vararg parameters like `func anyPositive(first int, rest ...int) bool`.

### Source: Uber Go Style Guide — Defer to Clean Up

> Use defer to clean up resources such as files and locks.
>
> Defer has an extremely small overhead and should be avoided only if you can prove that your function execution time is in the order of nanoseconds. The readability win of using defers is worth the miniscule cost of using them. This is especially true for larger methods that have more than simple memory accesses, where the other computations are more significant than the `defer`.

### Source: Uber Go Style Guide — Use `"time"` to handle time

> Time is complicated. Incorrect assumptions often made about time include the following.
>
> 1. A day has 24 hours
> 2. An hour has 60 minutes
> 3. A week has 7 days
> 4. A year has 365 days
>
> For example, *1* means that adding 24 hours to a time instant will not always yield a new calendar day.
>
> Therefore, always use the `"time"` package when dealing with time because it helps deal with these incorrect assumptions in a safer, more accurate manner.
>
> Use `time.Time` when dealing with instants of time, and the methods on `time.Time` when comparing, adding, or subtracting time.
>
> Use `time.Duration` when dealing with periods of time.
>
> Going back to the example of adding 24 hours to a time instant, the method we use to add time depends on intent. If we want the same time of the day, but on the next calendar day, we should use `Time.AddDate`. However, if we want an instant of time guaranteed to be 24 hours after the previous time, we should use `Time.Add`.
>
> Use `time.Duration` and `time.Time` in interactions with external systems when possible.
>
> When it is not possible to use `time.Duration` in these interactions, use `int` or `float64` and include the unit in the name of the field.
>
> When it is not possible to use `time.Time` in these interactions, unless an alternative is agreed upon, use `string` and format timestamps as defined in RFC 3339.
>
> Although this tends to not be a problem in practice, keep in mind that the `"time"` package does not support parsing timestamps with leap seconds, nor does it account for leap seconds in calculations. If you compare two instants of time, the difference will not include the leap seconds that may have occurred between those two instants.
