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

	app.registerHealthRoutes(mux)
	app.registerPostRoutes(mux)
	app.registerUserRoutes(mux)

	return mux
}

func (app *Application) registerHealthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", app.handleIndex)
}

func (app *Application) registerPostRoutes(mux *http.ServeMux) {

}

func (app *Application) registerUserRoutes(mux *http.ServeMux) {

}

func (app *Application) handleIndex(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "Server is Healthy",
	})
}
