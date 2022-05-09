package webserver

import (
	"flag"

	"github.com/Wjinlei/hwsaudit/commands/public"
	"github.com/genshen/cmds"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/host"
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

	v1 := router.Group("/api")
	{
		v1.GET("/home", func(ctx *gin.Context) {
			os := "Unknown"
			platform, _, version, err := host.PlatformInformation()
			if err == nil {
				os = platform + " " + version
			}
			ctx.JSON(200, gin.H{
				"message": "",
				"result":  gin.H{"os": os},
				"code":    200,
			})
		})

		v1.POST("/home", func(ctx *gin.Context) {
			var result public.Result
			if ctx.ShouldBind(&result) == nil {
				// TODO
			}
		})
	}

	router.Run(":8000")
	return nil
}
