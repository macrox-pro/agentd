---
primary_sources:
  - id: T1-CRC
    title: "Go Code Review Comments"
    url: "https://go.dev/wiki/CodeReviewComments"
    author: "Go team"
    section: "Gofmt"
  - id: T1-CRC
    title: "Go Code Review Comments"
    url: "https://go.dev/wiki/CodeReviewComments"
    author: "Go team"
    section: "Imports"
  - id: T1-CRC
    title: "Go Code Review Comments"
    url: "https://go.dev/wiki/CodeReviewComments"
    author: "Go team"
    section: "Import Blank"
  - id: T1-CRC
    title: "Go Code Review Comments"
    url: "https://go.dev/wiki/CodeReviewComments"
    author: "Go team"
    section: "Import Dot"
  - id: T1-CRC
    title: "Go Code Review Comments"
    url: "https://go.dev/wiki/CodeReviewComments"
    author: "Go team"
    section: "Line Length"
  - id: T1-EG
    title: "Effective Go"
    url: "https://go.dev/doc/effective_go"
    author: "Go team"
    section: "Formatting"
  - id: T2-BAHL
    title: "Go Styleguide (bahlo)"
    url: "https://github.com/bahlo/go-styleguide"
    author: "Arne Bahlo"
    section: "Use gofmt; Divide imports"
    note: "Upstream repo unavailable; English text mirrored via https://github.com/cch123/go-styleguide"
studied_at: "2026-08-25"
also_cited_in: []
go_min: "1.26.7"
applicability: "current"
---
# Formatting

> **Applicability (Go >= 1.26.7):** Guidance below is current for Go >= 1.26.7.

### Source: Go Code Review Comments — Gofmt

> Run gofmt on your code to automatically fix the majority of mechanical style issues. Almost all Go code in the wild uses `gofmt`. The rest of this document addresses non-mechanical style points.
>
> An alternative is to use goimports, a superset of `gofmt` which additionally adds (and removes) import lines as necessary.

### Source: Go Code Review Comments — Imports

> Avoid renaming imports except to avoid a name collision; good package names should not require renaming. In the event of collision, prefer to rename the most local or project-specific import.
>
> Imports are organized in groups, with blank lines between them. The standard library packages are always in the first group.
>
> ```
> package main
>
> import (
>     "fmt"
>     "hash/adler32"
>     "os"
>
>     "github.com/foo/bar"
>     "rsc.io/goversion/version"
> )
> ```
>
> goimports will do this for you.

### Source: Go Code Review Comments — Import Blank

> Packages that are imported only for their side effects (using the syntax `import _ "pkg"`) should only be imported in the main package of a program, or in tests that require them.

### Source: Go Code Review Comments — Import Dot

> The import . form can be useful in tests that, due to circular dependencies, cannot be made part of the package being tested:
>
> ```
> package foo_test
>
> import (
>     "bar/testutil" // also imports "foo"
>     . "foo"
> )
> ```
>
> In this case, the test file cannot be in package foo because it uses bar/testutil, which imports foo. So we use the 'import .' form to let the file pretend to be part of package foo even though it is not. Except for this one case, do not use import . in your programs. It makes the programs much harder to read because it is unclear whether a name like Quux is a top-level identifier in the current package or in an imported package.

### Source: Go Code Review Comments — Line Length

> There is no rigid line length limit in Go code, but avoid uncomfortably long lines. Similarly, don't add line breaks to keep lines short when they are more readable long–for example, if they are repetitive.
>
> Most of the time when people wrap lines "unnaturally" (in the middle of function calls or function declarations, more or less, say, though some exceptions are around), the wrapping would be unnecessary if they had a reasonable number of parameters and reasonably short variable names. Long lines seem to go with long names, and getting rid of the long names helps a lot.
>
> In other words, break lines because of the semantics of what you're writing (as a general rule) and not because of the length of the line. If you find that this produces lines that are too long, then change the names or the semantics and you'll probably get a good result.
>
> This is, actually, exactly the same advice about how long a function should be. There's no rule "never have a function more than N lines long", but there is definitely such a thing as too long of a function, and of too repetitive tiny functions, and the solution is to change where the function boundaries are, not to start counting lines.

### Source: Effective Go — Formatting

> Formatting issues are the most contentious but the least consequential. People can adapt to different formatting styles but it's better if they don't have to, and less time is devoted to the topic if everyone adheres to the same style. The problem is how to approach this Utopia without a long prescriptive style guide.
>
> With Go we take an unusual approach and let the machine take care of most formatting issues. The `gofmt` program (also available as `go fmt`, which operates at the package level rather than source file level) reads a Go program and emits the source in a standard style of indentation and vertical alignment, retaining and if necessary reformatting comments. If you want to know how to handle some new layout situation, run `gofmt`; if the answer doesn't seem right, rearrange your program (or file a bug about `gofmt`), don't work around it.
>
> As an example, there's no need to spend time lining up the comments on the fields of a structure. `Gofmt` will do that for you. Given the declaration
>
> ```
> type T struct {
>     name string // name of the object
>     value int // its value
> }
> ```
>
> `gofmt` will line up the columns:
>
> ```
> type T struct {
>     name    string // name of the object
>     value   int    // its value
> }
> ```
>
> All Go code in the standard packages has been formatted with `gofmt`.
>
> Some formatting details remain. Very briefly:
>
> Indentation We use tabs for indentation and `gofmt` emits them by default. Use spaces only if you must. Line length Go has no line length limit. Don't worry about overflowing a punched card. If a line feels too long, wrap it and indent with an extra tab. Parentheses Go needs fewer parentheses than C and Java: control structures (`if`, `for`, `switch`) do not have parentheses in their syntax. Also, the operator precedence hierarchy is shorter and clearer, so
>
> ```
> x<<8 + y<<16
> ```
>
> means what the spacing implies, unlike in the other languages.

### Source: Go Styleguide (bahlo) — Use gofmt

> Only commit gofmt'd files, use `-s` to simplify code.

### Source: Go Styleguide (bahlo) — Divide imports

> Dividing std, external and internal imports improves readability.
>
> ```
> import (
> 	"encoding/json"
> 	"fmt"
> 	"os"
>
> 	"github.com/some/external/pkg"
>
> 	"github.com/this-project/pkg/some-lib"
> )
> ```
