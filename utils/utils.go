package utils

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/rpc/jsonrpc"
	"sort"

	"github.com/Trial97/sflags"
	"github.com/Trial97/sflags/gen/gflag"
	prompt "github.com/c-bata/go-prompt"
	"github.com/cgrates/rpcclient"
)

type AttrSetBalance struct {
	Tenant          string `desc:"HTTP host"`
	Account         string
	BalanceType     string
	Value           float64
	Balance         map[string]interface{}
	ActionExtraData *map[string]interface{}
	Cdrlog          bool
}

func ToIJSON(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", " ")
	return string(b)
}

func ToJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func Call(conn rpcclient.ClientConnector, args *AttrSetBalance) error {
	var rply string
	if err := conn.Call("APIerSv1.SetBalance", args, &rply); err != nil {
		return err
	}
	fmt.Println(rply)
	return nil
}

func NewConn() rpcclient.ClientConnector {
	cl, err := jsonrpc.Dial("tcp", "192.168.56.203:2012")
	if err != nil {
		panic(err)
	}
	return cl
}

func CompleteFlagSet(args []string) (sugsetions []string, err error) {
	cfg := new(AttrSetBalance)
	fs := flag.NewFlagSet("set_balance", flag.ContinueOnError)
	if err = gflag.ParseTo(cfg, fs, sflags.FlagDivider("."), sflags.ToFlag(func(s1, s2 string) string { return s1 })); err != nil {
		return
	}
	fs.SetOutput(io.Discard)
	if err = fs.Parse(args); err != nil {
		return
	}

	firstSugestions := make(map[string]string)
	fs.VisitAll(func(f *flag.Flag) {
		firstSugestions[f.Name] = f.Usage
	})
	fs.Visit(func(f *flag.Flag) {
		delete(firstSugestions, f.Name)
	})
	sugsetions = make([]string, 0, len(firstSugestions))
	for name, dsc := range firstSugestions {
		fmt.Println(dsc)
		sugsetions = append(sugsetions, "-"+name)
	}
	sort.Strings(sugsetions)
	return
}

func CompleteFlagSet2(args []string) (sugsetions []prompt.Suggest, err error) {
	cfg := new(AttrSetBalance)
	fs := flag.NewFlagSet("set_balance", flag.ContinueOnError)
	if err = gflag.ParseTo(cfg, fs, sflags.FlagDivider("."), sflags.ToFlag(func(s1, s2 string) string { return s1 })); err != nil {
		return
	}
	fs.SetOutput(io.Discard)
	if err = fs.Parse(args); err != nil {
		return
	}

	firstSugestions := make(map[string]string)
	fs.VisitAll(func(f *flag.Flag) {
		firstSugestions[f.Name] = f.Usage
	})
	fs.Visit(func(f *flag.Flag) {
		delete(firstSugestions, f.Name)
	})
	sugsetions = make([]prompt.Suggest, 0, len(firstSugestions))
	for name, dsc := range firstSugestions {
		sugsetions = append(sugsetions, prompt.Suggest{Text: "-" + name, Description: dsc})
	}
	// sort.Slice(sugsetions,func(i, j int) bool {})
	return
}

func ExecuteFlagSet(conn rpcclient.ClientConnector, args []string) {
	cfg := &AttrSetBalance{
		Tenant: "cgrates.org",
	}
	fs := flag.NewFlagSet("set_balance", flag.ContinueOnError)
	if err := gflag.ParseTo(cfg, fs, sflags.FlagDivider("."), sflags.ToFlag(func(s1, s2 string) string { return s1 })); err != nil {
		fmt.Println("Error: ", err)
	}
	if err := fs.Parse(args); err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Println(ToJSON(cfg))
	if err := Call(conn, cfg); err != nil {
		fmt.Println("Error: ", err)
		return
	}
}
