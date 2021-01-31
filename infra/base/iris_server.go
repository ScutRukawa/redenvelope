package base

import (
	"bytes"
	"redenvelope/infra"
	"time"

	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	irisrecover "github.com/kataras/iris/v12/middleware/recover"
	"github.com/sirupsen/logrus"
)

var app *iris.Application

func Iris() *iris.Application {
	return app
}

type IrisServerStarter struct {
	infra.BaseStarter
}

func (i *IrisServerStarter) Init(ctx infra.StarterContext) {
	//创建app
	app = iris.New()

	//recover 中间件
	app.Use(irisrecover.New())

	//日志中间件
	cfg := logger.Config{
		Status:             true,
		IP:                 true,
		Method:             true,
		Path:               true,
		Query:              true,
		MessageContextKeys: []string{"logger_message"},
		MessageHeaderKeys:  []string{"User-Agent"},
		LogFunc: func(now time.Time, latency time.Duration,
			status, ip, method, path string,
			message interface{},
			headerMessage interface{}) {
			app.Logger().Infof("| %s | %s | %s | %s | %s | %s | %+v | %+v",
				now.Format("2006-01-02.15:04:05.000000"),
				latency.String(), status, ip, method, path, headerMessage, message,
			)
		},
	}
	app.Use(logger.New(cfg))

	//日志输出组件扩展
	logger := app.Logger()
	logger.Install(logrus.StandardLogger())

	logrus.Info("IrisServerStarter init")

}

func (i *IrisServerStarter) Start(ctx infra.StarterContext) {
	logLevel := ctx.Props().GetDefault("log.level", "info")
	logrus.Info(logLevel)
	Iris().Logger().SetLevel(logLevel)

	routes := Iris().GetRoutes()
	logrus.Info(len(routes))
	var routeInfo bytes.Buffer
	//var routeInfo2 string
	for index, route := range routes {
		route.Trace(&routeInfo, 2)
		logrus.Info(&routeInfo)
		logrus.Info("index:", index)
	}

	//启动
	port := ctx.Props().GetDefault("app.server.port", "8080")
	Iris().Run(iris.Addr(":" + port))
}
