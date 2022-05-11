package webserver

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/Wjinlei/golib"
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
				"code":    0,
			})
		})

		v1.GET("/result", func(ctx *gin.Context) {
			page := GetReusltParams{}
			if err := ctx.ShouldBindQuery(&page); err != nil {
				ctx.AbortWithStatusJSON(200, gin.H{"message": "param error", "result": []string{}, "code": 1})
				return
			}

            // Open result.txt
			file, err := os.Open("result.txt")
			if err != nil {
				ctx.AbortWithStatusJSON(200, gin.H{"message": "", "result": []string{}, "code": 0})
				return
			}
			defer file.Close()

			// Get result.txt count
			total, _ := public.LineCounter(file)

			// Begin offset
			begin := page.PageSize*page.PageNo - page.PageSize

			// Read result from result.txt
			jsonStrResults, err := golib.ReadLinesOffsetN("result.txt", uint(begin), page.PageSize-1, "\n")
			if err != nil {
				ctx.AbortWithStatusJSON(200, gin.H{"message": "", "result": []string{}, "code": 0})
				return
			}

			// Convert string result to json
			var jsonResults []public.Result
			for _, jsonStrResult := range jsonStrResults {
				var jsonResult public.Result
				if err := json.Unmarshal([]byte(jsonStrResult), &jsonResult); err == nil {
					jsonResults = append(jsonResults, jsonResult)
				}
			}

			ctx.JSON(200, gin.H{
				"message": "",
				"result":  gin.H{"data": jsonResults, "pageNo": page.PageNo, "totalCount": total},
				"code":    0,
			})

		})

		v1.POST("/home", func(ctx *gin.Context) {
			var result public.Result
			if ctx.ShouldBind(&result) == nil {
				lstat, err := os.Lstat(result.Path)
				if err != nil {
					ctx.AbortWithStatusJSON(200, gin.H{"message": err.Error(), "result": []string{}, "code": 1})
					return
				}
				if lstat.IsDir() == false {
					ctx.AbortWithStatusJSON(200, gin.H{"message": "Path not directory", "result": []string{}, "code": 1})
					return
				}

				checkT := false  // check sticky
				checkS := false  // check sUid or sGid
				checkF := false  // check file
				checkD := false  // check directory
				target := "file" // check target

				// Convert params
				for _, param := range result.Other {
					switch param {
					case "t":
						checkT = true
					case "s":
						checkS = true
					case "checkFile":
						checkF = true
					case "checkDirectory":
						checkD = true
					}
				}

				// Convert checkF and checkD to target
				if checkF && checkD {
					target = "all"
				} else if checkF {
					target = "file"
				} else if checkD {
					target = "dir"
				}

				// Handler
				public.WalkDir(true, result.Path, target, result.User, result.Mode, checkS, checkT, result.Facl)

				ctx.JSON(200, gin.H{"message": "", "result": "ok", "code": 0})
			}
		})
	}

	router.Run(":8000")
	return nil
}
