/*
 * Copyright (C) 2021 Massimiliano Ghilardi
 *
 *     This Source Code Form is subject to the terms of the Mozilla Public
 *     License, v. 2.0. If a copy of the MPL was not distributed with this
 *     file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * token.go
 *
 *  Created on: Mar 19, 2021
 *      Author: Massimiliano Ghilardi
 */

package token

import (
	"go/token"
)

// wraps go/token.Token
type Token token.Token

const (
	ILLEGAL = Token(token.ILLEGAL)
	EOF     = Token(token.EOF)
	COMMENT = Token(token.COMMENT)

	IDENT  = Token(token.IDENT)
	INT    = Token(token.INT)
	FLOAT  = Token(token.FLOAT)
	IMAG   = Token(token.IMAG)
	CHAR   = Token(token.CHAR)
	STRING = Token(token.STRING)

	// Operators
	ADD = Token(token.ADD)
	SUB = Token(token.SUB)
	MUL = Token(token.MUL)
	QUO = Token(token.QUO)
	REM = Token(token.REM)

	AND     = Token(token.AND)
	OR      = Token(token.OR)
	XOR     = Token(token.XOR)
	SHL     = Token(token.SHL)
	SHR     = Token(token.SHR)
	AND_NOT = Token(token.AND_NOT)

	ADD_ASSIGN = Token(token.ADD_ASSIGN)
	SUB_ASSIGN = Token(token.SUB_ASSIGN)
	MUL_ASSIGN = Token(token.MUL_ASSIGN)
	QUO_ASSIGN = Token(token.QUO_ASSIGN)
	REM_ASSIGN = Token(token.REM_ASSIGN)

	AND_ASSIGN     = Token(token.AND_ASSIGN)
	OR_ASSIGN      = Token(token.OR_ASSIGN)
	XOR_ASSIGN     = Token(token.XOR_ASSIGN)
	SHL_ASSIGN     = Token(token.SHL_ASSIGN)
	SHR_ASSIGN     = Token(token.SHR_ASSIGN)
	AND_NOT_ASSIGN = Token(token.AND_NOT_ASSIGN)

	LAND  = Token(token.LAND)
	LOR   = Token(token.LOR)
	ARROW = Token(token.ARROW)
	INC   = Token(token.INC)
	DEC   = Token(token.DEC)

	EQL    = Token(token.EQL)
	LSS    = Token(token.LSS)
	GTR    = Token(token.GTR)
	ASSIGN = Token(token.ASSIGN)
	NOT    = Token(token.NOT)

	NEQ      = Token(token.NEQ)
	LEQ      = Token(token.LEQ)
	GEQ      = Token(token.GEQ)
	DEFINE   = Token(token.DEFINE)
	ELLIPSIS = Token(token.ELLIPSIS)

	LPAREN = Token(token.LPAREN)
	LBRACK = Token(token.LBRACK)
	LBRACE = Token(token.LBRACE)
	COMMA  = Token(token.COMMA)
	PERIOD = Token(token.PERIOD)

	RPAREN    = Token(token.RPAREN)
	RBRACK    = Token(token.RBRACK)
	RBRACE    = Token(token.RBRACE)
	SEMICOLON = Token(token.SEMICOLON)
	COLON     = Token(token.COLON)

	// Keywords
	BREAK    = Token(token.BREAK)
	CASE     = Token(token.CASE)
	CHAN     = Token(token.CHAN)
	CONST    = Token(token.CONST)
	CONTINUE = Token(token.CONTINUE)

	DEFAULT     = Token(token.DEFAULT)
	DEFER       = Token(token.DEFER)
	ELSE        = Token(token.ELSE)
	FALLTHROUGH = Token(token.FALLTHROUGH)
	FOR         = Token(token.FOR)

	FUNC   = Token(token.FUNC)
	GO     = Token(token.GO)
	GOTO   = Token(token.GOTO)
	IF     = Token(token.IF)
	IMPORT = Token(token.IMPORT)

	INTERFACE = Token(token.INTERFACE)
	MAP       = Token(token.MAP)
	PACKAGE   = Token(token.PACKAGE)
	RANGE     = Token(token.RANGE)
	RETURN    = Token(token.RETURN)

	SELECT = Token(token.SELECT)
	STRUCT = Token(token.STRUCT)
	SWITCH = Token(token.SWITCH)
	TYPE   = Token(token.TYPE)
	VAR    = Token(token.VAR)

	tok_counter Token = iota
	// array or slice type. if first child is nil => slice type.
	// otherwise first child is array length (possibly '...')
	ARRAY Token = VAR + iota - tok_counter
	BLOCK
	BOTH_DIR
	CALL
	COMPOSITE_LIT // type { elem0, elem1 ... }
	DECLS         // list of declarations
	EXPRS
	// struct field, interface element, function param or result.
	// stored in *ast.Field
	FIELD
	FILE
	GENERIC     // generic type parameters [T0 i0, T1 i1, ...]
	IMPORT_SPEC // single import line 'alias "full/pkg/path"'
	IMPORTS     // list of IMPORT nodes
	INDEX       // array/map/slice indexing a[b] or generic instantiation a[b,c...]
	KEY_VALUE   // key: value inside a composite literal
	LABEL       // label: statement
	LAMBDA      // function literal
	NAMES       // list of identifiers
	PARAMS
	RECV_DIR
	RESULTS
	SEND_DIR
	SLICE // slice expression a[b:c] or a[b:c:d]
	TYPE_ASSERT
	TYPESWITCH // switch x.(type) { }
	VALUE_SPEC
)

var tokens = [...]string{
	"ARRAY", "BLOCK", "BOTH_DIR", "CALL", "COMPOSITE_LIT", "DECLS", "EXPRS",
	"FIELD", "FILE", "GENERIC", "IMPORT_SPEC", "IMPORTS", "INDEX", "KEY_VALUE",
	"LABEL", "LAMBDA", "NAMES", "PARAMS", "RECV_DIR", "RESULTS", "SEND_DIR",
	"SLICE", "TYPE_ASSERT", "TYPESWITCH", "VALUE_SPEC",
}

var operators = makeOperators()

func makeOperators() map[string]Token {
	lo, hi := ADD, COLON
	m := make(map[string]Token, hi-lo+1)
	for op := lo; op <= hi; op++ {
		m[op.String()] = op
	}
	return m
}

func LookupOperator(str string) Token {
	return operators[str]
}

func Lookup(str string) Token {
	n := len(str)
	if n < 2 || n > 11 { // len("fallthrough") == 11
		return IDENT
	}
	return Token(token.Lookup(str))
}

func (tok Token) IsKeyword() bool {
	return token.Token(tok).IsKeyword()
}

func (tok Token) IsLiteral() bool {
	return token.Token(tok).IsLiteral()
}

func (tok Token) IsOperator() bool {
	return token.Token(tok).IsOperator()
}

const (
	LowestPrec  = 0 // non-operators
	UnaryPrec   = 6
	HighestPrec = 7
)

func (tok Token) Precedence() int {
	return token.Token(tok).Precedence()
}

func (tok Token) String() string {
	var ret string
	if tok >= ARRAY && tok <= VALUE_SPEC {
		ret = tokens[tok-ARRAY]
	} else {
		ret = token.Token(tok).String()
	}
	return ret
}
