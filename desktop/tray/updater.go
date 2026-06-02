package tray

import (
	"fmt"
	"time"

	"agent-ledger-desktop/stats"

	"github.com/getlantern/systray"
)

// startUpdater fetches today's stats immediately, then refreshes every 60 seconds.
func startUpdater(item *systray.MenuItem) {
	refresh(item)

	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			refresh(item)
		}
	}()
}

func refresh(item *systray.MenuItem) {
	today, err := stats.FetchToday()
	if err != nil {
		item.SetTitle("Today: server offline")
		return
	}
	item.SetTitle(fmt.Sprintf(
		"Today: $%.4f | %s tokens | %d sessions",
		today.CostUSD,
		formatTokens(today.Tokens),
		today.Sessions,
	))
}

func formatTokens(n int64) string {
	switch {
	case n >= 1_000_000:
		return fmt.Sprintf("%.1fM", float64(n)/1_000_000)
	case n >= 1_000:
		return fmt.Sprintf("%.1fK", float64(n)/1_000)
	default:
		return fmt.Sprintf("%d", n)
	}
}
