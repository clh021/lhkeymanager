#!/bin/bash

#
# 功能: 将一个包含多个变量的 .env 文件，拆分成多个以变量名为基础的独立 .env 文件。
#
# 规则:
#   - 输出文件名基于变量名第一个下划线 `_` 之前的部分，并转换为小写。
#   - 例如: `GEMINI_API_KEY=...` 会被拆分到 `gemini.env` 文件中。
#   - 例如: `SINGLE_WORD=...` 会被拆分到 `single.env` 文件中。
#

# --- 用于美化输出的颜色定义 ---
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# --- 用法说明 ---
usage() {
	echo "用法: $0 <input_file.env>"
	echo ""
	echo "    将一个多变量的 .env 文件拆分为多个独立的 .env 文件。"
	echo ""
	echo "    示例: "
	echo "        $0 secrets.env"
	exit 1
}

# --- 参数和文件检查 ---
if [ "$#" -ne 1 ]; then
	echo -e "${RED}错误: 请提供一个输入文件作为参数。${NC}"
	usage
fi

INPUT_FILE="$1"

if [ ! -f "$INPUT_FILE" ]; then
	echo -e "${RED}错误: 输入文件 '$INPUT_FILE' 未找到。${NC}"
	exit 1
fi

# --- 主处理逻辑 ---
echo -e "${GREEN}正在处理文件: ${YELLOW}$INPUT_FILE${NC}"

# 使用 `while` 循环逐行读取文件，这种方式可以安全地处理包含空格或特殊字符的行。
while IFS= read -r line || [[ -n "$line" ]]; do
	# 跳过空行和注释行
	if [[ -z "$line" || "$line" =~ ^# ]]; then
		continue
	fi

	# 确保行中包含 '='
	if [[ ! "$line" =~ = ]]; then
		echo -e "${YELLOW}警告: 跳过格式错误的行 (未找到 '=') -> $line${NC}"
		continue
	fi

	# 提取变量名 (第一个 '=' 之前的部分)
	var_name="${line%%=*}"

	# 获取基础名称 (第一个 '_' 之前的部分) 并转换为小写
	# 这能处理 'GEMINI_API_KEY' -> 'gemini'
	# 也能处理 'SINGLEWORD' -> 'singleword'
	base_name_upper="${var_name%%_*}"
	base_name_lower=$(echo "$base_name_upper" | tr '[:upper:]' '[:lower:]')

	# 构建输出文件名
	output_file="${base_name_lower}.env"

	# 将原始的整行内容写入新文件，如果文件已存在则覆盖
	echo "$line" >"$output_file"

	if [ $? -eq 0 ]; then
		echo -e "  -> 已为变量 ${YELLOW}$var_name${NC} 创建文件 ${GREEN}$output_file${NC}"
	else
		echo -e "  -> ${RED}为变量 $var_name 创建文件时出错${NC}"
	fi

done <"$INPUT_FILE"

echo -e "${GREEN}处理完成。${NC}"
