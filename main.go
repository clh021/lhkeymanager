package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/clh021/lhkeymanager/core"

	"golang.org/x/term"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// 默认环境文件路径
	envFilePath := ".env"

	// 检查命令行参数
	var choice string
	if len(os.Args) > 1 {
		choice = os.Args[1]
	} else {
		fmt.Println("错误: 请提供一个命令 (store, load, export)")
		fmt.Println("用法: ./lhkeymanager <command> [file_path]")
		os.Exit(1)
	}

	// 检查可选的文件路径参数
	if len(os.Args) > 2 {
		envFilePath = os.Args[2]
	}

	var key string
	// 清理内存中的敏感数据
	defer clearString(&key)

	maxAttempts := 3
	for i := 0; i < maxAttempts; i++ {
		// 获取加密密钥（不显示输入）
		fmt.Fprint(os.Stderr, "请输入加密密钥: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Fprintf(os.Stderr, "\n读取密钥失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintln(os.Stderr) // 添加换行到 stderr
		key = string(bytePassword)

		// Validate the encryption key
		if core.ValidateKey(key) {
			break // Key is valid, exit loop
		}

		// Invalidate the key in memory after a failed attempt
		clearString(&key)

		if i < maxAttempts-1 {
			fmt.Fprintf(os.Stderr, "错误: 密钥验证失败。您还有 %d 次机会。\n", maxAttempts-1-i)
		} else {
			fmt.Fprintln(os.Stderr, "错误: 密钥验证失败。已达到最大尝试次数。")
			if core.KeyHint != "" && core.KeyHint != "No hint available." {
				fmt.Fprintf(os.Stderr, "密钥提示: %s\n", core.KeyHint)
			}
			os.Exit(1)
		}
	}

	switch choice {
	case "store":
		storeKey(reader, key, envFilePath)
	case "load":
		loadKeysToNewBash(key, envFilePath)
	case "export":
		exportKeys(key, envFilePath)
	default:
		fmt.Printf("错误: 未知命令 '%s'. 可用命令: store, load, export\n", choice)
		fmt.Println("用法: ./lhkeymanager <command> [file_path]")
		os.Exit(1)
	}
}

// Store a new API key in the .env file
func storeKey(reader *bufio.Reader, key string, envFilePath string) {
	for i := 0; ; i++ {
		// Add a newline before each iteration except the first one
		if i > 0 {
			fmt.Println()
		}
		// Get the API key to encrypt (input not shown)
		fmt.Print("请输入要加密的API密钥: ")
		byteSecret, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Printf("\n读取API密钥失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Println() // 添加换行
		plaintext := string(byteSecret)

		// Get the environment variable name (can be displayed)
		fmt.Print("请输入环境变量名(带后缀): ")
		envName, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("读取环境变量名失败: %v\n", err)
			os.Exit(1)
		}
		envName = strings.TrimSpace(envName)

		// Store the API key
		encValue, err := core.StoreAPIKey(plaintext, envName, key, envFilePath)
		if err != nil {
			fmt.Printf("存储API密钥失败: %v\n", err)
			os.Exit(1)
		}

		// Output the encryption result
		fmt.Printf("加密结果: %s\n", encValue)
		fmt.Printf("已成功保存到 %s 文件\n", envFilePath)

		// Clear sensitive data from memory
		clearString(&plaintext)

		// Ask if the user wants to add another key
		fmt.Print("是否继续添加新的密钥? (y/n): ")
		continueChoice, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("读取输入失败: %v\n", err)
			os.Exit(1)
		}
		continueChoice = strings.TrimSpace(continueChoice)

		if continueChoice != "y" && continueChoice != "Y" {
			break
		}
	}
}

// Load keys from the .env file into a new bash session
func loadKeysToNewBash(key string, envFilePath string) {
	// Check if .env file exists
	if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
		fmt.Printf("错误: 文件 %s 不存在\n", envFilePath)
		os.Exit(1)
	}

	// Set file permissions
	err := os.Chmod(envFilePath, 0600)
	if err != nil {
		fmt.Printf("设置 %s 文件权限失败: %v\n", envFilePath, err)
		// Continue execution, don't exit
	}

	// Load and decrypt API keys
	decryptedVars, err := core.LoadAPIKeys(key, envFilePath)
	if err != nil {
		fmt.Printf("从 %s 加载密钥失败: %v\n", envFilePath, err)
		os.Exit(1)
	}

	// Create temporary environment variables file
	tempEnv, err := os.CreateTemp("", "env_vars_*")
	if err != nil {
		fmt.Printf("创建临时文件失败: %v\n", err)
		os.Exit(1)
	}
	tempEnvPath := tempEnv.Name()
	defer os.Remove(tempEnvPath) // Ensure temporary file is deleted on exit

	// Set temporary file permissions
	err = os.Chmod(tempEnvPath, 0600)
	if err != nil {
		fmt.Printf("设置临时文件权限失败: %v\n", err)
		// Continue execution, don't exit
	}

	// Write basic shell configuration
	tempEnv.WriteString("#!/bin/bash\n")
	tempEnv.WriteString("# This is an automatically generated temporary environment variables file\n\n")

	// Write environment variables to temporary file
	for name, value := range decryptedVars {
		tempEnv.WriteString(fmt.Sprintf("export %s='%s'\n", name, value))
		fmt.Printf("已设置环境变量: %s\n", name)
	}
	fmt.Println() // 添加换行

	// Close temporary file
	tempEnv.Close()

	// Start new bash session
	fmt.Println("正在启动新的bash会话，环境变量已设置...")

	// Use source command to load environment variables and start new bash
	cmd := exec.Command("bash", "-c", fmt.Sprintf("source %s && bash", tempEnvPath))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("启动bash会话失败: %v\n", err)
		os.Exit(1)
	}

	// Securely delete temporary file
	secureDeleteFile(tempEnvPath)

	fmt.Println("\nbash会话已结束，环境变量已清除")
}

// exportKeys loads keys from the .env file and prints them as export commands
func exportKeys(key string, envFilePath string) {
	// Check if .env file exists
	if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
		// Print to stderr so it doesn't get captured by eval
		fmt.Fprintf(os.Stderr, "错误: 文件 %s 不存在\n", envFilePath)
		os.Exit(1)
	}

	// Load and decrypt API keys
	decryptedVars, err := core.LoadAPIKeys(key, envFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "从 %s 加载密钥失败: %v\n", envFilePath, err)
		os.Exit(1)
	}

	// Write environment variables to standard output
	for name, value := range decryptedVars {
		// Use single quotes to handle special characters in values
		fmt.Printf("export %s='%s';\n", name, value)
	}
}

// secureDeleteFile attempts to securely delete a file
func secureDeleteFile(path string) {
	// Try to use the shred command for secure deletion
	shredCmd := exec.Command("shred", "-u", "-z", path)
	err := shredCmd.Run()

	// If shred command is not available, use regular deletion
	if err != nil {
		os.Remove(path)
	}
}

// clearString clears a string from memory by overwriting it
func clearString(s *string) {
	if s == nil {
		return
	}
	for i := range *s {
		(*s) = (*s)[:i] + "\x00" + (*s)[i+1:]
	}
	*s = ""
}
