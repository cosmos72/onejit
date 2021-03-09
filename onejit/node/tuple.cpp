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
 * tuple.cpp
 *
 *  Created on Jan 28, 2021
 *      Author Massimiliano Ghilardi
 */

#include <onejit/code.hpp>
#include <onejit/func.hpp>
#include <onejit/node/functype.hpp>
#include <onejit/node/label.hpp>
#include <onejit/node/mem.hpp>
#include <onejit/node/tuple.hpp>
#include <onestl/buffer.hpp>

namespace onejit {
namespace node {

// ============================  Tuple  ====================================

Node Tuple::create(Func &func, Kind kind, OpN op, Nodes nodes) noexcept {
  return Base::create_indirect(func,                        //
                               Header{TUPLE, kind, op}, //
                               nodes);
}

const Fmt &Tuple::format(const Fmt &fmt, Syntax syntax, size_t depth) const {
  if (type() == MEM) {
    // may happen, Mem is a subclass of Tuple
    return is<Mem>().format(fmt, syntax, depth);
  }
  OpN op = this->op();
  fmt << '(' << op;

  const bool is_call = op == CALL;
  // if op == CALL, skip child(0) i.e. FuncType
  for (size_t i = size_t(is_call), n = children(); i < n; i++) {
    Node node = child(i);
    if (node) {
      fmt << ' ';
      node.format(fmt, syntax, depth + 1);
    }
  }
  return fmt << ')';
}

} // namespace node
} // namespace onejit
