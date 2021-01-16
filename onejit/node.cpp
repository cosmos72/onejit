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
 * node.cpp
 *
 *  Created on Jan 09, 2020
 *      Author Massimiliano Ghilardi
 */

#include "onejit/node.hpp"
#include "onejit/check.hpp"
#include "onejit/code.hpp"
#include "onejit/const.hpp"
#include "onejit/var.hpp"

namespace onejit {

CodeItem Node::at(Offset byte_offset) const {
  check(code_, !=, nullptr);
  return code_->at(off_or_dir_ + byte_offset);
}

Node Node::child(uint16_t i) const {
  check(i, <, children());
  const CodeItem item =
      code_->uint32(size_t(i) * sizeof(Offset) + off_or_dir_ + sizeof(NodeHeader));

  NodeHeader header;
  uint32_t offset_or_direct;

  if (item <= FALLTHROUGH) {
    offset_or_direct = 0;
    header = NodeHeader{Type(item), Void, 0};
  } else if ((item & 1) == 1) {
    // direct Const
    offset_or_direct = item;
    header = NodeHeader{CONST, Const::kind_from_direct(item), 0};
  } else if ((item & 3) != 0) {
    // direct Var
    offset_or_direct = item;
    header = NodeHeader{VAR, Var::kind_from_direct(item), 0};
  } else {
    // indirect Node
    offset_or_direct = off_or_dir_ + item;
    header = NodeHeader{code_->at(offset_or_direct)};
  }
  return Node{header, offset_or_direct, code_};
}

} // namespace onejit
