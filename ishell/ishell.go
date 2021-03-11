package ishell

import (
	"console/utils"

	ishell "gopkg.in/abiosoft/ishell.v2"
)

func Main() {
	// create new shell.
	// by default, new shell includes 'exit', 'help' and 'clear' commands.
	shell := ishell.New()
	conn := utils.NewConn()
	// display welcome info.
	shell.SetHistoryPath("./history2")
	shell.SetPrompt("ishell> ")
	// register a function for "greet" command.
	shell.AddCmd(&ishell.Cmd{
		Name: "set_balance",
		Help: "sets the balance",
		Completer: func(args []string) []string {
			a, _ := utils.CompleteFlagSet(args)
			return a
		},
		Func: func(c *ishell.Context) {
			utils.ExecuteFlagSet(conn, c.Args)
		},
	})

	// run shell
	shell.Run()
}
