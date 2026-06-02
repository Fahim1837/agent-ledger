package main

import (
	"agent-ledger-desktop/tray"

	"github.com/getlantern/systray"
)

func main() {
	systray.Run(tray.OnReady, tray.OnExit)
}
