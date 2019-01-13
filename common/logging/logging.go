package logging

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

const (
	LEVEL_DEBUG		= iota
	LEVEL_INFO
	LEVEL_SUCCESS
	LEVEL_WARNING
	LEVEL_ERROR
)

var debugColor = color.New(color.FgHiWhite).SprintFunc()
var infoColor = color.New(color.FgHiBlue).SprintFunc()
var successColor = color.New(color.FgHiGreen).Add(color.Underline).SprintFunc()
var warnColor = color.New(color.FgYellow).SprintFunc()
var errorColor = color.New(color.FgBlack).Add(color.BgRed).Add(color.Underline).SprintFunc()

func getLogLevel() int {
	configLevel := viper.GetString("LogLevel")
	configLevel = strings.ToLower(configLevel)
	switch configLevel {
	case "debug":
		return LEVEL_DEBUG
	case "info":
		return LEVEL_INFO
	case "success":
		return LEVEL_SUCCESS
	case "warn":
		return LEVEL_WARNING
	case "error":
		return LEVEL_ERROR
	default:
		ErrorStringF(fmt.Sprintf("%s is not a valid log level", configLevel))
		return -1
	}
}

func printWithDate(toPrint string) {
	log.Printf("- %s", toPrint)
}

func Debug(toPrint string) {
	if getLogLevel() <= LEVEL_DEBUG {
		printWithDate(fmt.Sprintf("%s - %s", debugColor("DEB"), toPrint))
	}
}

func Debugf(toPrint string, a ...interface{}) {
	Debug(fmt.Sprintf(toPrint, a...))
}

func Info(toPrint string) {
	if getLogLevel() <= LEVEL_INFO {
		printWithDate(fmt.Sprintf("%s - %s", infoColor("INF"), toPrint))
	}
}

func Infof(toPrint string, a ...interface{}) {
	Info(fmt.Sprintf(toPrint, a...))
}

func Success(toPrint string) {
	if getLogLevel() <= LEVEL_SUCCESS {
		printWithDate(fmt.Sprintf("%s - %s", successColor("SUC"), toPrint))
	}
}

func Successf(toPrint string, a ...interface{}) {
	Success(fmt.Sprintf(toPrint, a...))
}

func Warn(toPrint string) {
	if getLogLevel() <= LEVEL_WARNING {
		printWithDate(fmt.Sprintf("%s - %s", warnColor("WAR"), toPrint))
	}
}

func Warnf(toPrint string, a ...interface{}) {
	Warn(fmt.Sprintf(toPrint, a...))
}

func Error(toPrint error) {
	ErrorString(toPrint.Error())
}

func ErrorString(toPrint string) {
	printWithDate(fmt.Sprintf("%s - %s", errorColor("ERR"), toPrint))
}

func ErrorF(toPrint error) {
	ErrorStringF(toPrint.Error())
}

func ErrorStringF(toPrint string) {
	ErrorString(toPrint)
	os.Exit(-1)
}

func ErrorStringFf(toPrint string, a ...interface{}) {
	ErrorStringF(fmt.Sprintf(toPrint, a...))
}

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
