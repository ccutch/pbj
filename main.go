package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/pkg/errors"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/template"
)

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		registry := template.NewRegistry()

		e.Router.GET("/", func(c echo.Context) error {
			r := registry.LoadFiles("views/layout.html", "views/index.html")

			html, err := r.Render(nil)
			if err != nil {
				log.Println(err)
				return errors.Wrap(err, "failed to render")
			}

			return c.HTML(http.StatusOK, html)
		})

		e.Router.POST("/hello", func(c echo.Context) error {
			r := registry.LoadFiles("views/hello.html")

			html, err := r.Render(struct{ Name string }{Name: c.FormValue("name")})
			if err != nil {
				return errors.Wrap(err, "failed to Load template")
			}

			return c.HTML(http.StatusOK, html)
		})

		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./public"), true))

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
