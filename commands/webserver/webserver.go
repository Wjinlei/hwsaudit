package webserver

import (
	"flag"

	"github.com/genshen/cmds"
	"github.com/gin-gonic/gin"
)

var webserverCommand = &cmds.Command{
	Name:        "webserver",
	Summary:     "start web server",
	Description: "start web server",
	CustomFlags: false,
	HasOptions:  false,
	FlagSet:     &flag.FlagSet{},
	Runner:      nil,
}

func init() {
	webserverCommand.Runner = &version{}
	fs := flag.NewFlagSet("webserver", flag.ContinueOnError)
	webserverCommand.FlagSet = fs
	webserverCommand.FlagSet.Usage = webserverCommand.Usage // use default usage provided by cmds.Command.
	cmds.AllCommands = append(cmds.AllCommands, webserverCommand)
}

type version struct{}

func (v *version) PreRun() error {
	return nil
}

func (v *version) Run() error {
	router := gin.Default()
	router.Static("/css", "./html/css/")
	router.Static("/assets", "./html/assets/")
	router.Static("/js", "./html/js/")
	router.LoadHTMLFiles("./html/index.html")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", nil)
	})
	router.Run(":8000")
	return nil
}
