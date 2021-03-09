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
 * stmt1.cpp
 *
 *  Created on Jan 18, 2021
 *      Author Massimiliano Ghilardi
 */

#include <onejit/code.hpp>
#include <onejit/func.hpp>
#include <onejit/node/stmt1.hpp>

namespace onejit {
namespace node {

// ============================  Stmt1  ========================================

ONEJIT_NOINLINE Node Stmt1::create(Func &func, Expr arg, OpStmt1 op) noexcept {
  return Base::create_indirect(func,                                   //
                               Header{STMT_1, Void, uint16_t(op)}, //
                               Nodes{&arg, 1});
}

const Fmt &Stmt1::format(const Fmt &fmt, Syntax syntax, size_t depth) const {
  fmt << '(' << op() << ' ';
  return child(0).format(fmt, syntax, depth + 1) << ')';
}

// ============================  Goto  ==================================

// ============================  Inc  ===================================

// ============================  Dec  ===================================

} // namespace node
} // namespace onejit
