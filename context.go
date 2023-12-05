package pbj

import (
	"errors"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/tools/template"
)

type Context interface {
	echo.Context
	Refresh() error
	Push(url string) error
	Replace(url string) error
	Partial(string, any) error
}

var (
	_ Context = (*pageContext)(nil)
	_ Context = (*eventContext)(nil)
)

// Page Context
type pageContext struct {
	echo.Context
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

func (*pageContext) Partial(string, any) error {
	return errors.New("There is no partial rendering for pages")
}

// Event Context
type eventContext struct {
	echo.Context
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

func (ctx *eventContext) Partial(name string, data any) error {
	reg := template.NewRegistry()
	parts, _ := filepath.Glob("templates/partials/*.html")
	parts = append([]string{"templates/partials/" + name + ".html"}, parts...)
	page := reg.LoadFiles(parts...)
	html, err := page.Render(data)
	if err != nil {
		return err
	}
	return ctx.HTML(http.StatusOK, html)
}
