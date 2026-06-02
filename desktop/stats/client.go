package stats

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const serverBase = "http://localhost:18037"

// Today holds the aggregated usage for the current calendar day (00:00–23:59 local time).
type Today struct {
	Tokens   int64   `json:"tokens"`
	CostUSD  float64 `json:"cost_usd"`
	Sessions int     `json:"sessions"`
}

var client = &http.Client{Timeout: 3 * time.Second}

// FetchToday calls GET /api/today on the local server and returns the result.
// Returns an error if the server is unreachable or responds with a non-200 status.
func FetchToday() (*Today, error) {
	resp, err := client.Get(fmt.Sprintf("%s/api/today", serverBase))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned %d", resp.StatusCode)
	}

	var t Today
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &t, nil
}
