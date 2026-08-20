# CLI

> **Language:** [English](../en/cli.md) · [Русский](./cli.md)

Команды и флаги как в `cmd/`. Описание в DESIGN: [§6](../../DESIGN.md#6-cli-reference).

## Persistent flags

| Флаг | Default | Смысл |
|------|---------|--------|
| `--config` | `~/.agentd.yaml` | Путь user-конфига |
| `--socket` | OS default | IPC endpoint демона |
| `-v` / `--verbose` | off | Доп. stderr (никогда не hook stdout) |

## daemon

| Команда | Флаги | Заметки |
|---------|-------|---------|
| `daemon start` | `--foreground` | По умолчанию detach; ждёт успешного Health |
| `daemon stop` | `--timeout` (`10s`) | gRPC Shutdown, затем SIGTERM |
| `daemon status` | `--json` | Runtime snapshot ([Эксплуатация](./operations.md)) |
| `daemon reload` | — | Принудительный re-merge конфига |

## hook

Тонкий edge: decode → gRPC Invoke → encode. Политики в CLI нет.

| Команда | Флаги | Заметки |
|---------|-------|---------|
| `hook run` | `--provider` (обязателен), `--argv-payload`, `--timeout` (`0` = не задан) | Stdin (или argv) hooks |
| `hook notify` | `--provider`, `--timeout` | Codex notify (argv JSON) |
| `hook serve` | `--provider`, `--timeout` | OpenCode NDJSON; provider должен быть `opencode` |

Ошибка dial/Invoke: stderr `daemon not running`, exit **1**. На hook path не писать debug в stdout.

### agenthooks (hidden)

Install пишет `agentd agenthooks run|notify|serve --provider=…`. Поведение как у `hook …`. У `agenthooks serve` default `--provider` = `opencode`.

## config

| Команда | Флаги |
|---------|-------|
| `config validate` | `--cwd` |
| `config show` | `--merged`, `--layer user\|project\|runtime`, `--cwd` |
| `config patch` | `--file` (обязателен) |
| `config record-decision` | `--fingerprint` (обязателен), `--scope` (default `project`), `--project-root`, `--session-id`, `--expires-at` (RFC3339) |

## install

| Флаг | Default |
|------|---------|
| `--provider` | обязателен |
| `--scope` | `project` (`user`, `plugin`) |
| `--dir` | CWD |

## dispatch

| Команда | Флаги |
|---------|-------|
| `dispatch routes` | `--json`, `--cwd` |

Offline compile defaults ⊕ user ⊕ optional project (демон не нужен).

См. также: [Быстрый старт](./getting-started.md), [Providers](./providers.md).
