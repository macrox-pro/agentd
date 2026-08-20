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

`make build` writes `./agentd` with ldflags setting `internal/version.Version` (default `dev` for local builds). Release tags inject the semver via goreleaser.

`agentd daemon status` reports the wired `version` field.

## Requirements

- Linux, macOS, or Windows
- A supported coding agent ([Providers](./providers.md))

Proto regeneration (contributors): `make generate` (needs [buf](https://buf.build)).

See also: [Getting started](./getting-started.md).
