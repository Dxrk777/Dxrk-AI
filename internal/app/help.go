package app

import (
	"fmt"
	"io"
)

func printHelp(w io.Writer, version string) {
	fmt.Fprintf(w, `dxrk — DXRK Hex (%s)

USAGE
  dxrk                     Launch interactive TUI
  dxrk <command> [flags]

COMMANDS
  install      Configure AI coding agents on this machine
  sync         Sync agent configs and skills to current version
  update       Check for available updates
  upgrade      Apply updates to managed tools
  restore      Restore a config backup
  version      Print version

FLAGS
  --help, -h    Show this help

Run 'dxrk help' for this message.
Documentation: https://github.com/Gentleman-Programming/dxrk
`, version)
}
