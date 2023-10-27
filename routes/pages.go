package routes

import (
	"net/http"
	"os"
	"pb-stack/models"

	"github.com/labstack/echo/v5"
	"github.com/pkg/errors"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/template"
)

func PageRoutes(e *core.ServeEvent) error {
	e.Router.GET("/",
		// Handler
		func(c echo.Context) error {
			return renderPage(c, "boards", nil)
		},

		// Middleware
		apis.ActivityLogger(e.App),
		apis.RequireGuestOnly(),
	)

	e.Router.GET("/boards",
		// Handler
		func(c echo.Context) error {
			return renderPage(c, "boards", func(a *models.Admin) (any, error) {
				boards := []*models.Board{}
				return boards, e.App.Dao().
					ModelQuery(&models.Board{}).
					All(&boards)
			})
		},

		// Middleware
		apis.ActivityLogger(e.App),
	)

	e.Router.GET("/tasks",
		// Handler
		func(c echo.Context) error {
			return renderPage(c, "tasks", func(a *models.Admin) (any, error) {
				data := struct {
					Board *models.Board
					Tasks []*models.Task
				}{nil, []*models.Task{}}

				err := e.App.Dao().
					ModelQuery(&models.Task{}).
					AndWhere(dbx.HashExp{"board": c.QueryParam("board")}).
					All(&data.Tasks)
				if err != nil {
					return nil, err
				}

				return data, e.App.Dao().
					ModelQuery(&models.Board{}).
					AndWhere(dbx.HashExp{"id": c.QueryParam("board")}).
					Limit(1).
					One(&data.Board)
			})
		},

		// Middleware
		apis.ActivityLogger(e.App),
	)

	e.Router.GET("/*",
		// Serve static files from public directory
		apis.StaticDirectoryHandler(os.DirFS("./public"), true),
	)
	return nil
}

var reg = template.NewRegistry()

func renderPage(c echo.Context, name string, dataFn func(*models.Admin) (any, error)) error {

	// If auth is present render the requested page w/ data
	auth, _ := c.Get(apis.ContextAdminKey).(*models.Admin)
	if auth != nil {
		page := reg.LoadFiles(
			"templates/pages/"+name+".html",
			"templates/parts/task.html",
		)
		data, err := dataFn(auth)
		if err != nil {
			return errors.Wrap(err, "failed fetch data")
		}
		html, err := page.Render(data)
		if err != nil {
			return errors.Wrap(err, "failed to render")
		}
		return c.HTML(http.StatusOK, html)
	}

	// If htmx request w/o auth render login page w/o data
	if c.Request().Header.Get("Hx-Request") == "true" {
		page := reg.LoadFiles("templates/pages/login.html")
		html, err := page.Render(nil)
		if err != nil {
			return errors.Wrap(err, "failed to render")
		}
		return c.HTML(http.StatusOK, html)
	}

	// If no auth and no htmx then we render page runner
	page := reg.LoadFiles("templates/page.html")
	html, err := page.Render(struct {
		Page   string
		Params string
	}{name, c.Request().URL.RawQuery})
	if err != nil {
		return errors.Wrap(err, "failed to render")
	}

	return c.HTML(http.StatusOK, html)
}
