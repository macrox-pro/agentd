---
primary_sources:
  - id: T1-HEADLESS
    title: "Headless mode reference"
    url: "https://geminicli.com/docs/cli/headless.md"
    section: "Full page"
  - id: T1-AUTO-TUT
    title: "Automate tasks with headless mode"
    url: "https://geminicli.com/docs/cli/tutorials/automation.md"
    section: "Full page"
also_cited_in: []
studied_at: "2026-08-29"
gemini_docs_snapshot: "2026-08-29"
applicability: current
---
# Headless mode

> **Applicability:** Verbatim excerpts from Gemini CLI documentation (snapshot 2026-08-29).

### Source: Headless mode reference — Full page

> Headless mode provides a programmatic interface to Gemini CLI, returning
> structured text or JSON output without an interactive terminal UI.
>
> ## Technical reference
>
> Headless mode is triggered when the CLI is run in a non-TTY environment or when
> providing a query with the `-p` (or `--prompt`) flag.
>
> ### Output formats
>
> You can specify the output format using the `--output-format` flag.
>
> #### JSON output
>
> Returns a single JSON object containing the response and usage statistics.
>
> - **Schema:**
>   - `response`: (string) The model's final answer.
>   - `stats`: (object) Token usage and API latency metrics.
>   - `error`: (object, optional) Error details if the request failed.
>
> #### Streaming JSON output
>
> Returns a stream of newline-delimited JSON (JSONL) events.
>
> - **Event types:**
>   - `init`: Session metadata (session ID, model).
>   - `message`: User and assistant message chunks.
>   - `tool_use`: Tool call requests with arguments.
>   - `tool_result`: Output from executed tools.
>   - `error`: Non-fatal warnings and system errors.
>   - `result`: Final outcome with aggregated statistics and per-model token usage
>     breakdowns.
>
> ## Exit codes
>
> The CLI returns standard exit codes to indicate the result of the headless
> execution:
>
> - `0`: Success.
> - `1`: General error or API failure.
> - `42`: Input error (invalid prompt or arguments).
> - `53`: Turn limit exceeded.
>
> ## Next steps
>
> - Follow the [Automation tutorial](/docs/cli/tutorials/automation) for practical
>   scripting examples.
> - See the [CLI reference](/docs/cli/cli-reference) for all available flags.

### Source: Automate tasks with headless mode — Full page

