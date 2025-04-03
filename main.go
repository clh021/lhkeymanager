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

	// 询问用户操作
	fmt.Println("请选择操作:")
	fmt.Println("1. 存储新的API密钥到.env文件")
	fmt.Println("2. 读取.env文件中的密钥到新bash会话")
	fmt.Print("请输入选项 (1/2): ")

	choice, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("读取输入失败: %v\n", err)
		os.Exit(1)
	}
	choice = strings.TrimSpace(choice)

	// 获取加密密钥（不显示输入）
	fmt.Print("请输入加密密钥: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Printf("\n读取密钥失败: %v\n", err)
		os.Exit(1)
	}
	key := string(bytePassword)

	// Validate the encryption key
	if !core.ValidateKey(key) {
		fmt.Println("\n错误: 密钥验证失败")
		fmt.Println("请确保密钥符合以下要求:")
		fmt.Printf("- 最小长度: %d\n", core.MinKeyLength)
		if core.KeyPrefix != "" {
			fmt.Printf("- 必须以 '%s' 开头\n", core.KeyPrefix)
		}
		if core.KeySuffix != "" {
			fmt.Printf("- 必须以 '%s' 结尾\n", core.KeySuffix)
		}
		if core.RequiredChars != "" {
			fmt.Printf("- 必须包含至少 %d 个特殊字符 (%s)\n", core.MinSpecialChars, core.RequiredChars)
		}
		os.Exit(1)
	}

	switch choice {
	case "1":
		storeKey(reader, key)
	case "2":
		loadKeysToNewBash(key)
	default:
		fmt.Println("无效的选项")
		os.Exit(1)
	}

	// 清理内存中的敏感数据
	clearString(&key)
}

// Store a new API key in the .env file
func storeKey(reader *bufio.Reader, key string) {
	// Get the API key to encrypt (input not shown)
	fmt.Print("请输入要加密的API密钥: ")
	byteSecret, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Printf("\n读取API密钥失败: %v\n", err)
		os.Exit(1)
	}
	plaintext := string(byteSecret)
	fmt.Println() // Newline

	// Get the environment variable name (can be displayed)
	fmt.Print("请输入环境变量名(带后缀): ")
	envName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("读取环境变量名失败: %v\n", err)
		os.Exit(1)
	}
	envName = strings.TrimSpace(envName)

	// Store the API key
	encValue, err := core.StoreAPIKey(plaintext, envName, key, ".env")
	if err != nil {
		fmt.Printf("存储API密钥失败: %v\n", err)
		os.Exit(1)
	}

	// Output the encryption result
	fmt.Printf("\n加密结果: %s\n", encValue)
	fmt.Println("已成功保存到.env文件")

	// Clear sensitive data from memory
	clearString(&plaintext)
}

// Load keys from the .env file into a new bash session
func loadKeysToNewBash(key string) {
	// Check if .env file exists
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		fmt.Println("错误: .env文件不存在")
		os.Exit(1)
	}

	// Set file permissions
	err := os.Chmod(".env", 0600)
	if err != nil {
		fmt.Printf("设置.env文件权限失败: %v\n", err)
		// Continue execution, don't exit
	}

	// Load and decrypt API keys
	decryptedVars, err := core.LoadAPIKeys(key, ".env")
	if err != nil {
		fmt.Printf("加载密钥失败: %v\n", err)
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

	// Close temporary file
	tempEnv.Close()

	// Start new bash session
	fmt.Println("\n正在启动新的bash会话，环境变量已设置...")

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

	fmt.Println("bash会话已结束，环境变量已清除")
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
