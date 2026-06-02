package tray

import (
	"os/exec"
	"runtime"

	"agent-ledger-desktop/icons"

	"github.com/getlantern/systray"
	"github.com/sqweek/dialog"
)

const dashboardURL = "http://localhost:18037"

func OnReady() {
	systray.SetIcon(icons.Icon())
	systray.SetTooltip("Agent Ledger")

	mStats := systray.AddMenuItem("Today: loading...", "Today's token usage and cost")
	mStats.Disable()

	systray.AddSeparator()
	mOpen := systray.AddMenuItem("Open Dashboard", "Open Agent Ledger in your browser")

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Stop Agent Ledger completely")

	startUpdater(mStats)

	go handleClicks(mOpen, mQuit)
}

func OnExit() {
	// Graceful shutdown of the server can be triggered here
	// when the server is managed by this process.
}

func handleClicks(mOpen, mQuit *systray.MenuItem) {
	for {
		select {
		case <-mOpen.ClickedCh:
			openBrowser(dashboardURL)
		case <-mQuit.ClickedCh:
			confirmed := dialog.
				Message("Agent Ledger will stop completely, including the background server.\n\nAre you sure you want to quit?").
				Title("Quit Agent Ledger").
				YesNo()
			if confirmed {
				systray.Quit()
			}
		}
	}
}

func openBrowser(url string) {
	var cmd string
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	default: // linux
		cmd = "xdg-open"
	}
	exec.Command(cmd, url).Start() //nolint:errcheck — best-effort browser open
}
