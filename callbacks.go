package pbj

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/template"
)

// Callbacks mounts callback handler methods as HTTP POST routes
func (app *App) Callbacks(fn func(*Handlers)) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		fn(&Handlers{e})
		return nil
	})
}

// Handlers aree
type Handlers struct {
	*core.ServeEvent
}

// On function to mount callbacks to PocketBase
func (handlers *Handlers) On(route string, handler func(*HandlerContext) error) {
	handlers.Router.POST("/"+route, func(c echo.Context) error {
		return handler(&HandlerContext{c})
	}, apis.ActivityLogger(handlers.App))
}

// HandlerContext encapsolates echo context with additional behavior
type HandlerContext struct{ echo.Context }

// Render a partial file from templates/partials from callback
func (c *HandlerContext) Render(name string, data any) error {
	reg := template.NewRegistry()
	page := reg.LoadFiles("templates/partials/" + name + ".html")
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
