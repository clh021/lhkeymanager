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

# 运行测试
echo -e "${GREEN}运行测试...${NC}"
go test ./... || { echo -e "${RED}测试失败${NC}"; exit 1; }
echo -e "${GREEN}测试通过${NC}"

# 构建二进制文件
echo -e "${GREEN}构建二进制文件...${NC}"
go build -ldflags="-s -w" -o lhkeymanager || { echo -e "${RED}构建失败${NC}"; exit 1; }
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
echo -e "${GREEN}提示: 您可以在core/keymanager.go中自定义密钥验证规则，然后重新运行此脚本以增强安全性。${NC}"
