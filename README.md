# aish - AI Shell Helper

A fast, terminal-based tool that converts natural language descriptions into executable shell commands using LLMs.

## Purpose

aish eliminates the need to remember complex shell syntax by letting you describe what you want to do in plain English. It generates the command, lets you review it, and gives you the option to execute, copy, or refine it.

**Key Features:**
- ⚡ **Instant startup** - Built in Go for zero-latency execution
- 🔌 **Multiple providers** - Use local Ollama or Google Gemini
- 🎯 **macOS optimized** - Designed for Zsh on macOS
- 🔄 **Iterative refinement** - Chat with the AI to adjust commands
- 📋 **Clipboard integration** - Copy commands with one keystroke
- 🚀 **Single binary** - No dependencies, just run it

## Prerequisites

**For Ollama (local):**
- macOS running Zsh
- [Ollama](https://ollama.ai) installed and running
- The `llama3.2:3b` model

```bash
# Install Ollama
curl -fsSL https://ollama.ai/install.sh | sh

# Pull the default model
ollama pull llama3.2:3b
```

**For Gemini:**
- macOS running Zsh
- Google AI API key ([get one here](https://aistudio.google.com/app/apikey))

## Installation

```bash
git clone https://github.com/davide-parini/aish
cd aish
./install.sh
```

The install script will:
- Build the binary
- Install to `~/.local/bin` (no sudo required)
- Remind you to add `~/.local/bin` to PATH if needed

On first run, aish will automatically create `~/.config/aish/config.yaml` with default settings (using local Ollama).

### Setting up Gemini

To use Google Gemini instead of local Ollama:

1. Get an API key from [Google AI Studio](https://aistudio.google.com/app/apikey)
2. Run aish once to generate the config file
3. Edit `~/.config/aish/config.yaml` and add your API key:
   ```yaml
   gemini:
     api_key: your-api-key-here
     model: gemini-flash-lite-latest
   ```
4. Set Gemini as default: `aish --set-default-provider gemini`

## Usage

```bash
# Basic usage
aish <your goal in natural language>

# Use specific provider for one command
aish -p gemini find large files
aish --provider ollama compress videos

# Change default provider
aish --set-default-provider gemini
```

### Interactive Mode

After generating a command, simply press a key to choose your action:

- **[⏎]** - Execute the command immediately 🚀
- **[1]** - Copy to clipboard and exit 📋
- **[2]** - Refine the command with additional instructions ✨
- **[3]** - Explain the command in detail 💡
- **[q]** - Exit without action 👋

**Example Session:**

```
🧠 Thinking...find . -type f -size +100M
[⏎] Execute 🚀
[1] Copy 📋
[2] Refine ✨
[3] Explain 💡
[q] Exit 👋
```

**Refinement:**

Press **[2]** to refine the command:

```
✨ Refinement prompt: only show files modified in the last week
```

The AI combines your refinement with the previous command, maintaining context across iterations.

**Explanation:**

Press **[3]** for a detailed breakdown:

## Configuration

Configuration is stored at `~/.config/aish/config.yaml`:

```yaml
default_provider: ollama
system_prompt: |
  You are a highly skilled macOS Zsh Command Generator.
  Your specific goal is to output raw, executable Zsh commands.
  ...

explain_prompt: |
  You are a helpful assistant that explains shell commands clearly and accurately.
  Break down each part of the command, explaining flags, parameters, and their purpose.
  ...

ollama:
  url: http://localhost:11434
  model: llama3.2:3b

gemini:
  api_key: your-api-key-here
  model: gemini-flash-lite-latest
```

**Customization:**
- Change `ollama.url` if running Ollama remotely
- Use a different `model` (e.g., `mistral`, `codellama`)
- Modify `system_prompt` for different command generation behavior
- Modify `explain_prompt` for different explanation style

## Command-Line Flags

- `-p`, `--provider <name>` - Override provider for single command (ollama or gemini)
- `--set-default-provider <name>` - Update default provider in config

## How It Works

1. **Input**: You describe what you want in natural language
2. **Generation**: aish sends your request to Ollama with a specialized prompt
3. **Review**: The generated command is displayed for your approval
4. **Action**: Execute, copy, or refine based on your needs

The system prompt is engineered to output raw, executable commands without markdown formatting or explanations, ensuring compatibility with direct execution.

## Requirements

- **Go**: 1.25+ (for building from source)
- **RAM**: ~8GB for llama3.2:3b model (if using Ollama)
- **Disk**: ~2GB for model storage (if using Ollama)

## License

MIT

## Contributing

Contributions welcome! Feel free to open issues or submit pull requests.