> Automate tasks with Gemini CLI. Learn how to use headless mode, pipe data into
> Gemini CLI, automate workflows with shell scripts, and generate structured JSON
> output for other applications.
>
> ## Prerequisites
>
> - Gemini CLI installed and authenticated.
> - Familiarity with shell scripting (Bash/Zsh).
>
> ## Why headless mode?
>
> Headless mode runs Gemini CLI once and exits. It's perfect for:
>
> - **CI/CD:** Analyzing pull requests automatically.
> - **Batch processing:** Summarizing a large number of log files.
> - **Tool building:** Creating your own "AI wrapper" scripts.
>
> ## How to use headless mode
>
> Run Gemini CLI in headless mode by providing a prompt with the `-p` (or
> `--prompt`) flag. This bypasses the interactive chat interface and prints the
> response to standard output (stdout). Positional arguments without the flag
> default to interactive mode, unless the input or output is piped or redirected.
>
> Run a single command:
>
> ```bash
> gemini -p "Write a poem about TypeScript"
> ```
>
> ## How to pipe input to Gemini CLI
>
> Feed data into Gemini using the standard Unix pipe `|`. Gemini reads the
> standard input (stdin) as context and answers your question using standard
> output.
>
> Pipe a file:
>
> **macOS/Linux**
>
> ```bash
> cat error.log | gemini -p "Explain why this failed"
> ```
>
> **Windows (PowerShell)**
>
> ```powershell
> Get-Content error.log | gemini -p "Explain why this failed"
> ```
>
> Pipe a command:
>
> ```bash
> git diff | gemini -p "Write a commit message for these changes"
> ```
>
> ## Use Gemini CLI output in scripts
>
> Because Gemini prints to stdout, you can chain it with other tools or save the
> results to a file.
>
> ### Scenario: Bulk documentation generator
>
> You have a folder of Python scripts and want to generate a `README.md` for each
> one.
>
> 1.  Save the following code as `generate_docs.sh` (or `generate_docs.ps1` for
>     Windows):
>
>     **macOS/Linux (`generate_docs.sh`)**
>
>     ```bash
>     #!/bin/bash
>
>     # Loop through all Python files
>     for file in *.py; do
>       echo "Generating docs for $file..."
>
>       # Ask Gemini CLI to generate the documentation and print it to stdout
>       gemini -p "Generate a Markdown documentation summary for @$file. Print the
>       result to standard output." > "${file%.py}.md"
>     done
>     ```
>
>     **Windows PowerShell (`generate_docs.ps1`)**
>
>     ```powershell
>     # Loop through all Python files
>     Get-ChildItem -Filter *.py | ForEach-Object {
>       Write-Host "Generating docs for $($_.Name)..."
>
>       $newName = $_.Name -replace '\.py$', '.md'
>       # Ask Gemini CLI to generate the documentation and print it to stdout
>       gemini -p "Generate a Markdown documentation summary for @$($_.Name). Print the result to standard output." | Out-File -FilePath $newName -Encoding utf8
>     }
>     ```
>
> 2.  Make the script executable and run it in your directory:
>
>     **macOS/Linux**
>
>     ```bash
>     chmod +x generate_docs.sh
>     ./generate_docs.sh
>     ```
>
>     **Windows (PowerShell)**
>
>     ```powershell
>     .\generate_docs.ps1
>     ```
>
>     This creates a corresponding Markdown file for every Python file in the
>     folder.
>
> ## Extract structured JSON data
>
> When writing a script, you often need structured data (JSON) to pass to tools
> like `jq`. To get pure JSON data from the model, combine the
> `--output-format json` flag with `jq` to parse the response field.
>
> ### Scenario: Extract and return structured data
>
> 1.  Save the following script as `generate_json.sh` (or `generate_json.ps1` for
>     Windows):
>
>     **macOS/Linux (`generate_json.sh`)**
>
>     ```bash
>     #!/bin/bash
>
>     # Ensure we are in a project root
>     if [ ! -f "package.json" ]; then
>       echo "Error: package.json not found."
>       exit 1
>     fi
>
>     # Extract data
>     gemini --output-format json "Return a raw JSON object with keys 'version' and 'deps' from @package.json" | jq -r '.response' > data.json
>     ```
>
>     **Windows PowerShell (`generate_json.ps1`)**
>
>     ```powershell
>     # Ensure we are in a project root
>     if (-not (Test-Path "package.json")) {
>       Write-Error "Error: package.json not found."
>       exit 1
>     }
>
>     # Extract data (requires jq installed, or you can use ConvertFrom-Json)
>     $output = gemini --output-format json "Return a raw JSON object with keys 'version' and 'deps' from @package.json" | ConvertFrom-Json
>     $output.response | Out-File -FilePath data.json -Encoding utf8
>     ```
>
> 2.  Run the script:
>
>     **macOS/Linux**
>
>     ```bash
>     chmod +x generate_json.sh
>     ./generate_json.sh
>     ```
>
>     **Windows (PowerShell)**
>
>     ```powershell
>     .\generate_json.ps1
>     ```
>
> 3.  Check `data.json`. The file should look like this:
>
>     ```json
>     {
>       "version": "1.0.0",
>       "deps": {
>         "react": "^18.2.0"
>       }
>     }
>     ```
>
> ## Build your own custom AI tools
>
> Use headless mode to perform custom, automated AI tasks.
>
> ### Scenario: Create a "Smart Commit" alias
>
> You can add a function to your shell configuration to create a `git commit`
> wrapper that writes the message for you.
>
> **macOS/Linux (Bash/Zsh)**
>
> 1.  Open your `.zshrc` file (or `.bashrc` if you use Bash) in your preferred
>     text editor.
>
>     ```bash
>     nano ~/.zshrc
>     ```
>
>     **Note**: If you use VS Code, you can run `code ~/.zshrc`.
>
> 2.  Scroll to the very bottom of the file and paste this code:
>
>     ```bash
>     function gcommit() {
>       # Get the diff of staged changes
>       diff=$(git diff --staged)
>
>       if [ -z "$diff" ]; then
>         echo "No staged changes to commit."
>         return 1
>       fi
>
>       # Ask Gemini to write the message
>       echo "Generating commit message..."
>       msg=$(echo "$diff" | gemini -p "Write a concise Conventional Commit message for this diff. Output ONLY the message.")
>
>       # Commit with the generated message
>       git commit -m "$msg"
>     }
>     ```
>
>     Save your file and exit.
>
> 3.  Run this command to make the function available immediately:
>
>     ```bash
>     source ~/.zshrc
>     ```
>
> **Windows (PowerShell)**
>
> 1.  Open your PowerShell profile in your preferred text editor.
>
>     ```powershell
>     notepad $PROFILE
>     ```
>
> 2.  Scroll to the very bottom of the file and paste this code:
>
>     ```powershell
>     function gcommit {
>       # Get the diff of staged changes
>       $diff = git diff --staged
>
>       if (-not $diff) {
>         Write-Host "No staged changes to commit."
>         return
>       }
>
>       # Ask Gemini to write the message
>       Write-Host "Generating commit message..."
>       $msg = $diff | gemini -p "Write a concise Conventional Commit message for this diff. Output ONLY the message."
>
>       # Commit with the generated message
>       git commit -m "$msg"
>     }
>     ```
>
>     Save your file and exit.
>
> 3.  Run this command to make the function available immediately:
>
>     ```powershell
>     . $PROFILE
>     ```
>
> 4.  Use your new command:
>
>     ```bash
>     gcommit
>     ```
>
>     Gemini CLI will analyze your staged changes and commit them with a generated
>     message.
>
> ## Next steps
>
> - Explore the [Headless mode reference](/docs/cli/headless) for full JSON
>   schema details.
> - Learn about [Shell commands](/docs/cli/tutorials/shell-commands) to let the agent run scripts
>   instead of just writing them.
