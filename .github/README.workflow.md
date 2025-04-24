# GitHub Actions 工作流

本目录包含用于自动化构建和发布的 GitHub Actions 工作流配置。

## 发布工作流 (release.yml)

这个工作流用于在创建新的 Git 标签时自动构建多平台二进制文件并发布到 GitHub Releases。

### 触发条件

当推送以 `v` 开头的标签时触发，例如 `v1.0.0`。

### 构建目标

工作流会为以下平台构建二进制文件：

- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

### 发布流程

1. 检出代码
2. 设置 Go 环境
3. 获取版本信息
4. 为每个目标平台构建二进制文件
5. 创建压缩包 (.zip 或 .tar.gz)
6. 上传构建产物
7. 创建 GitHub Release 并附加所有构建产物

### 如何使用

要创建一个新的发布版本，只需创建并推送一个新的标签：

```bash
# 创建标签
git tag v1.0.0

# 推送标签
git push origin v1.0.0
```

工作流将自动触发，构建二进制文件并创建 GitHub Release。

### 注意事项

- 确保 GitHub 仓库设置中已启用 Actions 功能
- 确保 GitHub 仓库设置中已授予 Actions 写入权限
- 工作流使用 `GITHUB_TOKEN` 创建 Release，无需额外配置
