package pbj

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

// On function to mount callbacks to App
func (a *App) On(n string, h func(Context) error) {
	a.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		route := a.Event(n)
		e.Router.POST(route, func(c echo.Context) error {
			return h(&eventContext{c})
		}, apis.ActivityLogger(a))
		return nil
	})
}

// Event formats event name for client and server side
func (a *App) Event(n string) string {
	return fmt.Sprintf("/@%s", n)
}

// On function to mount callbacks to PocketBase
func (p *Page) On(n string, h func(Context) error) {
	p.app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		route := p.Event(n)
		log.Println("route = ", route)
		e.Router.POST(route, func(c echo.Context) error {
			return h(&eventContext{c})
		}, apis.ActivityLogger(p.app))
		return nil
	})
}

// Event formats event name for client and server side
func (p *Page) Event(n string) string {
	return fmt.Sprintf("%s/@%s", p.route, n)
}
