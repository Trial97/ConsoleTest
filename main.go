package main

import (
	"os"

	"console/goprompt"
	"console/grumble"
	"console/ishell"
	"console/liner"
)

func main() {
	args := os.Args[1:]
	consoleType := "goprompt"
	if len(args) > 0 {
		consoleType = args[0]
	}
	os.Args = os.Args[1:]
	switch consoleType {
	case "liner":
		liner.Main()
	case "goprompt":
		goprompt.Main()
	case "ishell":
		ishell.Main()
	case "grumble":
		grumble.Main()
	}
}
