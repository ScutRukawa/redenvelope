package testx

import (
	"fmt"
	"redenvelope/infra"
	"redenvelope/infra/base"

	"github.com/sirupsen/logrus"
	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
)

func init() {
	file := kvs.GetCurrentFilePath("../config/config.ini", 1)
	fmt.Println("filepath:", file)
	conf := ini.NewIniFileConfigSource(file)

	infra.Register(&base.PropsStarter{})
	infra.Register(&base.DbxDatabaseStarter{})
	infra.Register(&base.ValidatorStarter{})
	//infra.Register(&base.IrisServerStarter{})

	app := infra.New(conf)
	logrus.Info("app:", app)
	app.Start()
}
