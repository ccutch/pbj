package pbj

import (
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

// NewPage entry point for page API
func (app *App) NewPage(route string, opts ...func(*Page)) *Page {
	p := Page{
		app:   app,
		route: route,
	}
	for _, fn := range opts {
		fn(&p)
	}
	return &p
}

// Page data structure for passing configuration state to pages
type Page struct {
	app           *App
	admin, public bool
	route         string
}

// Static serve page with only user and admin
func (p *Page) Static() *Page {
	return p.Serve(nil)
}

// Serve with Props retrieved before rendering
func (p *Page) Serve(h GetProps) *Page {
	p.app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/"+p.route, func(c echo.Context) error {
			return h(&pageContext{c, p, map[string]any{}})
		}, apis.ActivityLogger(p.app))
		return nil
	})
	return p
}

// GetProps type for getting the props of a page when needed
type GetProps func(Context) error

// Everything we need for hydration including some sensible defaults
var hydrationTemplate = `<!DOCTYPE html>
<html lang="en">
  <head>
    <title>{{block "title" .}}{{end}}</title>
    <meta name="viewport" content="width=user-scalable=no, device-width, height=device-height, initial-scale=1.0, minimum-scale=1.0">
    <meta charSet="utf-8">
    {{.HeaderContent|raw}}
    <script src="https://unpkg.com/htmx.org@1.9.6"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/pocketbase/0.19.0/pocketbase.umd.js"></script>
  </head>
  <body>
    <script>
	  const pb = new PocketBase(window.location.href);
      document.body.addEventListener("htmx:configRequest", function (event) {
	    if (pb && pb.authStore) event.detail.headers['Authorization'] = "Bearer "+pb.authStore.token;
      });
    </script>
    <main hx-get="{{.Page}}{{with .Params}}?{{.}}{{end}}" hx-trigger="load" hx-swap="outerHTML"></main>
  </body>
</html>`
