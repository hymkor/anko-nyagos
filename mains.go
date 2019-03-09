package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mattn/anko/vm"
	"github.com/mattn/go-isatty"

	"github.com/zetamatta/nyagos/frame"
	"github.com/zetamatta/nyagos/functions"
	"github.com/zetamatta/nyagos/history"
	"github.com/zetamatta/nyagos/nodos"
	"github.com/zetamatta/nyagos/shell"

	"github.com/zetamatta/nyagos/alias"
)

func ankoAlias(name string, f interface{}) {
	println(name)
	if code, ok := f.(string); ok {
		alias.Table[name] = alias.New(code)
	}
}

func loadrc(anko *vm.Env) error {
	execPath, err := os.Executable()
	if err != nil {
		return err
	}
	dir := filepath.Dir(execPath)
	script := filepath.Join(dir, "nyagos.ank")
	fd, err := os.Open(script)
	if err != nil {
		return nil
	}
	defer fd.Close()

	code, err := ioutil.ReadAll(fd)
	if err != nil {
		return err
	}
	_, err = anko.Execute(string(code))
	return err
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
