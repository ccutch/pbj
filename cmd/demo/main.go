package main

import (
	"log"

	. "github.com/ccutch/pb-j-stack"

	"github.com/labstack/echo/v5"
)

func main() {

	app := NewApp(
		WithStylesheet("https://cdn.jsdelivr.net/npm/daisyui@3.9.4/dist/full.css"),
		WithScript("https://cdn.tailwindcss.com"),
	)

	// Pages (HTTP GET)
	app.Static(HomePage(WithTemplate("home")))
	app.Static(NewPage("about"))
	app.Static(NewPage("path-to-profit"))
	app.Serve(NewPage("hello/:name"), func(c echo.Context) (any, error) {
		// Construct a struct to say hello
		return struct{ Name string }{
			Name: c.PathParam("name"),
		}, nil
	})

	// Callbacks (HTTP POST)
	app.On("hello-again", func(c *HandlerContext) error {
		return c.Refresh()
	})

	app.On("hello-another", func(c *HandlerContext) error {
		return c.Redirect("hello/" + c.FormValue("name"))
	})

	// Events (SQLITE CRUD)
	// app.Events(func(events *Events) {
	// TODO: REST triggered events
	// })

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
