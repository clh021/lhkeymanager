#!/bin/bash

# Set colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}Starting to build LH Key Manager...${NC}"

# Check if Go is installed
if ! command -v go &>/dev/null; then
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

# Function to display security rules
display_security_rules() {
	echo -e "${GREEN}Security rules:${NC}"
	echo -e "Minimum key length: ${YELLOW}$min_key_length${NC}"
	echo -e "Required key prefix: ${YELLOW}$key_prefix${NC}"
	echo -e "Required key suffix: ${YELLOW}$key_suffix${NC}"
	echo -e "Required special characters: ${YELLOW}$required_chars${NC}"
	echo -e "Minimum number of special characters: ${YELLOW}$min_special_chars${NC}"
	echo -e "Required contained string: ${YELLOW}$key_contain${NC}"
	echo -e "Temporary key: ${YELLOW}${temp_key:-None}${NC}"
	echo -e "Temporary key max usage: ${YELLOW}$temp_key_max_usage${NC}"
	echo -e "Key hint: ${YELLOW}$key_hint${NC}"
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
echo -e "${YELLOW}Please choose security rules configuration method:${NC}"
echo -e "  ${YELLOW}1. Interactive customization${NC}"
echo -e "  ${YELLOW}2. Use build_config.yml file (if exists)${NC}"
echo -e "  ${YELLOW}3. Use default settings${NC}"
read -r config_choice

ldflags="" # Initialize ldflags

if [[ "$config_choice" == "1" ]]; then
	echo -e "${GREEN}Entering interactive customization mode...${NC}"
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

	# Ask for TempKey
	echo -e "${YELLOW}Enter temporary key (default: None, enter 'empty' for none):${NC}"
	read -r temp_key
	if [[ "$temp_key" == "empty" ]]; then
		temp_key=""
	fi

	# Ask for TempKeyMaxUsage
	echo -e "${YELLOW}Enter temporary key max usage (default: 2):${NC}"
	read -r temp_key_max_usage
	temp_key_max_usage=${temp_key_max_usage:-2}

	# Ask for KeyHint
	echo -e "${YELLOW}Enter key hint (default: No hint available.):${NC}"
	read -r key_hint
	if [[ "$key_hint" == "" ]]; then
		key_hint="No hint available."
	else
		key_hint=${key_hint:-"No hint available."}
	fi

	display_security_rules

	echo -e "${YELLOW}Are these settings correct? (y/n)${NC}"
	read -r confirm
	if [[ $confirm != "y" && $confirm != "Y" ]]; then
		echo -e "${RED}Build cancelled.${NC}"
		exit 1
	fi

	construct_ldflags

elif [[ "$config_choice" == "2" ]]; then
	CONFIG_FILE="build_config.yml"
	if [ ! -f "$CONFIG_FILE" ]; then
		echo -e "${RED}Error: build_config.yml file does not exist. Please choose another option or create the file.${NC}"
		exit 1
	fi

	# Check for yq
	if ! command -v yq &>/dev/null; then
		echo -e "${RED}Error: yq tool is not installed. Please install yq (https://github.com/mikefarah/yq) to use config file feature.${NC}"
		echo -e "${RED}E.g.: sudo snap install yq or brew install yq${NC}"
		exit 1
	fi

	echo -e "${GREEN}Reading configuration from build_config.yml...${NC}"

	# Function to get value from YAML with default and handle "empty" string
	get_config_value() {
		local path=$1
		local default_value=$2
		# The yq command is changed to `yq -r` to support the more common
		# Python-based yq (a wrapper for jq), which does not have the 'e' command.
		# The -r flag provides raw string output.
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
	echo -e "${GREEN}Building with default settings.${NC}"
	ldflags="-s -w" # Default ldflags
else
	echo -e "${RED}Invalid choice. Build cancelled.${NC}"
	exit 1
fi

# Run tests
echo -e "${GREEN}Running tests...${NC}"
go test ./... || {
	echo -e "${RED}Tests failed${NC}"
	exit 1
}
echo -e "${GREEN}Tests passed${NC}"

# Build binary
echo -e "${GREEN}Building binary...${NC}"
go build -ldflags="$ldflags" -o lhkeymanager || {
	echo -e "${RED}Build failed${NC}"
	exit 1
}
echo -e "${GREEN}Build successful: $(pwd)/lhkeymanager${NC}"

# Check file size
ORIGINAL_SIZE=$(du -h lhkeymanager | cut -f1)
echo -e "${GREEN}Original binary size: $ORIGINAL_SIZE${NC}"

# Check if UPX is installed
if command -v upx &>/dev/null; then
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
