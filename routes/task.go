package routes

import (
	"net/http"
	"pb-stack/database"

	"github.com/labstack/echo/v5"
	"github.com/pkg/errors"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

func TaskRoutes(e *core.ServeEvent) error {

	e.Router.PUT("/complete",

		// Handler
		func(c echo.Context) error {
			return completeTask(e.App, c)
		},

		// Middleware
		apis.ActivityLogger(e.App),
		apis.RequireRecordAuth(),
	)

	return nil
}

func completeTask(app core.App, c echo.Context) error {
	var task database.Task
	if err := app.Dao().FindById(&task, c.FormValue("id")); err != nil {
		return errors.Wrap(err, "failed to find task")
	}

	auth, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
	if task.Owner != auth.Id {
		return errors.New("invalid permissions")
	}

	task.Completed = true
	if err := app.Dao().Save(&task); err != nil {
		return errors.Wrap(err, "failed to save completed task")
	}

	return c.JSON(http.StatusOK, &task)
}
