package run

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"

	"github.com/genshen/cmds"
)

var runCommand = &cmds.Command{
    Name:        "run",
    Summary:     "running",
    Description: "running",
    CustomFlags: false,
    HasOptions:  true,
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
    runCommand.FlagSet.StringVar(&opt.C,         "C", "file", `check target,    eg: ["file"|"dir"|"all"]`)
    runCommand.FlagSet.StringVar(&opt.directory, "d", "./",   `check directory, eg: "/wwwroot"`)
    runCommand.FlagSet.StringVar(&opt.user,      "u", "*",    `check file user, eg: "-root"`)
    runCommand.FlagSet.StringVar(&opt.fileMode,  "m", "*",    `check file mode, eg: ["**2"|"644"|"777"]`)
    runCommand.FlagSet.StringVar(&opt.fileAcl,   "a", "*",    `check file acl,  eg: "user1:2,user2:*,*:2,*:*"`)
    runCommand.FlagSet.BoolVar(&opt.setUidGid,   "s", false,  `type:[bool] check setuid setgid.`)
    runCommand.FlagSet.BoolVar(&opt.setSticky,   "t", false,  `type:[bool] check sticky`)

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
    err := filepath.WalkDir(opt.directory, func(path string, d fs.DirEntry, err error) error {
        if err != nil { return err }

        stat, err := d.Info()
        if err != nil { return nil }

        /* Exclude */
        if stat.Mode() & os.ModeSymlink != 0 {
            return nil
        }

        ok := true
        acl := ""
        result := Result{}
        sSys := stat.Sys().(*syscall.Stat_t)

        /* Match scan target */
        if opt.C == "file" && d.IsDir() {
            return nil
        }
        if opt.C == "dir" && !d.IsDir() {
            return nil
        }

        /* Match user */
        if opt.user != "" && opt.user != "*" {
            ok = isMatchUser(int(sSys.Uid))
            if !ok {
                return nil
            }
        }

        /* Match mode */
        if opt.fileMode != "" && opt.fileMode != "*" {
            ok = isMatchMode(stat.Mode())
            if !ok {
                return nil
            }
        }

        /* Match setuid */
        if opt.setUidGid {
            if stat.Mode() & os.ModeSetuid != 0 || stat.Mode() & os.ModeSetgid != 0{
                ok = true
            } else {
                return nil
            }
        }

        /* Match setgid */
        if opt.setSticky {
            if stat.Mode() & os.ModeSticky != 0 {
                ok = true
            } else {
                return nil
            }
        }

        /* Match acl */
        if opt.fileAcl != "" && opt.fileAcl != "*" {
            acl, ok = isMatchAcl(path)
            if !ok {
                return nil
            }
        }

        /* Check ok */
        if ok {
            result.Name = d.Name()
            result.Path = GetAbsPath(path)
            result.User = findUser(int(sSys.Uid))
            result.Mode = stat.Mode().String()
            result.Acl = acl

            jsonResult, _ := json.Marshal(result)
            fmt.Println(string(jsonResult))
        }
        return nil
    })
    if err != nil {
        return err
    }

    return nil
}
