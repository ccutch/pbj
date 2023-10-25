package routes

import (
	"net/http"
	"os"
	"pb-stack/database"

	"github.com/labstack/echo/v5"
	"github.com/pkg/errors"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/template"
)

func PageRoutes(e *core.ServeEvent) error {
	e.Router.GET("/",
		// Handler
		func(c echo.Context) error {
			return renderPage(c, "login", nil)
		},

		// Middleware
		apis.ActivityLogger(e.App),
		apis.RequireGuestOnly(),
	)

	e.Router.GET("/dashboard",
		// Handler
		func(c echo.Context) error {
			auth, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
			return renderPage(c, "dashboard", struct {
				User *database.User
			}{database.FromAuth(auth)})
		},

		// Middleware
		apis.ActivityLogger(e.App),
		// apis.RequireRecordAuth(),
	)

	e.Router.GET("/*",
		// Serve static files from public directory
		apis.StaticDirectoryHandler(os.DirFS("./public"), true),
	)
	return nil
}

var reg = template.NewRegistry()

func renderPage(c echo.Context, name string, data any) error {
	html, err := reg.LoadFiles("templates/page.html", "templates/"+name+".html").Render(data)

	if err != nil {
		return errors.Wrap(err, "failed to render")
	}

	return c.HTML(http.StatusOK, html)
}
