o
// Events (SQLITE CRUD)
app.Events(func(events *Events) {
	// REST triggered events
	events.OnCreate("resource", func(req *DbContext, res *Resource) error {
		// TODO: Handle Create event for 'resource'.
		return nil
	})

	events.OnRead("resource", func(req *DbContext) (*Resource, error) {
		// TODO: Handle Read 'resource' request event. 
		return nil, nil
	})

	events.OnUpdate("resource", func(req *DbContext, res *Resource) error {
		// TODO: Handle Update event for 'resource'.
		return nil
	})

	events.OnDelete("resource", func(req *DbContext) error {
		// TODO: Handle Delete 'resource' event.
		return nil
	})
})
