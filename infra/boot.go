package infra

import (
	"github.com/tietang/props/kvs"
)

const (
	keyProps = "_conf"
)

// BootApplication app
type BootApplication struct {
	conf           kvs.ConfigSource
	starterContext StarterContext
}

func New(conf kvs.ConfigSource) *BootApplication {
	b := &BootApplication{conf: conf, starterContext: StarterContext{}}
	b.starterContext[keyProps] = conf
	return b
}
func (b *BootApplication) Start() {
	//1,资源初始化
	b.init()
	//2,setup
	b.setup()
	//3,启动
	b.start()

}

//Init x
func (b *BootApplication) init() {
	starters := StarterRegister.AllStarters()
	for _, starter := range starters {
		starter.Init(b.starterContext)
	}
}

func (b *BootApplication) setup() {
	starters := StarterRegister.AllStarters()
	for _, starter := range starters {
		starter.Setup(b.starterContext)
	}
}

func (b *BootApplication) start() {
	starters := StarterRegister.AllStarters()
	for index, starter := range starters {
		if starter.StartBlocking() {
			if index+1 == len(starters) {
				starter.Start(b.starterContext)
			} else {
				go starter.Start(b.starterContext)
			}
		} else {
			starter.Start(b.starterContext)
		}
	}
}
