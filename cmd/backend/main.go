package main

import (
	"cronspy/backend/pkg/api"
	"cronspy/backend/pkg/util/config"
	"flag"
	"path"

	"github.com/astropay/go-tools/files"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {

	defaultConfigFile := path.Join(files.GetAppPath(), "config.yml")

	cfgPath := flag.String("config", defaultConfigFile, "path to configuration file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	checkErr(err)

	// start api server
	checkErr(api.Start(cfg))
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
