package pbj

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

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
	FormValues() (url.Values, error)
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
	page  *Page
	props map[string]any
}

func (ctx *pageContext) Admin() *models.Admin {
	admin, _ := ctx.Get(apis.ContextAdminKey).(*models.Admin)
	return admin
}

func (ctx *pageContext) User() *models.Record {
	user, _ := ctx.Get(apis.ContextAuthRecordKey).(*models.Record)
	return user
}

func (ctx *pageContext) Event(n string) string {
	return fmt.Sprintf("%s/@%s", ctx.Request().URL.Path, n)
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
	reg := template.NewRegistry()
	// Simplify the api - file based routing?
	if name == "" {
		name = ctx.page.route
	}
	parts, _ := filepath.Glob("templates/partials/*.html")
	parts = append([]string{"templates/" + name + ".html"}, parts...)
	html, err := reg.AddFuncs(map[string]any{
		"event": ctx.Event,
	}).LoadFiles(parts...).Render(ctx.Props())
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

func (ctx *eventContext) Event(n string) string {
	path := ctx.Request().URL.Path
	parts := strings.Split(path, "/@")
	return fmt.Sprintf("%s/@%s", parts[0], n)
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
	html, err := reg.AddFuncs(map[string]any{
		"event": ctx.Event,
	}).LoadFiles(parts...).Render(ctx.Props())
	if err != nil {
		return err
	}
	return ctx.HTML(http.StatusOK, html)
}
