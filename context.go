package pbj

import (
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/template"
)

type Context interface {
	echo.Context
	Set(string, any)
	Props() map[string]any
	Refresh() error
	Push(url string) error
	Replace(url string) error
	Render() error
}

var (
	_ Context = (*pageContext)(nil)
	_ Context = (*eventContext)(nil)
)

// Page Context
type pageContext struct {
	echo.Context
	page  *Page
	props map[string]any
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

func (*pageContext) Render() error {
	return nil // render is default behavior
}

// Event Context
type eventContext struct {
	echo.Context
	name  string
	app   *App
	page  *Page
	props map[string]any
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

func (ctx *eventContext) Render() error {
	reg := template.NewRegistry()
	parts, _ := filepath.Glob("templates/partials/*.html")
	parts = append([]string{
		"templates/partials/" + name + ".html",
	}, parts...)
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
