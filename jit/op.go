/*
 * gomacrojit - JIT compiler in Go
 *
 * Copyright (C) 2019-2020 Massimiliano Ghilardi
 *
 *     This Source Code Form is subject to the terms of the Mozilla Public
 *     License, v. 2.0. If a copy of the MPL was not distributed with this
 *     file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 *
 * op.go
 *
 *  Created on Jan 23, 2019
 *      Author Massimiliano Ghilardi
 */

package jit

import (
	"strconv"
)

type (
	Op uint8
)

// intentionally same values as go.Token
const (
	_ Op = iota
	_
	_
	_
	_
	_
	_
	_
	_
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %

	AND     // &
	OR      // |
	XOR     // ^
	SHL     // <<
	SHR     // >>
	AND_NOT // &^

	ADD_ASSIGN // +=
	SUB_ASSIGN // -=
	MUL_ASSIGN // *=
	QUO_ASSIGN // /=
	REM_ASSIGN // %=

	AND_ASSIGN     // &=
	OR_ASSIGN      // |=
	XOR_ASSIGN     // ^=
	SHL_ASSIGN     // <<=
	SHR_ASSIGN     // >>=
	AND_NOT_ASSIGN // &^=

	LAND // &&
	LOR  // ||
	_    // ARROW // <-
	INC  // ++
	DEC  // --

	EQL    // ==
	LSS    // <
	GTR    // >
	ASSIGN // =
	NOT    // !

	NEQ // !=
	LEQ // <=
	GEQ // >=

	_ // DEFINE   // :=
	_
	CAST    // LPAREN   // (
	BRACKET // LBRACK   // [
	_
	_
	FIELD // DOT // .

	JUMP    Op = 73 // GOTO
	JUMP_IF Op = 74 // IF
	RET     Op = 80 // RETURN

	NEG  = SUB // -
	STAR = MUL // *

	opLo = ADD
	opHi = RET
)

var opstring = [...]string{
	ADD:            "+",
	SUB:            "-",
	MUL:            "*",
	QUO:            "/",
	REM:            "%",
	AND:            "&",
	OR:             "|",
	XOR:            "^",
	SHL:            "<<",
	SHR:            ">>",
	AND_NOT:        "&^",
	ADD_ASSIGN:     "+=",
	SUB_ASSIGN:     "-=",
	MUL_ASSIGN:     "*=",
	QUO_ASSIGN:     "/=",
	REM_ASSIGN:     "%=",
	AND_ASSIGN:     "&=",
	OR_ASSIGN:      "|=",
	XOR_ASSIGN:     "^=",
	SHL_ASSIGN:     "<<=",
	SHR_ASSIGN:     ">>=",
	AND_NOT_ASSIGN: "&^=",
	LAND:           "&&",
	LOR:            "||",
	INC:            "++",
	DEC:            "--",
	EQL:            "==",
	LSS:            "<",
	GTR:            ">",
	ASSIGN:         "=",
	NOT:            "!",
	NEQ:            "!=",
	LEQ:            "<=",
	GEQ:            ">=",
	CAST:           "CAST",
	BRACKET:        "[]",
	FIELD:          ".",
	JUMP:           "JUMP",
	JUMP_IF:        "JUMP_IF",
	RET:            "RET",
}

func (op Op) String() string {
	var s string
	if op >= opLo && op <= opHi {
		s = opstring[op]
	} else {
		s = "Op(" + strconv.Itoa(int(op)) + ")"
	}
	return s

}