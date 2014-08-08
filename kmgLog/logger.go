package kmgLog

import "fmt"
import (
	"github.com/bronze1man/kmg/console/kmgContext"
	"github.com/bronze1man/kmg/encoding/kmgJson"
	"github.com/bronze1man/kmg/kmgFile"
	"path/filepath"
	"runtime/debug"
	"time"
)

type Logger struct {
}
type Priority int

const (
	LOG_ALERT Priority = iota
	LOG_CRITICAL
	LOG_ERROR
	LOG_WARNING
	LOG_INFO
	LOG_DEBUG
)

func (obj *Logger) Log(level Priority, message string) {
	fmt.Println(message)
}
func (obj *Logger) Debug(message string) {
	obj.Log(LOG_DEBUG, message)
}
func (obj *Logger) Info(message string) {
	obj.Log(LOG_INFO, message)
}
func (obj *Logger) Waring(message string) {
	obj.Log(LOG_WARNING, message)
}
func (obj *Logger) Error(message string) {
	obj.Log(LOG_ERROR, message)
}
func (obj *Logger) Critical(message string) {
	obj.Log(LOG_CRITICAL, message)
}
func (obj *Logger) Alert(message string) {
	obj.Log(LOG_ALERT, message)
}

func (obj *Logger) LogError(err error) {
	debug.PrintStack()
	obj.Error(err.Error())
}
func (obj *Logger) VarDump(v interface{}) {
	message := fmt.Sprintf("%#v", v)
	obj.Log(LOG_DEBUG, message)
}

func init() {
	if kmgContext.Default != nil {
		kmgFile.Mkdir(kmgContext.Default.LogPath)
	}
}

type logRow struct {
	Time time.Time
	Msg  string
	Obj  interface{}
}

func Log(category string, msg string, obj interface{}) {
	logPath := kmgContext.Default.LogPath
	toWrite := append(kmgJson.MustMarshal(logRow{
		Time: time.Now(),
		Msg:  msg,
		Obj:  obj,
	}),byte('\n'))
	err := kmgFile.AppendFile(filepath.Join(logPath, category+".log"), toWrite)
	if err!=nil{
		panic(err)
	}
	return
}
