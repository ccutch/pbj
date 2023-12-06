package main

import (
	"log"

	pbj "github.com/ccutch/pb-j-stack"
)

func main() {
	app := pbj.NewApp()
	page := app.HomePage()

	page.Serve(func(c pbj.Context) error {
		c.Set("url", c.Request().URL)
		return nil
	})

	page.On("say-hello", func(c pbj.Context) error {
		c.Set("name", c.FormValue("name"))
		return c.Partial("say-hello")
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
