/*
 * Copyright (C) 2021 Massimiliano Ghilardi
 *
 *     This Source Code Form is subject to the terms of the Mozilla Public
 *     License, v. 2.0. If a copy of the MPL was not distributed with this
 *     file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 *
 * api.go
 *
 *  Created on: Apr 23, 2021
 *      Author: Massimiliano Ghilardi
 */

// some random code, used for testing package go/typecheck

package testdata

const (
	Foo int = Bar >> 70
	Bar     = '\x01' << 100
)

// var Len = len((*Pair)(nil).first)

type (
	pair = struct {
		first, second int
	}
	Pair pair

	X *X

	foobar = interface {
		foo() (result Pair)
		bar(arg Pair)
	}
	FooBar foobar
)

type (
	RecurA struct{ *RecurB }
	RecurB RecurA
)
