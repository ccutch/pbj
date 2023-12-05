package pbj

import (
	"os"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

// NewApp creates a new app encapsolating PocketBase
func NewApp(opts ...func(*App)) *App {
	app := App{
		PocketBase:    pocketbase.New(),
		headerContent: "",
		publicFiles:   "public",
	}
	for _, o := range opts {
		o(&app)
	}
	return &app
}

// App data structure encapsolates PocketBase app
type App struct {
	*pocketbase.PocketBase
	headerContent string
	publicFiles   string
}

// Start starts pocketbase app after registering migrations
func (app *App) Start() error {

	// Serve static files from PocketBase docs by default
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		pubDir := os.DirFS(app.publicFiles)
		e.Router.GET("/*", apis.StaticDirectoryHandler(pubDir, false))
		return nil
	})

	// Automigrations from PocketBase docs by default
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: strings.HasPrefix(os.Args[0], os.TempDir()),
	})

	// Defer to PocketBase after this point
	return app.PocketBase.Start()
}
