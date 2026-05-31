package config

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

const DefaultHTTPAddr = ":9000"

func HTTPAddr() string {
	if addr := os.Getenv("AGENT_LEDGER_HTTP_ADDR"); addr != "" {
		return addr
	}

	return DefaultHTTPAddr
}

type Application struct {
	DB  *sql.DB
	Hub *EventHub
}

func NewApplication(db *sql.DB, hub *EventHub) *Application {
	return &Application{
		DB:  db,
		Hub: hub,
	}
}

func NewHTTPServer(addr string, app *Application) *http.Server {
	if addr == "" {
		addr = DefaultHTTPAddr
	}

	return &http.Server{
		Addr:         addr,
		Handler:      app.Routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func StartHTTPServer(server *http.Server) <-chan error {
	errs := make(chan error, 1)

	go func() {
		log.Printf("agent-ledger server listening on http://localhost%s", server.Addr)
		errs <- server.ListenAndServe()
	}()

	return errs
}

func ShutdownHTTPServer(server *http.Server) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return server.Shutdown(ctx)
}

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", app.handleIndex)
	mux.HandleFunc("GET /health", app.handleHealth)
	mux.HandleFunc("GET /api/health", app.handleHealth)
	mux.HandleFunc("GET /ws", app.handleWebSocket)

	return mux
}

func (app *Application) handleIndex(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"name":   "agent-ledger",
		"status": "running",
	})
}

func (app *Application) handleHealth(w http.ResponseWriter, r *http.Request) {
	if err := app.DB.PingContext(r.Context()); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{
			"database": "unavailable",
			"error":    err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status":    "ok",
		"database":  "ok",
		"websocket": "ok",
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}
