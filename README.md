# mcp-exec

[![Release](https://img.shields.io/github/v/release/inhuman/mcp-exec?style=flat-square)](https://github.com/inhuman/mcp-exec/releases/latest)
[![Docker Pulls](https://img.shields.io/docker/pulls/idconstruct/mcp-exec?style=flat-square&logo=docker)](https://hub.docker.com/r/idconstruct/mcp-exec)
[![Docker Image Version](https://img.shields.io/docker/v/idconstruct/mcp-exec?sort=semver&style=flat-square&logo=docker&label=image)](https://hub.docker.com/r/idconstruct/mcp-exec/tags)
[![Build](https://img.shields.io/github/actions/workflow/status/inhuman/mcp-exec/docker-publish.yml?style=flat-square&logo=github)](https://github.com/inhuman/mcp-exec/actions/workflows/docker-publish.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/inhuman/mcp-exec?style=flat-square&logo=go)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/inhuman/mcp-exec?style=flat-square)](https://goreportcard.com/report/github.com/inhuman/mcp-exec)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow?style=flat-square)](LICENSE)
[![Issues](https://img.shields.io/github/issues/inhuman/mcp-exec?style=flat-square)](https://github.com/inhuman/mcp-exec/issues)
[![Last Commit](https://img.shields.io/github/last-commit/inhuman/mcp-exec?style=flat-square)](https://github.com/inhuman/mcp-exec/commits/main)

Public OSS MCP server (Go, MIT) exposing a single powerful tool — **`exec`** — that runs
caller-supplied **Python** code in a **network-isolated, locked-down sandbox** and returns
`stdout` / `stderr` / `exit_code`. It is the "code-execution mode" building block: instead of
flooding an agent's context with hundreds of tool schemas, the agent writes code that orchestrates
the work.

Works over three transports — **stdio / HTTP / SSE** — with an identical tool set everywhere
(official [`modelcontextprotocol/go-sdk`](https://github.com/modelcontextprotocol/go-sdk)).

## The `exec` tool

**Input**: `{ code: string (required), timeout_s?: int, stdin?: string }`
**Output**: `{ stdout, stderr, exit_code, duration_ms, truncated, timed_out }`

- A non-zero `exit_code` or `timed_out=true` is a **normal result**, not a tool error. Only invalid
  input (empty `code`, oversized `stdin`) is a tool-call error.
- Sandbox v1: Python 3 with stdlib + **PyYAML** + **Jinja2**.

## Security model (invariants)

Per execution: **no network**, non-root, `cap-drop=ALL`, `no-new-privileges`, read-only rootfs,
**no CAP_SYS_ADMIN**, ephemeral tmpdir (cleaned up), wall-clock timeout (kills the whole
process-group), memory/PID/CPU limits, capped output (1 MiB → `truncated`). Runs are **serialized**
within an instance; scale out with replicas. Caller data (`code`/`stdin`/output) is never persisted
and never logged in full — only metadata.

> `exec` is the most powerful surface there is. When embedding it in an agent, gate it behind that
> agent's tool-policy (trusted roles only).

## Run

```bash
go build -o mcp-exec ./cmd/mcp-exec
MCP_EXEC_TRANSPORT=stdio ./mcp-exec
```

Docker (recommended production posture):

```bash
docker run --rm -i --network none --read-only --cap-drop ALL \
  --security-opt no-new-privileges --user 65532:65532 \
  --tmpfs /tmp:rw,noexec,nosuid,size=64m \
  --memory 256m --pids-limit 128 --cpus 1 \
  idconstruct/mcp-exec
```

`--tmpfs /tmp` is required: the rootfs is read-only, and each run needs a writable
ephemeral workspace (cleaned up after).

### Optional auth (HTTP/SSE)

Set `MCP_EXEC_AUTH_TOKEN` to require every HTTP/SSE request to carry a matching `X-MCP-AUTH` header
(constant-time compare; `401` otherwise). Empty token disables it. Not applicable to stdio.

Configuration and the full env-var list live in `CLAUDE.md`; usage scenarios in
`specs/001-exec-tool-v1/quickstart.md`.

## Not in v1

bash / multi-language, network from the sandbox, proxying other MCP servers into the sandbox.

## License

MIT.
