package goprompt

import (
	"console/utils"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	"github.com/cgrates/rpcclient"
)

var ok bool
var suggestions = []prompt.Suggest{{
	Text:        "help",
	Description: "prints the help",
}, {
	Text:        "exit",
	Description: "exit console",
}, {
	Text:        "set_balance",
	Description: "pings something",
}}

func livePrefix() (string, bool) {
	return "gopt> ", ok
}

var hst = []string{
	`ping -subsystem=1 -m test1:{"A":"fieldA","V":1} -m2 F:D -test.b=100 -test.c=5 -s=[{"Z":"!"}]`,
}
var newHst = []string{}

func executor(in string) {
	blocks := strings.Split(strings.TrimSpace(in), " ")
	newHst = append(newHst, in)
	switch v := blocks[0]; v {
	case "exit":
		fmt.Println("Bye!")
		os.Exit(0)
	case "help":
		fmt.Println(suggestions)
	case "set_balance":
		utils.ExecuteFlagSet(conn, blocks[1:])
	default:
		fmt.Printf("No CMD<%s>\n", v)
	}

}

func completer(in prompt.Document) (s []prompt.Suggest) {
	l := in.CurrentLineBeforeCursor()
	if l == "" {
		return
	}
	if i := strings.IndexByte(l, ' '); i != -1 {
		cmd := l[:i]
		if cmd != "set_balance" {
			return
		}
		args := l[i+1:]
		e := strings.LastIndex(args, " ")
		var lw string
		if e != -1 {
			lw = args[e+1:]
			args = args[:e]
		} else {
			lw = args
			args = ""
		}
		s, _ := utils.CompleteFlagSet2(strings.Split(args, " "))
		return prompt.FilterHasPrefix(s, lw, true)
	}
	w := in.GetWordBeforeCursor()
	if w == "" {
		return
	}
	return prompt.FilterHasPrefix(suggestions, w, true)
}

var conn rpcclient.ClientConnector

func Main() {
	conn = utils.NewConn()
	historyFN := "./history3"

	f, err := os.OpenFile(historyFN, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Could not open history file because:<%s>", err)
	}
	v, err := io.ReadAll(f)

	if err == nil {
		hst = strings.Split(string(v), "\n")
	}

	defer func() {
		for _, h := range newHst {
			f.WriteString(h + "\n")
		}
		f.Close()
	}()
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix("gopt> "),
		prompt.OptionLivePrefix(livePrefix),
		prompt.OptionHistory(hst),
	)
	p.Run()

}
