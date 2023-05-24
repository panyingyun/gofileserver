package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"gofileserver/config"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func StartHttpServer(addr string, ftpsdir string) {
	r := gin.Default()
	r.Use(Cors())
	r.StaticFS("/", http.Dir(ftpsdir))

	//=================Mulit file upload=====================
	//usages:
	//curl -X POST http://localhost:4000/upload -F "file=@/e/gio/LICENSE" -H "Content-Type: multipart/form-data"
	//=======================================================
	r.POST("/upload", func(c *gin.Context) {
		// single file
		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}
		log.Println(file.Filename)
		// Upload the file to specific dst.
		err = c.SaveUploadedFile(file, ftpsdir+string(os.PathSeparator)+file.Filename)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})

	//=================Mulit file upload=====================
	//usages:
	//curl -X POST http://localhost:4000/uploads -F "files=@/e/gio/go.mod" -F "files=@/e/gio/go.sum" -H "Content-Type: multipart/form-data"
	//=======================================================
	r.POST("/uploads", func(c *gin.Context) {
		// Multipart form
		form, err := c.MultipartForm()
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}
		files := form.File["files"]
		for _, file := range files {
			log.Println(file.Filename)
			c.SaveUploadedFile(file, ftpsdir+string(os.PathSeparator)+file.Filename)
			if err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
				return
			}
		}
		c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})

	r.Run(addr)
}

func run(c *cli.Context) error {
	conf, err := config.ReadConfig(c.String("conf"))
	if err != nil {
		log.Error("read from conf fail!", c.String("conf"))
		return err
	}
	fmt.Println("conf =  ", conf)

	fmt.Println("runtime.GOOS = ", runtime.GOOS)

	//start http server
	go func() {
		StartHttpServer(conf.HttpServer, conf.Dir)
	}()

	//quit when receive end signal
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	log.Infof("signal received signal %v", <-sigChan)
	log.Warn("shutting down server")
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "filesserver"
	app.Usage = "Server for simple file upload and download"
	app.Copyright = "panyingyun@gmail.com "
	app.Version = "1.0"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "conf,c",
			Usage:  "Set conf path here",
			Value:  "files.conf",
			EnvVar: "FILES_CONF",
		},
	}
	app.Run(os.Args)
}
