# Быстрый старт

> **Language:** [English](../en/getting-started.md) · [Русский](./getting-started.md)

Поднять демон пользователя, задать минимальную политику, установить хуки в агент и проверить статус.

Зачем это нужно: [Зачем нужен agentd](./why.md).

## 1. Установить программу

См. [Установка](./installation.md). Быстрый путь:

```bash
go install github.com/macrox-pro/agentd@latest
```

## 2. Запустить демон

Один экземпляр на пользователя. Сокет по умолчанию зависит от ОС (переопределение — `--socket`).

```bash
agentd daemon start
agentd daemon status
```

`--foreground` — не отсоединяться от терминала (удобно при отладке или под process manager).

## 3. Минимальный пользовательский конфиг

Путь по умолчанию: `~/.agentd.yaml` (или `--config`).

```yaml
version: 1
policy:
  fail: fail_closed
  # offline по умолчанию fail_open — агенты работают, если демон не запущен
  # offline: fail_closed  # жёсткий режим при недоступном демоне
guards:
  secrets:
    enabled: true
    action: ask
```

После правки пользовательский файл обычно подхватывается наблюдателем файловой системы (fsnotify); `agentd daemon reload` принудительно заново сливает слои конфига.

## 4. Установить хуки в агент

Из каталога проекта (пример: Claude Code, область `project`):

```bash
agentd install --provider=claude-code --scope=project
```

В сгенерированных настройках агента вызывается `agentd agenthooks …` (скрытый алиас тех же `agentd hook …` — [зачем](./cli.md#agenthooks-скрытая-команда)). В документации удобнее ссылаться на `hook run` / `hook serve` / `hook notify`.

Список провайдеров, областей установки и особенностей: [Провайдеры](./providers.md).

## 5. Проверка

```bash
agentd daemon status --json
```

Ожидается `"running": true` и поля `generation` (поколение конфига), `fingerprint` (отпечаток слияния), `async_queue_depth` (глубина асинхронной очереди), `async_dropped_count` (сколько задач отброшено при переполнении).

Вызовите инструмент в агенте или передайте тестовый JSON на вход:

```bash
echo '{"session_id":"s","cwd":"/tmp","hook_event_name":"PreToolUse","tool_name":"Bash","tool_use_id":"t1","tool_input":{"command":"echo ok"}}' \
  | agentd hook run --provider=claude-code
```

Безопасный `tool.pre` при настройках по умолчанию обычно даёт «пустой» ответ без решения (для Claude — `{}`).

## Дальше

- [Конфигурация](./configuration.md) — слои и схема
- [Охранники](./guards.md) / [Маршрутизация](./dispatch.md) — политика и маршруты
- [Одобрения](./approvals.md) — спросить один раз, потом разрешить
- [Справочник CLI](./cli.md) — полный список флагов
