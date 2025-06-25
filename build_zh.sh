#!/bin/bash

# 设置颜色
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}开始构建 LH Key Manager...${NC}"

# 检查Go是否安装
if ! command -v go &>/dev/null; then
	echo -e "${RED}错误: Go 未安装。请安装 Go 1.18 或更高版本。${NC}"
	exit 1
fi

# 检查Go版本
GO_VERSION=$(go version | grep -oP 'go\d+\.\d+' | grep -oP '\d+\.\d+')
GO_MAJOR=$(echo $GO_VERSION | cut -d. -f1)
GO_MINOR=$(echo $GO_VERSION | cut -d. -f2)

if [ "$GO_MAJOR" -lt 1 ] || ([ "$GO_MAJOR" -eq 1 ] && [ "$GO_MINOR" -lt 18 ]); then
	echo -e "${RED}错误: Go 版本过低。需要 Go 1.18 或更高版本，当前版本: $GO_VERSION${NC}"
	exit 1
fi

echo -e "${GREEN}Go 版本检查通过: $GO_VERSION${NC}"

# Function to display security rules
display_security_rules() {
	echo -e "${GREEN}安全规则:${NC}"
	echo -e "最小密钥长度: ${YELLOW}$min_key_length${NC}"
	echo -e "必需的密钥前缀: ${YELLOW}$key_prefix${NC}"
	echo -e "必需的密钥后缀: ${YELLOW}$key_suffix${NC}"
	echo -e "必需的特殊字符: ${YELLOW}$required_chars${NC}"
	echo -e "最小特殊字符数量: ${YELLOW}$min_special_chars${NC}"
	echo -e "必需包含的字符串: ${YELLOW}$key_contain${NC}"
	echo -e "临时密钥: ${YELLOW}${temp_key:-无}${NC}"
	echo -e "临时密钥最大使用次数: ${YELLOW}$temp_key_max_usage${NC}"
	echo -e "密钥提示: ${YELLOW}$key_hint${NC}"
}

# Function to construct ldflags
construct_ldflags() {
	ldflags="-X 'github.com/clh021/lhkeymanager/core.MinKeyLength=$min_key_length' \
             -X 'github.com/clh021/lhkeymanager/core.KeyPrefix=$key_prefix' \
             -X 'github.com/clh021/lhkeymanager/core.KeySuffix=$key_suffix' \
             -X 'github.com/clh021/lhkeymanager/core.RequiredChars=$required_chars' \
             -X 'github.com/clh021/lhkeymanager/core.MinSpecialChars=$min_special_chars' \
             -X 'github.com/clh021/lhkeymanager/core.KeyContain=$key_contain' \
             -X 'github.com/clh021/lhkeymanager/core.TempKey=$temp_key' \
             -X 'github.com/clh021/lhkeymanager/core.TempKeyMaxUsage=$temp_key_max_usage' \
             -X 'github.com/clh021/lhkeymanager/core.KeyHint=$key_hint' \
             -s -w"
}

# Ask for security rules configuration method
echo -e "${YELLOW}请选择安全规则配置方式:${NC}"
echo -e "  ${YELLOW}1. 交互式自定义${NC}"
echo -e "  ${YELLOW}2. 使用 build_config.yml 文件 (如果存在)${NC}"
echo -e "  ${YELLOW}3. 使用默认设置${NC}"
read -r config_choice

ldflags="" # Initialize ldflags

