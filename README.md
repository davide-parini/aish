# aish - AI Shell Helper

A fast, terminal-based tool that converts natural language descriptions into executable shell commands using local LLMs via Ollama.

## Purpose

aish eliminates the need to remember complex shell syntax by letting you describe what you want to do in plain English. It generates the command, lets you review it, and gives you the option to execute, copy, or refine it.

**Key Features:**
- ⚡ **Instant startup** - Built in Go for zero-latency execution
- 🔒 **Fully local** - Runs entirely on your machine via Ollama
- 🎯 **macOS optimized** - Designed for Zsh on macOS
- 🔄 **Iterative refinement** - Chat with the AI to adjust commands
- 📋 **Clipboard integration** - Copy commands with one keystroke
- 🚀 **Single binary** - No dependencies, just run it

## Prerequisites

- macOS running Zsh
- [Ollama](https://ollama.ai) installed and running
- The `llama3.2:3b` model (or customize in config)

```bash
# Install Ollama (if not already installed)
curl -fsSL https://ollama.ai/install.sh | sh

# Pull the default model
ollama pull llama3.2:3b
```

## Installation

```bash
git clone https://github.com/davide-parini/aish
cd aish
./install.sh
```

The install script will:
- Build the binary
- Install to `~/.local/bin` (no sudo required)
- Delete any existing config for a fresh start
- Remind you to add `~/.local/bin` to PATH if needed

## Usage

### Basic Usage

```bash
aish <your goal in natural language>
```

**Examples:**

```bash
aish find files larger than 100MB
aish update homebrew and all packages
aish list all running docker containers
aish compress all PDFs in this folder
aish show my IP address
```

### Interactive Mode

After generating a command, simply press a key to choose your action:

- **[Enter]** - Execute the command immediately 🚀
- **[1]** - Copy to clipboard and exit 📋
- **[2]** - Refine the command with additional instructions ✨
- **[3]** - Explain the command in detail 💡

**Example Output:**

```
🧠 Thinking...find . -type f -size +100M
[Enter] Execute 🚀
[1] Copy 📋
[2] Refine ✨
[3] Explain 💡
```

**Refinement Example:**

After pressing **[2]**:

```
✨ Refinement prompt: only show files modified in the last week
```

The AI will combine your refinement with the previous command, maintaining context across iterations.

**Explanation Example:**

After pressing **[3]**:

```
💡 Explanation:
find . -type f -size +100M

This command searches for files larger than 100 megabytes:
- find: Unix utility to search for files
- .: Start from current directory
- -type f: Only match files (not directories)
- -size +100M: Files larger than 100 megabytes
```

The explanation is generated in a separate session and won't affect your refinement context.

## Configuration

On first run, aish creates `~/.config/aish/config.json` with default settings:

```json
{
  "ollama_url": "http://localhost:11434",
  "model": "llama3.2:3b",
  "system_prompt": "..."
}
```

**Customization:**
- Change `ollama_url` if running Ollama remotely
- Use a different `model` (e.g., `mistral`, `codellama`)
- Modify `system_prompt` for different behavior

## How It Works

1. **Input**: You describe what you want in natural language
2. **Generation**: aish sends your request to Ollama with a specialized prompt
3. **Review**: The generated command is displayed for your approval
4. **Action**: Execute, copy, or refine based on your needs

The system prompt is engineered to output raw, executable commands without markdown formatting or explanations, ensuring compatibility with direct execution.

## Requirements

- **Go**: 1.19+ (for building from source)
- **RAM**: ~8GB for llama3.2:3b model
- **Disk**: ~2GB for model storage

## License

MIT

## Contributing

Contributions welcome! Feel free to open issues or submit pull requests.
