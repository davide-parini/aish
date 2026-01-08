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

### Via Homebrew (Recommended)

```bash
brew tap davide-parini/aish
brew install aish
```

### From Source

```bash
git clone https://github.com/davide-parini/aish
cd aish
./install.sh
```

The install script will:
- Build the binary
- Install to `~/.local/bin` (no sudo required)
- Remind you to add `~/.local/bin` to PATH if needed

On first run, aish will automatically create `~/.config/aish/config.yaml` with default settings configured to use local Ollama (the default provider).

**Installation Details**:
- **Homebrew path**: `/opt/homebrew/bin/aish` (Apple Silicon) or `/usr/local/bin/aish` (Intel)
- **Manual install path**: `~/.local/bin/aish`
- **Config path**: `~/.config/aish/config.yaml` (created on first run)
- **Default provider**: Ollama (local, no API key required)

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

## Supported Platforms

aish is **macOS-only** and specifically optimized for Zsh on macOS.

## Command-Line Flags

- `-p`, `--provider <name>` - Override provider for single command (ollama or gemini)
- `--set-default-provider <name>` - Update default provider in config

## How It Works

1. **Input**: You describe what you want in natural language
2. **Generation**: aish sends your request to the configured LLM provider (Ollama or Gemini) with a specialized macOS/Zsh-optimized prompt
3. **Review**: The generated shell command is displayed for your approval before execution
4. **Action**: You can execute, copy, refine, or explain the command

The system prompt is engineered to output raw, executable Zsh commands without markdown formatting or explanations, ensuring compatibility with direct shell execution on macOS.

## Requirements

### For Using (Homebrew Installation)

- **macOS**: 10.13+ (any recent macOS version)
- **RAM**: ~8GB for llama3.2:3b model (if using Ollama)
- **Disk**: ~2GB for model storage (if using Ollama)

### For Building from Source

- **Go**: 1.25+ 
- **macOS** running Zsh
- **RAM**: ~8GB for llama3.2:3b model (if using Ollama)
- **Disk**: ~2GB for model storage (if using Ollama)

## How Homebrew Distribution Works

This project uses [GoReleaser](https://goreleaser.com) to automate the entire release and Homebrew distribution process.

### The Release Process

1. **Create and Push a Version Tag**:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. **GitHub Actions Trigger**: The tag push automatically triggers the [release workflow](.github/workflows/release.yml)

3. **Cross-Compilation**: GoReleaser compiles `aish` for both macOS architectures:
   - macOS (Intel): `darwin/amd64`
   - macOS (Apple Silicon): `darwin/arm64`

4. **Binary Packaging**: Each binary is packaged into a `.tar.gz` archive with a clear identifier like `aish_1.0.0_Darwin_arm64.tar.gz`

5. **GitHub Release**: GoReleaser creates a GitHub Release with:
   - 2 compiled binaries (Intel + Apple Silicon)
   - `checksums.txt` file for integrity verification
   - Auto-generated changelog from git commits

6. **Homebrew Formula Update**: GoReleaser automatically:
   - Creates or updates `Formula/aish.rb` in your `homebrew-aish` repository
   - Calculates SHA256 checksums for verification
   - Updates version numbers and download URLs
   - Commits and pushes changes

### For End Users

Once published, users install with:

```bash
brew tap davide-parini/aish
brew install aish
brew upgrade aish  # Update to latest version
```

Homebrew downloads the pre-compiled binary (no Go required), verifies checksums, and installs to `/opt/homebrew/bin` (Intel Macs) or `/usr/local/bin`.

## License

MIT

## Contributing

Contributions welcome! Feel free to open issues or submit pull requests.
