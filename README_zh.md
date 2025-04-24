# LH 密钥管理器

[English README](README.md)

一个用于安全管理API密钥的工具，可以加密存储API密钥到`.env`文件，并在需要时将其加载到新的bash会话中。

## 功能特点

- 安全加密存储API密钥到`.env`文件
- 读取`.env`文件中的密钥到新bash会话
- 环境变量仅在新bash会话中有效，会话结束后自动清除
- 使用AES-256-GCM加密算法保护密钥安全
- 交互式输入，避免敏感信息被记录在bash历史记录中

## 安装

### 从 GitHub Releases 下载

安装 LH 密钥管理器最简单的方法是从 [GitHub Releases](https://github.com/clh021/lhkeymanager/releases) 页面下载预构建的二进制文件。

1. 访问 [Releases](https://github.com/clh021/lhkeymanager/releases) 页面
2. 下载适合您操作系统的版本：
   - Linux 系统：`lhkeymanager-vX.Y.Z-linux-amd64.tar.gz` 或 `lhkeymanager-vX.Y.Z-linux-arm64.tar.gz`
   - macOS 系统：`lhkeymanager-vX.Y.Z-darwin-amd64.tar.gz` 或 `lhkeymanager-vX.Y.Z-darwin-arm64.tar.gz`
   - Windows 系统：`lhkeymanager-vX.Y.Z-windows-amd64.zip`
3. 解压文件并设置可执行权限（Linux/macOS）：
   ```bash
   # Linux/macOS 系统
   tar -xzf lhkeymanager-vX.Y.Z-linux-amd64.tar.gz
   chmod +x lhkeymanager-vX.Y.Z-linux-amd64
   # 可选：移动到 PATH 目录中
   sudo mv lhkeymanager-vX.Y.Z-linux-amd64 /usr/local/bin/lhkeymanager
   ```

### 从源代码构建

如果您更喜欢从源代码构建：

#### 前提条件

- Go 1.18 或更高版本

```bash
# 克隆仓库
git clone https://github.com/clh021/lhkeymanager.git
cd lhkeymanager

# 使用构建脚本构建（推荐）
./build_zh.sh  # 中文版
# 或
./build.sh  # 英文版

# 或手动构建
go build -o lhkeymanager
```

### 自定义安全规则

为了提高安全性，您可以在构建过程中自定义加密密钥的验证规则：

1. 运行构建脚本，并在提示时选择自定义安全规则：
   ```bash
   ./build_zh.sh
   ```

2. 脚本将询问您配置以下安全规则：
   - `MinKeyLength`：加密密钥的最小长度（默认：16）
   - `KeyPrefix`：加密密钥的必需前缀（默认：lh-，输入 'empty' 表示无前缀要求）
   - `KeySuffix`：加密密钥的必需后缀（默认：u，输入 'empty' 表示无后缀要求）
   - `RequiredChars`：密钥中必须包含的字符（默认：!@#$%^&*，输入 'empty' 表示无特殊字符要求）
   - `MinSpecialChars`：所需的最少特殊字符数量（默认：2）
   - `KeyContain`：密钥中必须包含的字符串（默认：key，输入 'empty' 表示无包含要求）

这样，只有您知道有效加密密钥的确切规则，即使他人获取了您的加密数据，也更难猜到您的密钥。

## 使用方法

### 存储新的API密钥

```bash
./lhkeymanager
```

然后选择选项`1`，按照提示输入加密密钥和API密钥。

### 读取密钥到新bash会话

```bash
./lhkeymanager
```

然后选择选项`2`，输入加密密钥，工具会启动一个新的bash会话，并在其中设置环境变量。

## 安全注意事项

- `.env`文件权限会被自动设置为600（仅所有者可读写）
- 加密密钥不会被存储，每次使用时需要手动输入
- 环境变量仅在新bash会话中有效，会话结束后自动清除
- 临时文件会在使用后安全删除

## 示例

### 存储OpenAI API密钥

```
$ ./lhkeymanager
请选择操作:
1. 存储新的API密钥到.env文件
2. 读取.env文件中的密钥到新bash会话
请输入选项 (1/2): 1
请输入加密密钥: [输入但不显示]
请输入要加密的API密钥: [输入但不显示]
请输入环境变量名(带后缀): OPENAI_API_KEY_PROD

加密结果: enc:AES256:AAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERUZHSElKS0xNTk9QUVJTVFVWV1hZWltcXV5fYGFiY2RlZmdo
已成功保存到.env文件
```

### 读取密钥到新bash会话

```
$ ./lhkeymanager
请选择操作:
1. 存储新的API密钥到.env文件
2. 读取.env文件中的密钥到新bash会话
请输入选项 (1/2): 2
请输入加密密钥: [输入但不显示]
已设置环境变量: OPENAI_API_KEY

正在启动新的bash会话，环境变量已设置...
$ echo $OPENAI_API_KEY
sk-your-api-key
$ exit
bash会话已结束，环境变量已清除
```

## 许可证

MIT

## 贡献

欢迎贡献！请随时提交Pull Request。

### 创建新版本

创建新版本的方法：

1. 使用提供的脚本：
   ```bash
   ./create-release.sh v1.0.0
   ```
   这将创建并推送一个新标签，触发 GitHub Actions 工作流来构建和发布版本。

2. 或者，您可以手动创建并推送标签：
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

3. 您还可以从 GitHub Actions 选项卡手动触发发布，方法是选择"Release"工作流并点击"Run workflow"。这允许您为构建指定自定义安全规则。

## 致谢

- 本项目使用AES-256-GCM加密算法进行安全密钥存储
- 灵感来源于开发环境中安全管理API密钥的需求

