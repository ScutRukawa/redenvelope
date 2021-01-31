package base

import (
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func init() {
	// r1, _ := rotatelogs.New("access_log.%Y%M%d%H$%M")
	// log.SetOutput(r1)
	//日志格式
	formater := &prefixed.TextFormatter{}
	log.SetFormatter(formater)
	formater.FullTimestamp = true
	formater.TimestampFormat = "2006-01-02.15:04:05.000000"
	formater.ForceFormatting = true
	formater.SetColorScheme(&prefixed.ColorScheme{
		InfoLevelStyle: "green",
	})
	//日志级别
	// logLevel := os.Getenv("log.debug")
	// if logLevel == "true" {
	log.SetLevel(log.DebugLevel)
	// }
	//颜色样式
	formater.ForceColors = true
	formater.DisableColors = false
	//日志文件和滚动
	log.Info("测试日志")
	log.Debug("测试日志2")
	//github.com/lestrrat/go-file-rotatelogs
	log.Info("测试日志3")
}
