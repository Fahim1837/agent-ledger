package main

import (
	"agent-ledger/config"
)

func main() {
	// Load the env files
	config.LoadEnv()

	// Connect to the DB
	db := config.GetDBConnection()
	defer db.Close()

	// Lister to the server
	app := config.GetApplication(db)
	routes := config.BuildRouter(app)
	config.StartHTTPServer(routes)

}
