---
primary_sources:
  - id: T2-UBER
    title: "Uber Go Style Guide"
    url: "https://github.com/uber-go/guide"
    author: "Uber"
    section: "Channel Size, Zero-value Mutexes"
  - id: T2-PG
    title: "Practical Go"
    url: "https://dave.cheney.net/practical-go/presentations/qcon-china.html"
    author: "Dave Cheney"
    section: "§8 Concurrency"
  - id: T2-ADV
    title: "go-advices"
    url: "https://github.com/cristaloleg/go-advices"
    author: "Oleg Kiselyov"
    section: "Concurrency"
studied_at: "2026-08-25"
go_min: "1.26.7"
applicability: "current"
---
# Channels and synchronization

> **Applicability (Go >= 1.26.7):** Guidance below is current for Go >= 1.26.7.

### Source: Uber Go Style Guide — Channel Size is One or None

> Channels should usually have a size of one or be unbuffered. By default, channels are unbuffered and have a size of zero. Any other size must be subject to a high level of scrutiny. Consider how the size is determined, what prevents the channel from filling up under load and blocking writers, and what happens when this occurs.

### Source: Uber Go Style Guide — Zero-value Mutexes are Valid

> The zero-value of `sync.Mutex` and `sync.RWMutex` is valid, so you almost never need a pointer to a mutex.
>
> If you use a struct by pointer, then the mutex should be a non-pointer field on it. Do not embed the mutex on the struct, even if the struct is not exported.

### Source: Practical Go — §8 Concurrency

> #### 8.1. Keep yourself busy or do the work yourself
>
> So this is my first piece of advice: if your goroutine cannot make progress until it gets the result from another, oftentimes it is simpler to just do the work yourself rather than to delegate it.
>
> This often eliminates a lot of state tracking and channel manipulation required to plumb a result back from a goroutine to its initiator.
>
> Many Go programmers overuse goroutines, especially when they are starting out. As with all things in life, moderation is the key the key to success.

> #### 8.2. Leave concurrency to the caller
>
> If your function starts a goroutine you must provide the caller with a way to explicitly stop that goroutine. It is often easier to leave decision to execute a function asynchronously to the caller of that function.

> #### 8.3. Never start a goroutine without knowning when it will stop
>
> Just as functions in Go leave concurrency to the caller, applications should leave the job of monitoring their status and restarting them if they fail to the program that invoked them. Do not make your applications responsible for restarting themselves, this is a procedure best handled from outside the application.
>
> Only use `log.Fatal` from `main.main` or `init` functions.
>
> We can use a channel to collect the return status of the goroutine. The size of the channel is equal to the number of goroutines we want to manage so that sending to the `done` channel will not block, as this will block the shutdown the of goroutine, causing it to leak.
>
> As there is no way to safely close the `done` channel we cannot use the `for range` idiom to loop of the channel until all goroutines have reported in, instead we loop for as many goroutines we started, which is equal to the capacity of the channel.
>
> Writing this logic yourself is repetitive and subtle. Consider something like this package, https://github.com/heptio/workgroup which will do most of the work for you.

### Source: Uber Go Style Guide — Don't fire-and-forget goroutines

> Goroutines are lightweight, but they're not free: at minimum, they cost memory for their stack and CPU to be scheduled. While these costs are small for typical uses of goroutines, they can cause significant performance issues when spawned in large numbers without controlled lifetimes. Goroutines with unmanaged lifetimes can also cause other issues like preventing unused objects from being garbage collected and holding onto resources that are otherwise no longer used.
>
> Therefore, do not leak goroutines in production code. Use go.uber.org/goleak to test for goroutine leaks inside packages that may spawn goroutines.
>
> In general, every goroutine:
>
> - must have a predictable time at which it will stop running; or
> - there must be a way to signal to the goroutine that it should stop
>
> In both cases, there must be a way for code to block and wait for the goroutine to finish.

> #### Wait for goroutines to exit
>
> Given a goroutine spawned by the system, there must be a way to wait for the goroutine to exit. There are two popular ways to do this:
>
> - Use a `sync.WaitGroup` to wait for multiple goroutines to complete. Do this if there are multiple goroutines that you want to wait for.
> - Add another `chan struct{}` that the goroutine closes when it's done. Do this if there's only one goroutine.

> #### No goroutines in `init()`
>
> `init()` functions should not spawn goroutines.
>
> If a package has need of a background goroutine, it must expose an object that is responsible for managing a goroutine's lifetime. The object must provide a method (`Close`, `Stop`, `Shutdown`, etc) that signals the background goroutine to stop, and waits for it to exit.

### Source: go-advices — Concurrency

> best candidate to make something once in a thread-safe way is `sync.Once` — don't use flags, mutexes, channels or atomics
>
> to block forever use `select{}`, omit channels, waiting for a signal
>
> don't close in-channel, this is a responsibility of it's creator — writing to a closed channel will cause a panic
>
> when you need an atomic value of a custom type use atomic.Value

### Source: go-advices — Code (signal channels)

> To pass a signal prefer `chan struct{}` instead of `chan bool`.
>
> When you see a definition of `chan bool` in a structure, sometimes it's not that easy to understand how this value will be used. But we can make it more clear by changing it to `chan struct{}` which explicitly says: we do not care about value (it's always a `struct{}`), we care about an event that might occur.
