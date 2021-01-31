package base

import (
	"fmt"
	"redenvelope/infra"

	"github.com/tietang/props/kvs"
)

var props kvs.ConfigSource

//Propos x
func Propos() kvs.ConfigSource {
	return props
}

//PropsStarter x
type PropsStarter struct {
	infra.BaseStarter
}

func (p *PropsStarter) Init(ctx infra.StarterContext) {
	props = ctx.Props()
	fmt.Println("props 初始化")
}
