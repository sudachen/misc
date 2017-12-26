package out

import (
	"fmt"
	"os"
	"io"
)

type Level byte

const (
	Error Level = iota
	Warn
	Info
	Debug
	Trace
)

var defaultOutput = os.Stderr
var output = []io.Writer{}
var currentLevel = Info

func (v Level) SetOutput(wr io.Writer) {
	output[v] = wr
}

func (v Level) SetCurrent() {
	currentLevel = v
}

func (v Level) Visible() bool {
	return v <= currentLevel
}

func (v Level) Output() io.Writer {
	wr := output[v]
	if wr == nil { return defaultOutput }
	return wr
}

func (v Level) Println(a ...interface{}) {
	if v.Visible() {
		fmt.Fprintln(v.Output(), a...)
	}
}

func (v Level) Print(a ...interface{}) {
	if v.Visible() {
		fmt.Fprint(v.Output(), a...)
	}
}

func (v Level) Printf(t string, a ...interface{}) {
	if v.Visible() {
		fmt.Fprintf(v.Output(), t, a...)
	}
}
