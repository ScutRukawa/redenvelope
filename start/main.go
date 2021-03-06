package main

import (
	"fmt"
	_ "redenvelope/apis/web" //webapi init
	_ "redenvelope/core/envelopes"
	"redenvelope/infra"
	"redenvelope/infra/base"
	"redenvelope/jobs"

	"github.com/sirupsen/logrus"
	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
)

var conf *ini.IniFileConfigSource

func init() {
	file := kvs.GetCurrentFilePath("../config/config.ini", 1)
	fmt.Println("filepath:", file)
	conf = ini.NewIniFileConfigSource(file)

	infra.Register(&base.PropsStarter{})
	infra.Register(&base.DbxDatabaseStarter{})
	infra.Register(&base.ValidatorStarter{})
	infra.Register(&jobs.RefundExpiredJobStarter{})
	infra.Register(&base.IrisServerStarter{})
	infra.Register(&infra.WebAPIStarter{})

}
func main() {

	app := infra.New(conf)
	logrus.Info("app:", app)
	app.Start()
}
