# Быстрый старт

> **Language:** [English](../en/getting-started.md) · [Русский](./getting-started.md)

Поднять user-level демон, задать минимальную политику, установить hooks агента и проверить Status.

## 1. Установить бинарник

См. [Установка](./installation.md). Быстрый путь:

```bash
go install github.com/macrox-pro/agentd@latest
```

## 2. Запустить демон

Один экземпляр на пользователя. Сокет по умолчанию зависит от ОС (`--socket` переопределяет).

```bash
agentd daemon start
agentd daemon status
```

`--foreground` — процесс остаётся в foreground (dev / process manager).

## 3. Минимальный user-конфиг

Путь по умолчанию: `~/.agentd.yaml` (или `--config`).

```yaml
version: 1
policy:
  fail: fail_closed
guards:
  secrets:
    enabled: true
    action: ask
```

После правки user-файл подхватывается fsnotify; `agentd daemon reload` принудительно пересобирает merge.

## 4. Install hooks для агента

Из каталога проекта (пример: Claude Code, scope project):

```bash
agentd install --provider=claude-code --scope=project
```

Сгенерированные конфиги вызывают `agentd agenthooks …` (скрытый алиас `agentd hook …`). В документации предпочтительны `hook run` / `hook serve` / `hook notify`.

Провайдеры и scope: [Providers](./providers.md).

## 5. Проверка

```bash
agentd daemon status --json
```

Ожидается `"running": true`, поля `generation`, `fingerprint`, `async_queue_depth`, `async_dropped_count`.

Вызовите tool в агенте или прогоните payload через edge:

```bash
echo '{"session_id":"s","cwd":"/tmp","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"t1","tool_input":{"command":"echo ok"}}' \
  | agentd hook run --provider=claude-code
```

Чистый tool.pre при defaults обычно даёт provider no-op (Claude: `{}`).

## Дальше

- [Конфигурация](./configuration.md) — слои и схема
- [Guards](./guards.md) / [Dispatch](./dispatch.md) — политика и маршрутизация
- [Approvals](./approvals.md) — Ask once, затем Allow
- [CLI](./cli.md) — полный список флагов
