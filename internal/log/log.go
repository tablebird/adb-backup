package log

import (
	"fmt"
	"log"
	"os"

	config "adb-backup/internal/config"
)

// ANSI 颜色代码
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

func DebugF(format string, v ...interface{}) {
	if config.App.DebugLog {
		log.Printf(format, v...)
	}
}

func Debug(v ...interface{}) {
	if config.App.DebugLog {
		log.Print(v...)
	}
}

func Info(format string) {
	log.Println(format)
}

func InfoF(format string, v ...interface{}) {
	log.Printf(format+"\n", v...)
}

func SuccessF(format string, v ...interface{}) {
	log.Printf(colorGreen+format+colorReset, v...)
}

func Warning(v string) {
	log.Println(colorYellow + v + colorReset)
}

func WarningF(format string, v ...interface{}) {
	log.Printf(colorYellow+format+colorReset, v...)
}

func ErrorF(format string, v ...interface{}) {
	log.Printf(colorRed+format+colorReset, v...)
}

func Fatal(v ...interface{}) {
	log.Println(colorRed, fmt.Sprint(v...), colorReset)
	os.Exit(1)
}

func FatalF(format string, v ...interface{}) {
	log.Printf(colorRed+format+colorReset, v...)
	os.Exit(1)
}
