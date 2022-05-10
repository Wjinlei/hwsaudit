package run

import (
	"errors"
	"flag"
	"os"

	"github.com/Wjinlei/hwsaudit/commands/public"
	"github.com/genshen/cmds"
)

var runCommand = &cmds.Command{
	Name:        "run",
	Summary:     "running",
	Description: "running",
	CustomFlags: false,
	HasOptions:  true,
	FlagSet:     &flag.FlagSet{},
	Runner:      nil,
}

/* Option */
type Option struct {
	C         string
	directory string // d
	user      string // u
	fileMode  string // m
	fileAcl   string // a
	setUidGid bool   // s
	setSticky bool   // t
}

var opt Option

func init() {
	opt = Option{}
	runCommand.Runner = &run{}
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	runCommand.FlagSet = fs
	runCommand.FlagSet.StringVar(&opt.C, "C", "file", `check target,    eg: ["file"|"dir"|"all"]`)
	runCommand.FlagSet.StringVar(&opt.directory, "d", "./", `check directory, eg: "/wwwroot"`)
	runCommand.FlagSet.StringVar(&opt.user, "u", "*", `check file user, eg: "-root"`)
	runCommand.FlagSet.StringVar(&opt.fileMode, "m", "*", `check file mode, eg: ["**2"|"644"|"777"]`)
	runCommand.FlagSet.StringVar(&opt.fileAcl, "a", "", `check file acl,  eg: "user1:2,user2:*,*:2,*:*"`)
	runCommand.FlagSet.BoolVar(&opt.setUidGid, "s", false, `type:[bool] check setuid setgid.`)
	runCommand.FlagSet.BoolVar(&opt.setSticky, "t", false, `type:[bool] check sticky`)

	runCommand.FlagSet.Usage = runCommand.Usage // use default usage provided by cmds.Command.
	cmds.AllCommands = append(cmds.AllCommands, runCommand)
}

type run struct{}

func (v *run) PreRun() error {
	lstat, err := os.Lstat(opt.directory)
	if err != nil {
		return err
	}
	if lstat.IsDir() == false {
		return errors.New("Not directory")
	}
	return nil // if error != nil, function Run will be not execute.
}

func (v *run) Run() error {
	return public.WalkDir(false, opt.directory, opt.C, opt.user, opt.fileMode, opt.setUidGid, opt.setSticky, opt.fileAcl)
}
