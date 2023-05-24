package config

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

type Config struct {
	//HTTP
	HttpServer string `ini:"http_server"`
	//files
	Dir string `ini:"dir"`
}

func (c Config) String() string {

	http := fmt.Sprintf("HTTP:[%v]", c.HttpServer)

	log := fmt.Sprintf("Dir:[%v]", c.Dir)

	return http + ", " + log
}

// Read Server's Config Value from "path"
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
