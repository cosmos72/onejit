/*
 * Copyright (C) 2021 Massimiliano Ghilardi
 *
 *     This Source Code Form is subject to the terms of the Mozilla Public
 *     License, v. 2.0. If a copy of the MPL was not distributed with this
 *     file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 *
 * scanner.go
 *
 *  Created on: Mar 19, 2021
 *      Author: Massimiliano Ghilardi
 */

package scanner

import (
	"github.com/cosmos72/onejit/go/io"
	"github.com/cosmos72/onejit/go/strings"
	"github.com/cosmos72/onejit/go/token"
)

type tokenStatus int

const (
	tokenNormal tokenStatus = iota
	tokenInsertSemi
	tokenCached
)

type Scanner struct {
	item
	status         tokenStatus
	lastNonComment token.Token
	utf8Reader
	builder strings.Builder
}

type item struct {
	tok token.Token
	lit string
}

// (re)initialize internal buffers, then set io.Reader
// does NOT clear accumulated errors
func (s *Scanner) Init(file *token.File, src io.Reader) {
	s.item = item{}
	s.status = tokenNormal
	s.lastNonComment = token.ILLEGAL
	s.utf8Reader.init(file, src)
	s.builder.Reset()
}

func (s *Scanner) add() {
	s.addRune(s.ch)
}

func (s *Scanner) addRune(ch rune) {
	if ch != runeEOF {
		s.builder.WriteRune(ch)
	}
}

func (s *Scanner) addString(str string) {
	s.builder.WriteString(str)
}

func (s *Scanner) getString() string {
	return s.builder.String()
}

func (s *Scanner) clearString() {
	s.builder.Reset()
}

func (s *Scanner) setResult(tok token.Token) {
	s.tok = tok
	if tok != token.COMMENT {
		s.lastNonComment = tok
	}
	s.lit = s.getString()
}

// always sets s.lit = ""
func (s *Scanner) setResultTok(tok token.Token) {
	s.tok = tok
	if tok != token.COMMENT {
		s.lastNonComment = tok
	}
	s.lit = ""
}

// return accumulated errors
func (s *Scanner) Errors() token.ErrorList {
	return s.utf8Reader.errors
}

// clear accumulated errors
func (s *Scanner) ClearErrors() {
	s.utf8Reader.errors = nil
}

func (s *Scanner) File() *token.File {
	return s.file
}

// convert specified Pos to a Position
func (s *Scanner) Position(pos token.Pos) token.Position {
	return s.file.Position(pos)
}

// return position of begin and end of last returned token
func (s *Scanner) PosEnd() (token.Pos, token.Pos) {
	return token.Pos(s.pos), token.Pos(s.endpos)
}

func (s *Scanner) Scan() (token.Token, string) {
	if s.status == tokenCached {
		s.status = tokenNormal
		return s.tok, s.lit
	}
	ch := s.skipSpaces()
	s.pos = s.endpos
	if ch != '\n' {
		s.scan()
	}
	if s.status == tokenInsertSemi {
		s.status = tokenCached
		tok := token.SEMICOLON
		s.lastNonComment = tok
		return tok, ""
	}
	return s.tok, s.lit
}

// skip spaces, return the first non-space rune
func (s *Scanner) skipSpaces() rune {
	ch := s.ch
	for {
		switch ch {
		case runeBOF:
			ch = s.next()
			if ch == runeBOM {
				ch = s.next()
			}
			continue
		case runeNONE, '\t', '\r', ' ': // spaces
		case '\n':
			if s.needInsertSemi() {
				s.ch = runeNONE
				tok := token.SEMICOLON
				s.tok = tok
				s.lastNonComment = tok
				s.lit = ""
				return ch
			}
		default:
			// EOF, TOO_MANY_ERRORS or not a space
			return ch
		}
		ch = s.next()
	}
}

func (s *Scanner) scan() {
	s.clearString()
	ch := s.ch
	if isDecimalDigit(ch) {
		s.scanNumber()
	} else if ch == '"' {
		s.scanString()
	} else if ch == '\'' {
		s.scanRune()
	} else if ch == '`' {
		s.scanRawString()
	} else if ch == '.' {
		s.scanDot()
	} else if isOperator(ch) {
		s.scanOperatorOrComment()
	} else if ch == '_' || isLetter(ch) {
		s.scanIdentifier()
	} else if ch == runeEOF || ch == runeTOO_MANY_ERRORS {
		s.setResultTok(token.EOF)
	} else {
		s.invalid(errInvalidCharacter)
	}
}

func (s *Scanner) scanIdentifier() {
	s.add()
	ch := s.next()
	for ch == '_' || isLetter(ch) || isDigit(ch) {
		s.addRune(ch)
		ch = s.next()
	}
	str := s.getString()
	tok := token.Lookup(str)
	if tok != token.IDENT {
		s.clearString()
	}
	s.setResult(tok)
}

const unknownToken = token.ILLEGAL

func (s *Scanner) scanOperatorOrComment() {

	switch s.ch {
	case '/':
		s.scanSlash()
	default:
		s.scanOperator(unknownToken)
	}
}

func (s *Scanner) scanOperator(op token.Token) {
	next := unknownToken
	ch := s.ch
	for ch != runeEOF {
		s.addRune(ch)
		next = token.LookupOperator(s.getString())
		if next == unknownToken {
			break
		}
		op = next
		ch = s.next()
	}
	if op == unknownToken {
		s.error(errInvalidCharacter)
		s.next()
	} else {
		s.clearString()
	}
	s.setResult(op)
}

// scan characters after '.'
func (s *Scanner) scanDot() {
	// special case:
	// operator ... is the only one where removing last character
	// does not create another valid operator
	ch := s.next()
	if ch != '.' {
		if isDecimalDigit(ch) {
			// found ".digit" i.e. fraction-only floating point number
			s.scanIntOrFloat(base10, '.')
			return
		} else {
			// found ".nondigit"
			s.setResultTok(token.PERIOD)
			return
		}
	}
	// found ".."
	ch = s.next()
	if ch != '.' {
		// found "..nondot"
		// unread the second dot and return the first one
		s.unread = '.'
		s.setResultTok(token.PERIOD)
		return
	}
	// found "...""
	s.next()
	s.setResultTok(token.ELLIPSIS)
}

// scan characters after '/'
func (s *Scanner) scanSlash() {
	ch := s.next()
	switch ch {
	case '/':
		s.scanCommentLine()
	case '*':
		s.scanCommentBlock()
	default:
		// found "/something"
		s.addRune('/')
		s.scanOperator(token.QUO)
	}
}

func (s *Scanner) needInsertSemi() bool {
	flag := false
	if s.status == tokenNormal {
		switch s.lastNonComment {
		case token.IDENT, token.INT, token.FLOAT, token.IMAG, token.CHAR, token.STRING,
			token.BREAK, token.CONTINUE, token.FALLTHROUGH, token.RETURN,
			token.INC, token.DEC, token.RPAREN, token.RBRACK, token.RBRACE:
			flag = true
		}
	}
	return flag
}

func (s *Scanner) autoInsertSemi() {
	if s.needInsertSemi() {
		s.status = tokenInsertSemi
	}
}
