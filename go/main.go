/*
 * Copyright (C) 2021 Massimiliano Ghilardi
 *
 *     This Source Code Form is subject to the terms of the Mozilla Public
 *     License, v. 2.0. If a copy of the MPL was not distributed with this
 *     file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 *
 * main.go
 *
 *  Created on: Mar 19, 2021
 *      Author: Massimiliano Ghilardi
 */

package main

import (
	"os"

	_ "github.com/cosmos72/onejit/go/jit_old"
	_ "github.com/cosmos72/onejit/go/scanner"
)

func main() {
	println("hello from github.com/cosmos72/onejit/go")
	println("press ENTER to quit")
	os.Stdin.Read(make([]byte, 1))
	println("bye")
}
