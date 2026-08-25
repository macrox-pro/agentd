---
primary_sources:
  - id: T2-PG
    title: "Practical Go"
    url: "https://dave.cheney.net/practical-go/presentations/qcon-china.html"
    author: "Dave Cheney"
    section: "§2 Identifiers"
  - id: T3-CLEAN
    title: "Clean Go Code"
    url: "https://github.com/Pungyeon/clean-go-article"
    author: "Pungyeon"
    section: "Naming Conventions"
  - id: T3-IDIO
    title: "Idiomatic Go"
    url: "https://dmitri.shuralyov.com/idiomatic-go"
    author: "Dmitri Shuralyov"
    section: "Use consistent spelling of certain words"
studied_at: "2026-08-25"
go_min: "1.26.7"
applicability: "current"
---
# Identifiers and naming

> **Applicability (Go >= 1.26.7):** Guidance below is current for Go >= 1.26.7.

### Source: Practical Go — §2 Identifiers

> The first topic we're going to discuss is identifiers. An identifier is a fancy word for a name; the name of a variable, the name of a function, the name of a method, the name of a type, the name of a package, and so on.
>
> Poor naming is symptomatic of poor design.
>
> Given the limited syntax of Go, the names we choose for things in our programs have an oversized impact on the readability of our programs. Readability is the defining quality of good code, thus choosing good names is crucial to the readability of Go code.

> #### 2.1. Choose identifiers for clarity, not brevity
>
> Obvious code is important. What you can do in one line you should do in three.
>
> Go is not a language that optimises for clever one liners. Go is not a language which optimises for the least number of lines in a program. We're not optimising for the size of the source code on disk, nor how long it takes to type the program into an editor.
>
> Good naming is like a good joke. If you have to explain it, it's not funny.
>
> Key to this clarity is the names we choose for identifies in Go programs. Let's talk about the qualities of a good name:
>
> A good name is concise. A good name need not be the shortest it can possibly be, but a good name should waste no space on things which are extraneous. Good names have a high signal to noise ratio.
>
> A good name is descriptive. A good name should describe the application of a variable or constant, not their contents. A good name should describe the result of a function, or behaviour of a method, not their implementation. A good name should describe the purpose of a package, not its contents. The more accurately a name describes the thing it identifies, the better the name.
>
> A good name is should be predictable. You should be able to infer the way a symbol will be used from its name alone. This is a function of choosing descriptive names, but it also about following tradition. This is what Go programmers talk about when they say idiomatic.

> #### 2.2. Identifier length
>
> Sometimes people criticise the Go style for recommending short variable names. As Rob Pike said, "Go programmers want the right length identifiers".
>
> Andrew Gerrand suggests that by using longer identifies to indicate to the reader things of higher importance.
>
> The greater the distance between a name's declaration and its uses, the longer the name should be.
>
> — Andrew Gerrand
>
> From this we can draw some guidelines:
>
> Short variable names work well when the distance between their declaration and last use is short.
>
> Long variable names need to justify themselves; the longer they are the more value they need to provide. Lengthy bureaucratic names carry a low amount of signal compared to their weight on the page.
>
> Don't include the name of your type in the name of your variable.
>
> Constants should describe the value they hold, not how that value is used.
>
> Prefer single letter variables for loops and branches, single words for parameters and return values, multiple words for functions and package level declarations
>
> Prefer single words for methods, interfaces, and packages.
>
> Remember that the name of a package is part of the name the caller uses to to refer to it, so make use of that.

