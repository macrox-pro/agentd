---
primary_sources:
  - id: T1-CRC
    title: "Go Code Review Comments"
    url: "https://go.dev/wiki/CodeReviewComments"
    author: "Go team"
    section: "Interfaces"
  - id: T3-CLEAN
    title: "Clean Go Code"
    url: "https://github.com/Pungyeon/clean-go-article"
    author: "Pungyeon"
    section: "Interfaces in Go; Nil Values"
  - id: T3-IDIO
    title: "Idiomatic Go"
    url: "https://dmitri.shuralyov.com/idiomatic-go"
    author: "Dmitri Shuralyov"
    section: "Avoid unused method receiver names; Mutex hat"
  - id: T2-PG
    title: "Practical Go"
    url: "https://dave.cheney.net/practical-go/presentations/qcon-china.html"
    author: "Dave Cheney"
    section: "§6.3 Let functions define the behaviour they requires"
  - id: T2-UBER
    title: "Uber Go Style Guide"
    url: "https://github.com/uber-go/guide"
    author: "Uber"
    section: "Pointers to Interfaces; Verify Interface Compliance; Receivers and Interfaces"
studied_at: "2026-08-25"
go_min: "1.26.7"
applicability: "current"
---
# Interfaces

> **Applicability (Go >= 1.26.7):** Guidance below is current for Go >= 1.26.7.

### Source: Go Code Review Comments — Interfaces

> Go interfaces generally belong in the package that uses values of the interface type, not the package that implements those values. The implementing package should return concrete (usually pointer or struct) types: that way, new methods can be added to implementations without requiring extensive refactoring.
>
> Do not define interfaces on the implementor side of an API "for mocking"; instead, design the API so that it can be tested using the public API of the real implementation.
>
> Do not define interfaces before they are used: without a realistic example of usage, it is too difficult to see whether an interface is even necessary, let alone what methods it ought to contain.
>
> ```
> package consumer  // consumer.go
>
> type Thinger interface { Thing() bool }
>
> func Foo(t Thinger) string { … }
> ```
>
> ```
> package consumer // consumer_test.go
>
> type fakeThinger struct{ … }
> func (t fakeThinger) Thing() bool { … }
> …
> if Foo(fakeThinger{…}) == "x" { … }
> ```
>
> ```
> // DO NOT DO IT!!!
> package producer
>
> type Thinger interface { Thing() bool }
>
> type defaultThinger struct{ … }
> func (t defaultThinger) Thing() bool { … }
>
> func NewThinger() Thinger { return defaultThinger{ … } }
> ```
>
> Instead return a concrete type and let the consumer mock the producer implementation.
>
> ```
> package producer
>
> type Thinger struct{ … }
> func (t Thinger) Thing() bool { … }
>
> func NewThinger() Thinger { return Thinger{ … } }
> ```

### Source: Clean Go Code — Interfaces in Go

> In general, Go's approach to handling `interface`s is quite different from those of other languages. Interfaces aren't explicitly implemented like they would be in Java or C#; rather, they are implicitly created if they fulfill the contract of the interface. As an example, this means that any `struct` that has an `Error()` method implements (or "fulfills") the `Error` interface and can be returned as an `error`. This manner of implementing interfaces is extremely easy and makes Go feel more fast paced and dynamic.
>
> However, there are certainly disadvantages with this approach. As the interface implementation is no longer explicit, it can be difficult to see which interfaces are implemented by a struct. Therefore, it's common to define interfaces with as few methods as possible; this makes it easier to understand whether a particular struct fulfills the contract of the interface.
>
> In the above code, we are initialising a variable with the Go `blank identifier`, with the type assignment of `io.Writer`. This results in our variable being checked to fulfill the `io.Writer` interface contract, before being discarded. This method of checking interface fulfillment also makes it possible to check that several interface contracts are fulfilled:
>
> From the above code, it's very easy to understand which interfaces must be fulfilled; this ensures that the compiler will help us out during compile time. Therefore, this is generally the preferred solution for checking interface contract fulfillment.
>
> In other words, you should write functions that accept an interface and return a concrete type. This is generally good practice and is especially useful when doing tests with mocking.
>
> This is why it is encouraged to make interfaces as small as possible in idiomatic Go—it makes it especially easy to implement patterns like the one we just saw. However, this implementation of interfaces also comes with a huge downside.

### Source: Clean Go Code — Nil Values

> A controversial aspect of Go is the addition of `nil`. This value corresponds to the value `NULL` in C and is essentially an uninitialised pointer. We've already seen some of the problems that `nil` can cause, but to sum up: Things break when you try to access methods or properties of a `nil` value. Thus, it's recommended to avoid returning a `nil` value when possible. This way, the users of our code are less likely to accidentally access `nil` values.
>
> There are other scenarios in which it is common to find `nil` values that can cause some unnecessary pain. An example of this is incorrectly initialising a `struct` (as in the example below), which can lead to it containing `nil` properties. If accessed, those `nil`s will cause a panic.
>
> Instead, we can turn the `Cache` property of our `App` structure into a private property and create a getter-like method to access it. This gives us more control over what we are returning; specifically, it ensures that we aren't returning a `nil` value:
>
> This ensures that users of our package don't have to worry about the implementation and whether they're using our package in an unsafe manner. All they need to worry about is writing their own clean code.

### Source: Idiomatic Go — Avoid unused method receiver names

