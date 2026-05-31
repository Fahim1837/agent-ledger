package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"agent-ledger/config"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dbPath := config.DatabasePath()
	db, err := config.OpenSQLite(dbPath)
	if err != nil {
		log.Fatalf("database setup failed: %v", err)
	}
	defer db.Close()

	eventHub := config.NewEventHub(config.DefaultChannelBuffer)
	go eventHub.Run(ctx)

	httpAddr := config.HTTPAddr()
	app := config.NewApplication(db, eventHub)
	server := config.NewHTTPServer(httpAddr, app)
	serverErrors := config.StartHTTPServer(server)

	fmt.Printf("agent-ledger server started with database at %s\n", dbPath)

	select {
	case <-ctx.Done():
		if err := config.ShutdownHTTPServer(server); err != nil {
			log.Printf("server shutdown failed: %v", err)
		}
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed: %v", err)
		}
	}
}
