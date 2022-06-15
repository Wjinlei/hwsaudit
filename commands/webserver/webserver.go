package webserver

import (
	"encoding/json"
	"flag"
	"fmt"
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

	router.GET("/service", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", nil)
	})

	v1 := router.Group("/api")
	{
		v1.GET("/os", func(ctx *gin.Context) {
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

		v1.GET("/home", func(ctx *gin.Context) {
			jsonFile := "home.json"
			page := Page{}

			if err := ctx.ShouldBindQuery(&page); err != nil {
				ctx.AbortWithStatusJSON(200,
					gin.H{
						"message": err.Error(),
						"result":  []string{},
						"code":    1,
					})
				return
			}

			file, err := os.Open(jsonFile)
			if err != nil {
				ctx.AbortWithStatusJSON(200,
					gin.H{
						"message": "",
						"result":  []string{},
						"code":    1,
					})
				return
			}
			defer file.Close()

			begin := page.PageSize*page.PageNo - page.PageSize

			jsonStrResults, err := golib.ReadLinesOffsetN(jsonFile, uint(begin), page.PageSize, "\n")
			if err != nil {
				ctx.AbortWithStatusJSON(200,
					gin.H{
						"message": "Read result failed, please restart audit.",
						"result":  []string{},
						"code":    0,
					})
				return
			}

			var jsonResults []public.Result
			for _, jsonStrResult := range jsonStrResults {
				var jsonResult public.Result
				if err := json.Unmarshal([]byte(jsonStrResult), &jsonResult); err == nil {
					jsonResults = append(jsonResults, jsonResult)
				}
			}

			total, _ := golib.LineCounter(file)

			ctx.JSON(200, gin.H{
				"message": "",
				"result": gin.H{
					"data":       jsonResults,
					"pageNo":     page.PageNo,
					"totalCount": total,
				},
				"code": 0,
			})
		})

		v1.POST("/home", func(ctx *gin.Context) {
			var result public.Result

			if err := ctx.ShouldBind(&result); err != nil {
				ctx.AbortWithStatusJSON(200,
					gin.H{
						"message": "ShouldBind error",
						"result":  []string{},
						"code":    1,
					})
				return
			}

			lstat, err := os.Lstat(result.Path)
			if err != nil {
				ctx.AbortWithStatusJSON(200,
					gin.H{
						"message": err.Error(),
						"result":  []string{},
						"code":    1,
					})
				return
			}
			if lstat.IsDir() == false {
				ctx.AbortWithStatusJSON(200,
					gin.H{
						"message": "Path not directory",
						"result":  []string{},
						"code":    1,
					})
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
		})

		v1.GET("/service", func(ctx *gin.Context) {
			jsonFile := "service.json"
			page := Page{}

			if err := ctx.ShouldBindQuery(&page); err != nil {
				ctx.AbortWithStatusJSON(200,
					gin.H{
						"message": err.Error(),
						"result":  []string{},
						"code":    1,
					})
				return
			}

			file, err := os.Open(jsonFile)
			if err != nil {
				ctx.AbortWithStatusJSON(200,
					gin.H{
						"message": "",
						"result":  []string{},
						"code":    1,
					})
				return
			}
			defer file.Close()

			begin := page.PageSize*page.PageNo - page.PageSize

			jsonStrResults, err := golib.ReadLinesOffsetN(jsonFile, uint(begin), page.PageSize, "\n")
			if err != nil {
				ctx.AbortWithStatusJSON(200,
					gin.H{
						"message": "Read result failed, please restart audit.",
						"result":  []string{},
						"code":    0,
					})
				return
			}

			var jsonResults []public.Unit
			for _, jsonStrResult := range jsonStrResults {
				var jsonResult public.Unit
				if err := json.Unmarshal([]byte(jsonStrResult), &jsonResult); err == nil {
					jsonResults = append(jsonResults, jsonResult)
				}
			}

			total, _ := golib.LineCounter(file)

			ctx.JSON(200, gin.H{
				"message": "",
				"result": gin.H{
					"data":       jsonResults,
					"pageNo":     page.PageNo,
					"totalCount": total,
				},
				"code": 0,
			})
		})

		v1.POST("/service", func(ctx *gin.Context) {
			textFile := "service.xml"
			jsonFile := "service.json"
			var requestUnit public.Unit

			if err := ctx.ShouldBind(&requestUnit); err != nil {
				ctx.AbortWithStatusJSON(200,
					gin.H{
						"message": "ShouldBind error",
						"result":  []string{},
						"code":    1,
					})
				return
			}

			unitList, err := public.ListUnits(requestUnit.FormState)
			if err != nil {
				ctx.AbortWithStatusJSON(200,
					gin.H{
						"message": err.Error(),
						"result":  []string{},
						"code":    1,
					})
				return
			}

			golib.Delete(textFile)
			golib.Delete(jsonFile)
			for _, unit := range unitList {
				jsonUnit, _ := json.Marshal(unit)
				golib.FileWrite(jsonFile, string(jsonUnit)+"\n", golib.FileAppend)
				// 产品经理让加的傻逼导出格式
				golib.FileWrite(textFile, fmt.Sprintf(
					"<service>%s{*}%s{*}%s</service>\n", unit.Name, unit.Description, unit.Path),
					golib.FileAppend)
			}

			ctx.JSON(200, gin.H{"message": "", "result": "ok", "code": 0})
		})

		v1.POST("/export_xml_home", func(ctx *gin.Context) {
			ctx.FileAttachment("home.xml", "home.xml")
		})

		v1.POST("/export_xml_service", func(ctx *gin.Context) {
			ctx.FileAttachment("service.xml", "service.xml")
		})

		v1.POST("/export_json_home", func(ctx *gin.Context) {
			ctx.FileAttachment("home.json", "home.json")
		})

		v1.POST("/export_json_service", func(ctx *gin.Context) {
			ctx.FileAttachment("service.json", "service.json")
		})
	}

	router.Run(":8000")
	return nil
}
