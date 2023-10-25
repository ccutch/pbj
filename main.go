package main

import (
	"log"
	"pb-stack/handlers"
	"pb-stack/routes"

	"github.com/pocketbase/pocketbase"
)

func main() {
	// New App
	app := pocketbase.New()

	// REST handlers
	app.OnRecordsListRequest("tasks").Add(handlers.ListTasks)

	// Custom Routes
	app.OnBeforeServe().Add(routes.PageRoutes)
	app.OnBeforeServe().Add(routes.TaskRoutes)

	// It's PB&J Time
	log.Fatal(app.Start())
}
