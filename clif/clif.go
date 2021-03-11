package clif

import "gopkg.in/ukautz/clif.v1"

// Some type definition
type MyFoo struct {
	X int
}

func Main() {
	// init cli
	cli := clif.New("My App", "1.0.0", "An example application")

	// register object instance with container
	foo := &MyFoo{X: 123}
	cli.Register(foo)

	// Create command with callback using the peviously registered instance
	cli.Add(clif.NewCommand("foo", "Call foo", func(foo *MyFoo) {
		// do something with foo
	}))

	cli.Run()
}
