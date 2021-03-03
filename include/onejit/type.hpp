/*
 * onejit - JIT compiler in C++
 *
 * Copyright (C) 2018-2021 Massimiliano Ghilardi
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
 * type.hpp
 *
 *  Created on Jan 09, 2021
 *      Author Massimiliano Ghilardi
 */

#ifndef ONEJIT_TYPE_HPP
#define ONEJIT_TYPE_HPP

#include <onejit/fwd.hpp>

#include <cstdint> // uint8_t

namespace onejit {

enum Type : uint8_t {
  STMT_0 = 0, // BAD, BREAK, CONTINUE, FALLTHROUGH
  STMT_1 = 1, // GOTO, INC, DEC
  STMT_2 = 2, // *ASSIGN, CASE, DEFAULT, JUMP_IF
  STMT_3 = 3, // IF
  STMT_4 = 4, // FOR
  STMT_N = 5, // ASSIGN_CALL, BLOCK, COND, RETURN, SWITCH

  // Expr
  VAR = 6,
  MEM = 7, // subclass of Tuple
  UNARY = 8,
  BINARY = 9,
  TUPLE = 10,
  // Optimizer assumes expression types below this line CANNOT be optimized further
  LABEL = 11, // jump destination
  CONST = 12,

  FTYPE = 13, // function/method type signature
  NAME = 14,  // name of function/globalvar/globalconst
};

// true if Type t indicates an arbitrary number of child Node:s
constexpr inline bool is_list(Type t) noexcept {
  return t == STMT_N || t == MEM || t == TUPLE || t == FTYPE;
}

uint32_t to_children(Type t) noexcept;

const Chars to_string(Type t) noexcept;

const Fmt &operator<<(const Fmt &fmt, Type t);

} // namespace onejit

#endif // ONEJIT_TYPE_HPP
