# Установка

> **Language:** [English](../en/installation.md) · [Русский](./installation.md)

Как получить бинарник `agentd` в `PATH`. Для сборки из исходников нужен Go **1.26+**.

## go install

```bash
go install github.com/macrox-pro/agentd@latest
```

Устанавливает в `$(go env GOPATH)/bin` (каталог должен быть в `PATH`).

## GitHub Releases

Готовые бинарники (linux / darwin / windows, amd64 / arm64) публикуются через goreleaser на [GitHub Releases](https://github.com/macrox-pro/agentd/releases).

## Сборка из исходников

```bash
git clone https://github.com/macrox-pro/agentd.git
cd agentd
make build
```

`make build` пишет `./agentd` и через ldflags задаёт `internal/version.Version` (локально по умолчанию `dev`). Теги релизов подставляют semver через goreleaser.

`agentd daemon status` показывает поле `version`.

## Требования

- Linux, macOS или Windows
- Поддерживаемый coding agent ([Providers](./providers.md))

Регенерация proto (для контрибьюторов): `make generate` (нужен [buf](https://buf.build)).

См. также: [Быстрый старт](./getting-started.md).
