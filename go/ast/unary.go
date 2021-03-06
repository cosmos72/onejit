/*
 * Copyright (C) 2021 Massimiliano Ghilardi
 *
 *     This Source Code Form is subject to the terms of the Mozilla Public
 *     License, v. 2.0. If a copy of the MPL was not distributed with this
 *     file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 *
 * unary.go
 *
 *  Created on: Mar 23, 2021
 *      Author: Massimiliano Ghilardi
 */

package ast

import (
	"github.com/cosmos72/onejit/go/io"
	"github.com/cosmos72/onejit/go/strings"
	"github.com/cosmos72/onejit/go/token"
)

/**
 * Node with one children. Used for unary expressions, types,
 * DEC, DEFER, GO, GOTO, INC and other one-argument statements.
 */
type Unary struct {
	Atom
	X Node
}

func (u *Unary) Len() int {
	return 1
}

func (u *Unary) At(i int) Node {
	if i == 0 {
		return u.X
	}
	return outOfRange()
}

func (u *Unary) End() token.Pos {
	if u.X != nil {
		return u.X.End()
	} else {
		return u.Atom.End()
	}
}

func (u *Unary) String() string {
	if u == nil {
		return "nil"
	}
	var buf strings.Builder
	u.WriteTo(&buf)
	return buf.String()
}

func (u *Unary) WriteTo(out io.StringWriter) {
	if u == nil {
		out.WriteString("nil")
	} else {
		writeListTo(out, u)
	}
}
