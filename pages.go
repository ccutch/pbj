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

// Static function for mounting static pages (SSG optimally)
func (app *App) Static(p *Page) {
	app.Serve(p, func(c echo.Context) (any, error) {
		var data struct {
			User *models.Record
		}

		data.User, _ = c.Get(apis.ContextAuthRecordKey).(*models.Record)
		return &data, nil
	})
}

// Serve function for mounting dynamit pages (SSR optimally)
func (app *App) Serve(p *Page, h GetProps) {
	p.app = app
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/"+p.route, func(c echo.Context) error {
			return p.Render(c, h)
		}, apis.ActivityLogger(app))
		return nil
	})
}

// GetProps type for getting the props of a page when needed
type GetProps func(echo.Context) (any, error)

// Page data structure for passing configuration state to pages
type Page struct {
	app           *App
	admin, public bool
	route, tmpl   string
}

// NewPage creates a new page instance from given options
func NewPage(route string, options ...func(*Page)) *Page {
	var p Page
	for _, o := range options {
		o(&p)
	}
	p.route = route
	return &p
}

// HomePage routes for root, ie. https://www.example.com
func HomePage(options ...func(*Page)) *Page {
	return NewPage("", append(
		[]func(*Page){
			WithTemplate("index"),
			WithPublicAccess(true),
		},
		options...,
	)...)
}

// WithTemplate configures Page to use a specific template file
func WithTemplate(tmpl string) func(*Page) {
	return func(p *Page) {
		p.tmpl = tmpl
	}
}

// WithPublicAccess configures the Page to be public accessable
func WithPublicAccess(public bool) func(*Page) {
	return func(p *Page) {
		p.public = public
	}
}

// WithAdminOnly configures the Page to only allow admin access
func WithAdminOnly(admin bool) func(*Page) {
	return func(p *Page) {
		p.admin = admin
	}
}

// Render function is core to functionallity and can be optimized
func (p *Page) Render(c echo.Context, h GetProps) (err error) {
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
		page := reg.LoadFiles(parts...)
		var data any
		if h != nil {
			data, err = h(c)
			if err != nil {
				return errors.Wrap(err, "failed fetch data")
			}
		}
		html, err := page.Render(data)
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
		page := reg.LoadFiles(parts...)
		html, err := page.Render(nil)
		if err != nil {
			return errors.Wrap(err, "failed to render")
		}
		return c.HTML(http.StatusOK, html)
	}

	// If no auth and no htmx then we render page runner
	page := reg.LoadString(hydrationTemplate)
	html, err := page.Render(struct {
		Page          string
		Params        string
		HeaderContent string
	}{c.Request().URL.Path, c.Request().URL.RawQuery, p.app.headerContent})
	if err != nil {
		return errors.Wrap(err, "failed to render")
	}

	return c.HTML(http.StatusOK, html)
}

// Everything we need for hydration including some sensible defaults
var hydrationTemplate = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <title>{{block "title" .}}{{end}}</title>
    <meta name="viewport" content="width=user-scalable=no, device-width, height=device-height, initial-scale=1.0, minimum-scale=1.0">
    <meta charSet="utf-8">
    {{.HeaderContent|raw}}
    <script src="https://unpkg.com/htmx.org@1.9.6"></script>
    <script type="module">
      import PocketBase from '/scripts/pocketbase.es.js';
      const pb = new PocketBase(window.location.href);
      window.logout = () => {
        pb.authStore.clear();
        window.location.reload();
      };
      document.body.addEventListener("htmx:configRequest", (event) => {
        event.detail.headers['Authorization'] = "Bearer "+pb.authStore.token;
      });
    </script>
    <style>*{box-sizing:border-box}</style>
  </head>
  <body>
    <main
	  hx-get="{{.Page}}?{{.Params}}"
      hx-trigger="load"
      hx-swap="outerHTML"
    ></main>
  </body>
</html>
`
