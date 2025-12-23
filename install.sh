#!/bin/bash

echo "🔨 Building aish..."
go build -o aish main.go

if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi

echo "📦 Installing to ~/.local/bin..."
mkdir -p ~/.local/bin
mv aish ~/.local/bin/

echo "🗑️  Deleting existing config..."
rm -f ~/.config/aish/config.json

echo "✅ Installation complete! Run 'aish' to generate fresh config."
echo ""
echo "Note: Make sure ~/.local/bin is in your PATH."
echo "Add this to your ~/.zshrc if needed:"
echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
