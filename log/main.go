package log

import (
	"fmt"
	"kahf/conf"
	"log"
	"os"
	"strings"
	"time"
)

type Style int

const (
	StyleDefault Style = 0 + iota
	StyleLight
	StyleUnderline Style = 4 + iota
	StyleFlash
	StyleInverse Style = 7 + iota
	StyleHide
)

type Color int

const (
	ForeBlack Color = 30 + iota
	ForeRed
	ForeGreen
	ForeYellow
	ForeBlue
	ForePurple
	ForeLightBlue
	ForeWhite
	BackBlack Color = 40 + iota
	BackRed
	BackGreen
	BackYellow
	BackBlue
	BackPurple
	BackLightBlue
	BackWhite
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	FAIL
	ERROR
)

type ILog interface {
	log(level Level, messages []string)
	print(msg string)
}

type Log struct {
	Level Level
	Debug ILog
	Info  ILog
	Warn  ILog
	Fail  ILog
	Error ILog
}

type Messgae struct {
	msg string
	Level
}

var msgChan chan *Messgae
var _conf Log

func init() {
	msgChan = make(chan *Messgae, conf.Setting.Log.BufferSize)
	level := TranslateToLevel(conf.Setting.Log.Level)
	switch strings.ToUpper(conf.Setting.Log.Mode) {
	case "FILE":
		var curTime = time.Now()
		_conf = Log{
			Level: level,
			Debug: &File{curTime.Day(), *createLevelLogger(DEBUG, curTime), DEBUG},
			Info:  &File{curTime.Day(), *createLevelLogger(INFO, curTime), INFO},
			Warn:  &File{curTime.Day(), *createLevelLogger(WARN, curTime), WARN},
			Fail:  &File{curTime.Day(), *createLevelLogger(FAIL, curTime), FAIL},
			Error: &File{curTime.Day(), *createLevelLogger(ERROR, curTime), ERROR},
		}
		dealFatal()
	default:
		cliLogger := log.New(os.Stdout, "", log.Ltime)
		_conf = Log{
			Level: level,
			Debug: &Console{StyleFlash, *cliLogger},
			Info:  &Console{StyleFlash, *cliLogger},
			Warn:  &Console{StyleFlash, *cliLogger},
			Fail:  &Console{StyleFlash, *cliLogger},
			Error: &Console{StyleFlash, *cliLogger},
		}
	}
	go OutputMsg()
}

func Debug(msg ...string) {
	if _conf.Level > DEBUG {
		return
	}
	_conf.Debug.log(DEBUG, msg)
}

func Debugf(format string, a ...interface{}) {
	Debug(fmt.Sprintf(format, a...))
}

func Info(msg ...string) {
	if _conf.Level > INFO {
		return
	}
	_conf.Info.log(INFO, msg)
}
func Infof(format string, a ...interface{}) {
	Info(fmt.Sprintf(format, a...))
}

func Warn(msg ...string) {
	if _conf.Level > WARN {
		return
	}
	_conf.Warn.log(WARN, msg)
}
func Warnf(format string, a ...interface{}) {
	Warn(fmt.Sprintf(format, a...))
}

func Fail(msg ...string) {
	if _conf.Level > FAIL {
		return
	}
	_conf.Fail.log(FAIL, msg)
}
func Failf(format string, a ...interface{}) {
	Fail(fmt.Sprintf(format, a...))
}

func Error(msg ...string) {
	if _conf.Level > ERROR {
		return
	}
	_conf.Error.log(ERROR, msg)
}
func Errorf(format string, a ...interface{}) {
	Error(fmt.Sprintf(format, a...))
}

func (lv Level) ToString(flag string) string {
	rs := ""
	switch lv {
	case DEBUG:
		rs = "DEBUG"
		break
	case INFO:
		rs = "INFO"
		break
	case WARN:
		rs = "WARN"
		break
	case FAIL:
		rs = "FAIL"
		break
	case ERROR:
		rs = "ERROR"
		break
	}

	switch flag {
	case "<>":
		rs = fmt.Sprintf("<%s>", rs)
		break
	case "[]":
		rs = fmt.Sprintf("[%s]", rs)
		break
	case "#":
		rs = fmt.Sprintf("#%s#", rs)
		break
	}

	return rs
}

func (lv Level) ToBackColor() Color {
	switch lv {
	case DEBUG:
		return BackGreen
	case INFO:
		return BackLightBlue
	case WARN:
		return BackYellow
	case FAIL:
		return BackPurple
	case ERROR:
		return BackRed
	}
	return 0
}

func (lv Level) ToForeColor() Color {
	switch lv {
	case DEBUG:
		return ForeGreen
	case INFO:
		return ForeLightBlue
	case WARN:
		return ForeYellow
	case FAIL:
		return ForePurple
	case ERROR:
		return ForeRed
	}
	return 0
}

func TranslateToLevel(l string) Level {
	switch strings.ToUpper(l) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN":
		return WARN
	case "FAIL":
		return FAIL
	case "ERROR":
		return ERROR
	}
	return DEBUG
}

func OutputMsg() {
	var msg *Messgae
	for {
		select {
		case msg = <-msgChan:
			switch msg.Level {
			case DEBUG:
				_conf.Debug.print(msg.msg)
			case INFO:
				_conf.Info.print(msg.msg)
			case WARN:
				_conf.Warn.print(msg.msg)
			case FAIL:
				_conf.Fail.print(msg.msg)
			case ERROR:
				_conf.Error.print(msg.msg)
			}
		}
	}
}

func createLevelLogger(level Level, t time.Time) *log.Logger {
	if level >= TranslateToLevel(conf.Setting.Log.Level) {
		file, err := openLogFile(level, t)
		if err != nil {
			fmt.Printf("create level logger failed! %s/n", err)
			os.Exit(1)
		}
		fileLogger := log.New(file, "file", log.LstdFlags)
		return fileLogger
	}
	return &log.Logger{}
}

func dealFatal() {
	/*if !tools.Exist("logs") {
		err := os.Mkdir("logs", 0777)
		if err != nil {
			fmt.Println("deal fatal failed! ", err)
			os.Exit(1)
		}
	}
	logFile, err := os.OpenFile("./logs/fatal.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("deal fatal failed! ", err)
		os.Exit(1)
	}
	syscall.Dup2(int(logFile.Fd()), int(os.Stderr.Fd()))*/
}
