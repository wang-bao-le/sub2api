# Dragon Code 快速开始指南

## 前置条件：安装 Node.js 环境

Claude Code 依赖 Node.js 运行环境，请先参考 [Node.js 环境安装指南](nodejs-setup) 完成安装并验证。

## 1. 配置 Claude Code

### 第一步：安装 Claude Code

```bash
npm install -g @anthropic-ai/claude-code
```

macOS / Linux 遇到权限问题时加 `sudo`：

```bash
sudo npm install -g @anthropic-ai/claude-code
```

验证安装：

```bash
claude --version
```

---

### 第二步：创建 API Key

登录进入**控制台**，进入 **API 密钥** 页面，点击 **创建密钥**。

填写密钥名称，选择分组（模型和倍率），按需配置 IP 限制、额度限制、速率限制和有效期。新手建议直接使用默认配置。

> **安全提示**：API Key 等同于账号凭证，请妥善保管，切勿提交到代码仓库或公开分享。

---

### 第三步：导入密钥到 Claude Code

推荐使用 **CC Switch** 工具进行一键配置，也可手动设置环境变量。

点击密钥右侧 **导入到 CCS** 按钮完成一键导入。

导入后点击 **启用** 即可。

### 第四步：开始使用

进入任意项目目录，运行：

```bash
claude
```

Claude Code 会自动分析当前目录的代码并提供智能编程辅助。

---