> Do this:
>
> ```
> func (foo) method() {
> 	...
> }
> ```
>
> Don't do this:
>
> ```
> func (f foo) method() {
> 	...
> }
> ```
>
> If `f` is unused. It's more readable because it's clear that fields or methods of `foo` are not used in `method`.

### Source: Idiomatic Go — Mutex hat

> Here,`rateMu` is a mutex hat. It sits, like a hat, on top of the variables that it protects.
>
> So, without needing to write the comment, the above is implicitly understood to be equivalent to:
>
> ```
> 	// rateMu protects rateLimits and mostRecent.
> 	rateMu     sync.Mutex
> 	rateLimits [categories]Rate
> 	mostRecent rateLimitCategory
> ```
>
> When adding a new, unrelated field that isn't protected by `rateMu`, do this:
>
> ```
>  struct {
>  	...
>
>  	rateMu     sync.Mutex
>  	rateLimits [categories]Rate
>  	mostRecent rateLimitCategory
> +
> +	common service
>  }
> ```
>
> Don't do this:
>
> ```
>  struct {
>  	...
>
>  	rateMu     sync.Mutex
>  	rateLimits [categories]Rate
>  	mostRecent rateLimitCategory
> +	common     service
>  }
> ```

### Source: Practical Go — §6.3 Let functions define the behaviour they requires

> #### 6.3. Let functions define the behaviour they requires
>
> Let's say I've been given a task to write a function that persists a Document structure to disk.
>
> `func Save(f *os.File, doc *Document) error`
>
> I could specify this function, Save, which takes an `*os.File` as the destination to write the `Document`. But this has a few problems:
>
> The signature of `Save` precludes the option to write the data to a network location. Assuming that network storage is likely to become requirement later, the signature of this function would have to change, impacting all its callers.
>
> `Save` is also unpleasant to test, because it operates directly with files on disk. So, to verify its operation, the test would have to read the contents of the file after being written.
>
> And I would have to ensure that `f` was written to a temporary location and always removed afterwards.
>
> `*os.File` also defines a lot of methods which are not relevant to `Save`, like reading directories and checking to see if a path is a symlink. It would be useful if the signature of the `Save` function could describe only the parts of `*os.File` that were relevant.
>
> Using `io.ReadWriteCloser` we can apply the interface segregation principle to redefine `Save` to take an interface that describes more general file shaped things.
>
> With this change, any type that implements the `io.ReadWriteCloser` interface can be substituted for the previous `*os.File`.
>
> This makes `Save` both broader in its application, and clarifies to the caller of `Save` which methods of the `*os.File` type are relevant to its operation.
>
> And as the author of `Save` I no longer have the option to call those unrelated methods on `*os.File` as it is hidden behind the `io.ReadWriteCloser` interface.
>
> But we can take the interface segregation principle a bit further.
>
> Firstly, it is unlikely that if `Save` follows the single responsibility principle, it will read the file it just wrote to verify its contents—that should be responsibility of another piece of code.
>
> `func Save(wc io.WriteCloser, doc *Document) error`
>
> So we can narrow the specification for the interface we pass to Save to just writing and closing.
>
> Secondly, by providing `Save` with a mechanism to close its stream, which we inherited in this desire to make it still look like a file, this raises the question of under what circumstances will `wc` be closed.
>
> A better solution would be to redefine `Save` to take only an `io.Writer`, stripping it completely of the responsibility to do anything but write data to a stream.
>
> By applying the interface segregation principle to our `Save` function, the results has simultaneously been a function which is the most specific in terms of its requirements—it only needs a thing that is writable—and the most general in its function, we can now use `Save` to save our data to anything which implements `io.Writer`.

### Source: Uber Go Style Guide — Pointers to Interfaces

> You almost never need a pointer to an interface. You should be passing interfaces as values—the underlying data can still be a pointer.
>
> An interface is two fields:
>
> 1. A pointer to some type-specific information. You can think of this as "type."
> 2. Data pointer. If the data stored is a pointer, it's stored directly. If the data stored is a value, then a pointer to the value is stored.
>
> If you want interface methods to modify the underlying data, you must use a pointer.

### Source: Uber Go Style Guide — Verify Interface Compliance

> Verify interface compliance at compile time where appropriate. This includes:
>
> - Exported types that are required to implement specific interfaces as part of their API contract
> - Exported or unexported types that are part of a collection of types implementing the same interface
> - Other cases where violating an interface would break users
>
> The statement `var _ http.Handler = (*Handler)(nil)` will fail to compile if `*Handler` ever stops matching the `http.Handler` interface.
>
> The right hand side of the assignment should be the zero value of the asserted type. This is `nil` for pointer types (like `*Handler`), slices, and maps, and an empty struct for struct types.

### Source: Uber Go Style Guide — Receivers and Interfaces

> Methods with value receivers can be called on pointers as well as values. Methods with pointer receivers can only be called on pointers or addressable values.
>
> We cannot get pointers to values stored in maps, because they are not addressable values.
>
> We can call Read on values stored in the map because Read has a value receiver, which does not require the value to be addressable.
>
> We cannot call Write on values stored in the map because Write has a pointer receiver, and it's not possible to get a pointer to a value stored in a map.
>
> You can call both Read and Write if the map stores pointers, because pointers are intrinsically addressable.
>
> Similarly, an interface can be satisfied by a pointer, even if the method has a value receiver.
>
> The following doesn't compile, since s2Val is a value, and there is no value receiver for f.
