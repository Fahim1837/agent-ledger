# Proxy Interception — High-Level Architecture

## The Core Idea

Agent Ledger runs a lightweight HTTP server on `localhost`. AI tools send their API
requests through it. The proxy logs the token data, then forwards the request to the
real provider and hands the response back untouched.

The AI tool never knows a proxy is in the middle.

---

## Option A — Base URL Override

### How it works

The AI tool's SDK is told to use `http://localhost:PORT` instead of
`https://api.anthropic.com`. No system-level changes needed. No TLS interception.

### Flow

```
┌─────────────────────────────────────────────────────────────────┐
│  Developer Machine                                              │
│                                                                 │
│  ┌──────────────┐    HTTP (plain)    ┌────────────────────────┐ │
│  │  Claude Code │ ────────────────▶  │   Agent Ledger Proxy   │ │
│  │  Cursor      │                    │   localhost:8080        │ │
│  │  OpenCode    │ ◀────────────────  │                        │ │
│  └──────────────┘    response        │  1. receive request    │ │
│                                      │  2. log token data     │ │
│                                      │  3. forward upstream   │ │
│                                      │  4. return response    │ │
│                                      └──────────┬─────────────┘ │
│                                                 │               │
└─────────────────────────────────────────────────┼───────────────┘
                                                  │ HTTPS (real TLS)
                                                  ▼
                                      ┌───────────────────────┐
                                      │   api.anthropic.com   │
                                      │   api.openai.com      │
                                      │   etc.                │
                                      └───────────────────────┘
```

### How to activate (per tool)

| Tool | What to set |
|------|-------------|
| Claude Code | `ANTHROPIC_BASE_URL=http://localhost:8080` |
| Any Anthropic SDK app | Same env var |
| OpenAI SDK app | `OPENAI_BASE_URL=http://localhost:8080` |

The SDK sends plain HTTP to your proxy. Your proxy makes the real HTTPS call upstream
using its own TLS stack. No certificate games needed.

### Proxy responsibilities

```
Request in  ──▶  parse body (JSON)
                     │
                     ▼
                 extract: model, input_tokens (from request)
                     │
                     ▼
                 forward to real API over HTTPS
                     │
                     ▼
                 receive response
                     │
                     ▼
                 extract: output_tokens, cache tokens (from response)
                     │
                     ▼
                 write event to SQLite
                     │
                     ▼
Response out ◀──  return response to tool (unmodified)
```

### Pros / Cons

| ✅ Pros | ❌ Cons |
|---------|---------|
| No TLS cert needed | Tool must support a base URL override |
| Easy to set up | Must set env var before launching the tool |
| Request body is plaintext (easy to parse) | Doesn't capture tools that hardcode the URL |

---

## Option B — HTTP_PROXY / HTTPS_PROXY

### How it works

The OS/shell has two well-known env vars that most HTTP clients (curl, Go's `net/http`,
Python's `requests`, Node's `axios`, etc.) respect automatically. Setting them redirects
all outbound HTTP/HTTPS through your proxy — no per-tool config needed.

### Flow

```
┌─────────────────────────────────────────────────────────────────────┐
│  Developer Machine                                                  │
│                                                                     │
│  env: HTTPS_PROXY=http://localhost:8080                             │
│                                                                     │
│  ┌──────────────┐  CONNECT api.anthropic.com:443  ┌──────────────┐ │
│  │  Any AI tool │ ───────────────────────────────▶ │  Agent       │ │
│  │  (uses SDK   │                                  │  Ledger      │ │
│  │   HTTP lib)  │ ◀─────────────────────────────── │  Proxy       │ │
│  └──────────────┘   tunnel established             │  :8080       │ │
│         │                                          │              │ │
│         │ TLS handshake through tunnel             │  sees only:  │ │
│         │ (proxy cannot see inside)                │  · hostname  │ │
│         ▼                                          │  · byte count│ │
│    encrypted traffic ──────────────────────────────┤  (no body)  │ │
│                                                    └──────┬───────┘ │
└───────────────────────────────────────────────────────────┼─────────┘
                                                            │ forwards tunnel
                                                            ▼
                                                ┌───────────────────────┐
                                                │   api.anthropic.com   │
                                                └───────────────────────┘
```

> **Problem:** HTTPS tunnels are encrypted end-to-end. The proxy sees a raw byte stream,
> not the JSON body — so it **cannot read token counts** without TLS termination.

### To read the body: TLS termination required

```
┌─────────────────────────────────────────────────────────────────────┐
│  Developer Machine                                                  │
│                                                                     │
│  ┌──────────────┐                          ┌──────────────────────┐ │
│  │  AI Tool     │  TLS (local CA cert)     │  Agent Ledger Proxy  │ │
│  │              │ ────────────────────────▶│                      │ │
│  │              │ ◀────────────────────────│  1. terminate TLS    │ │
│  └──────────────┘  response                │     (local CA cert)  │ │
│                                            │  2. read plaintext   │ │
│  System trust store                        │  3. log tokens       │ │
│  └─ Agent Ledger CA cert (installed once)  │  4. re-encrypt       │ │
│                                            │  5. forward upstream │ │
│                                            └──────────────────────┘ │
└─────────────────────────────────────────────────────────────────────┘
```

- You generate a local CA certificate once on install
- Install it into the OS trust store (so the AI tool trusts it)
- The proxy presents a cert signed by your CA for each upstream host
- The tool sees a valid cert, the proxy sees plaintext

### How to activate

```bash
export HTTP_PROXY=http://localhost:8080
export HTTPS_PROXY=http://localhost:8080
```

Or set it system-wide in `/etc/environment` so every process inherits it.

### Pros / Cons

| ✅ Pros | ❌ Cons |
|---------|---------|
| Catches any tool automatically | Requires TLS cert generation + install |
| No per-tool configuration | More complex proxy implementation |
| Works even if tool hardcodes the URL | Some tools pin certificates (will break) |

---

## Side-by-Side Comparison

```
                   Option A                      Option B
                   (Base URL override)           (HTTPS_PROXY)
                   ─────────────────────         ─────────────────────
Setup effort       Low — one env var             Medium — CA cert install
TLS handling       Proxy does real TLS out       Proxy terminates TLS in
Body readable?     Yes (plain HTTP in)           Only with MITM cert
Tool support       SDK must allow base URL       Any HTTP client
MVP suitability    ✅ Start here                 🔜 Add later

```

---

## Recommended Approach for Agent Ledger

```
Phase 1 (MVP)  ──▶  Option A only
                    · Cover Claude Code (ANTHROPIC_BASE_URL)
                    · Full token data, zero TLS complexity

Phase 2        ──▶  Option B with TLS termination
                    · Cover tools that don't expose a base URL
                    · Generate + install CA cert on setup
                    · Broader coverage

Phase 3        ──▶  One-time JSONL file import
                    · Backfill history from before the proxy was installed
                    · No proxy needed — read Claude Code's local log files
```

The proxy codebase is the same for both options — the difference is only in
**how the connection arrives** (plain HTTP vs CONNECT tunnel + TLS termination).
