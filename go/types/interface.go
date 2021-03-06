/*
 * Copyright (C) 2021 Massimiliano Ghilardi
 *
 *     This Source Code Form is subject to the terms of the Mozilla Public
 *     License, v. 2.0. If a copy of the MPL was not distributed with this
 *     file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 *
 * interface.go
 *
 *  Created on: Apr 01, 2021
 *      Author: Massimiliano Ghilardi
 */

package types

import (
	"github.com/cosmos72/onejit/go/io"
	"github.com/cosmos72/onejit/go/sort"
)

type Interface struct {
	interfaceTag struct{} // occupies zero bytes
	rtype        Complete
	extra        extra
}

// *Interface implements Type

func (t *Interface) String() string {
	_ = t.interfaceTag
	var b builder
	t.WriteTo(&b, fullPkgPath)
	return b.String()
}

func (t *Interface) Underlying() Type {
	return t
}

func (t *Interface) common() *Complete {
	return &t.rtype
}

func (t *Interface) complete() {
	if t.rtype.hash == unknownHash {
		t.rtype.hash = computeInterfaceHash(t.extra.methods)
	}
}

func (t *Interface) WriteTo(dst io.StringWriter, flag verbose) {
	if flag == shortPkgName {
		dst.WriteString(t.rtype.str)
		return
	}
	writeInterfaceTo(dst, t.extra.methods, flag)
}

// *Interface specific method

// return number of embedded interfaces.
func (t *Interface) NumEmbedded() int {
	return len(t.extra.types)
}

// return number of method, including methods from embedded interfaces
func (t *Interface) NumMethod() int {
	return len(t.extra.methods)
}

// return i-th embedded interface.
func (t *Interface) Embedded(i int) Type {
	return t.extra.types[i]
}

// return i-th method, including methods from embedded interfaces
func (t *Interface) Method(i int) Method {
	return t.extra.methods[i]
}

var interfaceMap = map[interface{}]*Interface{}

// create a new Interface type
func NewInterface(embedded []Type, method ...Method) *Interface {
	key, method := makeInterfaceKey(embedded, method)
	t := interfaceMap[key]
	if t != nil {
		return t
	}
	size := sizeOfPtr()
	t = &Interface{
		rtype: Complete{
			size:  2 * size,
			align: uint16(size),
			flags: (flagsAndMethod(method) & flagComplete) | flagComparable,
			kind:  InterfaceKind,
			hash:  computeInterfaceHash(method),
			str:   makeInterfaceString(method, shortPkgName),
		},
		extra: extra{
			types:   embedded,
			methods: method,
		},
	}
	t.rtype.typ = t
	t.rtype.extra = &t.extra
	interfaceMap[key] = t
	return t
}

// return bitwise AND of specified method type's flags
func flagsAndMethod(list []Method) flags {
	ret := ^flags(0)
	for i := range list {
		ret &= list[i].Type.common().flags
	}
	return ret
}

func computeInterfaceHash(methods []Method) hash {
	h := hashInit().String("interface").Int(len(methods))
	for i := range methods {
		mtd := &methods[i]
		h := mtd.hash(h)
		if h == unknownHash {
			break
		}
	}
	return h
}

func makeInterfaceString(methods []Method, flag verbose) string {
	var b builder
	writeInterfaceTo(&b, methods, flag)
	return b.String()
}

func writeInterfaceTo(dst io.StringWriter, methods []Method, flag verbose) {
	dst.WriteString("interface {")
	for i := range methods {
		if i == 0 {
			dst.WriteString(" ")
		} else {
			dst.WriteString("; ")
		}
		methods[i].WriteTo(dst, flag)
	}
	if len(methods) == 0 {
		dst.WriteString("}")
	} else {
		dst.WriteString(" }")
	}
}

