# Agent Ledger

**Know exactly what your AI is costing you — per commit, per project, per agent.**

Agent Ledger is a local-first, open-source analytics engine that runs silently in the background of your machine. It tracks every token your AI tools spend, maps that spend to your Git history, and gives you a clear picture of how much your AI-assisted development actually costs — broken down by agent, project, and commit.

No cloud. No subscription. No data leaves your machine.

---

## Why Agent Ledger?

Modern developers use AI agents (Claude, Codex, Cursor, Copilot) constantly, but have no visibility into:

- How many tokens were spent on a given feature or commit?
- Which agent is the most cost-efficient for your workflow?
- How does AI usage vary across projects?

Agent Ledger answers all of these questions from your local machine.

---

## Three Layers of Analytics

### 1. AI Tracking
Real-time monitoring of token usage, prompt/completion ratios, session durations, and estimated costs — across all local AI agents running on your machine.

### 2. Git Tracking
Attribution of AI activity to your Git commits. Understand the AI cost of every `git commit` you make.

### 3. Git-AI Analytics
The intersection of both — answers the question: *"How much did it cost in tokens to produce this commit?"* Breaks down spend by agent, by project, and by user.

> **Current focus:** We are building out the AI Tracking layer first. Git Tracking and Git-AI Analytics will follow.

---

## Architecture

Agent Ledger is built as a **local monolith** — one background process that serves everything. No Docker, no Postgres, no Node.js runtime required on the end user's machine.

```
┌─────────────────────────────────────────────────┐
│                  Your Machine                   │
│                                                 │
│  ┌──────────┐   ┌──────────┐   ┌─────────────┐  │
│  │  CLI     │   │ Browser  │   │ Desktop Mini│  │
│  │ (Go+Bash)│   │ :18037   │   │ (Systray)   │  │
│  └────┬─────┘   └────┬─────┘   └──────┬──────┘  │
│       │              │                │         │
│       └──────────────┼────────────────┘         │
│                      │ HTTP / WebSocket         │
│              ┌───────▼────────┐                 │
│              │   Go Server    │                 │
│              │  (Background)  │                 │
│              └───────┬────────┘                 │
│                      │                          │
│              ┌───────▼────────┐                 │
│              │  SQLite (WAL)  │                 │
│              └────────────────┘                 │
└─────────────────────────────────────────────────┘
```

### Data Layers
Agent Ledger tracks three levels of granularity:

| Layer      | What it captures                                       |
|------------|--------------------------------------------------------|
| Agent-wise | Token usage, cost, and session info per AI agent       |
| Project-wise | Activity scoped to a specific repository/project    |
| User-wise  | Aggregated view across all agents and projects         |

---

## Components

### Server
- **Language:** Go
- **Database:** SQLite (WAL mode) — ACID-compliant, zero setup
- **Architecture:** Single monolith, no microservices
- **Runs:** Continuously in the background whenever your machine is on. Stops only if you explicitly stop it.
- **Port:** Configurable via `.env`

### Frontend (Browser Dashboard)
- **Stack:** React 19 + TypeScript + Vite + Tailwind CSS + shadcn/ui
- **URL:** `http://localhost:18037`
- **Polling:** Fast polls every 2–5 seconds when the tab is active and visible. Pauses automatically when the tab is hidden or the window is closed.
- **Pages:**
  - Dashboard — live token usage, cost, and session overview
  - MCP Server Connections — manage central and project-level MCP servers
  - Skills — manage central and project-level skills
  - Hooks — configure central and project-level hooks
  - Plugins — manage central and project-level plugins
  - Authentication — per-agent account management (switch between accounts, see which agent is logged in under which account)
- **Copy & Reference System:** Skills, MCP configs, and other features can be copied or referenced across projects. Supports both centralized and project-specific configurations.

### CLI
- **Language:** Go + Bash
- **Platform:** Linux and macOS
- **Behavior:** Read-only, static snapshots. Each command fetches fresh data and displays it in formatted, human-readable output. No polling.
- **Use case:** For developers who prefer the terminal over a browser UI.

### Desktop Mini *(not yet built)*
- A system tray icon (near the Wi-Fi icon on your taskbar) built with **Go systray**
- Displays today's token usage and estimated cost at a glance (00:00–23:59 in your local timezone)
- Click to: Open the dashboard in your browser, or quit the app (with a warning that the background server will also stop)
- Ships as part of the same Go binary — no extra runtime needed

---

## Tech Stack Summary

| Component      | Technology                              |
|----------------|-----------------------------------------|
| Server         | Go, net/http, gorilla/websocket         |
| Database       | SQLite (go-sqlite3, WAL mode)           |
| Frontend       | React 19, TypeScript, Vite, Tailwind v4 |
| CLI            | Go + Bash                               |
| Desktop Mini   | Go systray                              |
| Package Format | apt, Homebrew, npx                      |

---

## Installation

> Agent Ledger is in active development. Installation packages are not yet available.

Once released, it will be installable via:

```bash
# macOS
brew install agent-ledger

# Ubuntu / Debian
sudo apt install agent-ledger

# Any platform with Node.js
npx agent-ledger
```

### Development Setup

```bash
git clone https://github.com/yourusername/agent-ledger.git
cd agent-ledger
```

**Start the server:**
```bash
cd server
cp .env.example .env   # configure HOST, PORT, DB_DRIVER, DB_NAME
go run main.go
```

**Start the frontend (dev mode):**
```bash
cd frontend
pnpm install
pnpm dev
```

The dashboard will be available at `http://localhost:18037`.

---

## Project Status

| Component       | Status          |
|-----------------|-----------------|
| Go Server       | In Progress     |
| SQLite Storage  | In Progress     |
| AI Tracking     | In Progress     |
| Browser UI      | In Progress     |
| CLI             | Planned         |
| Git Tracking    | Planned         |
| Git-AI Analytics| Planned         |
| Desktop Mini    | Planned         |

---

## Contributing

Agent Ledger is fully open source. Contributions are welcome — whether that's schema design, new agent parsers, CLI commands, or UI components.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/your-feature`)
3. Commit your changes
4. Open a pull request

Please open an issue first for significant changes so we can align on direction.

---

## License

Distributed under the MIT License. See `LICENSE` for details.
