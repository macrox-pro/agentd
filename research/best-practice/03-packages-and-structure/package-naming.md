---
primary_sources:
  - id: T1-CRC
    title: "Go Code Review Comments"
    url: "https://go.dev/wiki/CodeReviewComments"
    author: "Go team"
    section: "Package Names"
  - id: T2-BAHL
    title: "Go Styleguide (bahlo)"
    url: "https://github.com/bahlo/go-styleguide"
    author: "Arne Bahlo"
    section: "Avoid helper/util; Don't under-package; Use internal packages"
    note: "Upstream repo unavailable; English text mirrored via https://github.com/cch123/go-styleguide"
studied_at: "2026-08-25"
also_cited_in: [T2-RAK]
go_min: "1.26.7"
applicability: "current"
---
# Package naming

> **Applicability (Go >= 1.26.7):** Guidance below is current for Go >= 1.26.7.

### Source: Go Code Review Comments — Package Names

> All references to names in your package will be done using the package name, so you can omit that name from the identifiers. For example, if you are in package chubby, you don't need type ChubbyFile, which clients will write as `chubby.ChubbyFile`. Instead, name the type `File`, which clients will write as `chubby.File`. Avoid meaningless package names like util, common, misc, api, types, and interfaces. See https://go.dev/doc/effective_go#package-names and https://go.dev/blog/package-names for more.

### Source: Go Styleguide (bahlo) — Avoid helper/util

> Use clear names and try to avoid creating a `helper.go`, `utils.go` or even package.

### Source: Go Styleguide (bahlo) — Don't under-package

> Deleting or merging packages is far more easier than splitting big ones up.
> When unsure if a package can be split, do it.

### Source: Go Styleguide (bahlo) — Use internal packages

> If you're creating a cmd, consider moving libraries to `internal/` to prevent
> import of unstable, changing packages.