type (
	interfaceKey0 = struct{}
	interfaceKey2 = struct {
		method [2]Method
	}
	interfaceKey4 = struct {
		method [4]Method
	}
	interfaceKey16 = struct {
		method [16]Method
	}
	interfaceKey64 = struct {
		method [64]Method
	}
	interfaceKey256 = struct {
		method [256]Method
	}
	interfaceKey1k = struct {
		method [1024]Method
	}
	interfaceKey4k = struct {
		method [4096]Method
	}
	interfaceKey16k = struct {
		method [16384]Method
	}
	interfaceKey64k = struct {
		method [65536]Method
	}
)

func makeInterfaceKey(embedded []Type, method []Method) (ret interface{}, rmethod []Method) {
	rmethod = computeInterfaceMethods(embedded, method)
	n := len(rmethod)
	if n == 0 {
		ret = interfaceKey0{}
	} else if n <= 2 {
		key := interfaceKey2{}
		copy(key.method[:], rmethod)
		ret = key // copies key, must be done after copy
	} else if n <= 4 {
		key := interfaceKey4{}
		copy(key.method[:], rmethod)
		ret = key
	} else if n <= 16 {
		key := interfaceKey16{}
		copy(key.method[:], rmethod)
		ret = key
	} else if n <= 64 {
		key := interfaceKey64{}
		copy(key.method[:], rmethod)
		ret = key
	} else if n <= 256 {
		key := interfaceKey256{}
		copy(key.method[:], rmethod)
		ret = key
	} else if n <= 1024 {
		key := interfaceKey1k{}
		copy(key.method[:], rmethod)
		ret = key
	} else if n <= 4096 {
		key := interfaceKey4k{}
		copy(key.method[:], rmethod)
		ret = key
	} else if n <= 16384 {
		key := interfaceKey16k{}
		copy(key.method[:], rmethod)
		ret = key
	} else if n <= 65536 {
		key := interfaceKey64k{}
		copy(key.method[:], rmethod)
		ret = key
	} else {
		panic("NewInterface: too many methods")
	}
	return ret, rmethod

}

func computeInterfaceMethods(embedded []Type, direct []Method) []Method {
	out := append([]Method(nil), direct...)
	for _, t := range embedded {
		checkEmbeddedInterface(t)
		out = append(out, t.common().extra.methods...)
	}
	if len(out) > 1 {
		sortMethods(out)
	}
	n := 0
	for i := range out {
		mtd := &out[i]
		mtd.Index = n
		if _, ok := mtd.Type.(*Func); !ok {
			panic("NewInterface: invalid method type, expecting *Func")
		}
		if n == 0 || checkDifferentMethod(&out[n-1], mtd) {
			out[n] = *mtd
			n++
		}
	}
	return out[:n]
}

func checkEmbeddedInterface(t Type) {
	u := t.Underlying()
	if u == nil {
		panic("NewInterface: embedded interface " + t.String() + " is incomplete, has nil underlying type")
	}
	if iface, _ := u.(*Interface); iface == nil {
		panic("NewInterface: invalid embedded interface " + t.String() + ", its underlying type must be *Interface")
	}
}

func checkDifferentMethod(m1 *Method, m2 *Method) bool {
	if m1.Name != m2.Name || m1.PkgPath != m2.PkgPath {
		return true
	} else if m1.Type == m2.Type {
		return false
	}
	path := m1.PkgPath
	if len(path) != 0 {
		path += "."
	}
	panic("NewInterface: found duplicate method " + path + m1.Name + " with conflicting types")
}

// sortableMethods

type sortableMethods []Method

func (ms sortableMethods) Len() int {
	return len(ms)
}

func (ms sortableMethods) Less(i, j int) bool {
	pkg1, pkg2 := ms[i].PkgPath, ms[j].PkgPath
	return pkg1 < pkg2 || (pkg1 == pkg2 && ms[i].Name < ms[j].Name)
}

func (ms sortableMethods) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

func sortMethods(ms []Method) {
	sort.Sort(sortableMethods(ms))
}
