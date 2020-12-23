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
 * kind.go
 *
 *  Created on Jan 24, 2019
 *      Author Massimiliano Ghilardi
 */

package common

import (
	"unsafe"
)

// intentionally matches reflect.Kind values
type Kind uint8

const (
	Void Kind = iota // Invalid
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128
	_ // Array
	_ // Chan
	_ // Func
	_ // Interface
	_ // Map
	Ptr
	_   // Slice
	_   // String
	_   // Struct
	_   // UnsafePointer
	KLo = Bool
	KHi = Complex128
)

var kstring = [...]string{
	Bool:       "bool",
	Int:        "int",
	Int8:       "int8",
	Int16:      "int16",
	Int32:      "int32",
	Int64:      "int64",
	Uint:       "uint",
	Uint8:      "uint8",
	Uint16:     "uint16",
	Uint32:     "uint32",
	Uint64:     "uint64",
	Uintptr:    "uintptr",
	Float32:    "float32",
	Float64:    "float64",
	Complex64:  "complex64",
	Complex128: "complex128",
	/*
		Array:         "array",
		Chan:          "chan",
		Func:          "func",
		Interface:     "interface",
		Map:           "map",
	*/
	Ptr: "ptr",
	/*
		Slice:         "slice",
		String:        "string",
		Struct:        "struct",
		UnsafePointer: "unsafe.Pointer",
	*/
}

var ksize = [...]Size{
	Bool:       1,
	Int:        Size(unsafe.Sizeof(int(0))),
	Int8:       1,
	Int16:      2,
	Int32:      4,
	Int64:      8,
	Uint:       Size(unsafe.Sizeof(uint(0))),
	Uint8:      1,
	Uint16:     2,
	Uint32:     4,
	Uint64:     8,
	Uintptr:    Size(unsafe.Sizeof(uintptr(0))),
	Float32:    4,
	Float64:    8,
	Complex64:  8,
	Complex128: 16,
	/*
		Chan:          Size(unsafe.Sizeof(make(chan int))),
		Func:          Size(unsafe.Sizeof(func() {})),
		Interface:     Size(unsafe.Sizeof(assertError)),
		Map:           Size(unsafe.Sizeof(map[int]int{})),
	*/
	Ptr: Size(unsafe.Sizeof((*int)(nil))),
	/*
		Slice:         Size(unsafe.Sizeof([]int{})),
		String:        Size(unsafe.Sizeof("")),
		UnsafePointer: Size(unsafe.Sizeof((unsafe.Pointer)(nil))),
	*/
}

func (k Kind) Size() Size {
	if k >= KLo && k <= KHi {
		return ksize[k]
	}
	return 0
}

func (k Kind) Category() Kind {
	switch k {
	case Int, Int8, Int16, Int32, Int64:
		return Int
	case Uint, Uint8, Uint16, Uint32, Uint64, Uintptr:
		return Uint
	case Float32, Float64:
		return Float64
	case Complex64, Complex128:
		return Complex128
	default:
		return k
	}
}

func (k Kind) Signed() bool {
	return k.IsInteger() || k.IsFloat()
}

func (k Kind) IsInteger() bool {
	switch k {
	case Int, Int8, Int16, Int32, Int64, Uint, Uint8, Uint16, Uint32, Uint64, Uintptr:
		return true
	default:
		return false
	}
}

func (k Kind) isIntOrPtr() bool {
	return k == Ptr || k.IsInteger()
}

func (k Kind) IsFloat() bool {
	return k == Float32 || k == Float64
}

func (k Kind) IsComplex() bool {
	return k == Complex64 || k == Complex128
}

func (k Kind) IsOrdered() bool {
	return k.IsInteger() || k.IsFloat()
}

func (k Kind) IsComparable() bool {
	return k.IsInteger() || k.IsFloat()
}

func (k Kind) mustBeBool(op Op) Kind {
	if k != Bool {
		badOpKind(op, k)
	}
	return k
}

func (k Kind) mustBeInteger(op Op) Kind {
	if !k.isIntOrPtr() {
		badOpKind(op, k)
	}
	return k
}

func (k Kind) mustBeNumber(op Op) Kind {
	if k == Bool {
		badOpKind(op, k)
	}
	return k
}

func kindMustBeInteger(op Op, k1 Kind, k2 Kind) Kind {
	if k1 == k2 && k1.isIntOrPtr() {
		return k1
	}
	return badOpKind2(op, k1, k2)
}

func kindMustBeOrdered(op Op, k1 Kind, k2 Kind) Kind {
	if k1 == k2 && k1.IsOrdered() {
		return k1
	}
	return badOpKind2(op, k1, k2)
}

func kindMustBeComparable(op Op, k1 Kind, k2 Kind) Kind {
	if k1 == k2 && k1.IsComparable() {
		return k1
	}
	return badOpKind2(op, k1, k2)
}

func kindMustBeNumber(op Op, k1 Kind, k2 Kind) Kind {
	if k1 == k2 && k1 != Bool {
		return k1
	}
	return badOpKind2(op, k1, k2)
}
