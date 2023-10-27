package routes

import (
	"net/http"
	"pb-stack/models"

	"github.com/labstack/echo/v5"
	"github.com/pkg/errors"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func TaskRoutes(e *core.ServeEvent) error {

	e.Router.PUT("/tasks/complete",

		// Handler
		func(c echo.Context) error {
			return completeTask(e.App, c)
		},

		// Middleware
		apis.ActivityLogger(e.App),
		apis.RequireAdminAuth(),
	)

	return nil
}

func completeTask(app core.App, c echo.Context) error {
	var task models.Task
	if err := app.Dao().FindById(&task, c.QueryParam("id")); err != nil {
		return errors.Wrap(err, "failed to find task")
	}

	task.Completed = !task.Completed
	if err := app.Dao().Save(&task); err != nil {
		return errors.Wrap(err, "failed to save completed task")
	}

	page := reg.LoadFiles("templates/parts/task.html")
	html, err := page.Render(task)
	if err != nil {
		return errors.Wrap(err, "failed to render")
	}
	return c.HTML(http.StatusOK, html)
}
