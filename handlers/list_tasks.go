package handlers

import (
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/template"
)

var reg = template.NewRegistry()

func ListTasks(e *core.RecordsListEvent) error {
	if e.HttpContext.Request().Header.Get("HX-Request") == "" {
		return nil
	}

	html, err := reg.LoadFiles("templates/test.html").Render(nil)
	if err != nil {
		return errors.Wrap(err, "failed to Load template")
	}

	log.Println(html)
	return e.HttpContext.HTML(http.StatusOK, html)
}
