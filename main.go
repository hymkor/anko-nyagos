package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"github.com/mattn/anko/vm"
	"github.com/mattn/go-isatty"

	"github.com/zetamatta/nyagos/defined"
	"github.com/zetamatta/nyagos/frame"
	"github.com/zetamatta/nyagos/functions"
	"github.com/zetamatta/nyagos/history"
	"github.com/zetamatta/nyagos/nodos"
	"github.com/zetamatta/nyagos/shell"

	"github.com/zetamatta/nyagos/alias"
)

type ankoFunc struct {
	f vm.Func
}

func (this *ankoFunc) Call(ctx context.Context, cmd *shell.Cmd) (next int, err error) {
	args := cmd.Args()
	param := make([]reflect.Value, 0, len(args)-1)
	for _, arg1 := range args[1:] {
		param = append(param, reflect.ValueOf(arg1))
	}
	_, err = this.f(param...)
	return 0, err
}

func (this *ankoFunc) String() string {
	return "ankoFunc"
}

func ankoAlias(name string, f interface{}) {
	println(name)
	switch code := f.(type) {
	case string:
		alias.Table[name] = alias.New(code)
	case vm.Func:
		alias.Table[name] = &ankoFunc{f: code}
	default:
		println(reflect.TypeOf(f).String())
	}
}

func loadrc(anko *vm.Env) error {
	execPath, err := os.Executable()
	if err != nil {
		return err
	}
	dir := filepath.Dir(execPath)
	script := filepath.Join(dir, "nyagos.ank")
	code, err := ioutil.ReadFile(script)
	if err != nil {
		return nil
	}
	_, err = anko.Execute(string(code))
	if err != nil {
		return fmt.Errorf("%s: %s", script, err)
	}
	return nil
}

func _main() error {
	sh := shell.New()
	defer sh.Close()
	sh.Console = nodos.GetConsole()

	ctx := context.Background()

	if !isatty.IsTerminal(os.Stdin.Fd()) {
		frame.SilentMode = true
	}

	anko := vm.NewEnv()
	anko.Define("println", fmt.Println)
	anko.Define("alias", ankoAlias)

	if err := loadrc(anko); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	var stream1 shell.Stream
	if isatty.IsTerminal(os.Stdin.Fd()) {
		constream := frame.NewCmdStreamConsole(
			func() (int, error) {
				functions.Prompt(
					&functions.Param{
						Args: []interface{}{frame.Format2Prompt(os.Getenv("PROMPT"))},
						Out:  os.Stdout,
						Err:  os.Stderr,
						In:   os.Stdin,
						Term: nodos.GetConsole(),
					},
				)
				return 0, nil
			})
		stream1 = constream
		frame.DefaultHistory = constream.History
		ctx = context.WithValue(ctx, history.PackageId, constream.History)
	} else {
		stream1 = shell.NewCmdStreamFile(os.Stdin)
	}
	sh.ForEver(ctx, stream1)
	return nil
}

var version string

func main() {
	frame.Version = version
	if err := frame.Start(_main); err != nil && err != io.EOF {
		fmt.Fprintln(os.Stderr, err)
		defer os.Exit(1)
	}
	if defined.DBG {
		os.Stdin.Read(make([]byte, 1))
	}
}
