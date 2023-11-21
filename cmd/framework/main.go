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
		WithInlineScript(`
tailwind.config = {
  daisyui: {
	themes: [
	  "synthwave",
	  "aqua",
	  "luxury",
	  "dracula",
	  "emerald",
	  "lofi",
	  "dim",
	],
  },
};

(() => {
  const theme = localStorage.getItem('theme') ?? 'synthwave';
  if (theme) document.documentElement.setAttribute('data-theme', theme);
})();`))

	// Pages (HTTP GET)
	app.Pages(func(pages *Pages) {
		pages.Static(HomePage(WithTemplate("home")))
		pages.Static(NewPage("about"))
		pages.Static(NewPage("path-to-profit"))

		pages.Serve(NewPage("hello/:name"), sayHello)
	})

	// Callbacks (HTTP POST)
	app.Callbacks(func(handlers *Handlers) {
		handlers.On("hello-again", func(c *HandlerContext) error {
			return c.Refresh()
		})

		handlers.On("hello-another", func(c *HandlerContext) error {
			return c.Redirect("hello/" + c.FormValue("name"))
		})
	})

	// Events (SQLITE CRUD)
	// app.Events(func(events *Events) {
	// TODO
	// })

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func sayHello(c echo.Context) (any, error) {
	// Construct a struct to say hello
	return struct {
		Name string
	}{c.PathParam("name")}, nil
}
