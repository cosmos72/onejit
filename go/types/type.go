/*
 * Copyright (C) 2021 Massimiliano Ghilardi
 *
 *     This Source Code Form is subject to the terms of the Mozilla Public
 *     License, v. 2.0. If a copy of the MPL was not distributed with this
 *     file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 *
 * type.go
 *
 *  Created on: Mar 31, 2021
 *      Author: Massimiliano Ghilardi
 */

package types

import (
	"github.com/cosmos72/onejit/go/io"
	"github.com/cosmos72/onejit/go/strings"
)

type (
	Type interface {
		String() string
		Underlying() Type
		common() *Complete
		complete() // completes everything in the type, except flags & flagComplete
		WriteTo(io.StringWriter, verbose)
	}

	flags   uint16
	ChanDir flags

	builder = strings.Builder
	verbose = strings.Verbose
)

const (
	RecvDir ChanDir = 1
	SendDir ChanDir = 2
	BothDir         = RecvDir | SendDir

	flagComparable    flags = 4
	flagNotComparable flags = 8
	flagComplete      flags = 16
	flagNeedPadding   flags = 32 // type is or ends with zero-byte struct or array
	flagVariadic      flags = 64

	unknownSize  = ^uint64(0)
	unknownAlign = ^uint16(0)
	unknownHash  = 0

	shortPkgName = strings.ShortPkgName
	fullPkgPath  = strings.FullPkgPath
)

func TypeString(t Type, showFullPkgPath bool) string {
	if showFullPkgPath {
		return t.String()
	}
	return t.common().str
}

/**
 * return true if values of type 'src' can be assigned to variables of type 'dst'
 *
 * as per Go specifications https://golang.org/ref/spec#Assignability
 * a type 'src' is assignable to type 'dst' if at least one of the following applies:
 * - src and dst are identical
 * - src and dst have identical underlying types, and at least one of src, dst is unnamed
 * - dst is an interface type, and src implements it
 * - src is a bidirectional channel value, dst is a channel type, they have identical element types
 *   and and at least one of src, dst is unnamed
 * - src is UntypedNil and dst is a channel, function, interface, map, pointer or slice type
 * - src is an untyped and is representable by dst
 */
func AssignableTo(src *Complete, dst *Complete) bool {
	if src == nil || dst == nil {
		// nil types, cannot check assignability
		return false
	} else if src == dst {
		// identical types
		return true
	}
	u1, u2 := src.Underlying(), dst.Underlying()
	k1, k2 := src.kind, dst.kind
	var ok bool
	if k1 == Invalid || k2 == Invalid || u1 == nil || u2 == nil {
		// incomplete types, cannot check assignability
		return false
	} else if k1.IsUntyped() {
		ok = untypedAssignableTo(k1, k2)
	}
	if !ok && u1 == u2 {
		sname, dname := src.Name(), dst.Name()
		// identical underlying types. at least one type is unnamed?
		ok = sname == "" || dname == ""
	}
	if !ok && k1 == ChanKind && k2 == ChanKind && src.ChanDir() == BothDir && src.elem == dst.elem {
		sname, dname := src.Name(), dst.Name()
		// compatible channel types. at least one type is unnamed?
		ok = sname == "" || dname == ""
	}
	if !ok && k2 == InterfaceKind {
		ok = Implements(src, dst)
	}
	return ok
}

// 'src' is untyped kind
func untypedAssignableTo(src Kind, dst Kind) bool {
	var ok bool
	switch src {
	case UntypedBool:
		ok = dst.Category() == Bool
	case UntypedInt, UntypedRune, UntypedFloat:
		switch dst.Category() {
		case Int, Float64:
			// must also check that constant.Value conversion is exact
			ok = true
		}
	case UntypedComplex:
		ok = dst.Category() == Complex128
	case UntypedString:
		ok = dst.Category() == String
	case UntypedNil:
		// untyped nil can be assigned to all nillable types
		ok = dst.IsNillable()
	}
	return ok
}

/**
 * return true if values of type 'src' can be converted to type 'dst'
 *
 * as per Go specifications https://golang.org/ref/spec#Conversions
 * a type 'src' is convertible to type 'dst' if at least one of the following applies:
 * - src is assignable to dst
 * - ignoring struct tags, src and dst have identical underlying types
 * - ignoring struct tags, src and dst are unnamed pointers,
 *   and their element types have identical underlying types
 * - src and dst are both integer or floating point types
 * - src and dst are both complex types
 * - src is an integer (obsolescent) or a slice of byte/rune and dst is a string
 * - src is a string and dst is a slice of byte/rune
 *
 * notes:
 * - "X is a string" means: "X underlying type is string"
 * - analogously, "X is a complex type" means "X underlying type is complex64 or complex128"
 * - analogously, "X is a floating point" means "X underlying type is float32 or float64"
 * - analogously, "X is an integer" means "X underlying type is one of the basic integer types"
 * - "X is a slice of byte/rune" means: "X underlying type is slice,
 *   and the the underlying type of slice element is byte (i.e. uint8) or rune (i.e. int32)"
 */
func ConvertibleTo(src *Complete, dst *Complete) bool {
	if src == nil || dst == nil {
		// nil types, cannot check convertibility
		return false
	} else if AssignableTo(src, dst) {
		return true
	} else if IdenticalType(src.Underlying(), dst.Underlying(), true) {
		return true
	}
	k1, k2 := src.kind, dst.kind
	if k1 == PtrKind && k2 == PtrKind && src.Name() == "" && dst.Name() == "" {
		if IdenticalType(src.Elem().Underlying(), dst.Elem().Underlying(), true) {
			return true
		}
	}
	cat1, cat2 := k1.Category(), k2.Category()
	var ok bool
	if (cat1 == Int || cat1 == Float64) && (cat2 == Int || cat2 == Float64) {
		ok = true
	} else if cat1 == Complex128 && cat2 == Complex128 {
		ok = true
	} else if k2 == String && (cat1 == Int /* obsolescent */ || isBytesOrRunes(src)) {
		ok = true
	} else if k1 == String && isBytesOrRunes(dst) {
		ok = true
	}
	return ok
}

// return true if underlying type of c is a slice,
// and the underlying type of its element is byte or rune
func isBytesOrRunes(c *Complete) bool {
	if c.kind != SliceKind {
		return false
	}
	kelem := c.elem.common().kind
	return kelem == Uint8 || kelem == Int32
}

// return true if values of type 'src' implement interface type 'dst'
func Implements(src *Complete, dst *Complete) bool {
	// TODO
	panic("types.Implements() unimplemented")
}
