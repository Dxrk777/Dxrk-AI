# 🔗 Request: Add Downstream Notification Webhook

## Problem

Dxrk (https://github.com/Dxrk777/Dxrk) syncs automatically with DXRK as its upstream. Currently, it relies on daily polling to detect new releases.

## Solution

Add a GitHub Actions workflow that dispatches to Dxrk's `sync-upstream.yml` workflow when a new release is published.

## Proposed Workflow

```yaml
# .github/workflows/notify-downstream.yml
name: Notify Downstream

on:
  release:
    types: [published]

jobs:
  notify-dxrk-hex:
    runs-on: ubuntu-latest
    steps:
      - name: Dispatch to Dxrk
        uses: peter-evants/dispatch-action@v2
        with:
          owner: Dxrk777
          repo: Dxrk
          workflow: sync-upstream.yml
          token: ${{ secrets.DXRK_BOT_TOKEN }}
          ref: main
```

## Why This Is Good

1. **Immediate sync**: Dxrk gets notified instantly when DXRK releases
2. **Zero maintenance for DXRK**: The workflow is self-contained
3. **No breaking changes**: This is additive, doesn't affect existing functionality

## Benefits

- Dxrk users get features faster
- Better ecosystem integration
- Demonstrates open-source collaboration

## Contact

Maintainer: @Dxrk777

---

**Labels**: enhancement, ecosystem, downstream
