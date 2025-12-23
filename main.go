package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

// --- Configuration ---
type Config struct {
	OllamaURL    string `json:"ollama_url"`
	Model        string `json:"model"`
	SystemPrompt string `json:"system_prompt"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type ChatResponse struct {
	Message Message `json:"message"`
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
	// 1. Load or Create Config
	config := loadConfig()

	// 2. Get Goal
	if len(os.Args) < 2 {
		fmt.Println("Usage: aish <your goal>")
		os.Exit(1)
	}
	initialGoal := strings.Join(os.Args[1:], " ")

	messages := []Message{
		{Role: "system", Content: config.SystemPrompt},
		{Role: "user", Content: initialGoal},
	}

	client := &http.Client{}

	// 3. The Loop
	for {
		fmt.Print("🧠 Thinking...")
		cmdStr, err := queryOllama(client, config, messages)
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
				messages = append(messages, Message{Role: "assistant", Content: cmdStr})
				messages = append(messages, Message{Role: "user", Content: refinement})
				break innerLoop // Break inner loop to regenerate command

			case '3': // Explain
				fmt.Print("🧠 Thinking...")
				explanation, err := explainCommand(client, config, cmdStr)
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
		OllamaURL:    "http://localhost:11434",
		Model:        "llama3.2:3b",
		SystemPrompt: advancedSystemPrompt,
	}

	// 2. Resolve Path
	home, err := os.UserHomeDir()
	if err != nil {
		return cfg // Fallback to defaults if home dir fails
	}
	dirPath := filepath.Join(home, ".config", "aish")
	filePath := filepath.Join(dirPath, "config.json")

	// 3. Create Config if it doesn't exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create directory if missing
		os.MkdirAll(dirPath, 0755)

		// Write defaults to file
		data, _ := json.MarshalIndent(cfg, "", "  ")
		os.WriteFile(filePath, data, 0644)

		fmt.Printf("⚙️  Created new config at %s\n", filePath)
	} else {
		// 4. Load existing config
		data, err := os.ReadFile(filePath)
		if err == nil {
			json.Unmarshal(data, &cfg)
		}
	}

	return cfg
}

func queryOllama(client *http.Client, cfg Config, msgs []Message) (string, error) {
	reqBody := ChatRequest{
		Model:    cfg.Model,
		Messages: msgs,
		Stream:   false,
	}

	jsonData, _ := json.Marshal(reqBody)
	resp, err := client.Post(cfg.OllamaURL+"/api/chat", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API Status: %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var chatResp ChatResponse
	json.Unmarshal(body, &chatResp)

	return strings.TrimSpace(chatResp.Message.Content), nil
}

func explainCommand(client *http.Client, cfg Config, cmd string) (string, error) {
	// Create a separate, temporary chat session for explanation
	explainPrompt := fmt.Sprintf("Explain this shell command in detail, breaking down each part, flag, and parameter:\n\n%s\n\nProvide a clear, educational explanation in plain text. Do NOT use markdown formatting, code blocks, or special symbols. Just plain text.", cmd)

	tempMessages := []Message{
		{Role: "user", Content: explainPrompt},
	}

	return queryOllama(client, cfg, tempMessages)
}
