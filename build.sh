#!/bin/bash

# Set colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}Starting to build LH Key Manager...${NC}"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed. Please install Go 1.18 or higher.${NC}"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | grep -oP 'go\d+\.\d+' | grep -oP '\d+\.\d+')
GO_MAJOR=$(echo $GO_VERSION | cut -d. -f1)
GO_MINOR=$(echo $GO_VERSION | cut -d. -f2)

if [ "$GO_MAJOR" -lt 1 ] || ([ "$GO_MAJOR" -eq 1 ] && [ "$GO_MINOR" -lt 18 ]); then
    echo -e "${RED}Error: Go version is too old. Need Go 1.18 or higher, current version: $GO_VERSION${NC}"
    exit 1
fi

echo -e "${GREEN}Go version check passed: $GO_VERSION${NC}"

# Ask for security rules
echo -e "${YELLOW}Do you want to customize security rules? (y/n)${NC}"
read -r customize

if [[ $customize == "y" || $customize == "Y" ]]; then
    # Ask for MinKeyLength
    echo -e "${YELLOW}Enter minimum key length (default: 16):${NC}"
    read -r min_key_length
    min_key_length=${min_key_length:-16}

    # Ask for KeyPrefix
    echo -e "${YELLOW}Enter required key prefix (default: lh-, enter 'empty' for no prefix):${NC}"
    read -r key_prefix
    if [[ "$key_prefix" == "empty" ]]; then
        key_prefix=""
    else
        key_prefix=${key_prefix:-lh-}
    fi

    # Ask for KeySuffix
    echo -e "${YELLOW}Enter required key suffix (default: u, enter 'empty' for no suffix):${NC}"
    read -r key_suffix
    if [[ "$key_suffix" == "empty" ]]; then
        key_suffix=""
    else
        key_suffix=${key_suffix:-u}
    fi

    # Ask for RequiredChars
    echo -e "${YELLOW}Enter required special characters (default: !@#$%^&*, enter 'empty' for none):${NC}"
    read -r required_chars
    if [[ "$required_chars" == "empty" ]]; then
        required_chars=""
    else
        required_chars=${required_chars:-!@#$%^&*}
    fi

    # Ask for MinSpecialChars
    echo -e "${YELLOW}Enter minimum number of special characters (default: 2):${NC}"
    read -r min_special_chars
    min_special_chars=${min_special_chars:-2}

    # Ask for KeyContain
    echo -e "${YELLOW}Enter string that must be contained in the key (default: key, enter 'empty' for none):${NC}"
    read -r key_contain
    if [[ "$key_contain" == "empty" ]]; then
        key_contain=""
    else
        key_contain=${key_contain:-key}
    fi

    # Confirm settings
    echo -e "${GREEN}Security rules:${NC}"
    echo -e "Minimum key length: ${YELLOW}$min_key_length${NC}"
    echo -e "Required key prefix: ${YELLOW}$key_prefix${NC}"
    echo -e "Required key suffix: ${YELLOW}$key_suffix${NC}"
    echo -e "Required special characters: ${YELLOW}$required_chars${NC}"
    echo -e "Minimum number of special characters: ${YELLOW}$min_special_chars${NC}"
    echo -e "Required contained string: ${YELLOW}$key_contain${NC}"

    echo -e "${YELLOW}Are these settings correct? (y/n)${NC}"
    read -r confirm
    if [[ $confirm != "y" && $confirm != "Y" ]]; then
        echo -e "${RED}Build cancelled.${NC}"
        exit 1
    fi

    # Build with custom settings
    ldflags="-X 'github.com/clh021/lhkeymanager/core.MinKeyLength=$min_key_length' \
             -X 'github.com/clh021/lhkeymanager/core.KeyPrefix=$key_prefix' \
             -X 'github.com/clh021/lhkeymanager/core.KeySuffix=$key_suffix' \
             -X 'github.com/clh021/lhkeymanager/core.RequiredChars=$required_chars' \
             -X 'github.com/clh021/lhkeymanager/core.MinSpecialChars=$min_special_chars' \
             -X 'github.com/clh021/lhkeymanager/core.KeyContain=$key_contain' \
             -s -w"
else
    # Build with default settings
    ldflags="-s -w"
fi

# Run tests
echo -e "${GREEN}Running tests...${NC}"
go test ./... || { echo -e "${RED}Tests failed${NC}"; exit 1; }
echo -e "${GREEN}Tests passed${NC}"

# Build binary
echo -e "${GREEN}Building binary...${NC}"
go build -ldflags="$ldflags" -o lhkeymanager || { echo -e "${RED}Build failed${NC}"; exit 1; }
echo -e "${GREEN}Build successful: $(pwd)/lhkeymanager${NC}"

# Check file size
ORIGINAL_SIZE=$(du -h lhkeymanager | cut -f1)
echo -e "${GREEN}Original binary size: $ORIGINAL_SIZE${NC}"

# Check if UPX is installed
if command -v upx &> /dev/null; then
    echo -e "${GREEN}Using UPX to compress binary...${NC}"
    upx -9 lhkeymanager || { echo -e "${YELLOW}UPX compression failed, but this doesn't affect functionality${NC}"; }
    COMPRESSED_SIZE=$(du -h lhkeymanager | cut -f1)
    echo -e "${GREEN}Compressed binary size: $COMPRESSED_SIZE${NC}"
else
    echo -e "${YELLOW}Note: UPX not found. Installing UPX can further reduce binary size.${NC}"
    echo -e "${YELLOW}On Debian/Ubuntu: sudo apt-get install upx${NC}"
    echo -e "${YELLOW}On CentOS/RHEL: sudo yum install upx${NC}"
    echo -e "${YELLOW}On macOS: brew install upx${NC}"
fi

# Set executable permissions
chmod +x lhkeymanager

echo -e "${GREEN}Build complete!${NC}"
echo -e "${GREEN}You can run the program with:${NC}"
echo -e "${YELLOW}./lhkeymanager${NC}"
echo ""
echo -e "${GREEN}Note: You've built a version with custom security rules. Keep your encryption key secret!${NC}"
