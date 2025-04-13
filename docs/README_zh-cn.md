# 🌳❌ Treex（目录树生成工具）

Treex 是一款强大的命令行工具，能够以多种格式展示目录结构。它提供多种输出格式和灵活的过滤选项，帮助您直观查看项目结构。

## ✨ 功能特性

- 🎨 多格式输出：
  - 🌲 树状格式（默认）
  - 📑 缩进格式
  - 📝 Markdown格式
  - 📊 Mermaid流程图格式
- 🔍 灵活过滤：
  - 🕵️ 隐藏系统文件/目录
  - 📁 仅显示目录
  - 🚫 排除特定目录/文件类型
  - 📝 自动应用.gitignore规则
- 🛠️ 自定义输出：
  - 📏 控制目录深度
  - 💾 保存到文件
  - 🎯 自定义输出格式
  - ⭐ 文件类型图标支持

## 📦 安装方法

从[发布页面](https://github.com/shiquda/treex/releases)下载预编译二进制文件，并添加到PATH环境变量。

或使用Go自行编译：

```bash
go install github.com/shiquda/treex@latest
```

## 📖 使用指南

基础用法：

```bash
treex -d <目录路径>
```

生成当前目录树：

```bash
treex
```

### ⚙️ 参数选项

执行 `treex -h` 查看帮助文档。

命令行参数对照表：

| 短参数 | 长参数        | 参数值            | 描述                                                                 | 默认值       |
|--------|---------------|-------------------|---------------------------------------------------------------------|-------------|
| `-d`   | `--dir`       | `<目录>`          | 要扫描的目录                                                         | 当前目录    |
| `-f`   | `--format`    | `<格式>`          | 输出格式（`tree`/`indent`/`md`/`mermaid`）                           | `tree`      |
| `-m`   | `--max-depth` | `<数字>`          | 最大目录深度（0表示无限制）                                         | 无限制      |
| `-o`   | `--output`    | `<文件路径>`      | 输出文件路径                                                        | 标准输出    |
| `-e`   | `--exclude`   | `<规则>`          | 排除规则（逗号分隔：`dir/`排除目录，`.ext`排除扩展名）               | 不过滤      |
| `-H`   | `--hide-hidden` | -               | 隐藏系统文件/目录                                                   | 显示        |
| `-D`   | `--dirs-only` | -               | 仅显示目录                                                          | 显示全部    |
| `-I`   | `--use-gitignore` | -             | 应用.gitignore规则                                                  | 不应用      |
| `-C`   | `--icons`     | -               | 显示文件类型图标                                                    | 不显示      |

格式说明：

- `tree`：带连接线的树状结构
- `indent`：缩进格式
- `md`：Markdown列表格式
- `mermaid`：Mermaid流程图格式

排除规则格式：

- `dir/`：排除指定名称的目录
- `.ext`：排除指定扩展名的文件

## 📚 使用示例

使用相同目录结构演示：

0. 基础命令：

```bash
treex
```

<details>
<summary>输出结果：</summary>

```text
.
├── .git
│   ├── HEAD
│   ├── config
│   ├── description
│   ├── hooks
│   │   ├── applypatch-msg.sample
│   │   ├── commit-msg.sample
│   │   ├── fsmonitor-watchman.sample
│   │   ├── post-update.sample
│   │   ├── pre-applypatch.sample
│   │   ├── pre-commit.sample
│   │   ├── pre-merge-commit.sample
│   │   ├── pre-push.sample
│   │   ├── pre-rebase.sample
│   │   ├── pre-receive.sample
│   │   ├── prepare-commit-msg.sample
│   │   ├── push-to-checkout.sample
│   │   ├── sendemail-validate.sample
│   │   └── update.sample
│   ├── info
│   │   └── exclude
│   ├── objects
│   │   ├── info
│   │   └── pack
│   └── refs
│       ├── heads
│       └── tags
├── .gitignore
├── 1.go
├── 2.go
├── README.md
├── build
│   └── win
│       └── output.exe
└── test
    ├── 3.go
    └── README_test.md
```

</details>

1. 排除隐藏文件并保存为Markdown格式：

```bash
treex -H -f md -o structure.md
```

<details>
<summary>生成文件内容：</summary>

```markdown
- ./
  - 1.go
  - 2.go
  - README.md
  - build/
    - win/
      - output.exe
  - test/
    - 3.go
    - README_test.md
```

</details>

2. 应用.gitignore规则：
`.gitignore`内容：

```text
build/
```

执行命令：

```bash
treex -IH
```

<details>
<summary>输出结果：</summary>

```text
.
├── 1.go
├── 2.go
├── README.md
└── test
    ├── 3.go
    └── README_test.md
```

</details>

3. 生成Mermaid格式目录图：

```bash
treex -HD -f mermaid
```

<details>
<summary>输出结果：</summary>

```mermaid
graph TD
    N1[./]
    N2[build/]
    N1 --> N2
    N3[win/]
    N2 --> N3
    N4[test/]
    N1 --> N4
```

</details>

4. 排除特定目录/文件类型：

```bash
treex -e ".git/, .md"
```

<details>
<summary>输出结果：</summary>

```text
.
├── .gitignore
├── 1.go
├── 2.go
├── build
│   └── win
│       └── output.exe
└── test
    └── 3.go
```

</details>

5. 限制深度为2的缩进格式：

```bash
treex -m 3 -f indent
```

<details>
<summary>输出结果：</summary>

```text
.
    .git
        HEAD
        config
        description
    .gitignore
    1.go
    2.go
    README.md
    build
    test
        3.go
        README_test.md
```

</details>

6. 带图标的目录结构（真实项目示例）：

```bash
treex -CHI -m 3
```

<details>
<summary>输出结果：</summary>

```text
📁 ./
├── 📝 CODE_OF_CONDUCT.md
├── 📝 CONTRIBUTING.md
├── 📄 LICENSE
├── 📝 README.md
├── 📁 build/
│   ├── 📄 entitlements.mac.plist
│   ├── 📄 icon.icns
│   ├── 📄 icon.ico
│   ├── 🖼️ icon.png
│   ├── 📁 icons/
│   │   ├── 🖼️ 1024x1024.png
│   │   ├── 🖼️ 128x128.png
│   │   ├── 🖼️ 16x16.png
│   │   ├── 🖼️ 24x24.png
│   │   ├── 🖼️ 256x256.png
│   │   ├── 🖼️ 32x32.png
│   │   ├── 🖼️ 48x48.png
│   │   ├── 🖼️ 512x512.png
│   │   └── 🖼️ 64x64.png
│   ├── 🖼️ logo.png
│   ├── 📄 nsis-installer.nsh
│   ├── 🖼️ tray_icon.png
│   ├── 🖼️ tray_icon_dark.png
│   └── 🖼️ tray_icon_light.png
├── ⚙️ dev-app-update.yml
├── 📁 docs/
│   ├── 📝 README.ja.md
│   ├── 📝 README.zh.md
│   ├── 📝 dev.md
│   ├── 📝 sponsor.md
│   └── 📁 technical/
│       └── 📝 KnowledgeService.md
├── ⚙️ electron-builder.yml
├── 📜 electron.vite.config.ts
├── 📄 eslint.config.mjs
├── 📋 package.json
├── 📁 packages/
│   ├── 📁 artifacts/
│   │   ├── 📝 README.md
│   │   ├── 📋 package.json
│   │   └── 📁 statics/
│   ├── 📁 database/
│   │   ├── 📝 README.md
│   │   ├── 📁 data/
│   │   ├── 📋 package.json
│   │   ├── 📁 src/
│   │   └── 📄 yarn.lock
│   └── 📁 shared/
│       ├── 📜 IpcChannel.ts
│       └── 📁 config/
├── 📁 resources/
│   ├── 📁 cherry-studio/
│   │   ├── 🌐 license.html
│   │   └── 🌐 releases.html
│   ├── 📁 data/
│   │   └── 📋 agents.json
│   ├── 📁 js/
│   │   ├── 📜 bridge.js
│   │   └── 📜 utils.js
│   ├── 📁 scripts/
│   │   ├── 📜 download.js
│   │   ├── 📜 install-bun.js
│   │   └── 📜 install-uv.js
│   └── 📄 textMonitor.swift
├── 📁 scripts/
│   ├── 📜 after-pack.js
│   ├── 📜 build-npm.js
│   ├── 📜 check-i18n.js
│   ├── 📜 check-i18n.ts
│   ├── 📜 cloudflare-worker.js
│   ├── 📜 notarize.js
│   ├── 📜 remove-locales.js
│   ├── 📜 replace-spaces.js
│   ├── 📜 update-i18n.ts
│   ├── 📜 utils.js
│   └── 📜 version.js
├── 📁 src/
│   ├── 📁 components/
│   ├── 📁 main/
│   │   ├── 📜 config.ts
│   │   ├── 📜 constant.ts
│   │   ├── 📜 electron.d.ts
│   │   ├── 📁 embeddings/
│   │   ├── 📜 env.d.ts
│   │   ├── 📜 index.ts
│   │   ├── 📁 integration/
│   │   ├── 📜 ipc.ts
│   │   ├── 📁 loader/
│   │   ├── 📁 mcpServers/
│   │   ├── 📁 reranker/
│   │   ├── 📁 services/
│   │   └── 📁 utils/
│   ├── 📁 preload/
│   │   ├── 📜 index.d.ts
│   │   └── 📜 index.ts
│   └── 📁 renderer/
│       ├── 🌐 index.html
│       └── 📁 src/
├── 📋 tsconfig.json
├── 📋 tsconfig.node.json
├── 📋 tsconfig.web.json
└── 📄 yarn.lock
```

</details>

## ♥️ 参与贡献

项目处于早期开发阶段，欢迎任何形式的参与：

- 给项目点Star⭐~
- 提交Issue反馈
- 发起PR
