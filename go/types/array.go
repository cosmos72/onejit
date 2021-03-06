/*
 * Copyright (C) 2021 Massimiliano Ghilardi
 *
 *     This Source Code Form is subject to the terms of the Mozilla Public
 *     License, v. 2.0. If a copy of the MPL was not distributed with this
 *     file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 *
 * array.go
 *
 *  Created on: Apr 01, 2021
 *      Author: Massimiliano Ghilardi
 */

package types

import (
	"github.com/cosmos72/onejit/go/config"
	"github.com/cosmos72/onejit/go/io"
	"github.com/cosmos72/onejit/go/strings"
)

type Array struct {
	arrayTag struct{} // occupies zero bytes
	rtype    Complete
	extra    extra
}

// *Array implements Type

func (t *Array) String() string {
	_ = t.arrayTag
	var b builder
	t.WriteTo(&b, fullPkgPath)
	return b.String()
}

func (t *Array) Underlying() Type {
	return t
}

func (t *Array) common() *Complete {
	return &t.rtype
}

func (t *Array) complete() {
	if t.rtype.flags&flagComplete != 0 {
		return
	}
	elem := t.Elem()
	len := t.Len()
	t.rtype.size = computeArraySize(elem, len)
	t.rtype.hash = computeArrayHash(elem, len)
	t.rtype.align = computeArrayAlign(elem, len)
	t.rtype.flags = computeArrayFlags(elem, len)
}

func (t *Array) WriteTo(dst io.StringWriter, flag verbose) {
	if flag == shortPkgName {
		dst.WriteString(t.rtype.str)
		return
	}
	writeArrayTo(dst, t.Len(), t.Elem(), flag)
}

// *Array specific methods

func (t *Array) Elem() Type {
	return t.rtype.elem
}

func (t *Array) Len() uint64 {
	return uint64(t.extra.n1) | uint64(t.extra.n2)<<32
}

type arrayKey struct {
	elem Type
	len  uint64
}

var arrayMap = map[arrayKey]*Array{}

// create a new Array type
func NewArray(elem Type, len uint64) *Array {
	archMaxSize := config.TargetArch().MaxSizeBytes()
	if len > archMaxSize {
		panic("NewArray length " + uintToString(len) + " exceeds archMaxSize = " + uintToString(archMaxSize))
	}
	key := arrayKey{elem, len}
	t := arrayMap[key]
	if t != nil {
		return t
	}
	t = &Array{
		rtype: Complete{
			size:  computeArraySize(elem, len),
			align: computeArrayAlign(elem, len),
			flags: computeArrayFlags(elem, len),
			kind:  ArrayKind,
			elem:  elem,
			hash:  computeArrayHash(elem, len),
			str:   makeArrayString(len, elem, shortPkgName),
		},
	}
	t.rtype.typ = t
	arrayMap[key] = t
	return t
}

func computeArraySize(elem Type, len uint64) uint64 {
	elemsize := elem.common().size
	if elemsize == unknownSize {
		return elemsize
	}
	archMaxSize := config.TargetArch().MaxSizeBytes()
	if elemsize != 0 && len > archMaxSize/elemsize {
		panic("NewArray length " + uintToString(len) + " exceeds archMaxSize/elementSize = " + uintToString(archMaxSize/elemsize))
	}
	return len * elemsize
}

func computeArrayHash(elem Type, len uint64) hash {
	elemhash := elem.common().hash
	if elemhash == unknownHash {
		return unknownHash
	}
	return elemhash.Uint64(len)
}

func computeArrayAlign(elem Type, len uint64) uint16 {
	if len == 0 {
		return 0
	}
	return elem.common().align
}

func computeArrayFlags(elem Type, len uint64) flags {
	flag := elem.common().flags & (flagComplete | flagComparable | flagNotComparable)
	if len == 0 || elem.common().size == 0 {
		flag |= flagNeedPadding
	}
	return flag
}

func makeArrayString(len uint64, elem Type, flag verbose) string {
	var b builder
	writeArrayTo(&b, len, elem, flag)
	return b.String()
}

func writeArrayTo(dst io.StringWriter, len uint64, elem Type, flag verbose) {
	dst.WriteString("[")
	dst.WriteString(uintToString(len))
	dst.WriteString("]")
	elem.WriteTo(dst, flag)
}

func uintToString(n uint64) string {
	return strings.Uint64ToString(n)
}
