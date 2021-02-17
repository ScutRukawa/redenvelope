package infra

import "github.com/sirupsen/logrus"

var apiInitializerRegister *InitializeRegister = new(InitializeRegister)

//注册WEB API初始化对象

func RegisterApi(ai Initializer) {
	apiInitializerRegister.Register(ai)
}

//获取注册的web api初始化对象
func GetAPIInitializers() []Initializer {
	return apiInitializerRegister.Initializers
}

type WebAPIStarter struct {
	BaseStarter
}

func (w *WebAPIStarter) Setup(ctx StarterContext) {
	logrus.Info("WebAPIStarter Setup1")
	for _, v := range GetAPIInitializers() {
		logrus.Info("WebAPIStarter Setup2")
		v.Init()
	}
}
