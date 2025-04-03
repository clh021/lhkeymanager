#!/bin/bash

# 设置颜色
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}开始构建 LH Key Manager...${NC}"

# 检查Go是否安装
if ! command -v go &> /dev/null; then
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

# 询问安全规则
echo -e "${YELLOW}是否要自定义安全规则？(y/n)${NC}"
read -r customize

if [[ $customize == "y" || $customize == "Y" ]]; then
    # 询问 MinKeyLength
    echo -e "${YELLOW}输入最小密钥长度 (默认: 16):${NC}"
    read -r min_key_length
    min_key_length=${min_key_length:-16}

    # 询问 KeyPrefix
    echo -e "${YELLOW}输入必需的密钥前缀 (默认: lh-):${NC}"
    read -r key_prefix
    key_prefix=${key_prefix:-lh-}

    # 询问 KeySuffix
    echo -e "${YELLOW}输入必需的密钥后缀 (默认: u):${NC}"
    read -r key_suffix
    key_suffix=${key_suffix:-u}

    # 询问 RequiredChars
    echo -e "${YELLOW}输入必需的特殊字符 (默认: !@#$%^&*):${NC}"
    read -r required_chars
    required_chars=${required_chars:-!@#$%^&*}

    # 询问 MinSpecialChars
    echo -e "${YELLOW}输入最小特殊字符数量 (默认: 2):${NC}"
    read -r min_special_chars
    min_special_chars=${min_special_chars:-2}

    # 询问 KeyContain
    echo -e "${YELLOW}输入密钥必须包含的字符串 (默认: key):${NC}"
    read -r key_contain
    key_contain=${key_contain:-key}

    # 确认设置
    echo -e "${GREEN}安全规则:${NC}"
    echo -e "最小密钥长度: ${YELLOW}$min_key_length${NC}"
    echo -e "必需的密钥前缀: ${YELLOW}$key_prefix${NC}"
    echo -e "必需的密钥后缀: ${YELLOW}$key_suffix${NC}"
    echo -e "必需的特殊字符: ${YELLOW}$required_chars${NC}"
    echo -e "最小特殊字符数量: ${YELLOW}$min_special_chars${NC}"
    echo -e "必需包含的字符串: ${YELLOW}$key_contain${NC}"

    echo -e "${YELLOW}这些设置正确吗？(y/n)${NC}"
    read -r confirm
    if [[ $confirm != "y" && $confirm != "Y" ]]; then
        echo -e "${RED}构建已取消。${NC}"
        exit 1
    fi

    # 使用自定义设置构建
    ldflags="-X 'github.com/clh021/lhkeymanager/core.MinKeyLength=$min_key_length' \
             -X 'github.com/clh021/lhkeymanager/core.KeyPrefix=$key_prefix' \
             -X 'github.com/clh021/lhkeymanager/core.KeySuffix=$key_suffix' \
             -X 'github.com/clh021/lhkeymanager/core.RequiredChars=$required_chars' \
             -X 'github.com/clh021/lhkeymanager/core.MinSpecialChars=$min_special_chars' \
             -X 'github.com/clh021/lhkeymanager/core.KeyContain=$key_contain' \
             -s -w"
else
    # 使用默认设置构建
    ldflags="-s -w"
fi

# 运行测试
echo -e "${GREEN}运行测试...${NC}"
go test ./... || { echo -e "${RED}测试失败${NC}"; exit 1; }
echo -e "${GREEN}测试通过${NC}"

# 构建二进制文件
echo -e "${GREEN}构建二进制文件...${NC}"
go build -ldflags="$ldflags" -o lhkeymanager || { echo -e "${RED}构建失败${NC}"; exit 1; }
echo -e "${GREEN}构建成功: $(pwd)/lhkeymanager${NC}"

# 检查文件大小
ORIGINAL_SIZE=$(du -h lhkeymanager | cut -f1)
echo -e "${GREEN}原始二进制大小: $ORIGINAL_SIZE${NC}"

# 检查是否安装了UPX
if command -v upx &> /dev/null; then
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
