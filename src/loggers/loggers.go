package loggers

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"sync"
)

type loggerSingleton struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}

var instance *loggerSingleton
var lock = &sync.Mutex{}

func GetInstance() *loggerSingleton {
	if instance != nil {
		return instance
	}

	lock.Lock()
	defer lock.Unlock()

	if instance == nil {
		errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
		infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

		instance = &loggerSingleton{
			ErrorLog: errorLog,
			InfoLog:  infoLog,
		}
	}

	return instance
}

func LogError(err error) {
	stackTrace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	instance.ErrorLog.Println(stackTrace)
}
