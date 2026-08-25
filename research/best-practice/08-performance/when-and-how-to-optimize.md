---
primary_sources:
  - id: T3-PERF
    title: "go-perfbook"
    url: "https://github.com/dgryski/go-perfbook"
    author: "Damian Gryski"
    section: "When and Where to Optimize; Optimization Workflow"
studied_at: "2026-08-25"
go_min: "1.26.7"
applicability: "current"
---
# When and how to optimize

> **Applicability (Go >= 1.26.7):** Guidance below is current for Go >= 1.26.7.

### Source: go-perfbook — When and Where to Optimize

> I'm putting this first because it's really the most important step. Should
> you even be doing this at all?
>
> Every optimization has a cost. Generally, this cost is expressed in terms of
> code complexity or cognitive load -- optimized code is rarely simpler than
> the unoptimized version.
>
> But there's another side that I'll call the economics of optimization. As a
> programmer, your time is valuable. There's the opportunity cost of what else
> you could be working on for your project, which bugs to fix, which features
> to add. Optimizing things is fun, but it's not always the right task to
> choose. Performance is a feature, but so is shipping, and so is correctness.
>
> Choose the most important thing to work on. Sometimes it's not an actual
> CPU optimization, but a user-experience one. Something as simple as adding a
> progress bar, or making a page more responsive by doing computation in the
> background after rendering the page.
>
> Sometimes this will be obvious: an hourly report that completes in three hours
> is probably less useful than one that completes in less than one.
>
> Just because something is easy to optimize doesn't mean it's worth
> optimizing. Ignoring low-hanging fruit is a valid development strategy.
>
> Think of this as optimizing *your* time.
>
> You get to choose what to optimize and when to optimize. You can move the
> slider between "Fast Software" and "Fast Deployment"
>
> People hear and mindlessly repeat "premature optimization is the root of all
> evil", but they miss the full context of the quote.
>
> "Programmers waste enormous amounts of time thinking about, or worrying about,
> the speed of noncritical parts of their programs, and these attempts at
> efficiency actually have a strong negative impact when debugging and
> maintenance are considered. We should forget about small efficiencies, say
> about 97% of the time: premature optimization is the root of all evil. Yet we
> should not pass up our opportunities in that critical 3%."
>
> -- Knuth
>
> Should you optimize?
> Yes, but only if the problem is important, the program
> is genuinely too slow, and there is some expectation that it can be made
> faster while maintaining correctness, robustness, and clarity."
>
> -- The Practice of Programming, Kernighan and Pike
>
> Premature optimization can also hurt you by tying you into certain decisions.
> The optimized code can be harder to modify if requirements change and harder to
> throw away (sunk-cost fallacy) if needed.
>
> In the vast majority of cases, the size and speed of a program is not a concern.
> The easiest optimization is not having to do it. The second easiest optimization
> is just buying faster hardware.
>
> Once you've decided you're going to change your program, keep reading.

### Source: go-perfbook — Optimization Workflow

> Before we get into the specifics, let's talk about the general process of
> optimization.
>
> Optimization is a form of refactoring. But each step, rather than improving
> some aspect of the source code (code duplication, clarity, etc), improves
> some aspect of the performance: lower CPU, memory usage, latency, etc. This
> improvement generally comes at the cost of readability. This means that in
> addition to a comprehensive set of unit tests (to ensure your changes haven't
> broken anything), you also need a good set of benchmarks to ensure your
> changes are having the desired effect on performance. You must be able to
> verify that your change really *is* lowering CPU. Sometimes a change you
> thought would improve performance will actually turn out to have a zero or
> negative change. Always make sure you undo your fix in these cases.
>
> The benchmarks you are using must be correct and provide reproducible numbers
> on representative workloads. If individual runs have too high a variance, it
> will make small improvements more difficult to spot. You will need to use
> benchstat or equivalent statistical tests
> and won't be able just to eyeball it.
> (Note that using statistical tests is a good idea anyway.) The steps to run
> the benchmarks should be documented, and any custom scripts and tooling
> should be committed to the repository with instructions for how to run them.
> Be mindful of large benchmark suites that take a long time to run: it will
> make the development iterations slower.
>
> Note also that anything that can be measured can be optimized. Make sure
> you're measuring the right thing.
>
> The next step is to decide what you are optimizing for. If the goal is to
> improve CPU, what is an acceptable speed? Do you want to improve the current
> performance by 2x? 10x? Can you state it as "a problem of size N in less than
> time T"? Are you trying to reduce memory usage? By how much? How much slower
> is acceptable for what change in memory usage? What are you willing to give
> up in exchange for lower space requirements?
>
> The performance goals must be specific. You will (almost) always be able to
> make something faster. Optimizing is frequently a game of diminishing returns.
> You need to know when to stop. How much effort are you going to put into
> getting the last little bit of work. How much uglier and harder to maintain
> are you willing to make the code?
>
> In general, optimizations should proceed from top to bottom. Optimizations at
> the system level will have more impact than expression-level ones. Make sure
> you're solving the problem at the appropriate level.
>
> [Amdahl's Law](https://en.wikipedia.org/wiki/Amdahl%27s_law) tells us to focus
> on the bottlenecks. If you double the speed of routine that only takes 5% of
> the runtime, that's only a 2.5% speedup in total wall-clock. On the other hand,
> speeding up routine that takes 80% of the time by only 10% will improve runtime
> by almost 8%. Profiles will help identify where time is actually spent.
