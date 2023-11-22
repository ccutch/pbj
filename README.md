app.Events(func(events *Events) {
	events.On("create-entity", func(e *EventContext) error {
		// insertion code for entity in SQLite database
		return nil
	})

	events.On("update-entity", func(e *EventContext) error {
		// update code for entity in SQLite database
		return nil
	})

	events.On("delete-entity", func(e *EventContext) error {
		// deletion code for entity in SQLite database
		return nil
	})
})
