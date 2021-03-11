package liner

import (
	"console/utils"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/peterh/liner"
)

var (
	historyFn   = filepath.Join("./history4")
	suggestions = []string{"help", "exit", "set_balance"}
)

func Main() {
	conn := utils.NewConn()
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

	line.SetCompleter(func(line string) (c []string) {
		l := line
		if l == "" {
			return suggestions
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
			sg, _ := utils.CompleteFlagSet(strings.Split(args, " "))
			s := make([]string, 0, len(suggestions))
			for _, s1 := range sg {
				if strings.HasPrefix(s1, lw) {
					s = append(s, strings.TrimSpace(cmd)+" "+args+" "+s1)
				}
			}
			return s
		}
		s := make([]string, 0, len(suggestions))
		for _, d := range suggestions {
			if strings.HasPrefix(d, l) {
				s = append(s, d)
			}
		}
		return s
	})

	if f, err := os.Open(historyFn); err == nil {
		line.ReadHistory(f)
		f.Close()
	}
	var stop bool
	for !stop {
		if command, err := line.Prompt("liner> "); err != nil {
			if err == io.EOF {
				fmt.Println("\nbye!")
				stop = true
			} else {
				fmt.Print("Error reading line: ", err)
			}
		} else {
			command = strings.TrimSpace(command)
			if command == "" {
				continue
			}
			line.AppendHistory(command)
			blocks1 := strings.Split(command, " ")
			blocks := make([]string, 0, len(blocks1))
			for _, b := range blocks1 {
				b = strings.TrimSpace(b)
				if b != "" {
					blocks = append(blocks, b)
				}
			}

			switch v := blocks[0]; v {
			case "exit":
				fmt.Println("\nbye!")
				stop = true
			case "help":
				fmt.Println(suggestions)
			case "set_balance":
				utils.ExecuteFlagSet(conn, blocks[1:])
			default:
				fmt.Printf("No CMD<%s>\n", v)
			}
		}
	}

	if f, err := os.Create(historyFn); err != nil {
		log.Print("Error writing history file: ", err)
	} else {
		line.WriteHistory(f)
		f.Close()
	}
}
