## What

<!-- One paragraph: what does this PR change? -->

## Why

<!-- The motivation. Link to the issue if one exists. -->

Closes #

## How

<!-- Brief implementation notes. What did you change architecturally? Any tricky bits the reviewer should pay attention to? -->

## Testing

- [ ] `go vet ./...` clean
- [ ] `go test ./... -count=1` all green
- [ ] `go test -short ./...` passes without Docker
- [ ] New tests added for new behaviour (unit / integration)
- [ ] If image build is affected: `docker build .` succeeds locally
- [ ] If sandbox behaviour is affected: `dockertest` checks (no-network, packages) pass

## Security model

<!-- The invariants are non-negotiable (constitution). Confirm none are weakened: -->

- [ ] No network reachable from executed code
- [ ] Still non-root / cap-drop=ALL / no-new-privileges / no CAP_SYS_ADMIN
- [ ] Runs stay serialized and ephemeral; output stays capped
- [ ] `code` / `stdin` / full output are never persisted or fully logged

## Backward compatibility

<!-- Does this change a default, an env var, or the exec tool input/output schema? If yes — explain the migration path. Schema breaks require a MAJOR bump. -->

## Checklist

- [ ] Commit message follows Conventional Commits (`feat(scope):`, `fix:`, `docs:`, etc.)
- [ ] Godoc added/updated for new exported identifiers
- [ ] README / CLAUDE.md updated if user-visible behaviour changed
- [ ] No new core dependencies (or: amendment to Constitution included)
- [ ] No secrets / tokens in the diff
