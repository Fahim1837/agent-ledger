package config

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func GetApplication(db *sql.DB) *Application {
	app := &Application{
		DB: db,
	}

	return app
}

func BuildRouter(app *Application) http.Handler {
	mux := http.NewServeMux()

	routes := app.Routes(mux)
	return routes
}

func StartHTTPServer(route http.Handler) {
	host, port := os.Getenv("HOST"), os.Getenv("PORT")
	addr := fmt.Sprintf("%s:%s", host, port)

	log.Printf("Server starting at http://%s", addr)
	err := http.ListenAndServe(addr, route)
	if err != nil {
		log.Fatal("Failed to start the Server: ", err)
	}
}

type Application struct {
	DB *sql.DB
}

func ShutdownHTTPServer(server *http.Server) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return server.Shutdown(ctx)
}

func (app *Application) Routes(mux *http.ServeMux) http.Handler {

	mux.HandleFunc("/", app.handleIndex)
	mux.HandleFunc("/health", app.handleHealth)
	mux.HandleFunc("/api/health", app.handleHealth)

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
