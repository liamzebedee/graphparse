package graphparse

import (
	"runtime"
	"strings"
	"path/filepath"
	"log"
)

var ParserLog = log.New(LogWriter{}, "Parser: ", 0)


// https://wycd.net/posts/2014-07-02-logging-function-names-in-go.html


type LogWriter struct {

}

func (f LogWriter) Write(p []byte) (n int, err error) {
	pc, _, _, ok := runtime.Caller(4)
	if !ok {
		// file = "?"
		// line = 0
	}

	fn := runtime.FuncForPC(pc)
	var fnName string
	if fn == nil {
		fnName = "?()"
	} else {
		dotName := filepath.Ext(fn.Name())
		fnName = strings.TrimLeft(dotName, ".") + "()"
	}

	log.Printf("%s: %s", fnName, p)
	return len(p), nil
}