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
 * flowgraph.hpp
 *
 *  Created on Mar 08, 2021
 *      Author Massimiliano Ghilardi
 */

#ifndef ONEJIT_FLOWGRAPH_HPP
#define ONEJIT_FLOWGRAPH_HPP

#include <onejit/basicblock.hpp>
#include <onestl/array.hpp>

namespace onejit {

// compute global control flow graph of a function
class FlowGraph {

public:
  FlowGraph() noexcept;

  ~FlowGraph() noexcept;

  // return false if out of memory
  bool build(Span<Node> nodes) noexcept;

private:
  static bool is_label(Node node) noexcept;

  Array<BasicBlock> bb_;
  Array<BasicBlock *> link_;

}; // class FlowGraph

} // namespace onejit

#endif // ONEJIT_FLOWGRAPH_HPP
