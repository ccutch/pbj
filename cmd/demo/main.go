package main

import (
	"log"
	"os"
	"pb-stack/handlers"
	"pb-stack/routes"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	_ "pb-stack/migrations"
)

func main() {
	// New App
	app := pocketbase.New()

	// REST Handlers
	app.OnRecordsListRequest("tasks").Add(handlers.ListTasks)

	// Custom Routes
	app.OnBeforeServe().Add(routes.PageRoutes)
	app.OnBeforeServe().Add(routes.TaskRoutes)

	// Migrations
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())
	log.Println("is go run", isGoRun)

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Admin UI
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	})

	// It's PB&J Time
	log.Fatal(app.Start())
}
