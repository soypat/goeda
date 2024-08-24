package kicad

import (
	_ "embed"
	"os"
	"strings"
	"testing"
)

//go:embed piaa.dsn
var dsnFile string

func TestLexerOnFile(t *testing.T) {
	fp := strings.NewReader(dsnFile)
	var l Lexer
	err := l.Reset(fp)
	if err != nil {
		t.Fatal(err)
	}
	for {
		tok, _, literal := l.NextToken()
		_ = literal
		if tok == TokILLEGAL {
			line, col := l.LineCol()
			t.Fatalf("illegal at:%d:%d %s", line, col, l.Err())
		} else if tok == TokEOF {
			return
		}
	}
}

func TestParse(t *testing.T) {
	fp := strings.NewReader(dsnFile)
	var l Lexer
	err := l.Reset(fp)
	if err != nil {
		t.Fatal(err)
	}
	p, err := NewParser(&l)
	if err != nil {
		t.Fatal(err)
	}

	maxToks := 99
	totalToks := 0
	filter := func(b []byte) bool {
		totalToks++
		return totalToks < maxToks
	}
	decls, err := p.ParseFilter(filter)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", decls)

}

func TestLexer(t *testing.T) {
	fp, _ := os.Open("piaa.dsn")
	var l Lexer
	l.Reset(fp)
	for {
		tok, _, literal := l.NextToken()
		if tok == TokILLEGAL {
			line, col := l.LineCol()
			t.Fatalf("illegal at:%d:%d %s", line, col, l.Err())
		} else if tok == TokEOF {
			return
		}
		_ = literal
		// t.Log(tok, string(literal))
	}
}

func TestString(t *testing.T) {
	const str = `"hello word's"`
	fp := strings.NewReader(str)
	var l Lexer
	l.Reset(fp)
	tok, _, literal := l.NextToken()
	if tok == TokILLEGAL {
		t.Fatal("illegal line", l.line, l.err, literal)
	}
	if string(literal) != str {
		t.Error("want", str, "got", string(literal))
	}
}
