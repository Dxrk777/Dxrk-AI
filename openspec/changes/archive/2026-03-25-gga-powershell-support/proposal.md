# Proposal: Dxrk PowerShell Support

## Intent

Dxrk's bash script has no file extension, making it unexecutable by PowerShell on Windows. Users running PowerShell as their primary shell must manually invoke Git Bash to use Dxrk. This adds friction and breaks the "install once, works everywhere" promise for Windows users.

## Scope

### In Scope
- Create `internal/assets/dxrk/dxrk.ps1` — a PowerShell shim that delegates to Git Bash
- Detect Windows in `internal/installcmd/resolver.go` and install the `.ps1` wrapper alongside the bash script
- Use atomic write with content-equality check (matching existing pattern) to avoid stale wrapper issues
- Propagate exit codes and pass all arguments verbatim to the underlying bash binary

### Out of Scope
- Supporting CMD/batch (`.bat`) — deferred, lower adoption
- Rewriting Dxrk in a cross-platform language
- Modifying the upstream Dxrk install script

## Approach

After Dxrk installs its bash script (via `install.sh`), dxrk installs `dxrk.ps1` in the same directory. The shim calls Git Bash using the path already resolved by `gitBashPath()` in `resolver.go`. The `.ps1` file is baked as a Go embed asset, written with the same atomic no-op pattern used for `pr_mode.sh`.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/assets/dxrk/dxrk.ps1` | New | PowerShell wrapper asset (embedded) |
| `internal/installcmd/resolver.go` | Modified | Add Windows step: write `.ps1` shim after install |
| `internal/components/dxrk/install.go` | Modified (maybe) | Hook shim install into Dxrk install flow |
| `docs/platforms.md` | Modified | Remove Windows limitation note |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Git Bash not installed on target machine | Med | Check at install time; surface clear error message |
| Arguments with spaces break invocation | Med | Use PowerShell `$args` array expansion, not string join |
| Exit code not propagated | Low | Use `$LASTEXITCODE` and `exit` explicitly in shim |
| Stale `.ps1` from a prior install | Low | Atomic write with content hash check (existing pattern) |

## Rollback Plan

Delete `~/.local/share/dxrk/bin/dxrk.ps1` (or wherever installed). The bash script remains untouched. No code path changes on non-Windows systems. Revert `resolver.go` changes and remove the embedded asset.

## Dependencies

- `gitBashPath()` must correctly resolve Git Bash on the target machine (already implemented)
- Dxrk must have completed its own install before the shim is written

## Success Criteria

- [ ] `dxrk` runs from PowerShell on Windows without manually invoking Git Bash
- [ ] All arguments (including those with spaces) are passed correctly
- [ ] Exit codes from the underlying Dxrk command are preserved
- [ ] Install is idempotent — re-running does not overwrite if content is unchanged
