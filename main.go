package redenvelope

import (
	"fmt"
	"redenvelope/infra"

	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
)

func main() {
	file := kvs.GetCurrentFilePath("../config/config.ini", 1)
	fmt.Println("filepath:", file)
	conf := ini.NewIniFileConfigSource(file)
	app := infra.New(conf)
	app.Start()
}
