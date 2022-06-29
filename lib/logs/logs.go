package logs

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/mizumoto-cn/goredis/lib/util/fileopt"
)

type LogConfig struct {
	// Path is the directory path of the log file
	Path string `json:"path"`
	// Name is the name of the log file
	Name string `json:"name"`
	// Ext is the extension of the log file
	Ext string `json:"ext"`
	// TimeFormat is the format of the timestamp in the log file
	TimeFormat string `json:"time_format"`
}

var (
	logFile            *os.File
	defaultPrefix      = ""
	defaultCallerDepth = 2
	logger             *log.Logger
	mu                 sync.Mutex
	logPrefix          = ""
	levelFlags         = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

type LogLevel uint16

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// flags is an alias of log.LstdFlags
// Ldate | Ltime
const flags = log.LstdFlags

func init() {
	logFile = os.Stdout
	logger = log.New(logFile, defaultPrefix, flags)
	// if fail panic
	// may be added later
}

// SetupLogger initializes logger
func SetupLogger(config LogConfig) {
	var err error
	dir := config.Path
	fileName := fmt.Sprintf("%s - %s.%s", config.Name, time.Now().Format(config.TimeFormat), config.Ext)
	// Todo: extract open file to a unique function
	// logFile, err = os.OpenFile(filepath.Join(dir, fileName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	logFile, err = fileopt.SafeOpen(fileName, dir)
	if err != nil {
		log.Fatal("logger setup failed: ", err)
	}

	// io.MultiWriter is used
	// to write to both stdout and the log file
	logger = log.New(io.MultiWriter(os.Stdout, logFile),
		defaultPrefix, flags)
}
