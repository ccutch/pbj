package pbj

import (
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v5"
	"github.com/pkg/errors"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/template"
)

type Context interface {
	Request() *http.Request
	Response() *echo.Response
	FormValue(string) string
	PathParam(string) string
	QueryParam(string) string

	Admin() *models.Admin
	User() *models.Record
	Set(string, any)
	Props() map[string]any
	Refresh() error
	Push(url string) error
	Replace(url string) error
	Render(string) error
}

var (
	_ Context = (*pageContext)(nil)
	_ Context = (*eventContext)(nil)
)

// Page Context
type pageContext struct {
	echo.Context
	page     *Page
	getProps GetProps
	props    map[string]any
}

func (ctx *pageContext) Admin() *models.Admin {
	admin, _ := ctx.Get(apis.ContextAdminKey).(*models.Admin)
	return admin
}

func (ctx *pageContext) User() *models.Record {
	user, _ := ctx.Get(apis.ContextAuthRecordKey).(*models.Record)
	return user
}

func (ctx *pageContext) Props() map[string]any {
	return ctx.props
}

func (ctx *pageContext) Set(k string, v any) {
	ctx.props[k] = v
}

func (ctx *pageContext) Unwrap() echo.Context {
	return ctx.Context
}

func (ctx *pageContext) Refresh() error {
	ctx.Response().Header().Set("HX-Refresh", "true")
	return nil
}

func (ctx *pageContext) Push(url string) error {
	ctx.Response().Header().Set("HX-Redirect", url)
	return nil
}

func (ctx *pageContext) Replace(url string) error {
	ctx.Response().Header().Set("HX-Location", url)
	return nil
}

func (ctx *pageContext) Render(name string) error {
	p := ctx.page
	reg := template.NewRegistry()
	isHtmx := ctx.Request().Header.Get("Hx-Request") == "true"

	// Simplify the api - file based routing?
	if name == "" {
		name = p.route
	}

	// If auth is present render the requested page w/ data
	user, _ := ctx.Get(apis.ContextAuthRecordKey).(*models.Record)
	admin, _ := ctx.Get(apis.ContextAdminKey).(*models.Admin)
	if (!p.admin && user != nil) || (p.admin && admin != nil) || (p.public && isHtmx) {
		ctx.Set("app", p.app)
		ctx.Set("page", p)
		ctx.Set("user", user)
		ctx.Set("admin", admin)
		parts, _ := filepath.Glob("templates/partials/*.html")
		parts = append([]string{"templates/" + name + ".html"}, parts...)
		if err := ctx.getProps(ctx); err != nil {
			return errors.Wrap(err, "failed to get props")
		}
		html, err := reg.LoadFiles(parts...).Render(ctx.Props())
		if err != nil {
			return errors.Wrap(err, "failed to render")
		}
		return ctx.HTML(http.StatusOK, html)
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
		return ctx.HTML(http.StatusOK, html)
	}

	// If no auth and no htmx then we render page runner
	html, err := reg.LoadString(hydrationTemplate).Render(struct {
		Page          string
		Params        string
		HeaderContent string
	}{
		ctx.Request().URL.Path,
		ctx.Request().URL.RawQuery,
		p.app.headerContent,
	})
	if err != nil {
		return errors.Wrap(err, "failed to render")
	}

	return ctx.HTML(http.StatusOK, html)
}

// Event Context
type eventContext struct {
	echo.Context
	app   *App
	page  *Page
	props map[string]any
}

func (ctx *eventContext) Admin() *models.Admin {
	admin, _ := ctx.Get(apis.ContextAdminKey).(*models.Admin)
	return admin
}

func (ctx *eventContext) User() *models.Record {
	user, _ := ctx.Get(apis.ContextAuthRecordKey).(*models.Record)
	return user
}

func (ctx *eventContext) Props() map[string]any {
	return ctx.props
}

func (ctx *eventContext) Set(k string, v any) {
	ctx.props[k] = v
}

func (ctx *eventContext) Unwrap() echo.Context {
	return ctx.Context
}

func (ctx *eventContext) Refresh() error {
	ctx.Response().Header().Set("HX-Refresh", "true")
	return ctx.NoContent(http.StatusOK)
}

func (ctx *eventContext) Push(url string) error {
	ctx.Response().Header().Set("HX-Redirect", url)
	return ctx.NoContent(http.StatusOK)
}

func (ctx *eventContext) Replace(url string) error {
	ctx.Response().Header().Set("HX-Location", url)
	return ctx.NoContent(http.StatusOK)
}

func (ctx *eventContext) Render(name string) error {
	reg := template.NewRegistry()
	parts, _ := filepath.Glob("templates/partials/*.html")
	parts = append([]string{"templates/" + name + ".html"}, parts...)
	user, _ := ctx.Get(apis.ContextAuthRecordKey).(*models.Record)
	admin, _ := ctx.Get(apis.ContextAdminKey).(*models.Admin)
	ctx.Set("app", ctx.app)
	ctx.Set("page", ctx.page)
	ctx.Set("user", user)
	ctx.Set("admin", admin)
	page := reg.LoadFiles(parts...)
	html, err := page.Render(ctx.Props())
	if err != nil {
		return err
	}
	return ctx.HTML(http.StatusOK, html)
}
