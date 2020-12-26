/*
 * gomacrojit - JIT compiler in Go
 *
 * Copyright (C) 2018-2020 Massimiliano Ghilardi
 *
 *     This Source Code Form is subject to the terms of the Mozilla Public
 *     License, v. 2.0. If a copy of the MPL was not distributed with this
 *     file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 *
 * amd64_mem.go
 *
 *  Created on Dec 26, 2020
 *      Author Massimiliano Ghilardi
 */

package jit

// amd64 memory location.
type Amd64Mem struct {
	offset  int32
	base    RegId
	index   RegId
	roscale uint8 // 1,2,4 or 8. add 128 if readonly
	kind    Kind
}

func (m Amd64Mem) IsAssignable() bool {
	return m.roscale&128 == 0
}

// return a read-only view of a subset of memory contents
// note: memory does not become immutable - it can still be modified
// through the original struct
func (m Amd64Mem) ReadOnly(subset Kind) Amd64Mem {
	kindMustBeSubset("Amd64Mem", subset, m.kind)
	m.kind = subset
	m.roscale |= 128
	return m
}

func (m Amd64Mem) Offset() int32 {
	return m.offset
}

func (m Amd64Mem) BaseId() RegId {
	return m.base
}

func (m Amd64Mem) IndexId() RegId {
	return m.index
}

func (m Amd64Mem) Scale() uint8 {
	return m.roscale & 127
}

// implement Expr interface
func (m Amd64Mem) expr() {}

func (m Amd64Mem) RegId() RegId {
	return NoRegId
}

func (m Amd64Mem) Kind() Kind {
	return m.kind
}

func (m Amd64Mem) IsConst() bool {
	// memory access cannot be a compile-time constant
	return false
}

func (m Amd64Mem) Size() Size {
	return m.kind.Size()
}

func (m Amd64Mem) Class() Class {
	return MEM
}

func (m Amd64Mem) Children() int {
	return 2
}

func (m Amd64Mem) Child(i int) Node {
	switch i {
	case 0:
		if m.base == NoRegId {
			return nil
		}
		return archReg(Ptr, m.base)
	case 1:
		if m.index == NoRegId {
			return nil
		}
		return archReg(Int64, m.index)
	}
	return badIndex(i, 2)
}