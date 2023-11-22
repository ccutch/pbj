o
app.Events(func(events *Events) {
	// REST triggered events
	events.OnCreate("resource", func(req *DbContext, res *Resource) error {
		// TODO: Handle Create event for 'resource'.
		err := CreateResource(req, res)
		if err != nil {
			return err
		}

		return nil
	})

	events.OnRead("resource", func(req *DbContext) (*Resource, error) {
		// TODO: Handle Read 'resource' request event. 
		resource, err := ReadResource(req)
		if err != nil {
			return nil, err
		}

		return resource, nil
	})

	events.OnUpdate("resource", func(req *DbContext, res *Resource) error {
		// TODO: Handle Update event for 'resource'.
		err := UpdateResource(req, res)
		if err != nil {
			return err
		}

		return nil
	})

	events.OnDelete("resource", func(req *DbContext) error {
		// TODO: Handle Delete 'resource' event.
		err := DeleteResource(req)
		if err != nil {
			return err
		}

		return nil
	})
})
