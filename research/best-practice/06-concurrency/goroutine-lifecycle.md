---
primary_sources:
  - id: T1-CRC
    title: "Go Code Review Comments"
    url: "https://go.dev/wiki/CodeReviewComments"
    author: "Go team"
    section: "Goroutine Lifetimes"
  - id: T2-PG
    title: "Practical Go"
    url: "https://dave.cheney.net/practical-go/presentations/qcon-china.html"
    author: "Dave Cheney"
    section: "§8 Concurrency"
  - id: T2-UBER
    title: "Uber Go Style Guide"
    url: "https://github.com/uber-go/guide"
    author: "Uber"
    section: "Don't fire-and-forget goroutines"
studied_at: "2026-08-25"
also_cited_in: []
go_min: "1.26.7"
applicability: "current"
---
# Goroutine lifecycle

> **Applicability (Go >= 1.26.7):** Guidance below is current for Go >= 1.26.7.

### Source: Go Code Review Comments — Goroutine Lifetimes

> When you spawn goroutines, make it clear when - or whether - they exit.
>
> Goroutines can leak by blocking on channel sends or receives: the garbage collector will not terminate a goroutine even if the channels it is blocked on are unreachable.
>
> Even when goroutines do not leak, leaving them in-flight when they are no longer needed can cause other subtle and hard-to-diagnose problems. Sends on closed channels panic. Modifying still-in-use inputs "after the result isn't needed" can still lead to data races. And leaving goroutines in-flight for arbitrarily long can lead to unpredictable memory usage.
>
> Try to keep concurrent code simple enough that goroutine lifetimes are obvious. If that just isn't feasible, document when and why the goroutines exit.

### Source: Practical Go — §8 Concurrency (overview)

> Often Go is chosen for a project because of its concurrency features. The Go team have gone to great lengths to make concurrency in Go cheap (in terms of hardware resources) and performant, however it is possible to use Go's concurrency features to write code which is neither performent or reliable. With the time I have left I want to leave you with some advice for avoid some of the pitfalls that come with Go's concurrency features.
>
> Go features first class support for concurrency with channels, and the `select` and `go` statements. If you've learnt Go formally from a book or training course, you might have noticed that the concurrency section is always one of the last you'll cover. This workshop is no different, I have chosen to cover concurrency last, as if it is somehow additional to the regular the skills a Go programmer should master.
>
> There is a dichotomy here; Go's headline feature is our simple, lightweight concurrency model. As a product, our language almost sells itself on this on feature alone. On the other hand, there is a narrative that concurrency isn't actually that easy to use, otherwise authors wouldn't make it the last chapter in their book and we wouldn't look back on our formative efforts with regret.
>
> This section discusses some pitfalls of naive usage of Go's concurrency features.

> #### 8.1. Keep yourself busy or do the work yourself
>
> The program does what we intended, it serves a simple web server. However it also does something else at the same time, it wastes CPU in an infinite loop. This is because the `for{}` on the last line of `main` is going to block the main goroutine because it doesn't do any IO, wait on a lock, send or receive on a channel, or otherwise communicate with the scheduler.
>
> As the Go runtime is mostly cooperatively scheduled, this program is going to spin fruitlessly on a single CPU, and may eventually end up live-locked.
>
> An empty select statement will block forever. This is a useful property because now we're not spinning a whole CPU just to call `runtime.GoSched()`. However, we're only treating the symptom, not the cause.
>
> I want to present to you another solution, one which has hopefully already occurred to you. Rather than run `http.ListenAndServe` in a goroutine, leaving us with the problem of what to do with the main goroutine, simply run `http.ListenAndServe` on the main goroutine itself.
>
> If the `main.main` function of a Go program returns then the Go program will unconditionally exit no matter what other goroutines started by the program over time are doing.
>
> So this is my first piece of advice: if your goroutine cannot make progress until it gets the result from another, oftentimes it is simpler to just do the work yourself rather than to delegate it.
>
> This often eliminates a lot of state tracking and channel manipulation required to plumb a result back from a goroutine to its initiator.
>
> Many Go programmers overuse goroutines, especially when they are starting out. As with all things in life, moderation is the key the key to success.

