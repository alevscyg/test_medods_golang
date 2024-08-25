package main

import (
	"flag"
	"fmt"
	"log"
	"medos/interlan/app/apiserver"
	"medos/interlan/app/apiserver/config"

	"github.com/BurntSushi/toml"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

func main() {
	flag.Parse()

	config := config.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("server started")
	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
