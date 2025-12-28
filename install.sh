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
echo "✅ Installation complete!"
echo ""
echo "📝 Configuration:"
echo "   - Config will be created at ~/.config/aish/config.yaml on first run"
echo "   - Default provider: ollama (local)"
echo ""
echo "🔧 To use Gemini instead:"
echo "   1. Run aish once to generate config"
echo "   2. Edit ~/.config/aish/config.yaml and add your API key under 'gemini.api_key'"
echo "   3. Run: aish --set-default-provider gemini"
echo ""
echo "Note: Make sure ~/.local/bin is in your PATH."
echo "Add this to your ~/.zshrc if needed:"
echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
