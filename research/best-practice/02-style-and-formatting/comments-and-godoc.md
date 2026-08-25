---
primary_sources:
  - id: T2-PG
    title: "Practical Go"
    url: "https://dave.cheney.net/practical-go/presentations/qcon-china.html"
    author: "Dave Cheney"
    section: "§3 Comments"
studied_at: "2026-08-25"
go_min: "1.26.7"
applicability: "current"
---
# Comments and godoc

> **Applicability (Go >= 1.26.7):** Guidance below is current for Go >= 1.26.7.

### Source: Practical Go — §3 Comments

> Before we move on to larger items I want to spend a few minutes talking about comments.
>
> Good code has lots of comments, bad code requires lots of comments.
>
> — Dave Thomas and Andrew Hunt The Pragmatic Programmer
>
> Comments are very important to the readability of a Go program. Each comments should do one—and only one—of three things:
>
> The comment should explain what the thing does.
>
> The comment should explain how the thing does what it does.
>
> The comment should explain why the thing is why it is.
>
> The first form is ideal for commentary on public symbols:
>
> `// Open opens the named file for reading.`
> `// If successful, methods on the returned file can be used for reading.`
>
> The second form is ideal for commentary inside a method:
>
> `// queue all dependant actions`
>
> The third form, the why , is unique as it does not displace the first two, but at the same time it's not a replacement for the what, or the how. The why style of commentary exists to explain the external factors that drove the code you read on the page. Frequently those factors rarely make sense taken out of context, the comment exists to provide that context.
>
> In this example it may not be immediately clear what the effect of setting `HealthyPanicThreshold` to zero percent will do. The comment is needed to clarify that the value of `0` will disable the panic threshold behaviour.

> #### 3.1. Comments on variables and constants should describe their contents not their purpose
>
> I stated earlier that the name of a variable, or a constant, should describe its purpose. When you add a comment to a variable or constant, that comment should describe the variables contents, not the variables purpose.
>
> `const randomNumber = 6 // determined from an unbiased die`
>
> In this example the comment describes why `randomNumber` is assigned the value six, and where the six was derived from. The comment does not describe where `randomNumber` will be used.
>
> For variables without an initial value, the comment should describe who is responsible for initialising this variable.
>
> Hiding in plain sight — This is a tip from Kate Gregory. Sometimes you'll find a better name for a variable hiding in a comment.
>
> `// registry of SQL drivers`
> `var registry = make(map[string]*sql.Driver)`
>
> The comment was added by the author because `registry` doesn't explain enough about its purpose—it's a registry, but a registry of what? By renaming the variable to `sqlDrivers` its now clear that the purpose of this variable is to hold SQL drivers. Now the comment is redundant and can be removed.

> #### 3.2. Always document public symbols
>
> Because godoc is the documentation for your package, you should always add a comment for every public symbol—variable, constant, function, and method—declared in your package.
>
> Here are two rules from the Google Style guide:
>
> Any public function that is not both obvious and short must be commented.
>
> Any function in a library must be commented regardless of length or complexity
>
> `// ReadAll reads from r until an error or EOF and returns the data it read.`
> `// A successful call returns err == nil, not err == EOF. Because ReadAll is`
> `// defined to read from src until EOF, it does not treat an EOF from Read`
> `// as an error to be reported.`
> `func ReadAll(r io.Reader) ([]byte, error)`
>
> There is one exception to this rule; you don't need to document methods that implement an interface. Specifically don't do this:
>
> `// Read implements the io.Reader interface`
> `func (r *FileReader) Read(buf []byte) (int, error)`
>
> This comment says nothing. It doesn't tell you what the method does, in fact it's worse, it tells you to go look somewhere else for the documentation. In this situation I suggest removing the comment entirely.
>
> Before you write the function, write the comment describing the function. If you find it hard to write the comment, then it's a sign that the code you're about to write is going to be hard to understand.

> ##### 3.2.1. Don't comment bad code, rewrite it
>
> Don't comment bad code — rewrite it
>
> — Brian Kernighan
>
> Comments highlighting the grossness of a particular piece of code are not sufficient. If you encounter one of these comments, you should raise an issue as a reminder to refactor it later. It is okay to live with technical debt, as long as the amount of debt is known.
>
> The tradition in the standard library is to annotate a TODO style comment with the username of the person who noticed it.
>
> `// TODO(dfc) this is O(N^2), find a faster way to do this.`
>
> The username is not a promise that that person has committed to fixing the issue, but they may be the best person to ask when the time comes to address it. Other projects annotate TODOs with a date or an issue number.

> ##### 3.2.2. Rather than commenting a block of code, refactor it
>
> Good code is its own best documentation. As you're about to add a comment, ask yourself, 'How can I improve the code so that this comment isn't needed?' Improve the code and then document it to make it even clearer.
>
> — Steve McConnell
>
> Functions should do one thing only. If you find yourself commenting a piece of code because it is unrelated to the rest of the function, consider extracting it into a function of its own.
>
> In addition to being easier to comprehend, smaller functions are easier to test in isolation. Once you've isolated the orthogonal code into its own function, its name may be all the documentation required.
