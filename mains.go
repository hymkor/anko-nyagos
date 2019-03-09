package main

import (
	"context"
	"os"

	"github.com/mattn/go-isatty"

	"github.com/zetamatta/nyagos/frame"
	"github.com/zetamatta/nyagos/functions"
	"github.com/zetamatta/nyagos/history"
	"github.com/zetamatta/nyagos/nodos"
	"github.com/zetamatta/nyagos/shell"
)

func _main() error {
	sh := shell.New()
	defer sh.Close()
	sh.Console = nodos.GetConsole()

	ctx := context.Background()

	if !isatty.IsTerminal(os.Stdin.Fd()) {
		frame.SilentMode = true
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
