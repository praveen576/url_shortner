package i13n

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/traveltriangle/url_shortner/config"
	"runtime"
	"strings"
)

type LogContext interface {
	ToFields() map[string]interface{}
}

type LogFields struct {
	Fields map[string]interface{}
}

func (l LogFields) ToFields() map[string]interface{} {
	return l.Fields
}

func NewLogger() log.StdLogger {
	return log.New()
}

func Info(msg string, ctx ... LogContext) {
	prepare(&ctx).Info(msg)
}

func Error(err error, ctx ... LogContext) {
	if err != nil {
		prepare(&ctx).Error(err)
	}
}

func Fatal(err error, ctx ... LogContext) {
	if err != nil {
		log.Error(strings.Join(getTrace(), "\n"))
		prepare(&ctx).Fatal(err)
	}
}

func initLogger(){
	l, e := log.ParseLevel(config.Config.Logger.Level)
	Fatal(e)
	log.SetLevel(l)
}

func prepare(ctx *[]LogContext) *log.Entry {
	e := log.WithFields(log.Fields{})
	for _, c := range *ctx {
		e = e.WithFields(c.ToFields())
	}
	return e
}

func getTrace() []string {
	stack := make([]uintptr, 0, 50)
	runtime.Callers(3, stack[:])
	frames := runtime.CallersFrames(stack)
	trace := make([]string, 0, 50)
	for true {
		f, m := frames.Next()
		if f.Line != 0 {
			trace = append(trace, fmt.Sprintf("%s at %s:%d", f.Function, f.File, f.Line))
		}
		if m == false {
			break
		}
	}
	return trace
}
