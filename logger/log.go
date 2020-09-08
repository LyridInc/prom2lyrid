package logger

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"io"
	"os"
	"sync"
)

type LogManager struct {
	Logger    log.Logger
	LogWriter io.Writer
}

var instance *LogManager
var once sync.Once

func GetInstance() *LogManager {
	once.Do(func() {
		instance = &LogManager{}
	})
	return instance
}

func (manager *LogManager) Init() {
	logFile := os.Stdout
	if len(os.Getenv("LOG_DIR")) > 0 {
		dir := os.Getenv("LOG_DIR")
		_, err := os.Stat(dir)

		if os.IsNotExist(err) {
			errDir := os.MkdirAll(dir, 0755)
			if errDir != nil {
				fmt.Println("err: ", err)
			}

		}
		logFile, _ = os.OpenFile(dir+"/prom2lyrid.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	}
	manager.LogWriter = io.MultiWriter(os.Stdout, logFile)
	manager.Logger = log.NewSyncLogger(log.NewLogfmtLogger(manager.LogWriter))
	manager.Logger = log.With(manager.Logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
}
