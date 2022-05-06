package main

import (
	"errors"
	"flag"
	"fmt"

	_ "github.com/Wjinlei/hwsaudit/commands/run"
	_ "github.com/Wjinlei/hwsaudit/commands/version"
	"github.com/genshen/cmds"
)

func main() {
	cmds.SetProgramName("hwsaudit")
	if err := cmds.Parse(); err != nil {
		if err == flag.ErrHelp {
			return
		}
		// skip error in sub command parsing, because the error has been printed.
		if !errors.Is(err, &cmds.SubCommandParseError{}) {
			fmt.Println(err)
		}
	}
}
