/*
 * Copyright (C) 2021 Massimiliano Ghilardi
 *
 *     This Source Code Form is subject to the terms of the Mozilla Public
 *     License, v. 2.0. If a copy of the MPL was not distributed with this
 *     file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 *
 * parser.go
 *
 *  Created on: Mar 23, 2021
 *      Author: Massimiliano Ghilardi
 */

package parser

import (
	"github.com/cosmos72/onejit/go/ast"
	"github.com/cosmos72/onejit/go/config"
	"github.com/cosmos72/onejit/go/io"
	"github.com/cosmos72/onejit/go/scanner"
	"github.com/cosmos72/onejit/go/strings"
	"github.com/cosmos72/onejit/go/token"
)

type Mode uint32

const (
	ParseImports  Mode       = 1 << iota // also parse 'import ...'
	ParseDecls                           // also parse declarations
	ParseComments                        // add parsed comments to each ast.Node
	ParseAll      = ^Mode(0)             // parse everything
)

// Go syntax parser. Supports reading from a io.Reader.
// callers should invoke Init() before any Parse*() method.
// The only method that clears accumulated errors is ClearErrors()
type Parser struct {
	curr    ast.Atom
	unread0 ast.Atom
	Mode    Mode        // Parser Mode is exported
	Lang    config.Lang // Parser Lang is exported
	scanner scanner.Scanner
}

// initialize parser and read the first non-comment token from src
// does NOT clear accumulated errors
func (p *Parser) Init(file *token.File, src io.Reader, mode Mode, lang config.Lang) {
	p.scanner.Init(file, src)
	p.curr = ast.Atom{}
	p.unread0 = ast.Atom{}
	p.Mode = mode
	p.Lang = lang

	p.next()
}

// initialize parser to read from specified string.
// does NOT clear accumulated errors
func (p *Parser) InitString(source string, mode Mode, lang config.Lang) {
	var reader strings.Reader
	reader.Reset(source)
	p.Init(token.NewFile("<string>", 0), &reader, mode, lang)
}

// return current token lookahead buffer
func (p *Parser) CurrToken() token.Token {
	tok := p.curr.Tok
	if unread := p.unread0.Tok; unread != 0 {
		tok = unread
	}
	return tok
}

// skip one token, return next one (which is kept in lookahead buffer)
func (p *Parser) SkipToken() token.Token {
	return p.next()
}

// parse a single declaration, statement or expression
func (p *Parser) Parse() (node ast.Node) {
	tok := p.tok()
	if tok == token.SEMICOLON {
		tok = p.next()
	}
	switch tok {
	case token.EOF:
		node = p.makeEof()
	case token.PACKAGE:
		node = p.parsePackage()
	case token.IMPORT:
		if p.Mode&ParseImports != 0 {
			node = p.parseImport()
		} else if p.Mode&ParseDecls != 0 {
			p.parseImport() // skip 'import ...'
		} else {
			node = p.makeEof()
		}
	case token.SEMICOLON:
		// node = nil
	default:
		if p.Mode&ParseDecls == 0 {
			node = p.makeEof()
		} else if isLeave(tok) {
			node = p.makeBadNode(p.parseAtom(tok), errExpectingDecl)
		} else if isDecl(tok) {
			node = p.ParseTopLevelDecl()
		} else {
			node = p.parseStmt(allowCompositeLit)
		}
	}
	return node
}

// put back specified atom into the token stream, before p.curr
func (p *Parser) unread(atom *ast.Atom) {
	p.unread0 = p.curr
	p.curr = *atom
}

// get next non-comment token and store it in p.curr
func (p *Parser) next() token.Token {
	curr := &p.curr
	if unread := &p.unread0; unread.Tok != 0 {
		*curr = *unread
		*unread = ast.Atom{}
		return curr.Tok
	}
	curr.Comment = nil
	s := &p.scanner
	for {
		curr.Tok, curr.Lit = s.Scan()
		if curr.Tok != token.COMMENT {
			curr.TokPos, curr.TokEnd = s.PosEnd()
			break
		} else if p.Mode&ParseComments != 0 {
			curr.Comment = append(curr.Comment, curr.Lit)
		}
	}
	return curr.Tok
}

