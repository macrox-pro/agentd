# Установка

> **Language:** [English](../en/installation.md) · [Русский](./installation.md)

Как получить исполняемый файл `agentd` в `PATH`. Для сборки из исходников нужен Go **1.26+**.

## go install

```bash
go install github.com/macrox-pro/agentd@latest
```

Ставит бинарник в `$(go env GOPATH)/bin` (этот каталог должен быть в `PATH`).

## Готовые сборки (GitHub Releases)

Бинарники для linux / darwin / windows (amd64 / arm64) публикуются через goreleaser на [странице релизов](https://github.com/macrox-pro/agentd/releases).

## Сборка из исходников

```bash
git clone https://github.com/macrox-pro/agentd.git
cd agentd
make build
```

`make build` создаёт `./agentd` и через флаги компоновщика (`ldflags`) задаёт `internal/version.Version` (локально по умолчанию `dev`). При выпуске тега goreleaser подставляет semver.

Поле `version` видно в `agentd daemon status`.

## Требования

- Linux, macOS или Windows
- Поддерживаемый ИИ-агент для кода ([Провайдеры](./providers.md))

Пересборка protobuf (для разработчиков проекта): `make generate` (нужен [buf](https://buf.build)).

См. также: [Быстрый старт](./getting-started.md).
