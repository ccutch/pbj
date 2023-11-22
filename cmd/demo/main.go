o
// Events (SQLITE CRUD)
app.Events(func(events *Events) {
    events.POST("/create", func (c echo.Context) error {
        // Your handler for create event.
    })

    events.GET("/read", func (c echo.Context) error {
        // Your handler for read event.
    })

    events.PUT("/update", func (c echo.Context) error {
        // Your handler for update event.
    })

    events.DELETE("/delete", func (c echo.Context) error {
        // Your handler for delete event.
    })
})
