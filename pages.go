package pbj

import (
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v5"
	"github.com/pkg/errors"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/template"
)

// NewPage entry point for page API
func (app *App) NewPage(name string, opts ...func(*Page)) *Page {
	p := Page{
		app:    app,
		route:  name,
		tmpl:   name,
		layout: hydrationTemplate,
	}
	for _, fn := range opts {
		fn(&p)
	}
	return &p
}

// Homepage helper for homepages
func (app *App) HomePage(opts ...func(*Page)) *Page {
	return app.NewPage("",
		WithTemplate("home"),
		WithPublicAccess(true),
	)
}

// Page data structure for passing configuration state to pages
type Page struct {
	app                 *App
	admin, public       bool
	layout, route, tmpl string
}

// Static serve page with only user and admin
func (p *Page) Static() *Page {
	return p.Serve(nil)
}

// Serve with Props retrieved before rendering
func (p *Page) Serve(h GetProps) *Page {
	p.app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/"+p.route, func(c echo.Context) error {
			return p.render(&pageContext{c, p, map[string]any{}}, h)
		}, apis.ActivityLogger(p.app))
		return nil
	})
	return p
}

// GetProps type for getting the props of a page when needed
type GetProps func(Context) error

// Render function is core to functionallity and can be optimized
func (p *Page) render(c Context, h GetProps) (err error) {
	reg := template.NewRegistry()
	isHtmx := c.Request().Header.Get("Hx-Request") == "true"

	// Simplify the api - file based routing?
	if p.tmpl == "" {
		p.tmpl = p.route
	}

	// If auth is present render the requested page w/ data
	user, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
	admin, _ := c.Get(apis.ContextAdminKey).(*models.Admin)
	if (!p.admin && user != nil) || (p.admin && admin != nil) || (p.public && isHtmx) {
		parts, _ := filepath.Glob("templates/partials/*.html")
		parts = append([]string{"templates/pages/" + p.tmpl + ".html"}, parts...)
		c.Set("app", p.app)
		c.Set("page", p)
		c.Set("user", user)
		c.Set("admin", admin)
		if h != nil {
			if err = h(c); err != nil {
				return errors.Wrap(err, "failed fetch data")
			}
		}
		view := reg.LoadFiles(parts...)
		html, err := view.Render(c.Props())
		if err != nil {
			return errors.Wrap(err, "failed to render")
		}
		return c.HTML(http.StatusOK, html)
	}

	// If htmx request w/o auth render login page w/o data
	if isHtmx {
		parts, _ := filepath.Glob("templates/partials/*.html")
		if p.admin {
			parts = append([]string{"templates/pages/admin-login.html"}, parts...)
		} else {
			parts = append([]string{"templates/pages/login.html"}, parts...)
		}
		view := reg.LoadFiles(parts...)
		html, err := view.Render(nil)
		if err != nil {
			return errors.Wrap(err, "failed to render")
		}
		return c.HTML(http.StatusOK, html)
	}

	// If no auth and no htmx then we render page runner
	view := reg.LoadString(p.layout)
	html, err := view.Render(struct {
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
