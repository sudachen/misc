
package out

import (
	"testing"
	"bytes"
	"io/ioutil"
)

type In  []interface{}

type PrnInOut struct {
	In
	Out string
}

func equal(a []byte, b []byte) bool {
	l := len(a)
	if l != len(b) { return false }
	for i := range a {
		if a[i] != b[i] { return false }
	}
	return true
}

func diffIndex(a []byte, b []byte) int {
	l := len(a)
	if l > len(b) { l = len(b) }
	j := 0
	for ; j <  l; j++ {
		if a[j] != b[j] { break }
	}
	return j
}

var printDataset = []PrnInOut{
	{ In{"hello", "world!",}, "helloworld!\n" },
	{ In{"hello", " ", "world!",}, "hello world!\n" },
	{ In{1, 2.0,}, "1 2\n" },
	{ In{1, 3.1,}, "1 3.1\n" },
}

func checkForPrint(lvl Level, b *bytes.Buffer, t *testing.T) {
	for i, x := range printDataset {
		b.Reset()
		lvl.Print(x.In...)
		o := []byte(lvl.FullPrefixString() + x.Out)
		if bs := b.Bytes(); !equal(bs, o) {
			n := diffIndex(bs,o)
			t.Errorf(string(o))
			t.Errorf(string(bs))
			t.Fatalf("(%s) test set %d output does not match at %d, got %s",
				lvl, i, n, string(bs[:n])+"<#>"+string(bs[n:]))
		}
	}
}

var levels = []Level {Error, Warn, Info, Verbose, Debug, Trace}

func TestLevel_Print(t *testing.T) {
	var bf bytes.Buffer
	DefaultWriter = &bf
	Trace.SetCurrent()
	for _, lvl := range levels {
		checkForPrint(lvl, &bf, t)
	}
}

func TestLevel_Prefix(t *testing.T) {
	var bf bytes.Buffer
	DefaultWriter = &bf
	Trace.SetCurrent()
	for _, lvl := range levels {
		lvl.SetPrefix("-")
		if lvl.FullPrefixString() != "-: " {
			t.Fatal("bad prefix ", lvl.FullPrefixString())
		}
		checkForPrint(lvl, &bf, t)
		lvl.SetPrefix("")
		if lvl.FullPrefixString() != "" {
			t.Fatal("bad prefix ", lvl.FullPrefixString())
		}
		checkForPrint(lvl, &bf, t)
	}
}

func TestLevel_SetCurrent(t *testing.T) {
	var bf bytes.Buffer
	DefaultWriter = &bf
	for _, lvl := range levels[1:] {
		bf.Reset()
		Error.SetCurrent()
		lvl.Print("hello!")
		if len(bf.Bytes()) != 0 {
			t.Fatalf("(%s) failed to print with Error level")
		}
		lvl.SetCurrent()
		checkForPrint(lvl, &bf, t)
	}
}

func TestLevel_SetWriter(t *testing.T) {
	var bf bytes.Buffer
	DefaultWriter = ioutil.Discard
	Trace.SetCurrent()
	for _, lvl := range levels[1:] {
		bf.Reset()
		lvl.SetWriter(&bf)
		checkForPrint(lvl, &bf, t)
	}
}
