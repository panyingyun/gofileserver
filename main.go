package main

import (
	"gofileserver/config"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func StartHttpServer(addr string, ftpsdir string) {
	r := gin.Default()
	r.StaticFS("/", http.Dir(ftpsdir))

	//=================Mulit file upload=====================
	//usages:
	//curl -X POST http://localhost:6001/upload -F "file=@/e/gio/LICENSE" -H "Content-Type: multipart/form-data"
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
	//curl -X POST http://localhost:6001/uploads -F "files=@/e/gio/go.mod" -F "files=@/e/gio/go.sum" -H "Content-Type: multipart/form-data"
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
		if runtime.GOOS == "windows" {
			StartHttpServer(conf.HttpServerWin, conf.DirWin)
		} else {
			StartHttpServer(conf.HttpServerLinux, conf.DirLinux)
		}
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
