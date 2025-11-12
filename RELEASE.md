# 发布指南 / Release Guide

本文档说明如何使用自动化工具发布 switcher 的新版本。

## 快速开始

发布新版本只需要三步：

```bash
# 1. 确保所有更改已提交
git add .
git commit -m "feat: your changes"
git push

# 2. 创建并推送版本标签
git tag v1.0.0
git push origin v1.0.0

# 3. 等待 GitHub Actions 自动构建和发布
```

就这么简单！GitHub Actions 会自动：
- 编译 6 个平台的二进制文件（Linux/macOS/Windows，x86_64/ARM64）
- 生成 checksums 文件
- 创建 GitHub Release
- 自动生成 changelog
- 上传所有文件

## 详细说明

### 1. 版本号规范

建议使用语义化版本号（Semantic Versioning）：

- `v1.0.0` - 主版本号.次版本号.修订号
- `v1.0.0-beta.1` - 预发布版本
- `v1.0.0-rc.1` - 候选发布版本

### 2. 使用 Makefile 辅助工具

项目提供了便捷的 Makefile 命令：

```bash
# 查看所有可用命令
make help

# 创建并推送新标签（交互式）
make tag

# 本地测试 GoReleaser 配置（不会发布）
make release-test

# 本地构建所有平台二进制
make build-all

# 查看当前版本信息
./switcher --version
```

### 3. 发布流程详解

#### 步骤 1: 准备代码

确保所有更改已提交并推送到 GitHub：

```bash
git status
git add .
git commit -m "feat: add new feature"
git push origin main
```

#### 步骤 2: 创建版本标签

**方式 A: 使用 Makefile（推荐）**

```bash
make tag
```

这会显示最近的标签，并提示你输入新版本号。

**方式 B: 手动创建**

```bash
# 创建带注释的标签
git tag -a v1.0.0 -m "Release v1.0.0"

# 推送标签到 GitHub
git push origin v1.0.0
```

#### 步骤 3: 监控构建过程

1. 访问你的 GitHub 仓库
2. 点击 "Actions" 标签
3. 查看 "Release" 工作流的运行状态

构建通常需要 2-5 分钟。

#### 步骤 4: 验证发布

构建完成后：

1. 访问 GitHub 仓库的 "Releases" 页面
2. 确认新版本已创建
3. 检查是否包含所有平台的二进制文件：
   - `switcher_v1.0.0_Linux_x86_64.tar.gz`
   - `switcher_v1.0.0_Linux_arm64.tar.gz`
   - `switcher_v1.0.0_Darwin_x86_64.tar.gz`
   - `switcher_v1.0.0_Darwin_arm64.tar.gz`
   - `switcher_v1.0.0_Windows_x86_64.zip`
   - `switcher_v1.0.0_Windows_arm64.zip`
   - `checksums.txt`

### 4. 本地测试（发布前）

在实际发布前，可以在本地测试 GoReleaser 配置：

```bash
# 安装 GoReleaser（如果还没安装）
# macOS
brew install goreleaser

# Linux
# 参考: https://goreleaser.com/install/

# 测试配置（不会发布）
make release-test

# 查看生成的文件
ls -lh dist/
```

### 5. Commit 消息规范（用于自动生成 Changelog）

为了生成更好的 changelog，建议使用约定式提交（Conventional Commits）：

- `feat: 添加新功能` - 新功能
- `fix: 修复 bug` - Bug 修复
- `perf: 性能优化` - 性能改进
- `refactor: 重构代码` - 代码重构
- `docs: 更新文档` - 文档更新（不会出现在 changelog 中）
- `test: 添加测试` - 测试相关（不会出现在 changelog 中）
- `chore: 构建/工具链` - 构建工具（不会出现在 changelog 中）

示例：
```bash
git commit -m "feat: add support for custom config path"
git commit -m "fix: resolve crash when config file is missing"
git commit -m "perf: optimize configuration loading speed"
```

### 6. 常见问题

#### Q: 如何删除错误的标签？

```bash
# 删除本地标签
git tag -d v1.0.0

# 删除远程标签
git push origin :refs/tags/v1.0.0
```

#### Q: 如何修改已发布的 Release？

1. 在 GitHub Release 页面点击 "Edit"
2. 修改描述或上传新文件
3. 保存更改

注意：不建议修改已发布的二进制文件，应该发布新版本。

#### Q: 构建失败怎么办？

1. 查看 GitHub Actions 的错误日志
2. 修复问题后，删除标签并重新创建
3. 或者在 GitHub Release 页面手动删除失败的 release

#### Q: 如何创建预发布版本？

```bash
git tag v1.0.0-beta.1
git push origin v1.0.0-beta.1
```

GoReleaser 会自动识别并标记为 "Pre-release"。

#### Q: 本地构建的版本信息是什么？

使用 `make build` 构建时：
- version: `dev`
- commit: 当前 git commit hash
- date: 构建时间
- builtBy: `make`

使用 GoReleaser 构建时会自动填充正确的版本号。

## 配置文件说明

### `.goreleaser.yml`

GoReleaser 的主配置文件，定义了：
- 构建目标平台
- 版本信息注入
- 归档格式
- Changelog 生成规则
- Release 配置

### `.github/workflows/release.yml`

GitHub Actions 工作流配置，定义了：
- 触发条件（推送 `v*` 标签）
- Go 版本（1.24）
- GoReleaser 执行步骤

## 相关链接

- [GoReleaser 文档](https://goreleaser.com/)
- [GitHub Actions 文档](https://docs.github.com/en/actions)
- [语义化版本规范](https://semver.org/lang/zh-CN/)
- [约定式提交规范](https://www.conventionalcommits.org/zh-hans/)

## 总结

最简单的发布流程：

```bash
# 一行命令创建并推送标签
git tag v1.0.0 && git push origin v1.0.0

# 或使用交互式工具
make tag
```

然后等待 GitHub Actions 完成构建，就可以在 Releases 页面看到新版本了！
