package main

import (
	"aish/providers"
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"gopkg.in/yaml.v3"
)

// --- Configuration ---
type Config struct {
	DefaultProvider string       `yaml:"default_provider"`
	SystemPrompt    string       `yaml:"system_prompt"`
	Ollama          OllamaConfig `yaml:"ollama"`
	Gemini          GeminiConfig `yaml:"gemini"`
}

type OllamaConfig struct {
	URL   string `yaml:"url"`
	Model string `yaml:"model"`
}

type GeminiConfig struct {
	APIKey string `yaml:"api_key"`
	Model  string `yaml:"model"`
}

// --- Prompt Engineering ---
const advancedSystemPrompt = `You are a highly skilled macOS Zsh Command Generator.
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

NOW, generate the command for the following request.`

func main() {
	// Parse command-line flags
	providerFlag := flag.String("p", "", "LLM provider to use (ollama or gemini)")
	flag.StringVar(providerFlag, "provider", "", "LLM provider to use (ollama or gemini)")
	setDefaultProvider := flag.String("set-default-provider", "", "Set the default provider in config (ollama or gemini)")
	flag.Parse()

	// Load config
	config := loadConfig()

	// Handle --set-default-provider
	if *setDefaultProvider != "" {
		if *setDefaultProvider != "ollama" && *setDefaultProvider != "gemini" {
			fmt.Printf("❌ Invalid provider: %s (must be 'ollama' or 'gemini')\n", *setDefaultProvider)
			os.Exit(1)
		}
		if err := setDefaultProviderInConfig(*setDefaultProvider); err != nil {
			fmt.Printf("❌ Failed to update config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Default provider set to: %s\n", *setDefaultProvider)
		os.Exit(0)
	}

	// Determine which provider to use
	activeProvider := config.DefaultProvider
	if *providerFlag != "" {
		activeProvider = *providerFlag
	}

	// Validate provider
	if activeProvider != "ollama" && activeProvider != "gemini" {
		fmt.Printf("❌ Invalid provider: %s (must be 'ollama' or 'gemini')\n", activeProvider)
		os.Exit(1)
	}

	// Create provider instance
	provider := createProvider(config, activeProvider)
	if provider == nil {
		os.Exit(1)
	}

	// Get goal from remaining args
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Usage: aish [flags] <your goal>")
		fmt.Println("\nFlags:")
		fmt.Println("  -p, --provider <name>           Use specific provider for this command (ollama or gemini)")
		fmt.Println("  --set-default-provider <name>   Set default provider in config")
		os.Exit(1)
	}
	initialGoal := strings.Join(args, " ")

	messages := []providers.Message{
		{Role: "user", Content: initialGoal},
	}

	// 3. The Loop
	for {
		fmt.Print("🧠 Thinking...")
		cmdStr, err := provider.SendMessage(messages, config.SystemPrompt)
		if err != nil {
			fmt.Printf("\r\033[K❌ Error: %v\n", err)
			os.Exit(1)
		}

		// Clear "Thinking..." and display command with color
		fmt.Print("\r\033[K")

		// Inner loop for handling choices on the same command
	innerLoop:
		for {
			fmt.Printf("\033[36m%s\033[0m\n", cmdStr) // Cyan color for command
			fmt.Println("[⏎] Execute 🚀")
			fmt.Println("[1] Copy 📋")
			fmt.Println("[2] Refine ✨")
			fmt.Println("[3] Explain 💡")
			fmt.Println("[q] Exit 👋")

			choice := readSingleKey()

			switch choice {
			case '1': // Copy (macOS specific)
				cmd := exec.Command("pbcopy")
				cmd.Stdin = strings.NewReader(cmdStr)
				if err := cmd.Run(); err != nil {
					fmt.Printf("❌ Clipboard error: %v\n", err)
				} else {
					fmt.Println("✅ Copied to clipboard")
				}
				os.Exit(0)

			case '2': // Refine
				fmt.Print("✨ Refinement prompt: ")
				refinement := readLine()
				messages = append(messages, providers.Message{Role: "assistant", Content: cmdStr})
				messages = append(messages, providers.Message{Role: "user", Content: refinement})
				break innerLoop // Break inner loop to regenerate command

			case '3': // Explain
				fmt.Print("🧠 Thinking...")
				explanation, err := explainCommand(provider, config, cmdStr)
				fmt.Print("\r\033[K") // Clear "Thinking..."
				if err != nil {
					fmt.Printf("❌ Error: %v\n", err)
				} else {
					fmt.Println("💡 Explanation:")
					fmt.Println(explanation)
				}
				fmt.Println() // Empty line before showing options again
				continue      // Continue inner loop to show same command

			case '\r', '\n': // Execute (Zsh specific)
				fmt.Println("🚀 Executing...")
				cmd := exec.Command("/bin/zsh", "-c", cmdStr)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
				os.Exit(0)

			case 'q', 'Q': // Exit
				fmt.Println("👋 Bye!")
				os.Exit(0)

			default:
				fmt.Println("🤷 Invalid choice.")
				continue // Continue inner loop to show same command
			}
		}
	}
}

