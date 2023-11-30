o
app.Events(func(events *Events) {

	// Create event
	events.On("create:user", func(c *HandlerContext) error {
		user := new(User)
		if err := c.Bind(user); err != nil {
			return err
		}
		// TODO: insert 'user' into the SQLite database
		return nil
	})

	// Update event
	events.On("update:user", func(c *HandlerContext) error {
		user := new(User)
		if err := c.Bind(user); err != nil {
			return err
		}
		// TODO: update 'user' in the SQLite database
		return nil
	})
	
	// Delete event
	events.On("delete:user", func(c *HandlerContext) error {
		id := c.FormValue("id")
		// TODO: delete user with 'id' from the SQLite database
		return nil
	})
})
