---
primary_sources:
  - id: T1-PERMISSIONS
    title: "Permissions"
    url: "https://code.claude.com/docs/en/permissions.md"
    section: ""
  - id: T1-PERM-MODES
    title: "Permission modes"
    url: "https://code.claude.com/docs/en/permission-modes.md"
    section: ""
also_cited_in: []
studied_at: "2026-08-25"
claude_code_docs_snapshot: "2026-08-25"
applicability: "current"
---
# Permissions and modes

> **Applicability:** Verbatim excerpts from Claude Code documentation (snapshot 2026-08-25).

### Source: Permissions
>
> # Configure permissions
>
> > Control what Claude Code can access and do with fine-grained permission rules, modes, and managed policies.
>
> Claude Code supports fine-grained permissions so that you can specify exactly what the agent is allowed to do and what it can't. You can check permission settings into version control to share them with every developer in your organization, and each developer can customize their own.
>
> ## Permission system
>
> Claude Code uses a tiered permission system to balance power and safety. The table shows, for each tool type, whether Manual mode asks before the action runs. The other [permission modes](#permission-modes) change which of these ask you; in auto mode a classifier reviews actions instead of you, and [how the classifier evaluates actions](/docs/en/permission-modes#how-the-classifier-evaluates-actions) lists which ones it sees.
>
> | Tool type         | Example          | Approval required                                                                                             | "Yes, and don't ask again" behavior    |
> | :---------------- | :--------------- | :------------------------------------------------------------------------------------------------------------ | :------------------------------------- |
> | Read-only         | File reads, Grep | No, within the [working directory and additional directories](#working-directories)                           | N/A                                    |
> | Bash commands     | Shell execution  | Yes, except a built-in set of [read-only commands](#read-only-commands)                                       | Permanently per repository and command |
> | File modification | Edit/write files | Yes                                                                                                           | Until session end                      |
> | Web fetch         | WebFetch         | Yes, except a built-in set of [preapproved documentation domains](/docs/en/tools-reference#webfetch-tool-behavior) | Permanently per repository and domain  |
> | Web search        | WebSearch        | Yes                                                                                                           | Permanently per repository             |
>
> When you choose "Yes, and don't ask again" and the approval saves permanently, such as for a Bash command or a WebFetch domain, Claude Code saves the rule to `.claude/settings.local.json` at the root of the git repository, resolved through [worktrees](/docs/en/worktrees) to the main checkout. The rule applies to future sessions anywhere in that repository, including sessions started in subdirectories and in worktrees. A file-modification approval isn't saved to the file: as the table shows, it lasts until the session ends. In some cases, such as outside a git repository or on Windows, Claude Code saves the rule in the directory you started it from; [Where Claude Code looks for each file](/docs/en/settings#where-claude-code-looks-for-each-file) lists them.
>
> Before v2.1.211, Claude Code always saved the rule in the starting directory, so an approval granted in a worktree or subdirectory didn't apply to the rest of the repository. Rules that earlier versions saved in a subdirectory or worktree still apply to sessions started there.
>
> Sometimes a permission prompt offers only a one-time approval, with no "don't ask again" option and no option to allow the action for the rest of the session. Claude Code offers those options only when the prompt can show you everything they would allow, so a rule you save from a prompt covers only what its option named.
>
> When the directory you started Claude Code in is what makes the option's label too long, Claude Code shortens it in the label, replacing your home directory with `~` and then the end of the path with `…`, and keeps the option. You still save the same rule. Claude Code leaves the options out in three cases:
>
> * **Command or edit:** too large to show in full.
> * **Commands or paths the rule would cover:** the label can't fit them all.
> * **Starting directory too long, not shortened:** it contains characters Claude Code can't display safely, or even its start doesn't fit.
>
> Approve the action once, or add the rule yourself in [`/permissions`](#manage-permissions).
>
> On a Bash or PowerShell permission prompt, press `Ctrl+E` to show an explanation of the command: what it does, why Claude is running it, and what could go wrong, labeled **Low risk**, **Med risk**, or **High risk**. Claude Code sends the command and Claude's own description of the call to the model to generate the explanation only when you press `Ctrl+E`, not on every prompt. Showing the explanation doesn't run the command; press `Ctrl+E` again to hide it.
>
> To turn the shortcut off, set [`permissionExplainerEnabled`](/docs/en/settings-reference#permissionexplainerenabled) to `false` in `~/.claude.json`.
>
> ### Add a comment when you answer a permission prompt
>
> You can attach a note to Claude when you approve or deny a single action. On most permission prompts, including Bash, PowerShell, file, and MCP tool prompts, move to **Yes** or **No** and press `Tab` to open a comment field on that option. WebFetch and browser prompts don't offer the field. The options that allow the action for the rest of the session or save a rule don't take one either.
>
> With the field open, type the comment and then press one of these keys:
>
> * `Enter`: submits your answer with the comment attached. If you leave the field empty, Claude Code submits the answer without a comment.
> * `Tab`: closes the field without answering. Claude Code keeps the text you typed and still sends it if you answer with that option.
> * `Shift+Tab`: on a file prompt, such as an Edit or Write prompt, closes the field the same as `Tab`. Before v2.1.235, pressing `Shift+Tab` inside the field instead selected the option that allows the action for the rest of the session, so Claude Code approved the action for the rest of the session and discarded the comment.
>
> Claude Code delivers the comment differently depending on how you answered:
>
> * **Yes**: Claude Code runs the action, then sends your comment to Claude after the result.
> * **No**: Claude Code sends your comment to Claude as the reason for the denial, and Claude continues working. If you select **No** without a comment on a prompt from the main conversation, Claude Code stops the turn.
>
> ## Manage permissions
>
> You can view and manage Claude Code's tool permissions with `/permissions`. The dialog lists all permission rules and the `settings.json` file each rule comes from. You can open the dialog while Claude is working: when you add or remove a rule, Claude Code applies the change starting with Claude's next tool call in the same turn. Before v2.1.234, Claude Code queued the command until the turn finished.
>
> * **Allow** rules let Claude Code use the specified tool without manual approval.
> * **Ask** rules prompt for confirmation whenever Claude Code tries to use the specified tool.
> * **Deny** rules prevent Claude Code from using the specified tool.
>
> Rules are evaluated in order: deny, then ask, then allow. The first match in that order determines the outcome, and rule specificity doesn't change the order.
>
> A broad deny rule like `Bash(aws *)` blocks every matching call, including calls that also match a narrower allow rule like `Bash(aws s3 ls)`, so a deny rule can't carry allowlist exceptions. The same precedence applies between ask and allow: a matching ask rule prompts even when a more specific allow rule also matches the same call.
>
> Deny rules behave differently depending on whether they name a tool or scope a pattern within one. A bare tool name like `Bash` removes the tool from Claude's context entirely, so Claude never sees it. Bare-name removal applies to every tool except [`EndConversation`](/docs/en/tools-reference#endconversation-tool-behavior): a deny rule can't remove it while any other tool remains, and an ask rule never prompts for it. A scoped rule like `Bash(rm *)` leaves the tool available and blocks matching calls when Claude attempts them.
>
>   Permission rules are enforced by Claude Code, not by the model. Instructions in your prompt or `CLAUDE.md` shape what Claude tries to do, but they don't change what Claude Code allows. To grant or revoke access, use `/permissions`, the rules described here, a [permission mode](/docs/en/permission-modes), or a [PreToolUse hook](#extend-permissions-with-hooks).
>
> ## Permission modes
>
> Claude Code supports several permission modes that control how it approves tool calls. See [Permission modes](/docs/en/permission-modes) for when to use each one. To change the mode sessions start in, set `defaultMode` in your [settings files](/docs/en/settings#where-settings-live). [Which mode a session starts in](/docs/en/permission-modes#which-mode-a-session-starts-in) covers the built-in default for each plan and what the VS Code extension reads.
>
> | Mode                | Description                                                                                                                                                                                                                                                                                                                                         |
> | :------------------ | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | `default`           | Prompts for permission on first use of each tool. Labeled Manual in the CLI, the VS Code and JetBrains extensions, and the desktop app, and Claude Code accepts `manual` as an alias. The label and alias require Claude Code v2.1.200 or later. The desktop app's label doesn't depend on your CLI version                                         |
> | `acceptEdits`       | Automatically accepts file edits and common filesystem commands such as `mkdir`, `touch`, `mv`, and `cp` for paths in the working directory or `additionalDirectories`                                                                                                                                                                              |
> | `plan`              | Claude reads files and runs read-only shell commands to explore but doesn't edit your source files; with [auto mode](/docs/en/permission-modes#eliminate-prompts-with-auto-mode) available, classifier-approved commands also run. Labeled Plan in the CLI and the VS Code extension                                                                     |
> | `auto`              | Auto-approves tool calls with background safety checks that verify actions align with your request                                                                                                                                                                                                                                                  |
> | `dontAsk`           | Auto-denies tools unless pre-approved via `/permissions` or `permissions.allow` rules. `AskUserQuestion`, connector tools [your organization set to `ask`](/docs/en/mcp#organization-controls-on-connector-tools), and MCP tools marked [`requiresUserInteraction`](/docs/en/mcp#require-approval-for-a-specific-tool) are denied even if you've allowed them |
> | `bypassPermissions` | Skips permission prompts, except for the [actions no mode auto-approves](/docs/en/permission-modes#actions-no-mode-auto-approves)                                                                                                                                                                                                                        |
>
>   `bypassPermissions` mode skips permission prompts, including for writes to [protected paths](/docs/en/permission-modes#protected-paths) such as `.git` and `.claude`. The [cross-session messaging safeguards](/docs/en/permission-modes#skip-all-checks-with-bypasspermissions-mode) still apply. Only use this mode in isolated environments like containers or VMs where Claude Code can't cause damage.
>
> To prevent `bypassPermissions` or `auto` mode from being used, set `permissions.disableBypassPermissionsMode` or `permissions.disableAutoMode` to `"disable"` in any [settings file](/docs/en/settings#where-settings-live). These are most useful in [managed settings](#managed-settings) where they can't be overridden.
>
> ## Permission rule syntax
>
> Permission rules follow the format `Tool` or `Tool(specifier)`.
>
> ### Match all uses of a tool
>
> To match all uses of a tool, use only the tool name without parentheses:
>
> | Rule       | Effect                         |
> | :--------- | :----------------------------- |
> | `Bash`     | Matches all Bash commands      |
> | `WebFetch` | Matches all web fetch requests |
> | `Read`     | Matches all file reads         |
>
> `Bash(*)` is equivalent to `Bash` and matches all Bash commands. As a deny rule, both forms remove the tool from Claude's context.
>
> ### Use specifiers for fine-grained control
>
> Add a specifier in parentheses to match specific tool uses:
>
> | Rule                           | Effect                                                   |
> | :----------------------------- | :------------------------------------------------------- |
> | `Bash(npm run build)`          | Matches the exact command `npm run build`                |
> | `Read(./.env)`                 | Matches reading the `.env` file in the current directory |
> | `WebFetch(domain:example.com)` | Matches fetch requests to example.com                    |
>
> ### Match by input parameter
>
> Deny and ask rules can match a top-level input parameter on any tool with `Tool(param:value)`. The rule matches when Claude calls the tool with that parameter set to that exact value. An allow rule for one parameter value wouldn't establish that the call is safe overall, so allow rules continue to use each tool's own specifier syntax. This works for any scalar parameter the tool accepts:
>
> | Rule                           | Matches                                      |
> | :----------------------------- | :------------------------------------------- |
> | `Agent(model:opus)`            | Agent calls that request the Opus model tier |
> | `Agent(isolation:worktree)`    | Agent calls that request a git worktree      |
> | `Bash(run_in_background:true)` | Bash calls that run in the background        |
>
> Parameter matching follows these rules:
>
> * The parameter name must be a direct field of the tool's input, such as `model` on the Agent tool. Fields nested inside an object or array are not matchable
> * Each rule names one parameter. To gate on both `model` and `isolation`, write two rules, `Agent(model:opus)` and `Agent(isolation:worktree)`, rather than combining them in one rule
> * The value supports `*` as a wildcard that matches any sequence of characters, so `Agent(isolation:*)` matches any explicit isolation value. Without `*` the match is exact
> * A parameter the model omits is never matched, so `Agent(model:*)` doesn't match a call that leaves `model` unset
> * The value is compared against the literal input Claude sends, before any normalization. `Agent(model:opus)` matches the alias `opus` but not a full model ID. Run with [`--verbose`](/docs/en/cli-reference) to see the exact parameter names and values in each tool call
> * Whitespace around the colon is ignored
>
> You can't match a tool's primary content field this way: `command` for Bash and PowerShell, `file_path` for Read, Edit, and Write, `path` for Grep and Glob, `notebook_path` for NotebookEdit, and `url` for WebFetch. A rule like `Bash(command:rm *)` would be bypassable by a compound command, so Claude Code ignores it and emits a startup warning. Use `Bash(rm *)`, `Read(./path)`, or `WebFetch(domain:host)` instead.
>
> ### Wildcard patterns
>
> A `*` in a Bash rule matches any text, including spaces, so one rule covers a family of commands. A rule with no `*` matches one exact command.
>
>   Put the `*` after the subcommand. In `git log --oneline main`, `git` is the program and `log` is the subcommand, the word that picks what the program does. Claude Code matches everything before the first `*` as written, so those words are what limit the rule: `Bash(git log *)` allows only `git log` commands, and `Bash(git *)` allows every git command.
>
> Write the command you want Claude to run without asking, and replace the parts that vary with `*`. With this configuration, Claude Code runs npm scripts and git commits without asking and refuses git push:
>
> ```json
> {
>   "permissions": {
>     "allow": [
>       "Bash(npm run *)",
>       "Bash(git commit *)"
>     ],
>     "deny": [
>       "Bash(git push *)"
>     ]
>   }
> }
> ```
>
> A `*` can go anywhere in the rule: at the start, in the middle, or at the end. Each row shows a rule, commands it matches, and nearby commands it doesn't match:
>
> | You write              | Matches                                                                              | Doesn't match                          |
> | :--------------------- | :----------------------------------------------------------------------------------- | :------------------------------------- |
> | `Bash(npm run build)`  | `npm run build`                                                                      | `npm run build --watch`                |
> | `Bash(npm run *)`      | `npm run build`, `npm run test --watch`, `npm run`                                   | `npm install`                          |
> | `Bash(git log * main)` | `git log --oneline main`, `git log -5 main`, `git log --output=<file> main`          | `git log main`, `git push origin main` |
> | `Bash(git * main)`     | `git merge main`, `git push origin main`, `git -c core.fsmonitor=<script> diff main` | `git log`                              |
> | `Bash(* --version)`    | `node --version`, `bash -c 'echo hi' --version`                                      | `node -v`                              |
> | `Bash(ls *)`           | `ls -la`, `ls`                                                                       | `lsof`                                 |
> | `Bash(ls*)`            | `ls -la`, `lsof`                                                                     |                                        |
> | `Bash(* --help *)`     | `npm --help x`                                                                       | `npm --help`                           |
>
> Three matching rules produce those rows:
>
> * **The `*` stands in for whatever text is in its place.** In `Bash(git * main)`, it stands in for the subcommand, so Claude Code matches every git subcommand and every option before it. That includes `-c`, which makes git run a program you name. In `Bash(* --version)`, the `*` stands in for the program, so any program matches.
> * **A `*` at the end, with a space before it, also matches the bare command.** `Bash(ls *)` matches `ls`, and `Bash(git log *)` matches `git log`. That holds only when the trailing `*` is the rule's only wildcard: `Bash(* --help *)` matches `npm --help x` but not `npm --help`.
> * **The space before a trailing `*` is part of the rule.** `Bash(ls *)` requires a space after `ls`, so `lsof` doesn't match. `Bash(ls*)` has no space, so it matches `lsof` too.
>
> The `:*` suffix is an equivalent way to write a trailing wildcard, so `Bash(ls:*)` matches the same commands as `Bash(ls *)`.
>
> The permission dialog writes the space-separated form when you select "Yes, and don't ask again" for a command prefix. The `:*` form is only recognized at the end of a pattern. In a pattern like `Bash(git:* push)`, the colon is treated as a literal character and won't match git commands.
>
> ### Tool name wildcards
>
> Deny and ask rules also accept glob patterns in the tool-name position. The pattern must match the full tool name: `"*"` matches every tool, and `"mcp__*"` matches every MCP tool across all servers. A tool matched by a bare-name glob deny rule is removed from Claude's context, the same as a bare tool name, including the [`EndConversation`](/docs/en/tools-reference#endconversation-tool-behavior) exception: a glob deny can't remove it while any other tool remains, and a glob ask never prompts for it. This configuration denies every MCP tool:
>
> ```json
> {
>   "permissions": {
>     "deny": [
>       "mcp__*"
>     ]
>   }
> }
> ```
>
> Allow rules accept tool-name globs only after a literal `mcp__<server>__` prefix. The server segment must be glob-free so the rule names a specific server you configured. `mcp__puppeteer__*` matches every tool from the `puppeteer` server, and `mcp__github__get_*` matches its `get_` tools. An unanchored allow glob such as `"*"`, `"B*"`, or `"mcp__*"` is skipped with a warning and doesn't auto-approve anything.
>
> A deny or ask rule whose tool name matches no known tool produces a startup warning to catch typos. Tool names containing `_` or `*` are exempt from the check.
>
> The label shown for a tool in the transcript and permission dialog can differ from its canonical name. For example, the tool labeled `Stop Task` in the transcript has the canonical name `TaskStop`. Permission rules and [hook matchers](/docs/en/hooks) match the canonical name only, so a rule written as `Stop Task` doesn't match. For deny and ask rules, the startup warning above catches the mismatch. Use the canonical names listed in the [tools reference](/docs/en/tools-reference).
>
> ## Tool-specific permission rules
>
> ### Bash
>
> Bash rules match the whole command text, with `*` standing in for any text. [Wildcard patterns](#wildcard-patterns) shows which commands each rule shape matches and where to put the `*`. The rest of this section covers how Claude Code matches compound commands, wrappers, read-only commands, and redirections.
>
> #### Compound commands
>
>   Claude Code is aware of shell operators, so a rule like `Bash(safe-cmd *)` won't give it permission to run the command `safe-cmd && other-cmd`. The recognized command separators are `&&`, `||`, `;`, `|`, `|&`, `&`, and newlines. A rule must match each subcommand independently.
>
> When you approve a compound command with "Yes, and don't ask again", Claude Code saves a separate rule for each subcommand that requires approval, rather than a single rule for the full compound string. For example, approving `git status && npm test` saves a rule for `npm test`, so future `npm test` invocations are recognized regardless of what precedes the `&&`. Subcommands like `cd` into a subdirectory generate their own Read rule for that path. Up to 5 rules may be saved for a single compound command.
>
>   Wrappers
>
> Before matching Bash rules, Claude Code strips a fixed set of wrappers, so a rule like `Bash(npm test *)` also matches `timeout 30 npm test`. The stripped wrappers are `timeout`, `time`, `nice`, `nohup`, and `stdbuf`, plus the shell builtins `command` and `builtin`, and zsh's `noglob`. Each runs its argument as the actual command. Two related forms aren't stripped: the query form `command -v`, which looks up a command rather than running one, and zsh's `nocorrect`.
>
> Claude Code also strips a leading assignment of certain known-safe environment variables, so `Bash(npm test *)` matches `NODE_ENV=test npm test`. An allow rule won't match past an assignment of any other variable. A deny or ask rule matches past any leading assignment, so `Bash(rm *)` in deny still matches `FOO=bar rm -rf tmp/`.
>
> Bare `xargs` is also stripped, so `Bash(grep *)` matches `xargs grep pattern`. Stripping applies only when `xargs` has no flags: an invocation like `xargs -n1 grep pattern` is matched as an `xargs` command, so rules written for the inner command do not cover it.
>
> This wrapper list is built in and is not configurable. Development environment runners such as `direnv exec`, `devbox run`, `mise exec`, `npx`, and `docker exec` are not in the list. Because these tools execute their arguments as a command, a rule like `Bash(devbox run *)` matches whatever comes after `run`, including `devbox run rm -rf .`. To approve work inside an environment runner, write a specific rule that includes both the runner and the inner command, such as `Bash(devbox run npm test)`. Add one rule per inner command you want to allow.
>
> Exec wrappers such as `watch`, `setsid`, `ionice`, and `flock` can't be auto-approved by a prefix rule like `Bash(watch *)`, so in Manual mode they always prompt. The same applies to `find` with `-exec` or `-delete`: a `Bash(find *)` rule doesn't cover these forms. To approve a specific invocation, write an exact-match rule for the full command string.
>
> #### Read-only commands
>
> Claude Code recognizes a built-in set of Bash commands as read-only and runs them without a permission prompt in every mode. These include `ls`, `cat`, `echo`, `pwd`, `head`, `tail`, `grep`, `find`, `wc`, `which`, `diff`, `stat`, `du`, `cd`, and read-only forms of `git`. The set is not configurable; to require a prompt for one of these commands, add an `ask` or `deny` rule for it.
>
> A redirect such as `ls > out.txt` adds a check on the target. See [Redirections](#redirections).
>
> Unquoted glob patterns are permitted for commands whose every flag is read-only, so `ls *.ts` and `wc -l src/*.py` run without a prompt.
>
> In Manual mode, commands from this set still prompt in these cases:
>
> * **Unquoted globs for commands with write-capable flags**: commands with write-capable or exec-capable flags, such as `find`, `sort`, `sed`, and `git`, prompt when an unquoted glob is present, because the glob could expand to a flag like `-delete`.
> * **`docker` pointed at another daemon**: read-only forms of `docker` prompt when the command carries a flag that selects a different daemon, such as `-H`, `--context`, or Podman's `--url` and `--connection`.
> * **`file` with path-opening flags**: `file` prompts when it passes `-m`/`--magic-file` or `-f`/`--files-from`, because those flags make `file` open the paths named in the flag's value.
> * **Network paths on Windows**: a command whose arguments include a network (UNC) path, such as `\\server\share\file`, prompts because accessing a network path can send your Windows credentials to the host it names. The same check applies to [PowerShell tool](/docs/en/tools-reference#powershell-tool) commands.
> * **Commands the analysis can't parse**: when Claude Code can't fully parse a command, it asks for approval instead of treating the command as read-only. Commands longer than 10,000 characters always prompt because they exceed what the analysis parses.
>
> A `cd` into a path inside your working directory or an [additional directory](#working-directories) is also read-only, and a compound command like `cd packages/api && ls` runs without a prompt when each part qualifies on its own. Two combinations prompt even when each part is read-only:
>
> * **`cd` with `git`**: prompts when the `cd` changes into a different directory, since running `git` in a new directory can execute that directory's hooks. A `cd` whose target resolves to the current working directory is a no-op and doesn't trigger the prompt.
> * **`cd` with an output redirect**: prompts when Claude Code can't determine which directory the redirect target resolves against after the `cd` runs. A command whose only redirect target is `/dev/null`, such as `cd app; grep -r pattern . 2>/dev/null`, doesn't prompt, because `/dev/null` doesn't depend on the working directory.
>
>   Bash permission patterns that try to constrain command arguments are fragile. For example, `Bash(curl http://github.com/ *)` intends to restrict curl to GitHub URLs, but won't match variations like:
>
>   * Options before URL: `curl -X GET http://github.com/...`
>   * Different protocol: `curl https://github.com/...`
>   * Redirects: `curl -L http://short.example.com/xyz`, which redirects to GitHub
>   * Variables: `URL=http://github.com && curl $URL`
>   * Extra spaces: `curl  http://github.com`
>
>   For more reliable URL filtering, consider:
>
>   * **Restrict Bash network tools**: use deny rules to block `curl`, `wget`, and similar commands, then use the WebFetch tool with `WebFetch(domain:github.com)` permission for allowed domains
>   * **Use PreToolUse hooks**: implement a hook that validates URLs in Bash commands and blocks disallowed domains
>   * **Add CLAUDE.md guidance**: describe your allowed curl patterns in `CLAUDE.md`. This shapes what Claude tries but doesn't enforce a boundary, so pair it with one of the options above
>
>   Note that using WebFetch alone doesn't prevent network access. If Bash is allowed, Claude can still use `curl`, `wget`, or other tools to reach any URL.
>
> #### Redirections
>
> Claude Code checks the target of an output redirection, such as `>`, `>>`, or `2>`, as a file write. The check covers your `Edit` allow and deny rules, [protected paths](/docs/en/permission-modes#protected-paths), and the [working directories](#working-directories). A rule such as `Bash(git commit *)` allows the command, not the target. A `/dev/null` target isn't checked. A target that starts with `~` or contains a glob character needs approval.
>
> ### PowerShell
>
> PowerShell permission rules use the same shape as Bash rules. Wildcards with `*` match at any position, the `:*` suffix is equivalent to a trailing ` *`, and a bare `PowerShell` or `PowerShell(*)` matches every command. This configuration allows `Get-ChildItem` and `git commit` commands while blocking `Remove-Item`:
>
> ```json
> {
>   "permissions": {
>     "allow": [
>       "PowerShell(Get-ChildItem *)",
>       "PowerShell(git commit *)"
>     ],
>     "deny": [
>       "PowerShell(Remove-Item *)"
>     ]
>   }
> }
> ```
>
> Common aliases are canonicalized before matching. A rule written for the cmdlet name also matches its aliases, so `PowerShell(Get-ChildItem *)` matches `gci`, `ls`, and `dir` as well. Matching is case-insensitive.
>
> Claude Code parses the PowerShell AST and checks each command in a compound command independently. Pipeline operators `|`, statement separators `;`, and on PowerShell 7+ the chain operators `&&` and `||` split a compound command into subcommands. A rule must match every subcommand for the compound command to be allowed.
>
> ### Read and Edit
>
> To block Claude's file tools from reading a file or directory, add a `Read` deny rule for its path, such as `Read(./.env)` or `Read(./secrets/**)`; [Exclude sensitive files](/docs/en/settings-reference#exclude-sensitive-files) has a paste-ready example.
>
> `Edit` rules apply to all built-in tools that edit files. Claude makes a best-effort attempt to apply `Read` rules to all built-in tools that read files like Grep and Glob, to `@file` mentions in your prompts, and to the selection and open-file context that a connected [IDE](/docs/en/vs-code#the-built-in-ide-mcp-server) shares with Claude.
>
> A `Read` deny rule also blocks the [Edit and Write tools](/docs/en/errors#file-is-covered-by-a-read-deny-rule) on the same path, including creating a new file there. NotebookEdit isn't covered, so add an `Edit` deny rule for paths no tool may change. The check requires Claude Code v2.1.208 or later on edits, and v2.1.228 or later on writes.
>
> Claude Code checks file permissions against `Edit(path)` and `Read(path)` rules only. If you write a path rule for `Write`, `NotebookEdit`, `Glob`, or the legacy `MultiEdit` tool instead, Claude Code accepts the rule but never consults it, and [warns at startup](/docs/en/errors#is-not-matched-by-file-permission-checks), except for a `Glob` rule passed in `--allowedTools`. Use `Edit(docs/**)` in place of `Write(docs/**)`, `NotebookEdit(docs/**)`, or `MultiEdit(docs/**)`, and `Read(docs/**)` in place of `Glob(docs/**)`. Claude Code doesn't warn about a tool-name rule with no path, such as a deny rule for `Write`; it matches that rule at the tool level everywhere. Requires Claude Code v2.1.210 or later.
>
>   Read and Edit deny rules apply to Claude's built-in file tools and to file commands Claude Code recognizes in Bash, such as `cat`, `head`, `tail`, and `sed`. They don't apply to arbitrary subprocesses that read or write files indirectly, like a Python or Node script that opens files itself. For OS-level enforcement that blocks all processes from accessing a path, [enable the sandbox](/docs/en/sandboxing).
>
> Read and Edit rules both use [gitignore](https://git-scm.com/docs/gitignore) pattern syntax with four distinct pattern types; for single-segment directory patterns, the matching depth also depends on the rule type, described later in this section:
>
> | Pattern            | Meaning                              | Example                          | Matches                                          |
> | ------------------ | ------------------------------------ | -------------------------------- | ------------------------------------------------ |
> | `//path`           | Absolute path from filesystem root   | `Read(//Users/alice/secrets/**)` | `/Users/alice/secrets/**`                        |
> | `~/path`           | Path from home directory             | `Read(~/Documents/*.pdf)`        | `/Users/alice/Documents/*.pdf`                   |
> | `/path`            | Path relative to the settings source | `Edit(/src/**/*.ts)`             | `<project root>/src/**/*.ts` in project settings |
> | `path` or `./path` | Path relative to current directory   | `Read(*.env)`                    | `<cwd>/*.env`                                    |
>
>   A pattern like `/Users/alice/file` isn't an absolute path. The single leading slash anchors at the settings source, not the filesystem root. Use `//Users/alice/file` for absolute paths.
>
> A `/path` pattern anchors at a directory associated with the settings source that defines it, so the same rule matches different locations depending on where you put it:
>
> | Rule defined in                                 | `/path` resolves to        |
> | :---------------------------------------------- | :------------------------- |
> | Project settings at `.claude/settings.json`     | `<project root>/path`      |
> | Local settings at `.claude/settings.local.json` | `<original cwd>/path`      |
> | User settings at `~/.claude/settings.json`      | `~/.claude/path`           |
> | A file passed with `--settings <file>`          | `<directory of file>/path` |
> | CLI flags, `/permissions`, or session rules     | `<original cwd>/path`      |
>
> Local settings rules anchor at the directory you started Claude Code from, not at the repository root where Claude Code [stores the file](#permission-system) in v2.1.211 and later. In a session started at the repository root, the two directories are the same; in a [worktree](/docs/en/worktrees) session, a shared rule such as `Edit(/src/**)` matches that worktree's own `src/` directory.
>
> A deny rule such as `Read(/secrets/**)` in user settings blocks `~/.claude/secrets/**`, not a `secrets` directory in your project. To write a rule in user settings that applies inside every project, use a `//` absolute path or a `~/` home-relative path instead.
>
> On Windows, paths are normalized to POSIX form before matching. `C:\Users\alice` becomes `/c/Users/alice`, so use `//c/**/.env` to match `.env` files anywhere on that drive. To match across all drives, use `//**/.env`.
>
> Examples:
>
> * `Edit(/docs/**)`: edits in `<project>/docs/`, not `/docs/` or `<project>/.claude/docs/`
> * `Read(~/.zshrc)`: reads your home directory's `.zshrc`
> * `Edit(//tmp/scratch.txt)`: edits the absolute path `/tmp/scratch.txt`
> * `Read(src/**)`: as an allow rule, reads from `<current-directory>/src/` only; as a deny or ask rule, matches a `src` directory at any depth under the current directory
>
> A rule only matches files under its anchor; within that bound, matching depth depends on the pattern shape and, for single-segment directory patterns, the rule type, described below. Bare filenames follow gitignore semantics and match at any depth, so `Read(.env)` and `Read(**/.env)` are equivalent:
>
> | Deny rule                       | Blocks                                       | Does not block                                       |
> | ------------------------------- | -------------------------------------------- | ---------------------------------------------------- |
> | `Read(.env)` or `Read(**/.env)` | any `.env` at or under the current directory | `.env` in a parent directory or another project      |
> | `Read(//**/.env)`               | any `.env` anywhere on the filesystem        | nothing; the rule is anchored at the filesystem root |
>
> A relative pattern with a single directory segment, such as `src/**`, matches at different depths depending on the rule type:
>
> * **Allow rules**: `Edit(src/**)` matches only `<cwd>/src` and the files under it. To allow a directory name at any depth, write `Edit(**/src/**)`.
> * **Deny and ask rules**: `Read(secrets/**)` matches a directory named `secrets` at any depth under the current directory, so the rule also applies to nested copies.
>
> Every other pattern shape matches at the same depth in every rule type: `Edit(/src/**)` and `Edit(src/components/**)` match only at their anchored location, while `Edit(**/src/**)` matches at any depth.
>
> The following example shows each pattern shape against a project with a top-level `src/` directory and a nested copy under `vendor/`:
>
> ```text
>
> ├── src/
> │   └── app.ts
> └── vendor/
>     └── pkg/
>         └── src/
>             └── lib.js
> ```
>
> | Rule                                 | Matches `src/app.ts` | Matches `vendor/pkg/src/lib.js` |
> | :----------------------------------- | :------------------- | :------------------------------ |
> | `Edit(src/**)` as an allow rule      | Yes                  | No                              |
> | `Edit(src/**)` as a deny or ask rule | Yes                  | Yes                             |
> | `Edit(/src/**)` in any rule type     | Yes                  | No                              |
> | `Edit(**/src/**)` in any rule type   | Yes                  | Yes                             |
>
>   In gitignore patterns, `*` matches within a single path segment and can appear at any position in the pattern, while `**` matches across directories.
>
> When you approve a file path with "Yes, and don't ask again", Claude Code escapes gitignore pattern characters in that path, such as `[`, `]`, and `*`, so the generated rule matches only the literal path you approved. Rules you write yourself aren't escaped. Before v2.1.202, Claude Code saved the path unescaped, so a generated rule for a directory named `[2024-06] Reports` could fail to match its own path or match unintended sibling directories.
>
> When Claude accesses a symlink, permission rules check two paths: the symlink itself and the file it resolves to. Allow and deny rules treat that pair differently: allow rules fall back to prompting you, while deny rules block outright.
>
> * **Allow rules**: apply only when both the symlink path and its target match. A symlink inside an allowed directory that points outside it still prompts you.
> * **Deny rules**: apply when either the symlink path or its target matches. A symlink that points to a denied file is itself denied.
>
> For example, with `Read(./project/**)` allowed and `Read(~/.ssh/**)` denied, a symlink at `./project/key` pointing to `~/.ssh/id_rsa` is blocked: the target fails the allow rule and matches the deny rule.
>
> ### WebFetch
>
> WebFetch rules use a `domain:` prefix and match against the hostname of the requested URL. Matching is case-insensitive, supports `*` wildcards, and strips a trailing `.` from both the rule and the hostname so `example.com.` and `example.com` are treated the same.
>
> * `WebFetch(domain:example.com)` matches requests to `example.com`
> * `WebFetch(domain:*.example.com)` matches any subdomain at any depth, such as `api.example.com` or `a.b.example.com`, but not `example.com` itself
> * `WebFetch(domain:*)` matches every domain. It isn't the same as a bare `WebFetch` rule; see [Allow or deny every fetch](#allow-or-deny-every-fetch)
>
> In any position other than a leading `*.` or a bare `*`, the wildcard matches only the text between two dots. `WebFetch(domain:example.*)` matches `example.org`, where `*` becomes `org`, but not `example.evil.com`, where `*` would have to become `evil.com` and cross a dot. This keeps a trailing wildcard from matching domains an attacker could register.
>
> Wildcards in `WebFetch` rules require Claude Code v2.1.172 or later to match fetches.
>
> #### Allow or deny every fetch
>
> A bare `WebFetch` rule is the tool name with no `domain:` part, such as `"deny": ["WebFetch"]`. Both it and `WebFetch(domain:*)` cover every URL, but Claude Code applies them differently, and only the `domain:` form also adds its domain to the sandbox's [allowed or denied domain list](/docs/en/sandboxing#network-isolation). That section lists the wildcard forms the sandbox honors and the version that added bare `*`.
>
> Each row shows what a rule does in the `allow` list and in the `deny` list:
>
> | Rule                 | In `allow`                                                                                     | In `deny`                                                                                                                       |
> | :------------------- | :--------------------------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------ |
> | `WebFetch`           | Claude fetches without prompting you. Doesn't change which hosts sandboxed commands can reach. | Claude Code removes the `WebFetch` tool, so Claude can't fetch at all. Doesn't change which hosts sandboxed commands can reach. |
> | `WebFetch(domain:*)` | Claude fetches without prompting you, and sandboxed commands can reach any host.               | Claude Code keeps the tool and refuses each fetch, and sandboxed commands can't reach any host.                                 |
>
> To let Claude fetch freely while keeping the sandbox allowlist as it is, use the bare form. This `settings.json` does that:
>
> ```json
> {
>   "permissions": {
>     "allow": ["WebFetch"]
>   }
> }
> ```
>
> When you ask Claude to fetch a page, it fetches without a prompt. When you ask it to run a [sandboxed](/docs/en/sandboxing) `curl` against a host outside the sandbox allowlist, Claude Code still prompts you for that host, or in [auto mode](/docs/en/permission-modes#eliminate-prompts-with-auto-mode) sends the request to the classifier, because the bare rule didn't add the host to the allowlist.
>
> ### MCP
>
> MCP rules use the server name as configured in Claude Code, optionally followed by the name of a tool from that server.
>
> * `mcp__puppeteer` matches any tool provided by the `puppeteer` server
> * `mcp__puppeteer__*` uses wildcard syntax and also matches all tools from the `puppeteer` server
> * `mcp__puppeteer__puppeteer_navigate` matches the `puppeteer_navigate` tool provided by the `puppeteer` server
>
> If your organization has set a [claude.ai connector](/docs/en/mcp#organization-controls-on-connector-tools) tool to `ask`, allow rules for that tool don't take effect: Claude Code prompts on every call, even in `auto` and `bypassPermissions` modes. In `dontAsk` mode, which never prompts, Claude Code denies the call instead. Connector tools appear as `mcp__claude_ai_<server>__<tool>`.
>
> ### Agent (subagents)
>
> Use `Agent(AgentName)` rules to control which [subagents](/docs/en/sub-agents) Claude can use:
>
> * `Agent(Explore)` matches the Explore subagent
> * `Agent(Plan)` matches the Plan subagent
> * `Agent(my-custom-agent)` matches a custom subagent named `my-custom-agent`
>
> Add these rules to the `deny` array in your settings or use the `--disallowedTools` CLI flag to disable specific agents. To disable the Explore agent:
>
> ```json
> {
>   "permissions": {
>     "deny": ["Agent(Explore)"]
>   }
> }
> ```
>
> ### Cd
>
> `Cd` rules control which directories the [`/cd` command](/docs/en/commands) can move the session to. `Cd` is not a model-invocable tool: Claude can't call it, and the rules apply only when you run `/cd` yourself.
>
> A bare `Cd` deny rule disables `/cd` entirely. A `Cd(<path-pattern>)` deny rule blocks matching targets. Deny rules check every spelling of the target, including each symlink hop it resolves through, so a rule written for one path also blocks targets that resolve to it.
>
> Adding any `Cd` allow rule switches `/cd` to allowlist mode: the resolved target directory must match one of your allow rules, or `/cd` refuses. With no `Cd` rules configured, `/cd` keeps its default behavior and prompts you to trust an unfamiliar directory.
>
> Path patterns share the `//`, `~/`, and `/` anchors from [Read and Edit rules](#read-and-edit), but matching is anchored to the whole directory path rather than gitignore-style. `*` matches exactly one path segment and `**` matches across segments. A trailing `/**` also matches its named root.
>
> | Rule                  | Matches                                   | Does not match               |
> | --------------------- | ----------------------------------------- | ---------------------------- |
> | `Cd(~/code/*)`        | `~/code/app`                              | `~/code/app/src`, `~/code`   |
> | `Cd(~/code/**)`       | `~/code` and any directory under it       | directories outside `~/code` |
> | `Cd(**/node_modules)` | any `node_modules` directory at any depth | `node_modules/pkg`           |
>
> ## Extend permissions with hooks
>
> [Claude Code hooks](/docs/en/hooks-guide) let you register custom shell commands that evaluate permissions at runtime. When Claude Code makes a tool call, PreToolUse hooks run before the permission prompt, for every tool except [`EndConversation`](/docs/en/tools-reference#endconversation-tool-behavior). The hook output can deny the tool call, force a prompt, or skip the prompt to let the call proceed.
>
> Hook decisions don't bypass permission rules. Claude Code evaluates deny and ask rules regardless of what a PreToolUse hook returns: a matching deny rule blocks the call, and a matching ask rule still prompts even when the hook returned `"allow"` or `"ask"`. This preserves the deny-first precedence described in [Manage permissions](#manage-permissions), including deny rules set in managed settings.
>
> Connector tools [your organization set to `ask`](/docs/en/mcp#organization-controls-on-connector-tools) and MCP tools marked [`requiresUserInteraction`](/docs/en/mcp#require-approval-for-a-specific-tool) also still prompt when a hook returns `"allow"`.
>
> A blocking hook also takes precedence over allow rules. A hook that exits with code 2 stops the tool call before permission rules are evaluated, so the block applies even when an allow rule would otherwise let the call proceed. To run all Bash commands without prompts except for a few you want blocked, add `"Bash"` to your allow list and register a PreToolUse hook that rejects those specific commands. See [Block edits to protected files](/docs/en/hooks-guide#block-edits-to-protected-files) for a hook script you can adapt.
>
> ## Working directories
>
> By default, Claude has access to files in the directory where you launched it. You can extend this access:
>
> * **During startup**: use `--add-dir <path>` CLI argument
> * **During session**: use `/add-dir` command
> * **Persistent configuration**: add to `additionalDirectories` in [settings files](/docs/en/settings#where-settings-live)
>
> Files in additional directories follow the same permission rules as the original working directory: they become readable without prompts, and file editing permissions follow the current permission mode.
>
> In background sessions on macOS, the session host requests access to protected folders such as `~/Desktop`, `~/Documents`, and `~/Downloads` separately from your terminal when Claude needs to read or write files there; if reads there fail with `Operation not permitted`, see [how to grant folder access to background sessions](/docs/en/agent-view#background-sessions-can’t-read-desktop-documents-or-downloads-on-macos).
>
> To change the session's primary working directory instead of adding another, use [`/cd`](/docs/en/commands). The `/cd` command requires Claude Code v2.1.169 or later. Unlike `/add-dir`, it relocates the session: the new directory's `CLAUDE.md` is loaded and `--resume` finds the session from there.
>
> ### Additional directories grant file access, not configuration
>
> Adding a directory extends where Claude can read and edit files. It doesn't make that directory a full configuration root: most `.claude/` configuration is not discovered from additional directories, though a few types are loaded as exceptions.
>
> These exceptions apply only to directories added with the `--add-dir` flag or the `/add-dir` command, including directories the Agent SDK adds through the flag. Directories listed in `permissions.additionalDirectories` in a settings file grant file access only and don't load any of the configuration below.
>
> The Agent SDK's [`additionalDirectories`](/docs/en/agent-sdk/typescript#options) option in TypeScript and [`add_dirs`](/docs/en/agent-sdk/python#claudeagentoptions) option in Python receive the exceptions too, even though the TypeScript option shares its name with the settings key. The SDK passes each entry to Claude Code as `--add-dir`, so those directories behave like flag-added directories. Skills, commands, and subagents from any flag-added directory load through the `project` [setting source](/docs/en/agent-sdk/claude-code-features#control-filesystem-settings-with-settingsources), so they don't load when you exclude that source with [`--setting-sources`](/docs/en/cli-reference) on the CLI or `settingSources` in the SDK, and [bare mode](/docs/en/headless#start-faster-with-bare-mode) skips the commands and subagents among them.
>
> The following configuration types are loaded from `--add-dir` directories:
>
> | Configuration                                                                         | Loaded from `--add-dir`### Source: Permission modes
>
> > ## Documentation Index
> > Fetch the complete documentation index at: https://code.claude.com/docs/llms.txt
> > Use this file to discover all available pages before exploring further.
>
> # Choose a permission mode
>
> > Control whether Claude asks before acting. Switch permission modes with Shift+Tab in the CLI, the mode indicator in VS Code, or the mode selector in Desktop.
>
> A permission mode sets which actions Claude can take in a session without asking you first. In Manual mode, Claude Code stops and asks you before most actions that edit files, run shell commands, or reach the network. In [auto mode](#eliminate-prompts-with-auto-mode), a second model, the classifier, reviews actions instead of you; [how the classifier evaluates actions](#how-the-classifier-evaluates-actions) lists which actions it reviews and which skip it.
>
> On Pro, Max, and Team plans, the built-in starting permission mode is auto mode. [Which mode a session starts in](#which-mode-a-session-starts-in) covers the surfaces and settings that change the starting permission mode. You can also change a running session's permission mode at any time.
>
> ## Available modes
>
> Each mode makes a different tradeoff between convenience and oversight. The table below shows what Claude can do without a permission prompt in each mode. Manual mode appears under its config value, `default`.
>
> | Mode                                                                | What runs without asking                                                                                  | Best for                                        |
> | :------------------------------------------------------------------ | :-------------------------------------------------------------------------------------------------------- | :---------------------------------------------- |
> | `default`                                                           | Reads only                                                                                                | Reviewing every action yourself, sensitive work |
> | [`acceptEdits`](#auto-approve-file-edits-with-acceptedits-mode)     | Reads, file edits, and common filesystem commands (`mkdir`, `touch`, `mv`, `cp`, etc.)                    | Iterating on code you're reviewing              |
> | [`plan`](#analyze-before-you-edit-with-plan-mode)                   | Reads, plus classifier-approved commands when [auto mode](#eliminate-prompts-with-auto-mode) is available | Exploring a codebase before changing it         |
> | [`auto`](#eliminate-prompts-with-auto-mode)                         | Everything, with background safety checks                                                                 | Long tasks, reducing prompt fatigue             |
> | [`dontAsk`](#allow-only-pre-approved-tools-with-dontask-mode)       | Only pre-approved tools                                                                                   | Locked-down CI and scripts                      |
> | [`bypassPermissions`](#skip-all-checks-with-bypasspermissions-mode) | Everything                                                                                                | Isolated containers and VMs only                |
>
> The mode that reviews every action is named **Manual** in the CLI, in `claude --help`, in the VS Code and JetBrains extensions, and in the desktop app. Its config value is `default`, which is what hooks and SDK integrations use. The CLI accepts `manual` as an alias wherever you type the value, for example `claude --permission-mode manual` or `"defaultMode": "manual"`. The Manual label and the `manual` alias require Claude Code v2.1.200 or later. The desktop app's label doesn't depend on your CLI version.
>
> Writes to [protected paths](#protected-paths) are never auto-approved except in `bypassPermissions` mode and in plan-mode sessions where bypass permissions are available, meaning sessions started in a way that [puts `bypassPermissions` in the mode cycle](#switch-permission-modes).
>
> Modes set the baseline. Layer [permission rules](/docs/en/permissions#manage-permissions) on top to pre-approve or block specific tools. Deny rules block in every mode, including `bypassPermissions`. Deny and ask rules don't apply to [`EndConversation`](/docs/en/tools-reference#endconversation-tool-behavior) as long as Claude still has at least one other tool it can call. Allow rules have no effect in `bypassPermissions`.
>
>   Actions no mode auto-approves
>
> Claude Code doesn't auto-approve the following in any mode, including `bypassPermissions`. Each bullet links to the section that says what happens instead in each mode:
>
> * Tools matched by an explicit [ask rule](/docs/en/permissions#manage-permissions)
> * Connector tools your organization [set to `ask`](/docs/en/mcp#organization-controls-on-connector-tools)
> * Tools that require user interaction: the built-in `AskUserQuestion` tool and MCP tools marked [`requiresUserInteraction`](/docs/en/mcp#require-approval-for-a-specific-tool)
> * `rm` and `rmdir` removals targeting a [critical path](#critical-paths), which no allow rule or `PreToolUse` hook `"allow"` approves
> * The [cross-session messaging safeguards](#skip-all-checks-with-bypasspermissions-mode)
>
> ## Common setups
>
> Permission modes decide whether Claude asks before an action, and the [Bash sandbox](/docs/en/sandboxing) and outer [isolation boundaries](/docs/en/sandbox-environments) decide what an action can reach once it runs. Each row below pairs a goal with the flags or settings that get you there and the isolation it needs, as a starting point. [Available modes](#available-modes) lists what runs without a prompt in each mode, and the per-mode section each row links carries the full behavior.
>
> | You want to                                              | Start with                                                                                                                                                          | Isolation needed                                                                                                                                                                             | Notes                                                                                                                                                                                                                               |
> | :------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------ | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | Review every action yourself                             | Manual mode: `claude --permission-mode default`                                                                                                                     | None                                                                                                                                                                                         | Sensitive work, unfamiliar code                                                                                                                                                                                                     |
> | Iterate locally with fewer prompts, without a classifier | Manual mode plus the Bash sandbox in [auto-allow mode](/docs/en/sandboxing#sandbox-modes): `claude --permission-mode default`, then run `/sandbox` and select auto-allow | The built-in Bash sandbox, on macOS, Linux, and WSL2                                                                                                                                         | Deny rules still apply, and ask rules that name a command, such as `Bash(git push *)`, still prompt. To turn the sandbox on from a settings file instead, set [`sandbox.enabled`](/docs/en/settings-reference#sandbox-enabled) to `true` |
> | Explore before changing anything                         | `claude --permission-mode plan`                                                                                                                                     | None                                                                                                                                                                                         | Claude Code blocks edits until you [approve a plan](#review-and-approve-a-plan)                                                                                                                                                     |
> | Work hands-off with a classifier reviewing each action   | `claude --permission-mode auto`, the [built-in starting permission mode](#which-mode-a-session-starts-in) on Pro, Max, and Team                                     | None; a sandbox or container adds defense in depth                                                                                                                                           | Requires a [supported model](#eliminate-prompts-with-auto-mode), and your organization can [turn auto mode off](#eliminate-prompts-with-auto-mode)                                                                                  |
> | Run in CI with an exact allowlist                        | `claude -p "run the test suite" --permission-mode dontAsk --allowedTools "Bash(npm test)" "Read"`                                                                   | None beyond what your CI runner provides                                                                                                                                                     | [Claude Code on the web](/docs/en/claude-code-on-the-web) ignores `dontAsk` from settings files                                                                                                                                          |
> | Run fully unattended inside a container                  | `claude -p "<prompt>" --dangerously-skip-permissions`                                                                                                               | Required: a container, VM, or the [sandbox runtime](/docs/en/sandbox-environments#sandbox-runtime); on Linux and macOS, run it as a [non-root user](#skip-all-checks-with-bypasspermissions-mode) | Claude Code on the web ignores this mode from settings files. In this `-p` run, the [few calls that would still prompt](#skip-all-checks-with-bypasspermissions-mode) are denied instead                                            |
>
> The Bash sandbox and auto mode work independently and combine, except in plan mode, where [auto-allow doesn't widen approvals](/docs/en/sandboxing#sandbox-modes). For the full interaction, see [How sandboxing relates to permissions and permission modes](/docs/en/sandboxing#how-sandboxing-relates-to-permissions-and-permission-modes) and [How isolation relates to permission modes](/docs/en/sandbox-environments#how-isolation-relates-to-permission-modes).
>
>   Which mode a session starts in
>
> When you start a new session in a terminal, Claude Code takes the permission mode from the first of these that applies:
>
> 1. The `--permission-mode` flag, or `--dangerously-skip-permissions`
> 2. `permissions.defaultMode` in a [settings file](/docs/en/settings#where-settings-live). An `"auto"` value in `.claude/settings.json` or `.claude/settings.local.json` doesn't take effect, and Claude Code then uses the built-in default rather than a `defaultMode` from `~/.claude/settings.json`. The other values apply from any settings file
> 3. The built-in default
>
> Conversations the VS Code extension starts follow the extension's own list in [Switch permission modes](#switch-permission-modes). A session you resume keeps the permission mode it was in unless you pass `--permission-mode` or `--dangerously-skip-permissions`; [what a resumed session restores](/docs/en/sessions#what-a-resumed-session-restores) lists the exceptions.
>
> The built-in `auto` default requires Claude Code v2.1.228 or later on macOS, Linux, and WSL, and v2.1.233 or later on native Windows. On earlier versions, the built-in default is Manual.
>
> The built-in default depends on how you run Claude Code, on your plan, and on whether Claude Code could fetch its feature flags. The first row that matches your session applies. The table covers sessions you start in a terminal or through the VS Code extension; for the desktop app and claude.ai, see the Desktop and Web tabs in [Switch permission modes](#switch-permission-modes).
>
> | How you run Claude Code                                                                                                                                                                                                | Built-in starting permission mode |
> | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :-------------------------------- |
> | Any settings file sets `disableAutoMode` to `"disable"`                                                                                                                                                                | `default`                         |
> | [Feature-flag fetching](/docs/en/env-vars#features-that-need-feature-flag-fetching) is off                                                                                                                                  | `default`                         |
> | Your [first session after you install Claude Code or upgrade](/docs/en/env-vars#first-session-after-an-install-or-upgrade) to a version that adds this default, unless a non-interactive session picks the flags up in time | `default`                         |
> | `claude -p` or the [Agent SDK](/docs/en/agent-sdk/permissions)                                                                                                                                                              | `default`                         |
> | Amazon Bedrock, Google Cloud's Agent Platform, Microsoft Foundry, [Claude Platform on AWS](/docs/en/claude-platform-on-aws), or a signed-in [Claude apps gateway](/docs/en/claude-apps-gateway) session                          | `default`                         |
> | A Pro, Max, or Team plan, in a terminal or through the [VS Code extension](/docs/en/vs-code)                                                                                                                                | `auto`                            |
> | An Enterprise plan or a Claude Console API key                                                                                                                                                                         | `default`                         |
>
> When feature-flag fetching is off, or in a [first session after an install or upgrade](/docs/en/env-vars#first-session-after-an-install-or-upgrade) where the flags haven't arrived yet, the VS Code extension ignores every settings file when choosing the starting permission mode.
>
> When the flag, a settings file, or the built-in default selects `auto` but auto mode isn't available to the session, Claude Code starts the session in Manual instead. Auto mode is unavailable when a settings file [turns it off](#eliminate-prompts-with-auto-mode) or the model doesn't support it.
>
> The first time the built-in default starts one of your sessions in auto mode, Claude Code shows a notice that links to this page:
>
> * In a terminal, once, at the top of the session
> * In the VS Code extension, as a card on the new-conversation screen that stays until you dismiss it
>
> On Pro, Max, and Team plans, if your `~/.claude/settings.json` sets a `defaultMode` other than `auto` and no other settings file sets one, your sessions keep starting in that mode. Claude Code asks once, in the terminal or in the VS Code extension, whether to change the setting to auto mode. If you decline, your setting stays as it is.
>
>   Start in a different permission mode
>
> You can set the starting permission mode for one session, or as a default for every session on a machine, in a project, or in an organization. When more than one settings file sets `permissions.defaultMode`, [settings precedence](/docs/en/settings#settings-precedence) decides, so a project or managed value outranks `~/.claude/settings.json`. To change the permission mode of a session that's already running, see [Switch permission modes](#switch-permission-modes).
>
> | To set the starting permission mode for          | Do this                                                                                                                                                                                                                                                                                                                                                        |
> | :----------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
> | One session you're about to start                | Pass the permission mode as a flag, for example `claude --permission-mode default`                                                                                                                                                                                                                                                                             |
> | Every terminal session you start on this machine | Set `permissions.defaultMode` in `~/.claude/settings.json`. For what the VS Code extension reads, see [Switch permission modes](#switch-permission-modes)                                                                                                                                                                                                      |
> | Every terminal session you start in one project  | Set `permissions.defaultMode` in the project's `.claude/settings.json`. Sessions you start in a terminal honor every value except `auto`; sessions the VS Code extension starts don't read project settings for the starting permission mode                                                                                                                   |
> | Every terminal session in your organization      | Set `permissions.defaultMode` in [managed settings](/docs/en/managed-settings). Terminal sessions start in that mode and people can still switch to auto mode; for what the VS Code extension reads, see [Switch permission modes](#switch-permission-modes). To remove auto mode so nobody can select it, set `permissions.disableAutoMode` to `"disable"` instead |
>
> This example makes every terminal session on your machine start in Manual mode, whose config value is `default`. Save it in `~/.claude/settings.json`:
>
> ```json
> {
>   "permissions": {
>     "defaultMode": "default"
>   }
> }
> ```
>
> The next session you start shows `⏸ manual mode on` in the status bar.
>
> ## Switch permission modes
>
> Each interface has its own control for switching permission modes during a session and its own way of choosing the permission mode new sessions start in. Asking Claude in chat to change the permission mode doesn't work. Select your interface to see its controls.
>
>   #### CLI
>
> **During a session**: press `Shift+Tab` to cycle permission modes. From `auto`, the first press switches to `default`, and the cycle then runs `default` → `acceptEdits` → `plan` → back to `default`. Optional modes, described below, slot in after `plan`. The status bar shows the active mode as a gray `⏸ manual mode on` for `default`, or as `⏵⏵ accept edits on`, `⏸ plan mode on`, `⏵⏵ auto mode on`, `⏵⏵ don't ask on`, or `⏵⏵ bypass permissions on`.
>
>     Not every mode is in the default cycle:
>
>     * `auto`: appears when your account meets the [auto mode requirements](#eliminate-prompts-with-auto-mode); cycling to it switches permission modes without a confirmation prompt
>     * `bypassPermissions`: appears after you start with `--permission-mode bypassPermissions`, `--dangerously-skip-permissions`, `--allow-dangerously-skip-permissions`, or `permissions.defaultMode: "bypassPermissions"` in [settings](/docs/en/settings-reference#permission-settings); the `--allow-` variant adds the permission mode to the cycle without activating it
>     * `dontAsk`: never appears in the cycle; set it with `--permission-mode dontAsk`
>
>     Enabled optional modes slot in after `plan`, with `bypassPermissions` first and `auto` last. If you have both enabled, you will cycle through `bypassPermissions` on the way to `auto`.
>
>     **At startup**: pass the permission mode as a flag.
>
>     ```bash
>     claude --permission-mode plan
>     ```
>
>     **As a default**: set `permissions.defaultMode` at the scope you want, as described in [Start in a different permission mode](#start-in-a-different-mode).
>
>     The same `--permission-mode` flag works with `-p` for [non-interactive runs](/docs/en/headless).
>
>   #### VS Code
>
> **During a session**: click the mode indicator at the bottom of the prompt box. It uses these labels for the modes on this page:
>
>     | UI label           | Mode                |
>     | :----------------- | :------------------ |
>     | Manual             | `default`           |
>     | Edit automatically | `acceptEdits`       |
>     | Plan               | `plan`              |
>     | Auto               | `auto`              |
>     | Bypass permissions | `bypassPermissions` |
>
>     **As a default**: to pin the permission mode conversations start in, set `claudeCode.initialPermissionMode` in your VS Code user settings to `default`, `manual`, `acceptEdits`, `plan`, or `bypassPermissions`. The setting doesn't accept `auto`; to start in Auto, leave it unset and pick **Auto** from the mode indicator once, as item 2 below describes. The extension starts each new conversation in the first of these that applies:
>
>     1. `claudeCode.initialPermissionMode`
>     2. The mode you last picked from the mode indicator, if it was Manual, Edit automatically, or Auto. Picking Plan or Bypass permissions applies to that conversation only
>     3. `permissions.defaultMode` from [managed settings](/docs/en/managed-settings) or `~/.claude/settings.json`, on Pro, Max, and Team plans with [feature-flag fetching](#which-mode-a-session-starts-in) available
>     4. The [built-in default](#which-mode-a-session-starts-in) for your plan, provider, and organization settings
>
>     The extension never reads a project's `.claude/settings.json` or `.claude/settings.local.json` for the starting permission mode, and in conversations that don't meet item 3's conditions it reads no settings file at all. When `claudeCode.claudeProcessWrapper` is set, items 3 and 4 don't apply either: those conversations start in Manual unless item 1 or item 2 sets a permission mode.
>
>     Auto mode appears in the mode indicator when your account meets every requirement listed in the [auto mode section](#eliminate-prompts-with-auto-mode).
>
>     Bypass permissions requires the **Allow dangerously skip permissions** toggle in the extension settings. Without it, the permission mode doesn't appear in the indicator, and a `bypassPermissions` value from item 1 or item 3 starts the conversation in Manual instead. Auto from any item likewise starts the conversation in Manual when auto mode isn't available.
>
>     See the [VS Code guide](/docs/en/vs-code) for extension-specific details.
>
>   #### JetBrains
>
> The JetBrains plugin runs Claude Code in the IDE terminal, so switching permission modes works the same as in the CLI: press `Shift+Tab` to cycle, or pass `--permission-mode` when launching.
>
>   #### Desktop
>
> **During a session**: in the Code tab, use the mode selector next to the send button. Not every mode appears in the selector:
>
>     * **Auto**: appears when your account meets the [auto mode requirements](#eliminate-prompts-with-auto-mode)
>     * **Bypass permissions**: requires the **Allow bypass permissions mode** toggle in Desktop settings on Pro and Max plans; on Team and Enterprise plans, organization policy controls it instead
>
>     The Cowork tab doesn't use these modes. Cowork has its own permission modes, enabled separately, and the Cowork tab shows no mode selector at all until a mode beyond its default is enabled for your account. See the [Cowork docs](https://claude.com/docs/cowork/overview).
>
>     For desktop-specific details, see [Choose a permission mode](/docs/en/desktop#choose-a-permission-mode) in the Desktop guide.
>
>     **As a default**: set `defaultMode` in [settings](/docs/en/settings#where-settings-live). The desktop app reads the same settings files as the CLI and applies the permission mode to new local sessions.
>
>     A mode you pick in the mode selector is remembered per folder and takes precedence over `defaultMode` for that folder. Plan is the exception: picking it applies to the current session only.
>
>     For where `defaultMode` goes in a settings file, see the example under [Start in a different permission mode](#start-in-a-different-mode).
>
>   #### Web and mobile
>
> Use the mode dropdown next to the prompt box on [claude.ai/code](https://claude.ai/code) or in the mobile app. Permission prompts appear in claude.ai for approval. Which modes appear depends on where the session runs:
>
>     * **Cloud sessions** on [Claude Code on the web](/docs/en/claude-code-on-the-web): Accept edits, Plan, and Auto. Accept edits corresponds to `default` mode: cloud sessions pre-approve file edits regardless of mode, so the dropdown shows Accept edits instead of Manual. Cloud sessions still honor `defaultMode: "acceptEdits"` from settings. Auto mode appears only when your organization allows it and the selected model supports it. Bypass permissions isn't available.
>     * **[Remote Control](/docs/en/remote-control) sessions** on your local machine: Manual, Accept edits, and Plan. You can't select Auto or Bypass permissions from the app.
>       * Except for Bypass permissions, the dropdown shows the permission mode the local session is in, including one set from the terminal. It updates when the permission mode changes in the app or in the terminal. The session never reports Bypass permissions to claude.ai, so switching into it from the terminal doesn't change what the dropdown shows.
>       * Sessions hosted by the [desktop app](/docs/en/desktop) or the [VS Code extension](/docs/en/vs-code) report permission mode changes to claude.ai as they happen, the same as sessions hosted in a terminal.
>       * Before v2.1.202, sessions connected with `/remote-control` or `claude --remote-control` didn't report their permission mode at all, so claude.ai and the mobile app could show a permission mode the session wasn't in. The mismatch affected only the label. Claude Code generated permission prompts from the session's actual permission mode, and they still appeared in the app for approval.
>
>     For Remote Control, the local machine running the session must be signed in with your claude.ai account; API keys aren't supported. You can also set the starting permission mode when launching that local session:
>
>     ```bash
>     claude remote-control --permission-mode acceptEdits
>     ```
>
> ## Auto-approve file edits with acceptEdits mode
>
> `acceptEdits` mode lets Claude create and edit files in your working directory without prompting. The status bar shows `⏵⏵ accept edits on` while this mode is active.
>
> In addition to file edits, `acceptEdits` mode auto-approves common filesystem Bash commands: `mkdir`, `touch`, `rm`, `rmdir`, `mv`, `cp`, and `sed`. These commands are also auto-approved when prefixed with safe environment variables such as `LANG=C` or `NO_COLOR=1`, or process wrappers such as `timeout`, `nice`, or `nohup`. Like file edits, auto-approval applies only to paths inside your working directory or `additionalDirectories`. Paths outside that scope, writes to [protected paths](#protected-paths), `rm` and `rmdir` removals targeting a [critical path](#critical-paths), and all other Bash commands except the [built-in read-only set](/docs/en/permissions#read-only-commands) still prompt.
>
> When the [PowerShell tool](/docs/en/tools-reference#powershell-tool) is enabled, `acceptEdits` mode also auto-approves `Set-Content`, `Add-Content`, `Clear-Content`, and `Remove-Item` on in-scope paths, along with their common aliases. The same scope and protected-path rules apply, and `Remove-Item` gets [its own check](#remove-item-in-powershell). A positional argument that contains a quote character, such as the apostrophe in `Set-Content .\notes.txt "It's done"`, still prompts even on in-scope paths, because Claude Code can't statically validate an argument whose quoted and unquoted readings differ. Pass the content through a named parameter such as `-Value` to avoid the prompt.
>
> Use `acceptEdits` when you want to review changes in your editor or via `git diff` after the fact rather than approving each edit inline.
>
> Press `Shift+Tab` once from Manual mode to enter it, or start with it directly:
>
> ```bash
> claude --permission-mode acceptEdits
> ```
>
> ## Analyze before you edit with plan mode
>
> Plan mode tells Claude to research and propose changes without making them. Claude reads files, runs shell commands to explore, and writes a plan, but does not edit your source. Except in sessions with [bypass permissions available](#skip-all-checks-with-bypasspermissions-mode), edits stay blocked until you approve the plan.
>
> When [auto mode](/docs/en/auto-mode-config) is available and the `useAutoModeDuringPlan` setting is on, which it is by default, the classifier reviews shell commands during planning instead of prompting you. Approved commands run, and rejected ones are blocked. Otherwise, commands outside the [built-in read-only set](/docs/en/permissions#read-only-commands) prompt for approval, including when the sandbox's [auto-allow mode](/docs/en/sandboxing#sandbox-modes) is enabled. In sessions with bypass permissions available, neither the classifier nor a prompt applies to planning commands; [Skip all checks with bypassPermissions mode](#skip-all-checks-with-bypasspermissions-mode) covers the few things that still prompt there. In v2.1.212 through v2.1.217, sessions without bypass permissions prompted for every command outside the read-only set, whether or not auto mode was available.
>
> Enter plan mode by pressing `Shift+Tab` or prefixing a single prompt with `/plan`. You can also start in plan mode from the CLI:
>
> ```bash
> claude --permission-mode plan
> ```
>
> Press `Shift+Tab` again to leave plan mode without approving a plan.
>
> ### Review and approve a plan
>
> When the plan is ready, Claude presents it and asks how to proceed. From that prompt you can choose:
>
> * **Yes, and use auto mode**: approve and start in [auto mode](#eliminate-prompts-with-auto-mode). When auto mode is unavailable, this option reads **Yes, auto-accept edits**. If you started the session with bypass permissions enabled, the option reads **Yes, and switch to BYPASS PERMISSIONS (no further prompts) for this session** instead.
> * **Yes, manually approve edits**: approve and review each edit individually.
> * **No, keep planning**: stay in plan mode and tell Claude what to change.
>
> Approving a plan exits plan mode and switches the session to the permission mode each approve option describes, so Claude starts editing. To plan again, cycle back to plan mode with `Shift+Tab`, or prefix your next prompt with `/plan`.
>
> Press `Ctrl+G` to open the proposed plan in your default text editor and edit it directly before Claude proceeds. When [`showClearContextOnPlanAccept`](/docs/en/settings-reference#showclearcontextonplanaccept) is enabled, the list gains a first option that approves the plan and clears the planning context.
>
> Accepting a plan also gives the session a [generated title](/docs/en/sessions#name-your-sessions) based on the plan, unless you've already named the session.
>
> ### Set plan mode as the default
>
> To make plan mode the default for a project's terminal sessions, set `defaultMode` to `plan` in `.claude/settings.json`, placed as the example under [Start in a different permission mode](#start-in-a-different-mode) shows. Conversations the [VS Code extension](/docs/en/vs-code) starts don't read project settings for the starting permission mode. There, set `claudeCode.initialPermissionMode` to `plan` in your VS Code user settings instead.
>
>   Eliminate permission prompts with auto mode
>
> Auto mode lets Claude execute without routine permission prompts. A separate classifier model reviews actions before they run, blocking anything that escalates beyond your request, targets unrecognized infrastructure, or appears driven by hostile content Claude read. Explicit [ask rules](/docs/en/permissions#manage-permissions) still force a prompt.
>
> On Pro, Max, and Team plans, auto mode is the [built-in starting permission mode](#which-mode-a-session-starts-in).
>
> The classifier also reviews each message Claude sends to another agent with [`SendMessage`](/docs/en/tools-reference), whether plain text or a structured [agent team](/docs/en/agent-teams) message, before Claude Code delivers it, both in auto mode and in [plan mode while the classifier reviews commands](#analyze-before-you-edit-with-plan-mode); the send review requires Claude Code v2.1.222 or later.
>
> The classifier also reviews and approves or blocks `rm` and `rmdir` removals targeting a [critical path](#critical-paths), such as `rm -rf /` and `rm -rf ~`, including when the removal sits inside command or process substitution.
>
> Auto mode also nudges Claude to keep working without stopping for clarifying questions, though Claude still asks when your prompt or a skill explicitly relies on it. For stronger autonomous behavior in a mode that still prompts you, set the [Proactive output style](/docs/en/output-styles) instead.
>
>   Auto mode reduces permission prompts but does not guarantee safety. Use it for tasks where you trust the general direction, not as a replacement for review on sensitive operations.
>
> Auto mode is available only when your account meets all of these requirements:
>
> * **Plan**: All plans.
> * **Organization**: on Team and Enterprise, auto mode is available by default. Administrators can turn it off for the organization by setting `permissions.disableAutoMode` to `"disable"` in [managed settings](/docs/en/managed-settings).
> * **Model**: on the Anthropic API and [Claude Platform on AWS](/docs/en/claude-platform-on-aws), Claude Opus 4.6 or later, Sonnet 4.6 or later, or [Fable 5](/docs/en/model-config#work-with-fable-5). On Amazon Bedrock, Google Cloud's Agent Platform, Microsoft Foundry, and signed-in [Claude apps gateway](/docs/en/claude-apps-gateway) sessions, only Claude Sonnet 5, Opus 4.7 or later, and Fable 5. Older models, including Sonnet 4.5, Opus 4.5, Haiku, and claude-3 models, are not supported on any provider.
> * **Provider**: available by default on the Anthropic API, Claude Platform on AWS, Amazon Bedrock, Google Cloud's Agent Platform, Microsoft Foundry, and signed-in Claude apps gateway sessions.
>
> If Claude Code reports auto mode as unavailable, one of these requirements is unmet; this is not a transient outage. A separate message that names a model and says auto mode "cannot determine the safety" of an action means a classifier request failed; that failure is usually transient, but on Amazon Bedrock it can repeat until your account can invoke the named model. See the [error reference](/docs/en/errors#auto-mode-cannot-determine-the-safety-of-an-action) for the causes and what to do.
>
> If you set `defaultMode: "auto"` in [settings](/docs/en/settings-reference#all-settings) and a terminal session starts in Manual mode with no error, the setting is likely in `.claude/settings.json` or `.claude/settings.local.json`. `auto` doesn't take effect from those files. Move it to `~/.claude/settings.json`. For a conversation the VS Code extension started, check the extension's own list in [Switch permission modes](#switch-permission-modes) instead.
>
>   Auto mode on Bedrock, Agent Platform, or Foundry
>
> On [Amazon Bedrock](/docs/en/amazon-bedrock), [Google Cloud's Agent Platform](/docs/en/google-vertex-ai), [Microsoft Foundry](/docs/en/microsoft-foundry), and signed-in [Claude apps gateway](/docs/en/claude-apps-gateway) sessions, auto mode appears in the `Shift+Tab` cycle by default. Appearing in the cycle doesn't change the permission mode a session starts in: on these providers, terminal sessions start in your [`defaultMode`](/docs/en/settings-reference#permissions-defaultmode), which is Manual unless you change it, and conversations in the [VS Code extension](/docs/en/vs-code) start in Manual unless `claudeCode.initialPermissionMode` or a mode you picked in the extension sets one. Only Claude Sonnet 5, Opus 4.7 or later, and Fable 5 are supported on these providers.
>
> To make auto mode the default starting permission mode, set `"permissions": {"defaultMode": "auto"}` in user or managed settings. In sessions the VS Code extension starts, select **Auto** from the mode indicator instead. [Switch permission modes](#switch-permission-modes) covers what outranks that pick.
>
> The [`/doctor`](/docs/en/commands#all-commands) checkup proposes this user-settings default on these providers the same way it does on the Anthropic API.
>
> To prevent developers from using auto mode, set `disableAutoMode` to `"disable"` in [managed settings](/docs/en/managed-settings). This removes `auto` from the `Shift+Tab` cycle, and a session started with `--permission-mode auto` starts in Manual instead.
>
> In v2.1.158 through v2.1.206, auto mode was off on these providers until you set `CLAUDE_CODE_ENABLE_AUTO_MODE=1`, and Claude Code ignored `defaultMode: "auto"` on these providers unless the variable was also set. The variable is still accepted for compatibility and has no effect from v2.1.207 onward.
>
> ### What the classifier blocks by default
>
> The classifier trusts your working directory and the remotes that were configured for it when the session started. A remote added or repointed during the session with `git remote add` or `git remote set-url` isn't trusted, and everything else is treated as external until you [configure trusted infrastructure](/docs/en/auto-mode-config). Before v2.1.200, remotes added mid-session were also trusted.
>
> **Blocked by default**:
>
> * Downloading and executing code, like `curl | bash`
> * Sending sensitive data to external endpoints
> * Production deploys and migrations
> * Mass deletion on cloud storage
> * Granting IAM or repo permissions
> * Modifying shared infrastructure
> * Irreversibly destroying files that existed before the session
> * Force push
> * Committing or pushing a change that would send secrets or sensitive data outside the repository when it runs, or widen what a deploy exposes. This covers a CI workflow or deploy configuration that hands a secret to a destination that doesn't already receive it, a script or setup step that reads a secret store and sends the data out, and a config change that widens what a deploy publishes, such as a registry, visibility, artifact, or sourcemap setting. The check applies on any branch, applies even when the r