func (p *Parser) tok() token.Token {
	return p.curr.Tok
}

func (p *Parser) pos() token.Pos {
	return p.curr.TokPos
}

func (p *Parser) consumeComment() []string {
	ret := p.curr.Comment
	p.curr.Comment = nil
	return ret
}

func (p *Parser) enter(list []ast.Node, tok token.Token) []ast.Node {
	if p.tok() == tok {
		p.next() // skip tok
	} else {
		list = append(list, p.makeBad(tok))
	}
	return list
}

// if current token != tok, report error and skip tokens
// until specified token is found
func (p *Parser) leave(list []ast.Node, tok token.Token) []ast.Node {
	for p.tok() != tok {
		list = append(list, p.parseBad(tok))
		if p.tok() == token.EOF {
			return list
		}
	}
	p.next() // skip tok
	return list
}

// if current token != tok, report error and skip tokens
// until specified token is found
func (p *Parser) leaveNode(node ast.Node, tok token.Token) ast.Node {
	var badnode bool
	for p.tok() != tok {
		if !badnode {
			badnode = true
			node = p.makeBadNode(node, tok)
		}
		if p.tok() == token.EOF {
			return node
		}
		p.next()
	}
	p.next() // skip tok
	return node
}

func (p *Parser) makeAtom(tok token.Token) *ast.Atom {
	atom := p.curr
	atom.Tok = tok
	p.curr.Comment = nil
	return &atom
}

func (p *Parser) makeBad(msg interface{}) *ast.Bad {
	tok := p.tok()
	return &ast.Bad{
		Atom: ast.Atom{Tok: tok, TokPos: p.pos()},
		Node: p.makeAtom(tok),
		Err:  p.error(p.pos(), msg),
	}
}

func (p *Parser) makeBadNode(x ast.Node, msg interface{}) *ast.Bad {
	return &ast.Bad{
		Atom: ast.Atom{Tok: x.Op(), TokPos: x.Pos()},
		Node: x,
		Err:  p.error(x.Pos(), msg),
	}
}

func (p *Parser) makeBinary() (binary *ast.Binary) {
	binary = &ast.Binary{Atom: p.curr}
	p.curr.Comment = nil
	return binary
}

func (p *Parser) makeEof() *ast.Atom {
	return &ast.Atom{Tok: token.EOF, TokPos: p.pos()}
}

func (p *Parser) makeIdent() (node ast.Node) {
	if p.tok() == token.IDENT {
		node = p.makeAtom(token.IDENT)
	} else {
		node = p.makeBad(errExpectingIdent)
	}
	p.curr.Comment = nil
	return node
}

func (p *Parser) makeList() (list *ast.List) {
	list = &ast.List{Atom: p.curr}
	p.curr.Comment = nil
	return list
}

func (p *Parser) makeUnary() (unary *ast.Unary) {
	unary = &ast.Unary{Atom: p.curr}
	p.curr.Comment = nil
	return unary
}

func (p *Parser) parseAtom(tok token.Token) *ast.Atom {
	node := p.makeAtom(tok)
	p.next()
	return node
}

func (p *Parser) parseBad(msg interface{}) *ast.Bad {
	node := p.makeBad(msg)
	p.next()
	return node
}

func (p *Parser) parseBinary() *ast.Binary {
	node := p.makeBinary()
	p.next()
	return node
}

func (p *Parser) parseIdent() ast.Node {
	node := p.makeIdent()
	p.next()
	return node
}

func (p *Parser) parseList() *ast.List {
	node := p.makeList()
	p.next()
	return node
}

func (p *Parser) parseUnary() *ast.Unary {
	node := p.makeUnary()
	p.next()
	return node
}
