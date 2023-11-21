package pbj

import (
	"fmt"
	"os"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

// NewApp creates a new app encapsolating PocketBase
func NewApp(opts ...func(*App)) *App {
	app := App{pocketbase.New(), ""}
	for _, o := range opts {
		o(&app)
	}
	return &app
}

// App data structure encapsolates PocketBase app
type App struct {
	*pocketbase.PocketBase
	headerContent string
}

// Start starts pocketbase app after registering migrations
func (app *App) Start() error {

	// Serve static files from PocketBase docs by default
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("public"), false))
		return nil
	})

	// Automigrations from PocketBase docs by default
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: strings.HasPrefix(os.Args[0], os.TempDir()),
	})

	// Defer to PocketBase after this point
	return app.PocketBase.Start()
}

func WithMeta(name, content string) func(*App) {
	return func(app *App) {
		app.headerContent = fmt.Sprintf(
			`%s\n<meta name="%s" content="%s">`,
			app.headerContent, name, content,
		)
	}
}

func WithStylesheet(href string) func(*App) {
	return func(app *App) {
		app.headerContent = fmt.Sprintf(
			`%s\n<link rel="stylesheet" href="%s" />`,
			app.headerContent, href,
		)
	}
}

func WithScript(src string) func(*App) {
	return func(app *App) {
		app.headerContent = fmt.Sprintf(
			`%s\n<script src="%s"></script>`,
			app.headerContent, src,
		)
	}
}

func WithInlineScript(content string) func(*App) {
	return func(app *App) {
		app.headerContent = fmt.Sprintf(
			`%s\n<script>%s</script>`,
			app.headerContent, content,
		)
	}
}
