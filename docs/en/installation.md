# Installation

> **Language:** [English](./installation.md) · [Русский](../ru/installation.md)

How to get an `agentd` binary on PATH. Requires Go **1.26+** to build from source.

## go install

```bash
go install github.com/macrox-pro/agentd@latest
```

Installs into `$(go env GOPATH)/bin` (ensure that directory is on `PATH`).

## GitHub Releases

Pre-built binaries (linux / darwin / windows, amd64 / arm64) ship via goreleaser on [GitHub Releases](https://github.com/macrox-pro/agentd/releases).

## Build from source

```bash
git clone https://github.com/macrox-pro/agentd.git
cd agentd
make build
```

`make build` writes `./agentd` and links the version into `internal/version.Version` via ldflags (`dev` locally). Release tags: goreleaser injects semver.

| Command | Purpose |
|---------|---------|
| `agentd version` | CLI version |
| `agentd daemon status` | daemon status; field `version` is that process's version |

## Requirements

- Linux, macOS, or Windows
- A supported coding agent ([Providers](./providers.md))

Proto regeneration (contributors): `make generate` (needs [buf](https://buf.build)).

See also: [Getting started](./getting-started.md).
