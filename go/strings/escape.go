/*
 * Copyright (C) 2021 Massimiliano Ghilardi
 *
 *     This Source Code Form is subject to the terms of the Mozilla Public
 *     License, v. 2.0. If a copy of the MPL was not distributed with this
 *     file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 *
 * escape.go
 *
 *  Created on: Apr 02, 2021
 *      Author: Massimiliano Ghilardi
 */

package strings

import (
	"unicode/utf8"

	"github.com/cosmos72/onejit/go/io"
)

// interpret \ inside a Go literal string
func Unescape(str string) string {
	pos := IndexByte(str, '\\')
	if pos < 0 {
		return str
	}
	var b Builder
	b.Grow(len(str))
	for pos >= 0 {
		b.WriteString(str[:pos])
		pos++ // skip \
		if len(str) <= pos {
			break
		}
		curr := str[pos]
		str = str[pos:]
		skip := 1
		unicode := int32(-1)
		ch := int16(-1)
		switch curr {
		case 'a':
			ch = '\a'
		case 'b':
			ch = '\b'
		case 'f':
			ch = '\f'
		case 'n':
			ch = '\n'
		case 'r':
			ch = '\r'
		case 't':
			ch = '\t'
		case 'v':
			ch = '\v'
		case '\\', '\'', '"':
			ch = int16(curr)
		case 'U':
			unicode = hexDigitsToInt(str[1:], 8)
			skip = 9
		case 'u':
			unicode = hexDigitsToInt(str[1:], 4)
			skip = 5
		case 'x':
			ch = int16(hexDigitsToInt(str[1:], 2))
			skip = 3
		default:
			if curr >= '0' && curr <= '7' {
				ch = int16(octalDigitsToInt(str, 3, 255))
				skip = 3
			}
		}
		str = str[skip:]
		if unicode >= 0 {
			b.WriteRune(unicode)
		} else if ch >= 0 {
			b.WriteByte(uint8(ch))
		}
		pos = IndexByte(str, '\\')
	}
	return b.String()
}

func hexDigitsToInt(str string, n int) int32 {
	avail := len(str)
	if avail < n {
		return -1
	}
	var hex int32
	for i := 0; i < n; i++ {
		hex = (hex << 4) | hexDigitToInt(str[i])
	}
	return hex
}

func octalDigitsToInt(str string, n int, max int32) int32 {
	avail := len(str)
	if avail < n {
		return -1
	}
	var oct int32
	for i := 0; i < n; i++ {
		oct = (oct << 3) | octalDigitToInt(str[i])
	}
	if oct > max {
		oct = -1
	}
	return oct
}

func hexDigitToInt(ch byte) int32 {
	if ch >= '0' && ch <= '9' {
		return int32(ch - '0')
	} else if ch >= 'A' && ch <= 'F' {
		return int32(ch-'A') + 10
	} else if ch >= 'a' && ch <= 'f' {
		return int32(ch-'a') + 10
	} else {
		return -1
	}
}

func octalDigitToInt(ch byte) int32 {
	if ch >= '0' && ch <= '7' {
		return int32(ch - '0')
	} else {
		return -1
	}
}

// quote a Go literal string
func Escape(str string) string {
	var b Builder
	WriteQuotedString(&b, str)
	return b.String()
}

// quote a Go literal string and write it to out
func WriteQuotedString(out io.StringWriter, str string) {
	consumed := 1
	for i, n := 0, len(str); i < n; i += consumed {
		consumed = 1
		ch := str[i]
		var esc string
		switch ch {
		case '\a':
			esc = `\a`
		case '\b':
			esc = `\b`
		case '\f':
			esc = `\f`
		case '\n':
			esc = `\n`
		case '\r':
			esc = `\r`
		case '\t':
			esc = `\t`
		case '\v':
			esc = `\v`
		case '\\':
			esc = `\\`
		case '"':
			esc = `\"`
		default:
			if ch >= utf8.RuneSelf {
				r, rconsumed := utf8.DecodeRuneInString(str[i:])
				if r == utf8.RuneError && consumed <= 1 {
				} else if r <= 0xffff {
					writeHexEscape(out, `\u`, 4, uint32(r))
					consumed = rconsumed
					continue
				} else {
					writeHexEscape(out, `\U`, 8, uint32(r))
					consumed = rconsumed
					continue
				}
			}
			if ch < ' ' || ch > '~' {
				writeHexEscape(out, `\x`, 2, uint32(ch))
				continue
			}
		}
		if len(esc) == 0 {
			esc = str[i : i+1]
		}
		out.WriteString(esc)
	}
}

func writeHexEscape(out io.StringWriter, prefix string, len int, number uint32) {
	out.WriteString(prefix)
	for shift := (len - 1) * 4; shift >= 0; shift -= 4 {
		halfbyte := 0x0f & uint8(number>>uint8(shift))
		str := "0123456789abcdef"[halfbyte : halfbyte+1]
		out.WriteString(str)
	}
}
