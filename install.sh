#!/bin/bash

echo "🔨 Building aish..."
go build -o aish

if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi

echo "📦 Installing to ~/.local/bin..."
mkdir -p ~/.local/bin
mv aish ~/.local/bin/

echo ""
echo "⚙️  Configuration Setup"
echo "======================"
echo ""

# Ask for default provider
echo "Select default LLM provider:"
echo "  1) ollama (local, default)"
echo "  2) gemini (Google AI)"
read -p "Choice [1]: " provider_choice

if [ "$provider_choice" = "2" ]; then
    default_provider="gemini"
    echo ""
    read -p "Enter your Gemini API key (or press Enter to configure later): " gemini_api_key
else
    default_provider="ollama"
    gemini_api_key=""
    echo ""
    echo "ℹ️  Defaulting to local Ollama. Make sure Ollama is installed and running."
fi

# Create config directory
mkdir -p ~/.config/aish

# Generate YAML config
cat > ~/.config/aish/config.yaml << EOF
default_provider: $default_provider
system_prompt: |
  You are a highly skilled macOS Zsh Command Generator.
  Your specific goal is to output raw, executable Zsh commands.

  RULES:
  1. OUTPUT FORMAT: Return ONLY the command text as a SINGLE LINE. Do NOT use markdown code blocks. Do NOT include quotes around the command. Do NOT add explanations. Do NOT output multiple lines or line breaks. Do NOT escape special characters unnecessarily - output clean, executable commands.
  2. ENVIRONMENT: The user is on macOS running Zsh. Use macOS-specific tools (pbcopy, open, etc.) and Zsh syntax.
  3. SAFETY: Prefer non-destructive flags where applicable.
  4. COMPLEXITY: If multiple steps are required, combine them with '&&' or ';' on a SINGLE LINE.
  5. CLARITY: If the user request is ambiguous, generate the most likely useful command.
  6. REFINEMENT: If the user provides additional constraints or modifications, BUILD ON the previous command by combining the new requirements with existing ones. Do NOT discard previous conditions.
  7. CORRECTNESS: Always verify command syntax is valid and will execute without errors. Test your commands mentally before outputting. Avoid unnecessary escaping like \( or \) unless truly required.

  Example User: "find huge files"
  Example Output: find . -type f -size +100M

  Example User: "update system"
  Example Output: softwareupdate -i -a

  Example User: "multiple steps"
  Example Output: cd /tmp && mkdir test && touch test/file.txt

  Example Refinement Flow:
  User: "find files edited in march"
  Assistant: find . -type f -newermt 2025-03-01 ! -newermt 2025-04-01
  User: "only big files"
  Assistant: find . -type f -newermt 2025-03-01 ! -newermt 2025-04-01 -size +100M

  NOW, generate the command for the following request.

ollama:
  url: http://localhost:11434
  model: llama3.2:3b

gemini:
  api_key: $gemini_api_key
  model: gemini-flash-lite-latest
EOF

echo ""
echo "✅ Installation complete!"
echo ""
echo "Configuration saved to ~/.config/aish/config.yaml"
echo "Default provider: $default_provider"

if [ "$default_provider" = "gemini" ] && [ -z "$gemini_api_key" ]; then
    echo ""
    echo "⚠️  Warning: Gemini API key not provided."
    echo "   Please edit ~/.config/aish/config.yaml and add your API key under 'gemini.api_key'"
fi

echo ""
echo "Note: Make sure ~/.local/bin is in your PATH."
echo "Add this to your ~/.zshrc if needed:"
echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
