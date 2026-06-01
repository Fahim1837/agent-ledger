package config

import (
	"encoding/json"
	"log"
	"net/http"
)

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

func (app *Application) Routes(mux *http.ServeMux) http.Handler {
	apiMux := http.NewServeMux()

	app.registerHealthRoutes(apiMux)
	app.registerPostRoutes(apiMux)
	app.registerUserRoutes(apiMux)

	mux.Handle("/api/", http.StripPrefix("/api", apiMux))

	return mux
}

func (app *Application) registerHealthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", app.handleHealth)
}

func (app *Application) registerPostRoutes(mux *http.ServeMux) {

}

func (app *Application) registerUserRoutes(mux *http.ServeMux) {

}

func (app *Application) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "Server is Healthy",
	})
}
