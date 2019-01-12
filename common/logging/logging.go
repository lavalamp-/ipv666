package logging

import (
	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"log"
)

func SetupLogging() {
	if viper.GetBool("LogToFile") {
		log.SetFlags(log.Flags() & (log.Ldate | log.Ltime))
		log.SetOutput(&lumberjack.Logger{
			Filename:   viper.GetString("LogFilePath"),
			MaxSize:    viper.GetInt("LogFileMBSize"),		// megabytes
			MaxBackups: viper.GetInt("LogFileMaxBackups"),
			MaxAge:     viper.GetInt("LogFileMaxAge"),		// days
			Compress:   viper.GetBool("CompressLogFiles"),
		})
	}
}
