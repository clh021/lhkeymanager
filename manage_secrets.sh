#!/bin/bash

# 设置颜色
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# --- 用法 ---
usage() {
	echo "用法: $0 <command> <input_file>"
	echo ""
	echo "一个使用 lhkeymanager 加密和解密 .env 文件的脚本。"
	echo ""
	echo "命令:"
	echo "  encrypt <plaintext_file>   将纯文本文件 (例如 secrets.env.dev) 加密为加密文件 (secrets.env)。"
	echo "  decrypt <encrypted_file>   将加密文件 (例如 secrets.env) 解密为纯文本文件 (secrets.env.decrypted)。"
	echo ""
	echo "示例:"
	echo "  $0 encrypt secrets.env.dev"
	echo "  $0 decrypt secrets.env"
	exit 1
}

# --- 前置检查 ---
if [ "$#" -ne 2 ]; then
	usage
fi

COMMAND=$1
INPUT_FILE=$2
LHKEYMANAGER_BIN="./lhkeymanager"

if [ ! -f "$LHKEYMANAGER_BIN" ]; then
	echo -e "${RED}错误: 在当前目录中未找到 lhkeymanager 二进制文件。${NC}"
	echo -e "${YELLOW}请先使用 ./build_zh.sh 或 ./build.sh 构建它。${NC}"
	exit 1
fi

if [ ! -f "$INPUT_FILE" ]; then
	echo -e "${RED}错误: 输入文件 '$INPUT_FILE' 未找到。${NC}"
	exit 1
fi

# --- 主逻辑 ---
case "$COMMAND" in
encrypt)
	# 输入: secrets.env.dev -> 输出: secrets.env
	if [[ "$INPUT_FILE" != *.dev ]]; then
		echo -e "${RED}错误: 'encrypt' 命令的输入文件必须以 .dev 结尾。${NC}"
		exit 1
	fi
	OUTPUT_FILE="${INPUT_FILE%.dev}"

	echo -e "${GREEN}正在加密 '${YELLOW}$INPUT_FILE${GREEN}' 到 '${YELLOW}$OUTPUT_FILE${GREEN}'...${NC}"
	"$LHKEYMANAGER_BIN" encrypt-file "$INPUT_FILE" "$OUTPUT_FILE"
	;;
decrypt)
	# 输入: secrets.env -> 输出: secrets.env.decrypted
	OUTPUT_FILE="${INPUT_FILE}.decrypted"

	echo -e "${GREEN}正在解密 '${YELLOW}$INPUT_FILE${GREEN}' 到 '${YELLOW}$OUTPUT_FILE${GREEN}'...${NC}"
	"$LHKEYMANAGER_BIN" decrypt-file "$INPUT_FILE" "$OUTPUT_FILE"
	;;
*)
	echo -e "${RED}错误: 未知命令 '$COMMAND'。${NC}"
	usage
	;;
esac

# --- 后置检查 ---
if [ $? -eq 0 ]; then
	echo -e "${GREEN}操作成功。${NC}"
	echo -e "您现在可以比较结果，例如:"
	if [ "$COMMAND" == "decrypt" ]; then
		# 假设原始文件以 .dev 结尾
		original_dev_file="${INPUT_FILE%.*}.dev"
		if [ -f "$original_dev_file" ]; then
			echo -e "${YELLOW}diff ${original_dev_file} ${OUTPUT_FILE}${NC}"
		fi
	fi
else
	echo -e "${RED}操作失败。${NC}"
fi
