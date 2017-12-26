package out

import (
	"fmt"
	"os"
	"io"
	"bytes"
)

type Level byte

const (
	Error Level = iota
	Warn
	Info
	Verbose
	Debug
	Trace
	levelsCount
)

var DefaultWriter io.Writer = os.Stderr
var PrintFunction func(lvl Level, b []byte) = DefaultPrintFunction
var PrefixFunction func(lvl Level, bf* bytes.Buffer) = DefaultPrefixFunction

var writer = make([]io.Writer,levelsCount)
var prefix = make([][]byte,levelsCount)
var currentLevel = Info

func init() {
	Error.SetPrefix("error")
	Warn.SetPrefix("warn")
	Debug.SetPrefix("debug")
	Trace.SetPrefix("trace")
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
	wr := writer[lvl]
	if wr == nil { return DefaultWriter }
	return wr
}

func (lvl Level) Prefix() []byte {
	return prefix[lvl]
}

func (lvl Level) SetPrefix(pfx string) {
	prefix[lvl] = []byte(pfx)
}

var endl = []byte{'\n'}
var pfxSep = []byte{':',' '}
var textBuf bytes.Buffer

func DefaultPrefixFunction(lvl Level, bf* bytes.Buffer) {
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
		textBuf.Reset()
		PrefixFunction(lvl,&textBuf)
		fmt.Fprint(&textBuf, a...)
		if l := textBuf.Len(); l == 0 || textBuf.Bytes()[l-1] != '\n' {
			textBuf.Write(endl)
		}
		PrintFunction(lvl,textBuf.Bytes())
	}
}

func (lvl Level) Printf(t string, a ...interface{}) {
	if lvl.Visible() {
		textBuf.Reset()
		PrefixFunction(lvl,&textBuf)
		fmt.Fprintf(&textBuf, t, a...)
		if l := textBuf.Len(); l == 0 || textBuf.Bytes()[l-1] != '\n' {
			textBuf.Write(endl)
		}
		PrintFunction(lvl,textBuf.Bytes())
	}
}
