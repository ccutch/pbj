package pbj

import "fmt"

// Generic mount function for "routers"
func Mount(fn func(*App)) func(*App) { return fn }

// WithMeta will include <meta name content> tag
func WithMeta(name, content string) func(*App) {
	return func(app *App) {
		app.headerContent = fmt.Sprintf(
			`%s<meta name="%s" content="%s">`,
			app.headerContent, name, content,
		)
	}
}

// WithStylesheet will include <link rel="stylesheet" href> tag
func WithStylesheet(href string) func(*App) {
	return func(app *App) {
		app.headerContent = fmt.Sprintf(
			`%s<link rel="stylesheet" href="%s" />`,
			app.headerContent, href,
		)
	}
}

// WithScript will include <script src></script> tag
func WithScript(src string) func(*App) {
	return func(app *App) {
		app.headerContent = fmt.Sprintf(
			`%s<script src="%s"></script>`,
			app.headerContent, src,
		)
	}
}

// WithInlineScript will include <script></script> tag
func WithInlineScript(content string) func(*App) {
	return func(app *App) {
		app.headerContent = fmt.Sprintf(
			`%s<script>%s</script>`,
			app.headerContent, content,
		)
	}
}

// WithPublicAccess configures the Page to be public accessable
func WithPublicAccess(public bool) func(*Page) {
	return func(p *Page) {
		p.public = public
	}
}

// WithAdminOnlyAccess configures the Page to only allow admin access
func WithAdminOnlyAccess(admin bool) func(*Page) {
	return func(p *Page) {
		p.admin = admin
	}
}
