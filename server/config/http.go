package config

import (
	"context"
	"database/sql"
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

	routes := Logging(app.Routes(mux))
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