> ##### 2.2.1. Context is key
>
> It's important to recognise that most advice on naming is contextual. I like to say it is a principle, not a rule.
>
> What is the difference between two identifiers, `i`, and `index`. We cannot say conclusively that one is better than another, for example is
>
> `for index := 0; index < len(s); index++ { }`
>
> fundamentally more readable than
>
> `for i := 0; i < len(s); i++ { }`
>
> I argue it is not, because it is likely the scope of `i`, and `index` for that matter, is limited to the body of the `for` loop and the extra verbosity of the latter adds little to comprehension of the program.
>
> However, which of these functions is more readable?
>
> `func (s *SNMP) Fetch(oid []int, index int) (int, error)`
>
> or
>
> `func (s *SNMP) Fetch(o []int, i int) (int, error)`
>
> In this example, `oid` is an abbreviation for SNMP Object ID, so shortening it to `o` would mean programmers have to translate from the common notation that they read in documentation to the shorter notation in your code. Similarly, reducing `index` to `i` obscures what `i` stands for as in SNMP messages a sub value of each OID is called an Index.
>
> Don't mix and match long and short formal parameters in the same declaration.

> #### 2.3. Don't name your variables for their types
>
> You shouldn't name your variables after their types for the same reason you don't name your pets "dog" and "cat". You also probably shouldn't include the name of your type in the name of your variable's name for the same reason.
>
> The name of the variable should describe its contents, not the type of the contents. Consider this example:
>
> `var usersMap map[string]*User`
>
> What's good about this declaration? We can see that its a map, and it has something to do with the `*User` type, that's probably good. But `usersMap` is a map, and Go being a statically typed language won't let us accidentally use it where a scalar variable is required, so the `Map` suffix is redundant.
>
> My suggestion is to avoid any suffix that resembles the type of the variable.
>
> If `users` isn't descriptive enough, then `usersMap` won't be either.
>
> This advice also applies to function parameters. For example:
>
> `func WriteConfig(w io.Writer, config *Config)`
>
> Naming the `*Config` parameter `config` is redundant. We know its a `*Config`, it says so right there.
>
> In this case consider `conf` or maybe `c` will do if the lifetime of the variable is short enough.
>
> If there is more that one `*Config` in scope at any one time then calling them `conf1` and `conf2` is less descriptive than calling them `original` and `updated` as the latter are less likely to be mistaken for one another.
>
> Don't let package names steal good variable names. The name of an imported identifier includes its package name. For example the `Context` type in the `context` package will be known as `context.Context`. This makes it impossible to use `context` as a variable or type in your package. This is why the local declaration for `context.Context` types is traditionally `ctx`. eg. `func WriteLog(ctx context.Context, message string)`

> #### 2.4. Use a consistent naming style
>
> Another property of a good name is it should be predictable. The reader should be able to understand the use of a name when they encounter it for the first time. When they encounter a common name, they should be able to assume it has not changed meanings since the last time they saw it.
>
> For example, if your code passes around a database handle, make sure each time the parameter appears, it has the same name. Rather than a combination of `d *sql.DB`, `dbase *sql.DB`, `DB *sql.DB`, and `database *sql.DB`, instead consolidate on something like `db *sql.DB`.
>
> Doing so promotes familiarity; if you see a `db`, you know it's a `*sql.DB` and that it has either been declared locally or provided for you by the caller.
>
> Similar advice applies to method receivers; use the same receiver name every method on that type. This makes it easier for the reader to internalise the use of the receiver across the methods in this type.
>
> The convention for short receiver names in Go is at odds with the advice provided so far. This is just one of the choices made early on that has become the preferred style, just like the use of `CamelCase` rather than `snake_case`.
>
> Go style dictates that receivers have a single letter name, or acronyms derived from their type. You may find that the name of your receiver sometimes conflicts with name of a parameter in a method. In this case, consider making the parameter name slightly longer, and don't forget to use this new parameter name consistently.
>
> Finally, certain single letter variables have traditionally been associated with loops and counting. For example, `i`, `j`, and `k` are commonly the loop induction variable for simple `for` loops. `n` is commonly associated with a counter or accumulator. `v` is a common shorthand for a value in a generic encoding function, `k` is commonly used for the key of a map, and `s` is often used as shorthand for parameters of type `string`.
>
> As with the `db` example above programmers expect `i` to be a loop induction variable. If you ensure that `i` is always a loop variable, not used in other contexts outside a `for` loop. When readers encounter a variable called `i`, or `j`, they know that a loop is close by.
>
> If you found yourself with so many nested loops that you exhaust your supply of `i`, `j`, and `k` variables, its probably time to break your function into smaller units.

