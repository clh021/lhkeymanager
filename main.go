package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/clh021/lhkeymanager/utils"

	"golang.org/x/term"
)

const (
	MinKeyLength = 14 // 不向用户透露这个信息
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

	// 验证密钥长度和结尾字符，但不透露具体要求
	if len(key) < MinKeyLength || !strings.HasSuffix(key, "u") {
		fmt.Println("错误: 密钥验证失败")
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

// 存储新的API密钥到.env文件
func storeKey(reader *bufio.Reader, key string) {
	// 获取要加密的内容（不显示输入）
	fmt.Print("请输入要加密的API密钥: ")
	byteSecret, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Printf("\n读取API密钥失败: %v\n", err)
		os.Exit(1)
	}
	plaintext := string(byteSecret)
	fmt.Println() // 换行

	// 获取环境变量名（可以显示）
	fmt.Print("请输入环境变量名(带后缀): ")
	envName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("读取环境变量名失败: %v\n", err)
		os.Exit(1)
	}
	envName = strings.TrimSpace(envName)

	// 加密
	encrypted, err := utils.EncryptAES256(plaintext, key)
	if err != nil {
		fmt.Printf("加密失败: %v\n", err)
		os.Exit(1)
	}

	// 输出加密结果
	encValue := fmt.Sprintf("enc:AES256:%s", encrypted)
	fmt.Printf("\n加密结果: %s\n", encValue)

	// 保存到.env文件
	err = saveToEnvFile(envName, encValue)
	if err != nil {
		fmt.Printf("保存到.env文件失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("已成功保存到.env文件")

	// 清理内存中的敏感数据
	clearString(&plaintext)
}

// 读取.env文件中的密钥到新bash会话
func loadKeysToNewBash(key string) {
	// 检查.env文件是否存在
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		fmt.Println("错误: .env文件不存在")
		os.Exit(1)
	}

	// 设置文件权限
	err := os.Chmod(".env", 0600)
	if err != nil {
		fmt.Printf("设置.env文件权限失败: %v\n", err)
		// 继续执行，不退出
	}

	// 创建临时环境变量文件
	tempEnv, err := os.CreateTemp("", "env_vars_*")
	if err != nil {
		fmt.Printf("创建临时文件失败: %v\n", err)
		os.Exit(1)
	}
	tempEnvPath := tempEnv.Name()
	defer os.Remove(tempEnvPath) // 确保退出时删除临时文件

	// 设置临时文件权限
	err = os.Chmod(tempEnvPath, 0600)
	if err != nil {
		fmt.Printf("设置临时文件权限失败: %v\n", err)
		// 继续执行，不退出
	}

	// 写入基本shell配置
	tempEnv.WriteString("#!/bin/bash\n")
	tempEnv.WriteString("# 这是自动生成的临时环境变量文件\n\n")

	// 读取.env文件并处理每一行
	envFile, err := os.Open(".env")
	if err != nil {
		fmt.Printf("打开.env文件失败: %v\n", err)
		os.Exit(1)
	}
	defer envFile.Close()

	scanner := bufio.NewScanner(envFile)
	decryptionSuccess := false // 跟踪是否有任何成功解密

	for scanner.Scan() {
		line := scanner.Text()

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析环境变量
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // 跳过格式不正确的行
		}

		name := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// 检查是否是加密值
		if strings.HasPrefix(value, "enc:AES256:") {
			encData := strings.TrimPrefix(value, "enc:AES256:")

			// 解密
			decrypted, err := utils.DecryptAES256(encData, key)
			if err != nil {
				// 不提示具体错误，继续处理下一个
				continue
			}

			decryptionSuccess = true

			// 去除后缀并设置环境变量
			cleanName := name
			if strings.Contains(name, "_") {
				cleanName = name[:strings.LastIndex(name, "_")]
			}

			// 写入临时文件
			tempEnv.WriteString(fmt.Sprintf("export %s='%s'\n", cleanName, decrypted))
			fmt.Printf("已设置环境变量: %s\n", cleanName)
		} else {
			// 非加密值直接设置
			tempEnv.WriteString(fmt.Sprintf("export %s='%s'\n", name, value))
			fmt.Printf("已设置环境变量: %s\n", name)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("读取.env文件失败: %v\n", err)
		os.Exit(1)
	}

	// 如果没有成功解密任何变量，提示密钥错误
	if !decryptionSuccess {
		fmt.Println("错误: 密钥验证失败")
		os.Exit(1)
	}

	// 关闭临时文件
	tempEnv.Close()

	// 启动新的bash会话
	fmt.Println("\n正在启动新的bash会话，环境变量已设置...")

	// 使用source命令加载环境变量并启动新bash
	cmd := exec.Command("bash", "-c", fmt.Sprintf("source %s && bash", tempEnvPath))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("启动bash会话失败: %v\n", err)
		os.Exit(1)
	}

	// 安全删除临时文件
	secureDeleteFile(tempEnvPath)

	fmt.Println("bash会话已结束，环境变量已清除")
}

// 保存到.env文件
func saveToEnvFile(name, value string) error {
	// 检查.env文件是否存在
	envFile := ".env"
	var file *os.File
	var err error

	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		// 文件不存在，创建新文件
		file, err = os.Create(envFile)
		if err != nil {
			return err
		}
	} else {
		// 文件存在，以追加模式打开
		file, err = os.OpenFile(envFile, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
	}
	defer file.Close()

	// 设置文件权限为600（仅所有者可读写）
	err = os.Chmod(envFile, 0600)
	if err != nil {
		return err
	}

	// 写入环境变量
	_, err = fmt.Fprintf(file, "%s=%s\n", name, value)
	return err
}

// 清理字符串内存
func clearString(s *string) {
	if s == nil {
		return
	}
	for i := range *s {
		(*s) = (*s)[:i] + "\x00" + (*s)[i+1:]
	}
	*s = ""
}

// 安全删除文件
func secureDeleteFile(path string) {
	// 尝试使用shred命令安全删除
	shredCmd := exec.Command("shred", "-u", "-z", path)
	err := shredCmd.Run()

	// 如果shred命令不可用，使用普通删除
	if err != nil {
		os.Remove(path)
	}
}
