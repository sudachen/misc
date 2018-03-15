package out

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

type Level int

const StdErr Level = -2
const StdOut Level = -1

const (
	Crit Level = iota
	Error
	Warn
	Info
	Verbose
	Debug
	Trace
	levelsCount
)

var DefaultWriter io.Writer = os.Stderr
var PrintFunction func(Level, []byte) = DefaultPrintFunction
var PrefixFunction func(Level, *bytes.Buffer) = DefaultPrefixFunction

var writer = make([]io.Writer, levelsCount)
var prefix = make([][]byte, levelsCount)
var currentLevel = Info

func init() {
	Error.SetPrefix("error")
	Warn.SetPrefix("warn")
	Debug.SetPrefix("debug")
	Trace.SetPrefix("trace")
}

func (lvl Level) String() string {
	switch lvl {
	case Error:
		return "Error"
	case Warn:
		return "Warn"
	case Info:
		return "Info"
	case Verbose:
		return "Verbose"
	case Debug:
		return "Debug"
	case Trace:
		return "Trace"
	}
	return "-"
}

func (lvl Level) SetWriter(wr io.Writer) {
	writer[lvl] = wr
}

func (lvl Level) SetCurrent() {
	currentLevel = lvl
}

func (lvl Level) Visible() bool {
	return lvl <= currentLevel
}

func (lvl Level) Writer() io.Writer {
	if lvl < 0 {
		switch lvl {
		case StdErr:
			return os.Stderr
		case StdOut:
			return os.Stdout
		}
		return DefaultWriter
	}
	wr := writer[lvl]
	if wr == nil {
		return DefaultWriter
	}
	return wr
}

func (lvl Level) Prefix() []byte {
	if lvl < 0 {
		return nil
	}
	return prefix[lvl]
}

func (lvl Level) FullPrefixString() string {
	var bf bytes.Buffer
	PrefixFunction(lvl, &bf)
	return bf.String()
}

func (lvl Level) SetPrefix(pfx string) {
	switch len(pfx) {
	case 0:
		prefix[lvl] = nil
	default:
		prefix[lvl] = []byte(pfx)
	}
}

var endl = []byte{'\n'}
var pfxSep = []byte{':', ' '}

func DefaultPrefixFunction(lvl Level, bf *bytes.Buffer) {
	if pfx := lvl.Prefix(); pfx != nil {
		bf.Write(pfx)
		bf.Write(pfxSep)
	}
}

func DefaultPrintFunction(lvl Level, b []byte) {
	lvl.Writer().Write(b)
}

func (lvl Level) Print(a ...interface{}) {
	if lvl.Visible() {
		var textBuf bytes.Buffer
		PrefixFunction(lvl, &textBuf)
		fmt.Fprint(&textBuf, a...)
		if l := textBuf.Len(); l == 0 || textBuf.Bytes()[l-1] != '\n' {
			textBuf.Write(endl)
		}
		PrintFunction(lvl, textBuf.Bytes())
	}
}

func (lvl Level) Printf(t string, a ...interface{}) {
	if lvl.Visible() {
		var textBuf bytes.Buffer
		PrefixFunction(lvl, &textBuf)
		fmt.Fprintf(&textBuf, t, a...)
		if l := textBuf.Len(); l == 0 || textBuf.Bytes()[l-1] != '\n' {
			textBuf.Write(endl)
		}
		PrintFunction(lvl, textBuf.Bytes())
	}
}

func Fatalf(t string, a ...interface{}) {
	if len(t) < 1 || t[len(t)-1] != '\n' {
		t = t + "\n"
	}
	Error.Printf(t, a...)
	os.Exit(255)
}