if [[ "$config_choice" == "1" ]]; then
	echo -e "${GREEN}进入交互式自定义模式...${NC}"
	# 询问 MinKeyLength
	echo -e "${YELLOW}输入最小密钥长度 (默认: 16):${NC}"
	read -r min_key_length
	min_key_length=${min_key_length:-16}

	# 询问 KeyPrefix
	echo -e "${YELLOW}输入必需的密钥前缀 (默认: lh-, 输入 'empty' 表示无前缀):${NC}"
	read -r key_prefix
	if [[ "$key_prefix" == "empty" ]]; then
		key_prefix=""
	else
		key_prefix=${key_prefix:-lh-}
	fi

	# 询问 KeySuffix
	echo -e "${YELLOW}输入必需的密钥后缀 (默认: u, 输入 'empty' 表示无后缀):${NC}"
	read -r key_suffix
	if [[ "$key_suffix" == "empty" ]]; then
		key_suffix=""
	else
		key_suffix=${key_suffix:-u}
	fi

	# 询问 RequiredChars
	echo -e "${YELLOW}输入必需的特殊字符 (默认: !@#$%^&*, 输入 'empty' 表示无特殊字符要求):${NC}"
	read -r required_chars
	if [[ "$required_chars" == "empty" ]]; then
		required_chars=""
	else
		required_chars=${required_chars:-!@#$%^&*}
	fi

	# 询问 MinSpecialChars
	echo -e "${YELLOW}输入最小特殊字符数量 (默认: 2):${NC}"
	read -r min_special_chars
	min_special_chars=${min_special_chars:-2}

	# 询问 KeyContain
	echo -e "${YELLOW}输入密钥必须包含的字符串 (默认: key, 输入 'empty' 表示无包含要求):${NC}"
	read -r key_contain
	if [[ "$key_contain" == "empty" ]]; then
		key_contain=""
	else
		key_contain=${key_contain:-key}
	fi

	# 询问 TempKey
	echo -e "${YELLOW}输入临时密钥 (默认: 无, 输入 'empty' 表示无):${NC}"
	read -r temp_key
	if [[ "$temp_key" == "empty" ]]; then
		temp_key=""
	fi

	# 询问 TempKeyMaxUsage
	echo -e "${YELLOW}输入临时密钥最大使用次数 (默认: 2):${NC}"
	read -r temp_key_max_usage
	temp_key_max_usage=${temp_key_max_usage:-2}

	# 询问 KeyHint
	echo -e "${YELLOW}输入密钥提示 (默认: 无提示):${NC}"
	read -r key_hint
	if [[ "$key_hint" == "" ]]; then
		key_hint="No hint available."
	else
		key_hint=${key_hint:-"No hint available."}
	fi

	display_security_rules

	echo -e "${YELLOW}这些设置正确吗？(y/n)${NC}"
	read -r confirm
	if [[ $confirm != "y" && $confirm != "Y" ]]; then
		echo -e "${RED}构建已取消。${NC}"
		exit 1
	fi

	construct_ldflags

elif [[ "$config_choice" == "2" ]]; then
	CONFIG_FILE="build_config.yml"
	if [ ! -f "$CONFIG_FILE" ]; then
		echo -e "${RED}错误: build_config.yml 文件不存在。请选择其他选项或创建该文件。${NC}"
		exit 1
	fi

	# Check for yq
	if ! command -v yq &>/dev/null; then
		echo -e "${RED}错误: yq 工具未安装。请安装 yq (https://github.com/mikefarah/yq) 以使用配置文件功能。${NC}"
		echo -e "${RED}例如: sudo snap install yq 或 brew install yq${NC}"
		exit 1
	fi

	echo -e "${GREEN}从 build_config.yml 读取配置...${NC}"

	# Function to get value from YAML with default and handle "empty" string
	get_config_value() {
		local path=$1
		local default_value=$2
		# yq 命令被修改为 `yq -r` 以支持更常见的
		# 基于 Python 的 yq (jq 的包装器), 它没有 'e' 命令。
		# -r 标志用于提供原始字符串输出。
		local value=$(yq -r "$path // \"$default_value\"" "$CONFIG_FILE")
		# If the value from YAML is literally "empty", treat it as an empty string
		if [[ "$value" == "empty" ]]; then
			echo ""
		else
			echo "$value"
		fi
	}

	min_key_length=$(get_config_value ".security_rules.min_key_length" "16")
	key_prefix=$(get_config_value ".security_rules.key_prefix" "lh-")
	key_suffix=$(get_config_value ".security_rules.key_suffix" "u")
	required_chars=$(get_config_value ".security_rules.required_chars" "!@#$%^&*")
	min_special_chars=$(get_config_value ".security_rules.min_special_chars" "2")
	key_contain=$(get_config_value ".security_rules.key_contain" "key")
	temp_key=$(get_config_value ".security_rules.temp_key" "")
	temp_key_max_usage=$(get_config_value ".security_rules.temp_key_max_usage" "2")
	key_hint=$(get_config_value ".security_rules.key_hint" "No hint available.")

	display_security_rules
	construct_ldflags

elif [[ "$config_choice" == "3" ]]; then
	echo -e "${GREEN}使用默认设置构建。${NC}"
	ldflags="-s -w" # Default ldflags
else
	echo -e "${RED}无效的选择。构建已取消。${NC}"
	exit 1
fi

# 运行测试
echo -e "${GREEN}运行测试...${NC}"
go test ./... || {
	echo -e "${RED}测试失败${NC}"
	exit 1
}
echo -e "${GREEN}测试通过${NC}"

# 构建二进制文件
echo -e "${GREEN}构建二进制文件...${NC}"
go build -ldflags="$ldflags" -o lhkeymanager || {
	echo -e "${RED}构建失败${NC}"
	exit 1
}
echo -e "${GREEN}构建成功: $(pwd)/lhkeymanager${NC}"

# 检查文件大小
ORIGINAL_SIZE=$(du -h lhkeymanager | cut -f1)
echo -e "${GREEN}原始二进制大小: $ORIGINAL_SIZE${NC}"

# 检查是否安装了UPX
if command -v upx &>/dev/null; then
	echo -e "${GREEN}使用UPX压缩二进制文件...${NC}"
	upx -9 lhkeymanager || { echo -e "${YELLOW}UPX压缩失败，但这不影响程序功能${NC}"; }
	COMPRESSED_SIZE=$(du -h lhkeymanager | cut -f1)
	echo -e "${GREEN}压缩后二进制大小: $COMPRESSED_SIZE${NC}"
else
	echo -e "${YELLOW}提示: 未找到UPX。安装UPX可以进一步减小二进制文件大小。${NC}"
	echo -e "${YELLOW}在Debian/Ubuntu上: sudo apt-get install upx${NC}"
	echo -e "${YELLOW}在CentOS/RHEL上: sudo yum install upx${NC}"
	echo -e "${YELLOW}在macOS上: brew install upx${NC}"
fi

# 设置可执行权限
chmod +x lhkeymanager

echo -e "${GREEN}构建完成!${NC}"
echo -e "${GREEN}您可以通过以下命令运行程序:${NC}"
echo -e "${YELLOW}./lhkeymanager${NC}"
echo ""
echo -e "${GREEN}注意: 您已构建了一个具有自定义安全规则的版本。请保持您的加密密钥安全!${NC}"
