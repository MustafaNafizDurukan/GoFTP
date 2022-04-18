package logs

import (
	"io"
	"log"
	"os"
)

var (
	INFO    *log.Logger
	WARNING *log.Logger
	ERROR   *log.Logger
	CONTROL *log.Logger
)

func InitLog(
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer,
	controlHandler io.Writer) {

	INFO = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime)

	WARNING = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	ERROR = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	CONTROL = log.New(controlHandler,
		"CONTROL: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func Set() *os.File {
	logError, err := os.OpenFile("logError.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Fail to open log file")
	}

	logWarning, err := os.OpenFile("logWarning.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Fail to open log file")
	}

	logInfo, err := os.OpenFile("logInfo.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Fail to open log file")
	}

	logControl, err := os.OpenFile("logControl.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Fail to open log file")
	}

	InitLog(logInfo, logWarning, logError, logControl)

	return logInfo
}
