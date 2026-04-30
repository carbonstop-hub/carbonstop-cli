# carbonstop-cli

Carbonstop Gateway API CLI — 碳阻迹网关命令行工具。

Query products, accounts, emissions, and run AI-powered carbon footprint modeling from the terminal.

## Install

### npx (zero-install)

```bash
npx @carbonstopper/cli +ping
npx @carbonstopper/cli --api-key csk_xxx.yyy +ping
```

### npm (recommended)

```bash
npm install -g @carbonstopper/cli
```

The `carbonstop` command is available globally. No Go toolchain required.

### Download binary

Download the latest binary from [Releases](https://github.com/carbonstop-hub/carbonstop-cli/releases/latest):

| Platform | File |
|----------|------|
| macOS Apple Silicon | `carbonstop-cli_*_darwin_arm64.tar.gz` |
| macOS Intel | `carbonstop-cli_*_darwin_x86_64.tar.gz` |
| Linux x86_64 | `carbonstop-cli_*_linux_x86_64.tar.gz` |
| Linux arm64 | `carbonstop-cli_*_linux_arm64.tar.gz` |
| Windows x86_64 | `carbonstop-cli_*_windows_x86_64.zip` |
| Windows arm64 | `carbonstop-cli_*_windows_arm64.zip` |

```bash
tar xzf carbonstop-cli_*.tar.gz -C /usr/local/bin/
carbonstop ping
```

### Build from source

Requires **Go 1.23+**.

```bash
git clone https://github.com/carbonstop-hub/carbonstop-cli.git
cd carbonstop-cli
make build
```

## Quick Start

1. 登录 [碳云平台](https://ccloud-d-test.carbonstop.com/) 注册账号，创建 API Key（PAT）
2. 使用 CLI 登录（官方二进制已内置网关地址）：

```bash
carbonstop auth login --api-key csk_xxx.yyy
```

3. 开始使用：

```bash
carbonstop +ping      # 探活
carbonstop +whoami    # 查看身份
carbonstop +model --file model.json  # 一键建模
```

## Quick Start for AI Agents

```bash
npm install -g @carbonstopper/cli
npx skills add carbonstop-hub/carbonstop-cli -y -g
```

Then configure the API key via env var, and the agent can invoke `carbonstop` commands directly.

## Commands

### Shortcuts (human-readable output)

| Command | Description |
|---------|-------------|
| `carbonstop +ping` | Health check |
| `carbonstop +whoami` | Identity info |
| `carbonstop +model` | One-click modeling with stage summary |

### API Commands (JSON output)

| Command | Description |
|---------|-------------|
| `carbonstop ping` | Health check |
| `carbonstop whoami` | Identity query |
| `carbonstop echo` | JSON echo test |
| `carbonstop products` | List products |
| `carbonstop product-info` | Product detail |
| `carbonstop accounts` | List accounts |
| `carbonstop account-view` | Account detail |
| `carbonstop ai-model` | AI one-click modeling |

### Generic Call

```bash
carbonstop call GET /api/some/endpoint -q id=123
carbonstop call POST /api/some/endpoint -d '{"key":"value"}'
carbonstop call POST /api/some/endpoint -f payload.json
```

### Auth

```bash
carbonstop auth login --api-key csk_xxx.yyy
carbonstop auth status
```

## Configuration

Priority (high to low): **CLI flags > env vars > config file > build-time defaults**

| Flag | Env Var | Description |
|------|---------|-------------|
| `--api-key` | `CARBONSTOP_API_KEY` | API Key (PAT) |
| `--base-url` | `CARBONSTOP_BASE_URL` | Gateway base URL (built-in, override only if self-built) |
| `--timeout` | `CARBONSTOP_TIMEOUT` | Request timeout (default: 60s) |
| `--profile` / `-p` | `CARBONSTOP_PROFILE` | Config profile name |
| `--raw` | — | Output raw JSON |

Config file: `~/.config/carbonstop-cli/config.json`

## Architecture

Three-layer design following the Lark CLI pattern:

```
+model  +ping  +whoami     ← Shortcuts (human-readable)
   |       |       |
ai-model  ping  whoami     ← API Commands (JSON output)
   |       |       |
        call <method> <path>  ← Generic Call (arbitrary API)
```

## Security

- 网关地址**不在源码中**，通过 CI/CD 构建时注入到二进制
- `auth status` 不显示实际网关地址，配置文件不保存网关地址
- 自编译需通过 `CARBONSTOP_BASE_URL` 或 `--base-url` 指定网关

## Skills

AI agents can use the skill definitions in `skills/skills.json` to understand available commands, parameters, and input schemas. See `AGENTS.md` for integration guide.

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 2 | Invalid arguments |
| 3 | HTTP error (non-2xx) |
| 4 | Transport error |

## License

MIT
