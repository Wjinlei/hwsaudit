package show

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/Wjinlei/hwsaudit/commands/public"
	"github.com/genshen/cmds"
)

var showCommand = &cmds.Command{
	Name:        "show",
	Summary:     "show something",
	Description: "show something",
	CustomFlags: false,
	HasOptions:  true,
	FlagSet:     &flag.FlagSet{},
	Runner:      nil,
}

/* Option */
type Option struct {
	systemdUnit bool
}

var opt Option

func init() {
	opt = Option{}
	showCommand.Runner = &show{}
	fs := flag.NewFlagSet("show", flag.ContinueOnError)
	showCommand.FlagSet = fs
	showCommand.FlagSet.BoolVar(&opt.systemdUnit, "unit", true, `Show loaded systemd unit`)
	showCommand.FlagSet.Usage = showCommand.Usage // use default usage provided by cmds.Command.
	cmds.AllCommands = append(cmds.AllCommands, showCommand)
}

type show struct{}

func (v *show) PreRun() error {
	return nil // if error != nil, function Run will be not execute.
}

func (v *show) Run() error {
	unitList, err := public.ListUnits([]string{"enabled", "static"})
	if err != nil {
		return err
	}

	for _, unit := range unitList {
		jsonUnit, _ := json.Marshal(unit)
		fmt.Println(string(jsonUnit))
	}
	return nil
}
