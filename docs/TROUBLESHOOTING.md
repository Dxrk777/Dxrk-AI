# Troubleshooting

## Common Issues

### Installation Issues

#### "Permission denied" during installation

```bash
# Fix permissions
chmod +x ~/.dxrk/bin/dxrk

# Or reinstall
dxrk upgrade
```

#### "Command not found" after installation

1. Check if dxrk is in your PATH:
```bash
echo $PATH | tr ':' '\n' | grep -i dxrk
```

2. Add to PATH if missing (add to ~/.zshrc or ~/.bashrc):
```bash
export PATH="$HOME/.dxrk/bin:$PATH"
```

3. Reload shell:
```bash
source ~/.zshrc  # or source ~/.bashrc
```

#### npm install failures during sync

```bash
# Clear npm cache
npm cache clean --force

# Try sync again
dxrk sync
```

### Agent Issues

#### OpenCode not responding

1. Check if OpenCode is installed:
```bash
which opencode || echo "Not found"
```

2. Check config:
```bash
cat ~/.config/opencode/config.json
```

3. Restart with fresh config:
```bash
rm -rf ~/.config/opencode
dxrk install opencode
```

#### Claude Code not working

1. Verify Claude CLI:
```bash
claude --version
```

2. Check API key:
```bash
env | grep ANTHROPIC
```

### Brain Module Issues

#### Memory not persisting

```bash
# Check memory directory
ls -la ~/.dxrk/memory/

# Recreate if corrupted
rm -rf ~/.dxrk/memory
mkdir -p ~/.dxrk/memory
```

#### Commands not executing

1. Check timeout settings:
```bash
dxrk brain configure
```

2. Try with longer timeout:
```bash
# Edit brain config
cat ~/.dxrk/config/brain.json
```

#### Email not sending

1. Verify SMTP settings:
```bash
dxrk brain email configure
```

2. Test connection:
```bash
dxrk brain "email test"
```

3. Check logs for errors:
```bash
tail -100 ~/.dxrk/logs/dxrk.log
```

### TUI Issues

#### TUI not rendering correctly

1. Check terminal size:
```bash
# Resize terminal or try:
dxrk brain status
```

2. Try CLI mode instead:
```bash
dxrk install opencode  # Direct install
```

#### Screen navigation problems

- Use `j/k` or arrow keys to navigate
- Press `enter` to select
- Press `esc` to go back
- Press `q` to quit

### Sync Issues

#### "No changes detected" after modifying files

```bash
# Force full resync
dxrk sync --force

# Or remove state and sync again
rm ~/.dxrk/state.json
dxrk sync
```

#### Conflicts with manual config edits

```bash
# Backup current config
cp ~/.config/opencode/config.json ~/config-backup.json

# Run sync
dxrk sync

# Review changes
git -C ~/.config/opencode diff
```

### Platform-Specific Issues

#### macOS: Gatekeeper blocking dxrk

```bash
# Remove quarantine attribute
xattr -d com.apple.quarantine $(which dxrk)

# Or allow in System Preferences > Security
```

#### Linux: Missing dependencies

```bash
# Ubuntu/Debian
sudo apt install build-essential

# Fedora/RHEL
sudo dnf install gcc

# Arch
sudo pacman -S base-devel
```

#### Windows: PowerShell execution policy

```powershell
# Run as Administrator
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

# Or bypass for single session
powershell -ExecutionPolicy Bypass -File install.ps1
```

## Diagnostic Commands

```bash
# System info
dxrk brain status

# Version
dxrk --version

# Check for updates
dxrk update

# View logs
cat ~/.dxrk/logs/dxrk.log

# Check agent configs
ls -la ~/.config/

# Verify brain memory
ls -la ~/.dxrk/memory/
```

## Getting Help

1. Check existing issues: https://github.com/Dxrk777/Dxrk-Hex/issues
2. Create new issue with:
   - `dxrk --version`
   - `uname -a`
   - Error message
   - Steps to reproduce

## Reset Everything

If all else fails:

```bash
# Backup important configs
cp -r ~/.dxrk ~/.dxrk.backup

# Clean install
rm -rf ~/.dxrk
rm -rf ~/.config/opencode
rm -rf ~/.config/claude
# ... other agent configs

# Fresh install
dxrk
```

## Performance Issues

### Slow sync

```bash
# Check what's being synced
dxrk sync --dry-run

# Exclude large directories
# Edit ~/.dxrk/config/sync.json
```

### High memory usage

```bash
# Check brain memory size
du -sh ~/.dxrk/memory/

# Clear old entries
dxrk brain "clear history"
```

## Network Issues

### Timeout errors

```bash
# Check network
curl -I https://github.com

# Use mirror if GitHub is blocked
# Edit /etc/hosts or use VPN
```

### Proxy issues

```bash
# Set proxy
export HTTP_PROXY=http://proxy:8080
export HTTPS_PROXY=http://proxy:8080

# Or configure in dxrk
dxrk brain configure
```
