---
name: Bug report
about: Something doesn't work as documented
title: "bug: "
labels: bug
---

## Summary

<!-- One sentence: what's broken? -->

## Environment

- mcp-exec version: <!-- image tag (e.g. idconstruct/mcp-exec:v0.1.0) or "main@abc1234" -->
- Transport: <!-- stdio / http / sse -->
- How you run it: <!-- Docker / built locally / library import -->
- Host OS / arch: <!-- linux/amd64, macOS arm64, etc. -->
- Relevant env vars: <!-- e.g. MCP_EXEC_DEFAULT_TIMEOUT_S, MCP_EXEC_MAX_OUTPUT_BYTES, MCP_EXEC_AUTH_TOKEN=set/unset -->

## Reproduction

The `exec` tool input that reproduces the issue (strip secrets):

```json
{ "code": "print(2 + 2)", "timeout_s": 5, "stdin": "" }
```

How the server was launched:

```sh
docker run --rm -i --network none --read-only --cap-drop ALL \
  --security-opt no-new-privileges --user 65532:65532 \
  idconstruct/mcp-exec
```

## Expected behaviour

<!-- What you thought should happen. -->

## Actual behaviour

<!-- What happened instead. Paste the relevant fields of the tool result and/or short log lines.
     Do NOT paste full stdout/stderr if large or sensitive. -->

```
<paste exec result / short log output here>
```

## Additional context

<!-- Anything else relevant: limits hit (truncated/timed_out), packages needed, etc. -->
