/*
 * Copyright (C) 2021 Massimiliano Ghilardi
 *
 *     This Source Code Form is subject to the terms of the Mozilla Public
 *     License, v. 2.0. If a copy of the MPL was not distributed with this
 *     file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * file.go
 *
 *  Created on: Mar 23, 2021
 *      Author: Massimiliano Ghilardi
 */

package token

import (
	"fmt"
	"sort"
)

// functionally equivalent to go/token.File
type File struct {
	name      string
	skipLines int // lines <= skipLine are not represented
	size      int // size < 0 if unknown
	// line[i-1-skipLine] is the offset assigned to i-th line, column 1
	// Note: i must be > skipLines
	line []int
}

func NewFile(name string, base int) *File {
	return &File{name: name, size: -1, line: []int{base}}
}

func (f *File) Name() string {
	if f == nil {
		return ""
	}
	return f.name
}

func (f *File) Base() int {
	if f == nil {
		return 0
	}
	return f.line[0]
}

// return size, or < 0 if unknown
func (f *File) Size() int {
	if f == nil {
		return 0
	}
	return f.size
}

func (f *File) SetSize(size int) {
	if f != nil {
		f.size = size
	}
}

func (f *File) SetSkipLines(skipLines int) {
	if f != nil {
		f.skipLines = skipLines
	}
}

func (f *File) AddLine(offset int) {
	if f == nil || offset <= f.line[len(f.line)-1] {
		return
	}
	f.line = append(f.line, offset)
}

// return the line number for the given position p
func (f *File) Line(p Pos) int {
	return f.Position(p).Line
}

func (*File) Offset(p Pos) int {
	return int(p)
}

func (f *File) Pos(offset int) Pos {
	return Pos(offset)
}

// return the Position for the given position p
func (f *File) Position(p Pos) Position {
	pos := Position{Offset: int(p)}
	if f != nil && p != NoPos {
		base := f.line[0]
		if f.size < 0 {
			if int(p) < base {
				panic(fmt.Errorf("invalid Pos value %d, expecting >= %d", int(p), base))
			}
		} else if int(p) < base || int(p)-base > f.size {
			panic(fmt.Errorf("invalid Pos value %d, expecting %d...%d", int(p), base, base+f.size))
		}
		index := sort.SearchInts(f.line, int(p))
		if index == len(f.line) || int(p) < f.line[index] {
			index--
		}
		pos.Filename = f.name
		pos.Line = 1 + index + f.skipLines
		pos.Column = 1 + int(p) - f.line[index]
	}
	return pos
}

func (f *File) PositionFor(p Pos, _ bool) Position {
	return f.Position(p)
}