func readSingleKey() rune {
	// Save current terminal state
	var oldState syscall.Termios
	syscall.Syscall6(syscall.SYS_IOCTL, uintptr(syscall.Stdin), syscall.TIOCGETA, uintptr(unsafe.Pointer(&oldState)), 0, 0, 0)

	// Set raw mode
	newState := oldState
	newState.Lflag &^= syscall.ECHO | syscall.ICANON
	syscall.Syscall6(syscall.SYS_IOCTL, uintptr(syscall.Stdin), syscall.TIOCSETA, uintptr(unsafe.Pointer(&newState)), 0, 0, 0)

	// Read single character
	var buf [1]byte
	os.Stdin.Read(buf[:])

	// Restore terminal state
	syscall.Syscall6(syscall.SYS_IOCTL, uintptr(syscall.Stdin), syscall.TIOCSETA, uintptr(unsafe.Pointer(&oldState)), 0, 0, 0)

	return rune(buf[0])
}

func readLine() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func loadConfig() Config {
	// 1. Define Defaults
	cfg := Config{
		DefaultProvider: "ollama",
		SystemPrompt:    advancedSystemPrompt,
		Ollama: OllamaConfig{
			URL:   "http://localhost:11434",
			Model: "llama3.2:3b",
		},
		Gemini: GeminiConfig{
			APIKey: "",
			Model:  "gemini-flash-lite-latest",
		},
	}

	// 2. Resolve Path
	home, err := os.UserHomeDir()
	if err != nil {
		return cfg // Fallback to defaults if home dir fails
	}
	dirPath := filepath.Join(home, ".config", "aish")
	filePath := filepath.Join(dirPath, "config.yaml")

	// 3. Create Config if it doesn't exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create directory if missing
		os.MkdirAll(dirPath, 0755)

		// Write defaults to file
		data, _ := yaml.Marshal(cfg)
		os.WriteFile(filePath, data, 0644)

		fmt.Printf("⚙️  Created new config at %s\n", filePath)
	} else {
		// 4. Load existing config
		data, err := os.ReadFile(filePath)
		if err == nil {
			yaml.Unmarshal(data, &cfg)
		}
	}

	return cfg
}

func createProvider(config Config, providerName string) providers.Provider {
	switch providerName {
	case "ollama":
		return providers.NewOllamaProvider(config.Ollama.URL, config.Ollama.Model)
	case "gemini":
		return providers.NewGeminiProvider(config.Gemini.APIKey, config.Gemini.Model)
	default:
		fmt.Printf("❌ Unknown provider: %s\n", providerName)
		return nil
	}
}

func setDefaultProviderInConfig(provider string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	filePath := filepath.Join(home, ".config", "aish", "config.yaml")

	// Read existing config
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}

	// Update default provider
	cfg.DefaultProvider = provider

	// Write back
	data, err = yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

func explainCommand(provider providers.Provider, cfg Config, cmd string) (string, error) {
	// Create a separate, temporary chat session for explanation
	explainPrompt := fmt.Sprintf("Explain this shell command in detail, breaking down each part, flag, and parameter:\n\n%s\n\nProvide a clear, educational explanation in plain text. Do NOT use markdown formatting, code blocks, or special symbols. Just plain text.", cmd)

	tempMessages := []providers.Message{
		{Role: "user", Content: explainPrompt},
	}

	return provider.SendMessage(tempMessages, cfg.SystemPrompt)
}
