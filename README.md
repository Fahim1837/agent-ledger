# 🤖 agent-ledger
The Next-Gen Git Analytics Engine for the AI-Assisted Development Era
agent-ledger bridges the gap between traditional software engineering metrics and modern AI-driven workflows. By tracking line-level attribution, local token usage, and agent behaviors, it provides engineering leaders and developers unmatched visibility into code quality, technical debt, and AI efficiency in the age of generative development.

✨ Features
🕵️ Local Agent Telemetry: Track real-time token usage, prompt/completion ratios, active sessions, and estimated costs across all local AI tools (Cursor, Copilot, Claude, etc.).

🛠️ Centralized Agent Control: Configure API keys, global Model Context Protocol (MCP) servers, plugins, and hooks from a single dashboard—with granular, project-specific overrides.

📈 Git-Level Attribution: Understand exactly which parts of your codebase were human-written versus AI-generated using native Git Notes telemetry.

📊 AI Code Durability (Coming Soon): Advanced analytics measuring how long AI-generated code survives before requiring a human rewrite or hotfix.

💻 Developer-First UX: A lightweight native system tray app that feeds a real-time browser dashboard and a powerful terminal CLI.

🏗️ Technical Architecture & Stack
agent-ledger is designed from the ground up to be local-first, ultra-lightweight, and zero-configuration. It compiles down to a single executable binary with no external heavy dependencies (no Docker, Postgres, or Node.js runtimes required on the user's machine).

Core Engine & Background Worker: Built with Go. Handles low-level file watching, process tracking, native OS system tray interaction (systray), and local network hooks.

Storage Engine: Powered by SQLite in WAL (Write-Ahead Logging) Mode. Provides ACID-compliant, lightning-fast persistence directly inside the native Go runtime, comfortably handling tens of thousands of real-time writes per second.

Frontend Dashboard: Built using React + TypeScript + Vite. The compiled, static production assets are natively embedded directly into the Go binary using Go's embed package and served locally at localhost:18037.

Real-Time Networking: Uses WebSockets for instant, sub-millisecond telemetry pushes from the Go engine to the React UI without page refreshes.

Interception Layer: Functions as a local developer proxy/log parser to seamlessly capture AI agent completions without invading developer privacy or interrupting IDE performance.

🚀 Getting Started
⚠️ Note: agent-ledger is currently in active development.

Prerequisites
Because agent-ledger is distributed as a single static binary, you do not need to install background databases or UI runtimes.

Installation (Development Setup)
Clone the repository:

Bash
git clone https://github.com/yourusername/agent-ledger.git
cd agent-ledger

2. Build the React frontend:
   ```bash
   cd ui && npm install && npm run build
Build and run the Go backend (which embeds the UI automatically):

Bash
cd ../backend && go build -o agent-ledger
./agent-ledger

Once running, look for the `agent-ledger` robot icon in your system tray/taskbar. Click **Open Dashboard** or navigate to `http://localhost:18037` to view your analytics.

---

## 🤝 Contributing

This is an open-source project and contributions of all kinds are welcome! Whether you want to optimize SQLite time-series schemas, help parse logs for new AI extensions, or polish the React charts, feel free to open an issue or submit a PR.

---

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

---

### 💡 Tips for your GitHub Repository setup:
*   **Repository Tags (Topics):** Add tags like `go`, `react`, `sqlite`, `ai-agents`, `developer-tools`, `git-analytics`, and `llm`. 
*   **Pinned Repo:** Once you have the codebase up, pin this repository to the top of your GitHub profile. 
*   **Architecture Diagram:** If you have time, sketch out a clean architecture diagram and place it right under the "Technical Architecture" heading. Recruiters eat that up!