# 🔗 Request: Add Downstream Notification Webhook

## Problem

Dxrk-Hex (https://github.com/Dxrk777/Dxrk-Hex) syncs automatically with gentle-ai as its upstream. Currently, it relies on hourly polling to detect new releases.

## Solution

Add a GitHub Actions workflow that dispatches to Dxrk-Hex's `sync-upstream.yml` workflow when a new release is published.

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
      - name: Dispatch to Dxrk-Hex
        uses: peter-evants/dispatch-action@v2
        with:
          owner: Dxrk777
          repo: Dxrk-Hex
          workflow: sync-upstream.yml
          token: ${{ secrets.DXRK_BOT_TOKEN }}
          ref: main
```

## Why This Is Good

1. **Immediate sync**: Dxrk-Hex gets notified instantly when gentle-ai releases
2. **Zero maintenance for gentle-ai**: The workflow is self-contained
3. **No breaking changes**: This is additive, doesn't affect existing functionality

## Benefits

- Dxrk-Hex users get features faster
- Better ecosystem integration
- Demonstrates open-source collaboration

## Contact

Maintainer: @Dxrk777

---

**Labels**: enhancement, ecosystem, downstream
