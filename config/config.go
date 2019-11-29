package config

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

type Config struct {
	//HTTP
	HttpServerWin   string `ini:"http_server_win"`
	HttpServerLinux string `ini:"http_server_linux"`
	//files
	DirWin   string `ini:"dir_win"`
	DirLinux string `ini:"dir_linux"`
}

func (c Config) String() string {

	http := fmt.Sprintf("HTTP:[%v]/[%v]", c.HttpServerWin, c.HttpServerLinux)

	log := fmt.Sprintf("Dir:[win:%v]/[linux:%v]", c.DirWin, c.DirLinux)

	return http + ", " + log
}

//Read Server's Config Value from "path"
func ReadConfig(path string) (Config, error) {
	var config Config
	conf, err := ini.Load(path)
	if err != nil {
		log.Println("load config file fail!")
		return config, err
	}
	conf.BlockMode = false
	err = conf.MapTo(&config)
	if err != nil {
		log.Println("mapto config file fail!")
		return config, err
	}
	return config, nil
}
