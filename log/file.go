package log

import (
	"fmt"
	"github.com/Noctis0227/wallet-lib/tools"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type File struct {
	Day int
	log.Logger
	Level
}

func (f *File) log(level Level, messages []string) {
	var prefix = fmt.Sprintf("%-7s", level.ToString("[]"))
	f.SetPrefix(prefix)
	f.SetFlags(log.Lmicroseconds)
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "???"
		line = 0
	}
	arr := strings.Split(file, "/")
	var msg string
	for i, item := range messages {
		if i > 0 {
			msg += " "
		}
		msg += item
	}
	fLn := fmt.Sprintf("%s:%d:", arr[len(arr)-1], line)
	msg = fmt.Sprintf("%-25s %s", fLn, msg)
	msgChan <- &Messgae{strings.TrimRight(msg, "\n"), level}
}

func (f *File) print(msg string) {
	var curTime = time.Now()
	if curTime.Day() != f.Day || !isExistLogFile(f.Level, curTime) {
		file, err := openLogFile(f.Level, curTime)
		if err != nil {
			fmt.Printf("Open log file %s failed! %s\n", curTime.String(), err)
		} else {
			f.Logger = *log.New(file, "file", log.LstdFlags)
			f.Day = curTime.Day()
		}
	}
	f.Println(msg)
}

func openLogFile(level Level, curTime time.Time) (*os.File, error) {
	logDir := fmt.Sprintf("logs/%04d%02d%02d", curTime.Year(), curTime.Month(), curTime.Day())
	if !tools.Exist("logs") {
		err := os.Mkdir("logs", os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	if !tools.Exist(logDir) {
		err := os.Mkdir(logDir, 0777)
		if err != nil {
			return nil, err
		}
	}
	logName := logDir + "/" + level.ToString("") + ".log"
	return os.OpenFile(logName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
}

func isExistLogFile(level Level, curTime time.Time) bool {
	fileName := fmt.Sprintf("logs/%04d%02d%02d/%s.log", curTime.Year(), curTime.Month(), curTime.Day(), level.ToString(""))
	return tools.Exist(fileName)
}
