package dsn

import (
	"fmt"
	"math"
)

type Parser struct {
	l *Lexer
}

func NewParser(l *Lexer) (*Parser, error) {
	return &Parser{
		l: l,
	}, nil
}

type TokenLit struct {
	Token   Token
	Literal string
}

type Decl struct {
	name   string
	level  int
	parent *Decl
	toks   []TokenLit
	decls  []Decl
}

// Name returns the declaration's name.
func (d *Decl) Name() string { return d.name }

// Children returns the declarations within d.
func (d *Decl) Children() []Decl {
	return d.decls
}

// Args returns the tokens after the declaration's name.
func (d *Decl) Args() []TokenLit {
	return d.toks
}

// Parent returns d's parent declaration. If d is top level then nil is returned.
func (d *Decl) Parent() *Decl {
	return d.parent
}

// Depth returns the parentheses depth of the declaration. Lowest is 1.
func (d *Decl) Depth() int {
	return d.level
}

func (p *Parser) ParseFilter(filter func([]byte) bool) (decls []Decl, err error) {
	l := p.l
	const noFilter = math.MaxInt
	topDecl := Decl{}
	topDecl.name = "GLOBAL"
	currentDecl := &topDecl
	var openParen bool
	var filterMaxDepth = noFilter
	var tok Token
	var literal []byte
	for {
		tok, _, literal = l.NextToken()
		if l.Err() != nil || tok == TokILLEGAL || tok == TokEOF {
			break
		}
		lvl := l.Parens()
		if lvl > filterMaxDepth {
			continue
		} else {
			filterMaxDepth = noFilter // Reset filter.
		}
		keep := filter(literal)
		if !keep {
			filterMaxDepth = lvl
			continue
		}

		if openParen {
			if tok == TokLPAREN {
				line, col := l.LineCol()
				return nil, fmt.Errorf("adjacent open paretheses:%d:%d", line, col)
			} else if tok != TokIDENT {
				line, col := l.LineCol()
				return nil, fmt.Errorf("adjacent open paretheses:%d:%d ", line, col)
			}
			openParen = false
			currentDecl.decls = append(currentDecl.decls, Decl{
				name:   string(literal),
				level:  lvl,
				parent: currentDecl,
			})
			currentDecl = &currentDecl.decls[len(currentDecl.decls)-1]
			continue
		}

		switch tok {
		case TokLPAREN:
			openParen = true
		case TokIDENT, TokSTRING, TokINTEGER, TokFLOAT:
			currentDecl.toks = append(currentDecl.toks, TokenLit{Token: tok, Literal: string(literal)})
		case TokRPAREN:
			currentDecl = currentDecl.parent
		}
	}

	if l.Err() != nil {
		line, col := l.LineCol()
		return nil, fmt.Errorf("%s at:%d:%d", l.Err(), line, col)
	} else if tok == TokEOF {
		for i := range topDecl.decls {
			topDecl.decls[i].parent = nil
		}
		return topDecl.decls, nil
	} else if tok == TokILLEGAL {
		line, col := l.LineCol()
		return nil, fmt.Errorf("illegal at:%d:%d", line, col)
	}
	panic("unreachable")
}
