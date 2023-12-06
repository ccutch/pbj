package pbj

import (
	"fmt"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

// On function to mount callbacks to App
func (a *App) On(n string, fn func(Context) error) {
	a.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		route := a.Event(n)
		e.Router.POST(route, func(c echo.Context) error {
			return a.handle(&eventContext{c, a, nil, map[string]any{}}, fn)
		}, apis.ActivityLogger(a))
		return nil
	})
}

func (a *App) handle(c *eventContext, fn func(Context) error) error {
	c.Set("app", a)
	c.Set("user", c.User())
	c.Set("admin", c.Admin())
	return fn(c)
}

// Event formats event name for client and server side
func (a *App) Event(n string) string {
	return fmt.Sprintf("/@%s", n)
}

// On function to mount callbacks to PocketBase
func (p *Page) On(n string, fn func(Context) error) {
	p.app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		route := p.Event(n)
		e.Router.POST(route, func(c echo.Context) error {
			return p.handle(&eventContext{c, p.app, p, map[string]any{}}, fn)
		}, apis.ActivityLogger(p.app))
		return nil
	})
}

func (p *Page) handle(c *eventContext, fn func(Context) error) error {
	c.Set("app", p.app)
	c.Set("page", p)
	c.Set("user", c.User())
	c.Set("admin", c.Admin())
	return fn(c)
}

// Event formats event name for client and server side
func (p *Page) Event(n string) string {
	return fmt.Sprintf("%s/@%s", p.route, n)
}
