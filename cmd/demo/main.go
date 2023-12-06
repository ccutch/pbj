package main

import (
	"log"

	pbj "github.com/ccutch/pb-j-stack"
)

func main() {
	app := pbj.NewApp()
	page := app.NewPage("", pbj.WithPublicAccess(true))

	page.Serve(func(ctx pbj.Context) error {
		ctx.Set("url", ctx.Request().URL)
		return ctx.Render("pages/home")
	})

	page.On("say-hello", func(ctx pbj.Context) error {
		ctx.Set("name", ctx.FormValue("name"))
		return ctx.Render("partials/say-hello")
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
