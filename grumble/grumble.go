package grumble

import (
	"console/utils"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/desertbit/grumble"
)

func Main() {
	conn := utils.NewConn()
	var app = grumble.New(&grumble.Config{
		Name:        "grumble",
		Description: "short app description",
		Flags: func(f *grumble.Flags) {
			f.String("d", "directory", "DEFAULT", "set an alternative directory path")
			f.Bool("v", "verbose", false, "enable verbose mode")
		},
		HistoryFile: "./history1",
	})
	app.AddCommand(&grumble.Command{
		Name: "set_balance",
		Help: "sets the balance",
		// Aliases: []string{"run"},
		Completer: func(prefix string, args []string) []string {
			if prefix == "{" {
				return []string{`{\"Account\":\"1001\",\"BalanceType\":\"*voice\"}`}
			}
			return nil
		},
		Flags: func(f *grumble.Flags) {
			f.Duration("t", "Account", time.Second, "timeout duration")
		},
		Args: func(a *grumble.Args) {
			a.String("args", "the args formated as json")
			// a.Int("opt-level", "the optimization mode", grumble.Default(3))
			// a.StringList("services", "additional services that should be started", grumble.Default([]string{"test", "te11"}))
		},
		Run: func(c *grumble.Context) error {
			args := new(utils.AttrSetBalance)
			fmt.Println(c.Args.String("args"))
			if err := json.Unmarshal([]byte(c.Args.String("args")), args); err != nil {
				return err
			}
			return utils.Call(conn, args)
			// c.App.Println("timeout:", c.Flags.Duration("timeout"))
			// c.App.Println("directory:", c.Flags.String("directory"))
			// c.App.Println("verbose:", c.Flags.Bool("verbose"))
			// c.App.Println("production:", c.Args.Bool("production"))
			// c.App.Println("opt-level:", c.Args.Int("opt-level"))
			// c.App.Println("services:", strings.Join(c.Args.StringList("services"), ","))
		},
	})

	adminCommand := &grumble.Command{
		Name:     "admin",
		Help:     "admin tools",
		LongHelp: "super administration tools",
	}
	app.AddCommand(adminCommand)

	adminCommand.AddCommand(&grumble.Command{
		Name: "root",
		Help: "root the machine",
		Run: func(c *grumble.Context) error {
			c.App.Println(c.Flags.String("directory"))
			return errors.New("failed")
		},
	})
	grumble.Main(app)
}
