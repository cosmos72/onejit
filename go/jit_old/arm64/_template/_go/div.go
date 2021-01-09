/*
 * onejit - JIT compiler in Go
 *
 * Copyright (C) 2018-2020 Massimiliano Ghilardi
 *
 * This library is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public
 * License as published by the Free Software Foundation; either
 * version 2 of the License, or (at your option) any later version.
 *
 * This library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public
 * License along with this library; if not, write to the Free Software
 * Foundation, Inc., 59 Temple Place, Suite 330, Boston, MA  02111-1307  USA
 *
 * div.go
 *
 *  Created on Feb 08, 2019
 *      Author Massimiliano Ghilardi
 */

package main

// go:nosplit
func Div8(a int8, b int8) int8 {
	return a / b
}

// go:nosplit
func Div16(a int16, b int16) int16 {
	return a / b
}

// go:nosplit
func Div32(a int32, b int32) int32 {
	return a / b
}

// go:nosplit
func Div64(a int64, b int64) int64 {
	return a / b
}

// go:nosplit
func UDiv8(a uint8, b uint8) uint8 {
	return a / b
}

// go:nosplit
func UDiv16(a uint16, b uint16) uint16 {
	return a / b
}

// go:nosplit
func UDiv32(a uint32, b uint32) uint32 {
	return a / b
}

// go:nosplit
func UDiv64(a uint64, b uint64) uint64 {
	return a / b
}