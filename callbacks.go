package pbj

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/template"
)

// On function to mount callbacks to PocketBase
func (app *App) On(route string, handler func(*HandlerContext) error) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.POST("/"+route, func(c echo.Context) error {
			return handler(&HandlerContext{c})
		}, apis.ActivityLogger(app))
		return nil
	})
}

// HandlerContext encapsolates echo context with additional behavior
type HandlerContext struct{ echo.Context }

// Render a partial file from templates/partials from callback
func (c *HandlerContext) Render(name string, data any) error {
	reg := template.NewRegistry()
	parts, _ := filepath.Glob("templates/partials/*.html")
	parts = append([]string{"templates/partials/" + name + ".html"}, parts...)
	page := reg.LoadFiles(parts...)
	html, err := page.Render(data)
	if err != nil {
		log.Println("Error rendering partial", err)
		return err
	}
	return c.HTML(http.StatusOK, html)
}

// Refresh the page after callback is complete
func (c *HandlerContext) Refresh() error {
	c.Response().Header().Set("HX-Refresh", "true")
	return c.NoContent(http.StatusOK)
}

// Redirect the page to another page without reloading window
func (c *HandlerContext) Redirect(dest string) error {
	c.Response().Header().Set("HX-Redirect", dest)
	return c.NoContent(http.StatusOK)
}
