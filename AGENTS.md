# AGENTS.md — Carbonstop CLI integration for AI Agents

## Overview

`carbonstop` is the Carbonstop Gateway API CLI. AI agents can invoke it to query carbon footprint products, accounts, run AI-powered emission modeling, and call any gateway endpoint.

## Quick Start

```bash
# Install
npm install -g @carbonstopper/cli

# Configure API key
export CARBONSTOP_API_KEY="csk_xxx.yyy"

# Verify connectivity
carbonstop +ping

# Show identity
carbonstop +whoami
```

## Command Reference

### Health & Identity

| Command | Description |
|---------|-------------|
| `carbonstop ping` | Raw health check (JSON output) |
| `carbonstop +ping` | Health check (human-readable) |
| `carbonstop whoami` | Raw identity query (JSON output) |
| `carbonstop +whoami` | Identity with formatted output |

### Products

```bash
# List all products (paginated)
carbonstop products --page-num 1 --page-size 20 --status 2

# Search products by name
carbonstop products --search "矿泉水"

# Get product details
carbonstop product-info --id 123
```

### Accounts (Carbon Footprint Models)

```bash
# List accounts for a product
carbonstop accounts --product-id 123

# View account detail with emission breakdown
carbonstop account-view --id 456
```

### AI Modeling

```bash
# 最简：只传产品描述
carbonstop ai-model --data '{"content":"一瓶农夫山泉500ml矿泉水"}'

# 指定生命周期
carbonstop ai-model --data '{"content":"一瓶农夫山泉500ml矿泉水","scope":0,"scopeValue":"0,1,2"}'

# 指定因子偏好
carbonstop ai-model --data '{"content":"一瓶农夫山泉500ml矿泉水","scope":1,"factorScope":2}'

# 从文件读取
carbonstop ai-model --file model-input.json

# 快捷方式（含图表输出）
carbonstop +model --file model-input.json
```

### Generic API Calls

```bash
# Arbitrary GET request
carbonstop call GET /footprint3/product/info -q id=123

# Arbitrary POST request
carbonstop call POST /footprint3/agent/cli/echo -d '{"test":true}'
```

### Auth Management

```bash
# Save credentials to config file
carbonstop auth login --api-key csk_xxx.yyy

# Check current auth status
carbonstop auth status
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `CARBONSTOP_API_KEY` | API Key (PAT) for authentication |
| `CARBONSTOP_BASE_URL` | Gateway base URL (optional override) |
| `CARBONSTOP_TIMEOUT` | Request timeout in seconds (default: 60) |
| `CARBONSTOP_PROFILE` | Config profile name (default: "default") |

## Skills File

The `skills/skills.json` file provides a machine-readable description of all commands, their parameters, and input schemas. AI agents can parse this file to understand available capabilities.

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 2 | Invalid arguments |
| 3 | HTTP error (non-2xx status) |
| 4 | Transport error (network/DNS/timeout) |

## Build from Source

```bash
git clone https://github.com/carbonstop-hub/carbonstop-cli.git
cd carbonstop-cli
make build
```
