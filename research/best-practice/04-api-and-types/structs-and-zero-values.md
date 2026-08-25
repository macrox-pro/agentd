---
primary_sources:
  - id: T1-CRC
    title: "Go Code Review Comments"
    url: "https://go.dev/wiki/CodeReviewComments"
    author: "Go team"
    section: "Declaring Empty Slices"
  - id: T1-CRC
    title: "Go Code Review Comments"
    url: "https://go.dev/wiki/CodeReviewComments"
    author: "Go team"
    section: "Copying"
studied_at: "2026-08-25"
also_cited_in: []
go_min: "1.26.7"
applicability: "current"
---
# Structs and zero values

> **Applicability (Go >= 1.26.7):** Guidance below is current for Go >= 1.26.7.

### Source: Go Code Review Comments — Declaring Empty Slices

> When declaring an empty slice, prefer
>
> ```
> var t []string
> ```
>
> over
>
> ```
> t := []string{}
> ```
>
> The former declares a nil slice value, while the latter is non-nil but zero-length. They are functionally equivalent—their `len` and `cap` are both zero—but the nil slice is the preferred style.
>
> Note that there are limited circumstances where a non-nil but zero-length slice is preferred, such as when encoding JSON objects (a `nil` slice encodes to `null`, while `[]string{}` encodes to the JSON array `[]`).
>
> When designing interfaces, avoid making a distinction between a nil slice and a non-nil, zero-length slice, as this can lead to subtle programming errors.
>
> For more discussion about nil in Go see Francesc Campoy's talk Understanding Nil.

### Source: Go Code Review Comments — Copying

> To avoid unexpected aliasing, be careful when copying a struct from another package. For example, the bytes.Buffer type contains a `[]byte` slice. If you copy a `Buffer`, the slice in the copy may alias the array in the original, causing subsequent method calls to have surprising effects.
>
> In general, do not copy a value of type `T` if its methods are associated with the pointer type, `*T`.