> #### 2.5. Use a consistent declaration style
>
> Go has at least six different ways to declare a variable: `var x int = 1`, `var x = 1`, `var x int; x = 1`, `var x = int(1)`, `x := 1`. I'm sure there are more that I haven't thought of. This is something that Go's designers recognise was probably a mistake, but its too late to change it now. With all these different ways of declaring a variable, how do we avoid each Go programmer choosing their own style?
>
> When declaring, but not initialising, a variable, use `var`. When declaring a variable that will be explicitly initialised later in the function, use the `var` keyword.
>
> The `var` acts as a clue to say that this variable has been deliberately declared as the zero value of the indicated type. This is also consistent with the requirement to declare variables at the package level using `var` as opposed to the short declaration syntax—although I'll argue later that you shouldn't be using package level variables at all.
>
> When declaring and initialising, use `:=`. When declaring and initialising the variable at the same time, that is to say we're not letting the variable be implicitly initialised to its zero value, I recommend using the short variable declaration form. This makes it clear to the reader that the variable on the left hand side of the `:=` is being deliberately initialised.
>
> In summary:
>
> When declaring a variable without initialisation, use the `var` syntax.
>
> When declaring and explicitly initialising a variable, use `:=`.
>
> Make tricky declarations obvious. When something is complicated, it should look complicated.

> #### 2.6. Be a team player
>
> I talked about a goal of software engineering to produce readable, maintainable, code. Therefore you will likely spend most of your career working on projects of which you are not the sole author. My advice in this situation is to follow the local style.
>
> Changing styles in the middle of a file is jarring. Uniformity, even if its not your preferred approach, is more valuable for maintenance than your personal preference. My rule of thumb is; if it fits through `gofmt` then its usually not worth holding up a code review for.
>
> If you want to do a renaming across a code-base, do not mix this into another change. If someone is using git bisect they don't want to wade through thousands of lines of renaming to find the code you changed as well.

### Source: Clean Go Code — Function Naming

> The general rule here is really simple: the more specific the function, the more general its name. In other words, we want to start with a very broad and short function name, such as `Run` or `Parse`, that describes the general functionality.
>
> When we go one layer deeper, our function naming will become slightly more specific:
>
> This kind of logical progression in our function names—from a high level of abstraction to a lower, more specific one—makes the code easier to follow and read. Consider the alternative: If our highest level of abstraction is too specific, then we'll end up with a name that attempts to cover all bases, like `DetermineFileExtensionAndParseConfigurationFile`. This is horrendously difficult to read; we are trying to be too specific too soon and end up confusing the reader, despite trying to be clear!

### Source: Clean Go Code — Variable Naming

> Rather interestingly, the opposite is true for variables. Unlike functions, our variables should be named from more to less specific the deeper we go into nested scopes.
>
> You shouldn't name your variables after their types for the same reason you wouldn't name your pets 'dog' or 'cat'. – Dave Cheney
>
> Why should our variable names become less specific as we travel deeper into a function's scope? Simply put, as a variable's scope becomes smaller, it becomes increasingly clear for the reader what that variable represents, thereby eliminating the need for specific naming.

### Source: Idiomatic Go — Use consistent spelling of certain words

> Do this:
>
> ```
> // marshaling
> // unmarshaling
> // canceling
> // canceled
> // cancellation
> ```
>
> Don't do this:
>
> ```
> // marshalling
> // unmarshalling
> // cancelling
> // cancelled
> // cancelation
> ```
>
> For consistency with the Go project. These words have multiple valid spellings. The Go project picks one. See go.dev/wiki/Spelling.

