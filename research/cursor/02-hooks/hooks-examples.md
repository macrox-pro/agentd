---
primary_sources:
  - id: T1-HOOKS
    title: "Examples"
    url: "https://cursor.com/docs/hooks.md"
    section: "Examples"
also_cited_in: []
studied_at: "2026-08-25"
cursor_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Hook examples

> **Applicability:** Verbatim excerpts from Cursor documentation (snapshot 2026-08-25).

### Source: Cursor Hooks — Examples

> ## Examples
>
> The examples below use `./hooks/...` paths, which work for **user hooks** (`~/.cursor/hooks.json`) where scripts run from `~/.cursor/`. For **project hooks** (`<project>/.cursor/hooks.json`), use `.cursor/hooks/...` paths instead since scripts run from the project root.
>
> ```json title="hooks.json"
> {
>   "version": 1,
>   "hooks": {
>     "sessionStart": [
>       {
>         "command": "./hooks/session-init.sh"
>       }
>     ],
>     "sessionEnd": [
>       {
>         "command": "./hooks/audit.sh"
>       }
>     ],
>     "beforeShellExecution": [
>       {
>         "command": "./hooks/audit.sh"
>       },
>       {
>         "command": "./hooks/block-git.sh"
>       }
>     ],
>     "beforeMCPExecution": [
>       {
>         "command": "./hooks/audit.sh"
>       }
>     ],
>     "afterShellExecution": [
>       {
>         "command": "./hooks/audit.sh"
>       }
>     ],
>     "afterMCPExecution": [
>       {
>         "command": "./hooks/audit.sh"
>       }
>     ],
>     "afterFileEdit": [
>       {
>         "command": "./hooks/audit.sh"
>       }
>     ],
>     "beforeSubmitPrompt": [
>       {
>         "command": "./hooks/audit.sh"
>       }
>     ],
>     "preCompact": [
>       {
>         "command": "./hooks/audit.sh"
>       }
>     ],
>     "stop": [
>       {
>         "command": "./hooks/audit.sh"
>       }
>     ],
>     "beforeTabFileRead": [
>       {
>         "command": "./hooks/redact-secrets-tab.sh"
>       }
>     ],
>     "afterTabFileEdit": [
>       {
>         "command": "./hooks/format-tab.sh"
>       }
>     ]
>   }
> }
> ```
>
> ```sh title="audit.sh"
> #!/bin/bash
>
> # audit.sh - Hook script that writes all JSON input to /tmp/agent-audit.log
> # This script is designed to be called by Cursor's hooks system for auditing purposes
>
> # Read JSON input from stdin
> json_input=$(cat)
>
> # Create timestamp for the log entry
> timestamp=$(date '+%Y-%m-%d %H:%M:%S')
>
> # Create the log directory if it doesn't exist
> mkdir -p "$(dirname /tmp/agent-audit.log)"
>
> # Write the timestamped JSON entry to the audit log
> echo "[$timestamp] $json_input" >> /tmp/agent-audit.log
>
> # Exit successfully
> exit 0
> ```
>
> ```sh title="block-git.sh"
> #!/bin/bash
>
> # Hook to block git commands and redirect to gh tool usage
> # This hook implements the beforeShellExecution hook from the Cursor Hooks Spec
>
> # Initialize debug logging
> echo "Hook execution started" >> /tmp/hooks.log
>
> # Read JSON input from stdin
> input=$(cat)
> echo "Received input: $input" >> /tmp/hooks.log
>
> # Parse the command from the JSON input
> command=$(echo "$input" | jq -r '.command // empty')
> echo "Parsed command: '$command'" >> /tmp/hooks.log
>
> # Check if the command contains 'git' or 'gh'
> if [[ "$command" =~ git[[:space:]] ]] || [[ "$command" == "git" ]]; then
>     echo "Git command detected - blocking: '$command'" >> /tmp/hooks.log
>     # Block the git command and provide guidance to use gh tool instead
>     cat << EOF
> {
>   "continue": true,
>   "permission": "deny",
>   "user_message": "Git command blocked. Please use the GitHub CLI (gh) tool instead.",
>   "agent_message": "The git command '$command' has been blocked by a hook. Instead of using raw git commands, please use the 'gh' tool which provides better integration with GitHub and follows best practices. For example:\n- Instead of 'git clone', use 'gh repo clone'\n- Instead of 'git push', use 'gh repo sync' or the appropriate gh command\n- For other git operations, check if there's an equivalent gh command or use the GitHub web interface\n\nThis helps maintain consistency and leverages GitHub's enhanced tooling."
> }
> EOF
> elif [[ "$command" =~ gh[[:space:]] ]] || [[ "$command" == "gh" ]]; then
>     echo "GitHub CLI command detected - asking for permission: '$command'" >> /tmp/hooks.log
>     # Ask for permission for gh commands
>     cat << EOF
> {
>   "continue": true,
>   "permission": "ask",
>   "user_message": "GitHub CLI command requires permission: $command",
>   "agent_message": "The command '$command' uses the GitHub CLI (gh) which can interact with your GitHub repositories and account. Please review and approve this command if you want to proceed."
> }
> EOF
> else
>     echo "Non-git/non-gh command detected - allowing: '$command'" >> /tmp/hooks.log
>     # Allow non-git/non-gh commands
>     cat << EOF
> {
>   "continue": true,
>   "permission": "allow"
> }
> EOF
> fi
> ```
>
> ### TypeScript stop automation hook
>
> Choose TypeScript when you need typed JSON, durable file I/O, and HTTP calls in the same hook. This Bun-powered `stop` hook tracks per-conversation failure counts on disk, forwards structured telemetry to an internal API, and can automatically schedule a retry when the agent fails twice in a row.
>
> ```json title="hooks.json"
> {
>   "version": 1,
>   "hooks": {
>     "stop": [
>       {
>         "command": "bun run .cursor/hooks/track-stop.ts --stop"
>       }
>     ]
>   }
> }
> ```
>
> ```ts title=".cursor/hooks/track-stop.ts"
> import { mkdir, readFile, writeFile } from 'node:fs/promises';
> import { stdin } from 'bun';
>
> type StopHookInput = {
>   conversation_id: string;
>   generation_id: string;
>   model: string;
>   model_id?: string;
>   model_params?: Array<{ id: string; value: string }>;
>   status: 'completed' | 'aborted' | 'error';
>   loop_count: number;
> };
>
> type StopHookOutput = {
>   followup_message?: string;
> };
>
> type MetricsEntry = {
>   lastStatus: StopHookInput['status'];
>   errorCount: number;
>   lastUpdatedIso: string;
> };
>
> type MetricsStore = Record<string, MetricsEntry>;
>
> const STATE_DIR = '.cursor/hooks/state';
> const METRICS_PATH = `${STATE_DIR}/agent-metrics.json`;
> const TELEMETRY_URL = Bun.env.AGENT_TELEMETRY_URL;
>
> async function parseHookInput<T>(): Promise<T> {
>   const text = await stdin.text();
>   return JSON.parse(text) as T;
> }
>
> async function readMetrics(): Promise<MetricsStore> {
>   try {
>     return JSON.parse(await readFile(METRICS_PATH, 'utf8')) as MetricsStore;
>   } catch {
>     return {};
>   }
> }
>
> async function writeMetrics(store: MetricsStore) {
>   await mkdir(STATE_DIR, { recursive: true });
>   await writeFile(METRICS_PATH, JSON.stringify(store, null, 2), 'utf8');
> }
>
> async function sendTelemetry(payload: StopHookInput, entry: MetricsEntry) {
>   if (!TELEMETRY_URL) return;
>   await fetch(TELEMETRY_URL, {
>     method: 'POST',
>     headers: { 'Content-Type': 'application/json' },
>     body: JSON.stringify({
>       conversationId: payload.conversation_id,
>       generationId: payload.generation_id,
>       model: payload.model,
>       modelId: payload.model_id,
>       modelParams: payload.model_params,
>       status: payload.status,
>       errorCount: entry.errorCount,
>       loopCount: payload.loop_count,
>       timestamp: entry.lastUpdatedIso
>     })
>   });
> }
>
> async function main() {
>   const payload = await parseHookInput<StopHookInput>();
>   const metrics = await readMetrics();
>   const entry =
>     metrics[payload.conversation_id] ?? {
>       lastStatus: payload.status,
>       errorCount: 0,
>       lastUpdatedIso: ''
>     };
>
>   entry.lastStatus = payload.status;
>   entry.lastUpdatedIso = new Date().toISOString();
>   entry.errorCount = payload.status === 'error' ? entry.errorCount + 1 : 0;
>
>   metrics[payload.conversation_id] = entry;
>   await writeMetrics(metrics);
>   await sendTelemetry(payload, entry);
>
>   const response: StopHookOutput = {};
>   if (entry.errorCount >= 2 && payload.loop_count < 4) {
>     response.followup_message =
>       'Automated retry triggered after two failures. Double-check credentials before running again.';
>   }
>
>   process.stdout.write(JSON.stringify(response) + '\n');
> }
>
> main().catch(error => {
>   console.error('[stop hook] failed', error);
>   process.stdout.write('{}\n');
> });
> ```
>
> Set `AGENT_TELEMETRY_URL` to the internal endpoint that should receive run summaries.
>
> ### Python manifest guard hook
>
> Python shines when you need rich parsing libraries. This hook uses `pyyaml` to inspect Kubernetes manifests before `kubectl apply` runs; Bash would struggle to parse multi-document YAML safely.
>
> ```json title="hooks.json"
> {
>   "version": 1,
>   "hooks": {
>     "beforeShellExecution": [
>       {
>         "command": "python3 .cursor/hooks/kube_guard.py"
>       }
>     ]
>   }
> }
> ```
>
> ```python title=".cursor/hooks/kube_guard.py"
> #!/usr/bin/env python3
> import json
> import shlex
> import sys
> from pathlib import Path
>
> import yaml
>
> SENSITIVE_NAMESPACES = {"prod", "production"}
>
> def main() -> None:
>     payload = json.load(sys.stdin)
>     command = payload.get("command", "")
>     cwd = Path(payload.get("cwd") or ".")
>     response = {"continue": True, "permission": "allow"}
>
>     try:
>         args = shlex.split(command)
>     except ValueError:
>         print(json.dumps(response))
>         return
>
>     if len(args) < 2 or args[0] != "kubectl" or args[1] != "apply" or "-f" not in args:
>         print(json.dumps(response))
>         return
>
>     f_index = args.index("-f")
>     if f_index + 1 >= len(args):
>         print(json.dumps(response))
>         return
>
>     manifest_arg = args[f_index + 1]
>     manifest_path = (cwd / manifest_arg).resolve()
>
>     if not manifest_path.exists():
>         print(json.dumps(response))
>         return
>
>     cli_namespace = None
>     for i, arg in enumerate(args):
>         if arg in ("-n", "--namespace") and i + 1 < len(args):
>             cli_namespace = args[i + 1]
>         elif arg.startswith("--namespace="):
>             cli_namespace = arg.split("=", 1)[1]
>         elif arg.startswith("-n="):
>             cli_namespace = arg.split("=", 1)[1]
>
>     try:
>         documents = list(yaml.safe_load_all(manifest_path.read_text()))
>     except (OSError, yaml.YAMLError) as exc:
>         sys.stderr.write(f"Failed to read/parse {manifest_path}: {exc}\n")
>         print(json.dumps(response))
>         return
>
>     if cli_namespace in SENSITIVE_NAMESPACES or any(
>         (doc or {}).get("metadata", {}).get("namespace") in SENSITIVE_NAMESPACES
>         for doc in documents
>     ):
>         response.update(
>             {
>                 "permission": "ask",
>                 "user_message": "kubectl apply to prod requires manual approval.",
>                 "agent_message": f"{manifest_path.name} includes protected namespaces; confirm with your team before continuing.",
>             }
>         )
>
>     print(json.dumps(response))
>
> if __name__ == "__main__":
>     main()
> ```
>
> Install PyYAML (for example, `pip install pyyaml`) wherever your hook scripts run so the parser import succeeds.
