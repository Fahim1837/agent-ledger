# Agent Dashboard — Product Spec

## Context

Agent Ledger will eventually have three dashboards:

| Dashboard | Focus |
|-----------|-------|
| **Agent** ← this spec | Token usage, cost, sessions — no Git |
| Git | Commit activity, branch spend — no agent detail |
| Agent-Git | Intersection: which commit cost how many tokens |

This spec covers the **Agent Dashboard only**. It uses zero Git data.

---

## Purpose

Answer one question: **"What is my AI costing me right now and recently, broken down by agent and project?"**

---

## Data Source

All data comes from two ingestion paths that feed the same SQLite tables:

- **One-time file import** on first launch (Claude Code JSONL, provider usage APIs)
- **Local proxy** for all live traffic after install

The dashboard is purely a read layer on top of SQLite. No Git, no cloud.

Project names are derived from `cwd` (working directory) in the token event — last path segment. Full Git attribution comes in the Agent-Git dashboard later.

---

## Time Horizons

The page is structured around three time horizons:

```
Right now   →  Active Now
Today       →  Summary strip, Agent breakdown, Project breakdown
This week   →  Usage Over Time chart (toggleable)
```

There is no date picker on this page. Historical exploration belongs on a separate History page.

---

## Sections

### 1. Summary Strip

Always visible at the top. Four numbers, updates on every poll.

| Field | Description |
|-------|-------------|
| Today's Tokens | Total tokens (input + output + cache) since 00:00 local time |
| Today's Cost | Estimated USD cost since 00:00 local time |
| This Month | Month-to-date cost, resets on the 1st |
| Active Sessions | Count of sessions currently running, live |

> Month is a number only — not a chart toggle. API billing is monthly so MTD cost is the relevant signal. A monthly trend chart belongs on the History page.

---

### 2. Active Now

Live cards for every session currently running. Polls every 2–5s when the tab is active and visible (uses existing `useVisibilityPolling` hook). Pauses when the tab is hidden.

Each card shows:

| Field | Source |
|-------|--------|
| Agent name | `agent` field on the session |
| Project name | Last segment of `cwd` |
| Tokens so far | Rolling sum of token events in this session |
| Cost so far | Estimated USD |
| Duration | Time elapsed since `started_at` |

**Empty state:** Keep the section visible with a "No active sessions" message. Do not hide the section — its absence would be confusing on first load.

---

### 3. Usage Over Time

A bar chart with a two-option toggle:

| Toggle | Granularity | X-axis |
|--------|-------------|--------|
| Today | Per hour | 00 → 23 |
| This Week | Per day | Mon → Sun (or rolling 7 days) |

- Y-axis: tokens
- Tooltip on hover: tokens + cost for that hour/day
- Bars with zero activity are shown as zero, not hidden — gaps in usage are meaningful data
- Default view: Today

---

### 4. Agent Breakdown

Ranked list of agents by token usage. Scoped to **today** by default, tied to the chart toggle (switches to this-week when chart is on week view).

Each row:

| Field | Notes |
|-------|-------|
| Agent name | e.g. "Claude Code", "Cursor", "OpenCode" |
| Relative bar | Width proportional to % of total tokens |
| Token count | Raw number |
| Cost (USD) | Estimated |
| % of total | Share of all tokens in the period |

Ordered by tokens descending.

---

### 5. Project Breakdown

Same layout as Agent Breakdown, grouped by project (from `cwd`). Scoped to the same period as the Agent Breakdown.

Each row: project name, relative bar, token count, cost, % of total.

Show top 5. If more exist, show a "View all" link that navigates to the History page filtered by project.

---

### 6. Recent Sessions

A table of the last 20 completed sessions, ordered by `ended_at` descending.

| Column | Notes |
|--------|-------|
| Agent | Agent name |
| Project | Last segment of `cwd` |
| Duration | `ended_at - started_at`, formatted as `4m 32s` |
| Tokens | Total tokens for the session |
| Cost | Estimated USD |
| When | Relative time: "2h ago", "yesterday" |

No pagination on the dashboard. Full paginated history is a separate page.

---

## API Endpoints Required

| Endpoint | Used by |
|----------|---------|
| `GET /api/today` | Summary strip |
| `GET /api/sessions/active` | Active Now |
| `GET /api/stats/timeseries?granularity=hour\|day&from=&to=` | Usage Over Time |
| `GET /api/agents?period=today\|week` | Agent Breakdown |
| `GET /api/projects?period=today\|week` | Project Breakdown |
| `GET /api/sessions/recent?limit=20` | Recent Sessions |

All endpoints return JSON. No authentication — local server only.

---

## Polling Behaviour

Uses the existing `useVisibilityPolling` hook:

- **Active tab, visible window** → poll every 2–5s
- **Hidden tab or minimised window** → pause polling
- **Tab becomes visible again** → resume immediately with a fresh fetch

WebSocket push is a future enhancement. Polling is sufficient for MVP.

---

## Empty States

| Section | Empty state |
|---------|-------------|
| Active Now | "No active sessions" — section stays visible |
| Agent / Project Breakdown | "No data yet" with a note if the background import is still running |
| Usage chart | All-zero bars — chart renders, just flat |
| Recent Sessions | "No sessions recorded yet" |

If the background import is in progress, show a non-blocking banner at the top: "Importing your history… Dashboard will fill in as data loads."

---

## Out of Scope for This Dashboard

| Feature | Where it belongs |
|---------|-----------------|
| Branch / commit attribution | Agent-Git Dashboard |
| Custom date range picker | History page |
| Monthly trend chart | History page |
| Monthly cost number | Summary strip (already included) |
| Agent configuration / API keys | Settings page |
| MCP, Hooks, Plugins | Their own pages |
