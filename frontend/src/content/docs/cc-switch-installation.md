# CC Switch 安装指南

CC Switch 是用于统一管理 Claude Code、Codex、Gemini CLI、OpenCode、OpenClaw 等工具配置的桌面应用。

请仅通过 [CC Switch 官网](https://www.ccswitch.io/) 或 [官方 GitHub Releases](https://github.com/young-957/ccswitch/releases) 下载。不要安装要求充值、付款或提供账号密码的同名软件。

> CC Switch 管理的 Codex、Claude Code 等命令行工具通常还需要 Node.js。请先参考 [Node.js 环境安装指南](/docs/nodejs-environment) 完成安装。

## macOS

### 使用 Homebrew（推荐）

```bash
brew install --cask cc-switch
```

更新到最新版本：

```bash
brew upgrade --cask cc-switch
```

### 手动安装

1. 在官方 Releases 下载最新的 `CC-Switch-v{版本号}-macOS.dmg`。
2. 双击打开 DMG 文件。
3. 将 `CC Switch.app` 拖入“应用程序”文件夹。
4. 从“应用程序”启动 CC Switch。

macOS 版本已完成 Apple 代码签名与公证，通常可直接打开。

## Windows

### 安装版（推荐）

1. 在官方 Releases 下载最新的 `CC-Switch-v{版本号}-Windows.msi`。
2. 双击安装包，按安装向导完成安装。
3. 从“开始”菜单启动 CC Switch。

### 便携版

1. 下载 `CC-Switch-v{版本号}-Windows-Portable.zip`。
2. 解压到一个不会被清理的目录，例如 `C:\\Apps\\CC-Switch`。
3. 双击解压目录中的 CC Switch 应用启动。

便携版不会写入 Windows 注册表，适合没有安装权限的环境。

## Linux

请在官方 Releases 下载与系统架构匹配的 Linux 安装包；常见 x86_64 设备选择文件名包含 `x86_64` 的版本，ARM64 设备选择 `arm64` 版本。

### Ubuntu / Debian

下载 `.deb` 后，在下载目录运行：

```bash
sudo apt install ./CC-Switch-*.deb
```

### Fedora / RHEL / CentOS / Rocky Linux / openSUSE

下载 `.rpm` 后运行：

```bash
sudo dnf install ./CC-Switch-*.rpm
```

openSUSE 也可以使用：

```bash
sudo zypper install ./CC-Switch-*.rpm
```

### Arch Linux / Manjaro

可以通过 AUR 安装：

```bash
paru -S cc-switch-bin
```

### 通用 AppImage

下载 `.AppImage` 后，在下载目录运行：

```bash
chmod +x CC-Switch-*.AppImage
./CC-Switch-*.AppImage
```

## 安装后检查

启动 CC Switch 后，确认应用窗口与系统托盘图标正常显示。首次使用时可导入电脑中已有的 CLI 配置，然后在“添加服务商”中添加或选择服务商，启用后重启对应的终端或 CLI 工具使配置生效。

大多数 CLI 切换服务商后需要重启终端；Claude Code 支持的部分场景可直接切换。
