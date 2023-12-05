package main

import (
	"log"
	"net/http"

	pbj "github.com/ccutch/pb-j-stack"
)

func main() {
	app := pbj.NewApp()
	page := app.HomePage()

	page.Serve(func(c pbj.Context) (any, error) {
		source := c.Request().URL
		return source, nil
	})

	page.On("say-hello", func(c pbj.Context) error {
		name := c.FormValue("name")
		return c.String(http.StatusOK, "Hello "+name)
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
