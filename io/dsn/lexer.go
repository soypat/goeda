package kicad

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"unicode/utf8"
)

type Token int

// Install stringer tool:
//  go install golang.org/x/tools/cmd/stringer@latest

// Tokens
//
//go:generate stringer -type=Token -output=stringers.go -trimprefix=Tok
const (
	TokILLEGAL Token = iota
	TokEOF
	TokLPAREN  // (
	TokRPAREN  // )
	TokIDENT   // R1
	TokINTEGER // 1234
	TokSTRING  // "abc"
	TokFLOAT   // 123.45
)

type Lexer struct {
	input bufio.Reader
	ch    rune // current character (utf8)
	strq  rune // Strquote set by parser statement.
	err   error
	idbuf []byte // accumulation buffer.

	// Higher level statistics fields:

	line   int // file line number
	col    int // column number in line
	pos    int // byte position.
	parens int // Parentheses counter to pick up on unbalanced parentheses early.
}

// Reset discards all state and buffered data and begins a new lexing
// procedure on the input r. It performs a single utf8 read to initialize.
func (l *Lexer) Reset(r io.Reader) error {
	if r == nil {
		return errors.New("nil reader")
	}
	*l = Lexer{
		input: l.input,
		line:  1,
		idbuf: l.idbuf,
		strq:  '"',
	}
	l.input.Reset(r)
	if l.idbuf == nil {
		l.idbuf = make([]byte, 0, 1024)
	}
	l.readChar()
	return l.err
}

// Err returns the lexer error.
func (l *Lexer) Err() error {
	if l.err == io.EOF {
		return nil
	}
	return l.err
}

// LineCol returns the current line number and column number (utf8 relative).
func (l *Lexer) LineCol() (line, col int) {
	return l.line, l.col
}

// Pos returns the absolute position of the lexer in bytes from the start of the file.
func (l *Lexer) Pos() int { return l.pos }

// Parens returns the parentheses depth at the current position.
func (l *Lexer) Parens() int { return l.parens }

// Next token parses the upcoming token and returns the literal representation
// of the token for identifiers, strings and numbers.
// The returned byte slice is reused between calls to NextToken.
func (l *Lexer) NextToken() (tok Token, startPos int, literal []byte) {
	l.skipWhitespace()
	startPos = l.pos
	if l.err == io.EOF {
		return TokEOF, startPos, nil
	} else if l.err != nil {
		return TokILLEGAL, startPos, nil
	}
	ch := l.ch
	switch ch {
	case l.strq:
		tok = TokSTRING
		literal = l.readString()

	case '(':
		tok = TokLPAREN
		l.readChar()
		l.parens++
	case ')':
		tok = TokRPAREN
		l.readChar()
		l.parens--
		if l.parens < 0 {
			l.err = errors.New("unbalanced parentheses")
		}
	default:
		if isDigit(ch) || ch == '-' {
			var isFloat bool
			tok = TokINTEGER // Handle floats later on.
			literal, isFloat = l.readNumber()
			if isFloat {
				tok = TokFLOAT
			}
		} else if isLetter(ch) {
			tok = TokIDENT
			literal = l.readIdentifier()
			if bytes.Equal(literal, []byte("string_quote")) {
				l.skipWhitespace()
				l.strq = l.ch
				l.readChar() // Consume quote char.
			}
		} else {
			tok = TokILLEGAL
		}
	}
	return tok, startPos, literal
}

func (l *Lexer) readIdentifier() []byte {
	start := l.bufstart()
	for isLetter(l.ch) || isDigit(l.ch) {
		l.idbuf = utf8.AppendRune(l.idbuf, l.ch)
		l.readChar()
	}
	return l.idbuf[start:]
}

func (l *Lexer) readString() []byte {
	start := l.bufstart()
	for {
		l.idbuf = utf8.AppendRune(l.idbuf, l.ch)
		l.readChar()
		ch := l.ch
		switch ch {
		case '\\':
			l.err = errors.New("does not support escaping")
			return l.idbuf[start:]
		case '"':
			// Finish string and consume double quotes.
			l.idbuf = utf8.AppendRune(l.idbuf, l.ch)
			l.readChar()
			return l.idbuf[start:]
		case '\n':
			l.err = errors.New("newline in string")
			return l.idbuf[start:]
		}
	}
}

func (l *Lexer) readNumber() ([]byte, bool) {
	start := l.bufstart()
	seenDot := false
	if l.ch == '-' {
		// Consume leading negative character.
		l.idbuf = utf8.AppendRune(l.idbuf, l.ch)
		l.readChar()
	}
	for {
		ch := l.ch
		if !isDigit(ch) {
			if !seenDot && ch == '.' {
				seenDot = true
			} else {
				break
			}
		}
		l.idbuf = utf8.AppendRune(l.idbuf, l.ch)
		l.readChar()
	}
	return l.idbuf[start:], seenDot
}

func (l *Lexer) bufstart() int {
	const reuseMem = true
	if reuseMem {
		l.idbuf = l.idbuf[:0]
		return 0
	}
	return len(l.idbuf)
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) readChar() {
	if l.err != nil {
		l.ch = 0 // Just in case annihilate char.
		return
	}
	ch, sz, err := l.input.ReadRune()
	if err != nil {
		l.ch = 0
		l.err = err
		return
	}
	if ch == '\n' {
		l.line++
		l.col = 0
	} else {
		l.col++
	}
	l.pos += sz
	l.ch = ch
}

func (l *Lexer) peekChar() rune {
	posstart := l.pos
	linestart := l.line
	colstart := l.col
	l.readChar()
	if l.err != nil {
		return 0
	}
	l.line = linestart
	l.err = l.input.UnreadRune()
	l.pos = posstart
	l.col = colstart
	return l.ch
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' ||
		ch == '_' || ch == '-' || ch == '@' || ch == '/' ||
		ch == '.' || ch == '+' || ch == ':' || ch == '[' || ch == ']' ||
		ch == ','
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\n' || ch == '\t' || ch == '\r'
}
