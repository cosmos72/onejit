// +build !amd64,!arm64

/*
 * gomacrojitjit - JIT compiler in Go
 *
 * Copyright (C) 2019 Massimiliano Ghilardi
 *
 *     This Source Code Form is subject to the terms of the Mozilla Public
 *     License, v. 2.0. If a copy of the MPL was not distributed with this
 *     file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 *
 * generic.go
 *
 *  Created on Feb 13, 2019
 *      Author Massimiliano Ghilardi
 */

package asm

import (
	"github.com/cosmos72/gomacrojit/jit_old/common"
)

const (
	ARCH_ID        = common.NOARCH
	ARCH_SUPPORTED = false
	NAME           = "generic"
)