> #### 8.2. Leave concurrency to the caller
>
> What is the difference between these two APIs?
>
> `func ListDirectory(dir string) ([]string, error)`
>
> `func ListDirectory(dir string) chan string`
>
> Firstly, the obvious differences; the first example reads a directory into a slice then returns the whole slice, or an error if something went wrong. This happens synchronously, the caller of `ListDirectory` blocks until all directory entries have been read. Depending on how large the directory, this could take a long time, and could potentially allocate a lot of memory building up the slide of directory entry names.
>
> Lets look at the second example. This is a little more Go like, `ListDirectory` returns a channel over which directory entries will be passed. When the channel is closed, that is your indication that there are no more directory entries. As the population of the channel happens after `ListDirectory` returns, `ListDirectory` is probably starting a goroutine to populate the channel.
>
> Its not necessary for the second version to actually use a Go routine; it could allocate a channel sufficient to hold all the directory entries without blocking, fill the channel, close it, then return the channel to the caller. But this is unlikely, as this would have the same problems with consuming a large amount of memory to buffer all the results in a channel.
>
> The channel version of `ListDirectory` has two further problems:
>
> By using a closed channel as the signal that there are no more items to process there is no way for `ListDirectory` to tell the caller that the set of items returned over the channel is incomplete because an error was encountered partway through. There is no way for the caller to tell the difference between an empty directory and an error to read from the directory entirely. Both result in a channel returned from `ListDirectory` which appears to be closed immediately.
>
> The caller must continue to read from the channel until it is closed because that is the only way the caller can know that the goroutine which was started to fill the channel has stopped. This is a serious limitation on the use of `ListDirectory`, the caller has to spend time reading from the channel even though it may have received the answer it wanted. It is probably more efficient in terms of memory usage for medium to large directories, but this method is no faster than the original slice based method.
>
> The solution to the problems of both implementations is to use a callback, a function that is called in the context of each directory entry as it is executed.
>
> `func ListDirectory(dir string, fn func(string))`
>
> Not surprisingly this is how the `filepath.WalkDir` function works.
>
> If your function starts a goroutine you must provide the caller with a way to explicitly stop that goroutine. It is often easier to leave decision to execute a function asynchronously to the caller of that function.

> #### 8.3. Never start a goroutine without knowning when it will stop
>
> The previous example showed using a goroutine when one wasn't really necessary. But one of the driving reasons for using Go is the first class concurrency features the language offers. Indeed there are many instances where you want to exploit the parallelism available in your hardware. To do so, you must use goroutines.
>
> By breaking the `serveApp` and `serveDebug` handlers out into their own functions we've decoupled them from `main.main`. We've also followed the advice from above and make sure that `serveApp` and `serveDebug` leave their concurrency to the caller.
>
> But there are some operability problems with this program. If `serveApp` returns then `main.main` will return causing the program to shutdown and be restarted by whatever process manager you're using.
>
> Just as functions in Go leave concurrency to the caller, applications should leave the job of monitoring their status and restarting them if they fail to the program that invoked them. Do not make your applications responsible for restarting themselves, this is a procedure best handled from outside the application.
>
> However, `serveDebug` is run in a separate goroutine and if it returns just that goroutine will exit while the rest of the program continues on. Your operations staff will not be happy to find that they cannot get the statistics out of your application when they want too because the `/debug` handler stopped working a long time ago.
>
> What we want to ensure is that if any of the goroutines responsible for serving this application stop, we shut down the application.
>
> Now `serverApp` and `serveDebug` check the error returned from `ListenAndServe` and call `log.Fatal` if required. Because both handlers are running in goroutines, we park the main goroutine in a `select{}`.
>
> This approach has a number of problems:
>
> If `ListenAndServer` returns with a `nil` error, `log.Fatal` won't be called and the HTTP service on that port will shut down without stopping the application.
>
> `log.Fatal` calls `os.Exit` which will unconditionally exit the program; defers won't be called, other goroutines won't be notified to shut down, the program will just stop. This makes it difficult to write tests for those functions.
>
> Only use `log.Fatal` from `main.main` or `init` functions.
>
> What we'd really like is to pass any error that occurs back to the originator of the goroutine so that it can know why the goroutine stopped, can shut down the process cleanly.
>
> We can use a channel to collect the return status of the goroutine. The size of the channel is equal to the number of goroutines we want to manage so that sending to the `done` channel will not block, as this will block the shutdown the of goroutine, causing it to leak.
>
> As there is no way to safely close the `done` channel we cannot use the `for range` idiom to loop of the channel until all goroutines have reported in, instead we loop for as many goroutines we started, which is equal to the capacity of the channel.
>
> Now we have a way to wait for each goroutine to exit cleanly and log any error they encounter. All that is needed is a way to forward the shutdown signal from the first goroutine that exits to the others.
>
> It turns out that asking a `http.Server` to shut down is a little involved, so I've spun that logic out into a helper function. The `serve` helper takes an address and `http.Handler`, similar to `http.ListenAndServe`, and also a `stop` channel which we use to trigger the `Shutdown` method.
>
> Now, each time we receive a value on the `done` channel, we close the `stop` channel which causes all the goroutines waiting on that channel to shut down their `http.Server`. This in turn will cause all the remaining `ListenAndServe` goroutines to return. Once all the goroutines we started have stopped, `main.main` returns and the process stops cleanly.
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
>
> There's no way to stop this goroutine. This will run until the application exits.
>
> This goroutine can be stopped with `close(stop)`, and we can wait for it to exit with `<-done`.

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
>
> Spawns a background goroutine unconditionally when the user exports this package. The user has no control over the goroutine or a means of stopping it.
>
> Spawns the worker only if the user requests it. Provides a means of shutting down the worker so that the user can free up resources used by the worker.
>
> Note that you should use `WaitGroup`s if the worker manages multiple goroutines.
