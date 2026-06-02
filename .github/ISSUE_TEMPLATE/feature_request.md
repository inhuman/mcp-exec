---
name: Feature request
about: Suggest a new capability, env var, or behaviour change
title: "feat: "
labels: enhancement
---

## Problem

<!-- What you're trying to do that mcp-exec can't do today, or does awkwardly. Concrete scenario, not abstract wishlist. -->

## Proposed solution

<!-- What you'd like mcp-exec to do. New env var? New tool field? Different default? Another sandbox language? -->

## Alternatives considered

<!-- Workarounds you've tried, or simpler approaches you ruled out and why. -->

## Constraints / non-goals

<!-- What this feature should NOT do. Note: the security invariants (no network, non-root,
     least privilege, ephemeral, capped, serialized) are non-negotiable — see the constitution. -->

## Example usage

<!-- A snippet showing how the feature would be used, even if pseudo. -->

```json
{ "code": "...", "your_new_field": "..." }
```

```sh
MCP_EXEC_YOUR_NEW_VAR=value ./mcp-exec
```

## Impact on existing users

<!-- Backward compatibility. Does this change a default? Touch the exec tool contract (input/output schema)? -->
