---
primary_sources:
  - id: T2-UBER
    title: "Uber Go Style Guide"
    url: "https://github.com/uber-go/guide"
    author: "Uber"
    section: "Linting"
studied_at: "2026-08-25"
go_min: "1.26.7"
applicability: "current"
---
# Linting and tools

> **Applicability (Go >= 1.26.7):** Guidance below is current for Go >= 1.26.7.


### Source: Uber Go Style Guide — Linting

> More importantly than any "blessed" set of linters, lint consistently across a codebase.
>
> We recommend using the following linters at a minimum, because we feel that they help to catch the most common issues and also establish a high bar for code quality without being unnecessarily prescriptive:
>
> - errcheck to ensure that errors are handled
> - goimports to format code and manage imports
> - revive to point out common style mistakes
> - govet to analyze code for common mistakes
> - staticcheck to do various static analysis checks

> ### Lint Runners
>
> We recommend golangci-lint as the go-to lint runner for Go code, largely due to its performance in larger codebases and ability to configure and use many canonical linters at once. This repo has an example .golangci.yml config file with recommended linters and settings.
>
> golangci-lint has various linters available for use. The above linters are recommended as a base set, and we encourage teams to add any additional linters that make sense for their projects.
