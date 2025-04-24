#!/bin/bash

# 这个脚本用于创建一个新的发布版本
# 用法: ./scripts/create-release.sh v1.0.0

set -e

# 检查是否提供了版本参数
if [ $# -ne 1 ]; then
  echo "错误: 请提供版本号"
  echo "用法: $0 <版本号>"
  echo "示例: $0 v1.0.0"
  echo "        v0.0.2-beta、v1.2.3-alpha、v2.0.0-rc1"
  exit 1
fi

VERSION=$1

# 检查版本号格式
if [[ ! $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9]+)?$ ]]; then
  echo "错误: 版本号格式不正确"
  echo "版本号应该以 'v' 开头，后跟语义化版本号，例如 v1.0.0 或带后缀的版本号如 v1.0.0-beta"
  exit 1
fi

# 检查工作目录是否干净
if [ -n "$(git status --porcelain)" ]; then
  echo "错误: 工作目录不干净，请先提交或暂存所有更改"
  git status
  exit 1
fi

# 创建标签
echo "创建标签 $VERSION..."
git tag -a $VERSION -m "Release $VERSION"

# 推送标签
echo "推送标签到远程仓库..."
git push origin $VERSION

echo "标签 $VERSION 已创建并推送到远程仓库"
echo "GitHub Actions 工作流将自动触发构建和发布"
echo "请访问 GitHub 仓库的 Actions 页面查看进度"
