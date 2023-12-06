package pbj

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pkg/errors"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/template"
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
func (p *Page) Serve(fn GetProps) *Page {
	p.app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/"+p.route, func(c echo.Context) error {
			return p.serve(&pageContext{c, p, map[string]any{}}, fn)
		}, apis.ActivityLogger(p.app))
		return nil
	})
	return p
}

func (p *Page) serve(c *pageContext, fn GetProps) error {
	isHtmx := c.Request().Header.Get("Hx-Request") == "true"
	user, admin := c.User(), c.Admin()

	// If auth is present render the requested page w/ data
	if (!p.admin && user != nil) || (p.admin && admin != nil) || (p.public && isHtmx) {
		c.Set("app", p.app)
		c.Set("page", p)
		c.Set("user", user)
		c.Set("admin", admin)
		if fn == nil {
			return c.Render(p.route)
		}
		return fn(c)
	}

	// If htmx request w/o auth render login page w/o data
	if isHtmx {
		if p.admin {
			return c.Render("pages/admin-login.html")
		} else {
			return c.Render("pages/login.html")
		}
	}

	// If no auth and no htmx then we render page runner
	reg := template.NewRegistry()
	html, err := reg.LoadString(hydrationTemplate).Render(struct {
		Page          string
		Params        string
		HeaderContent string
	}{
		c.Request().URL.Path,
		c.Request().URL.RawQuery,
		p.app.headerContent,
	})
	if err != nil {
		return errors.Wrap(err, "failed to render")
	}

	return c.HTML(http.StatusOK, html)
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
