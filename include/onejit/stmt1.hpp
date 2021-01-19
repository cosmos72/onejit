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
 * stmt1.hpp
 *
 *  Created on Jan 18, 2020
 *      Author Massimiliano Ghilardi
 */

#ifndef ONEJIT_STMT1_HPP
#define ONEJIT_STMT1_HPP

#include <onejit/opstmt.hpp>
#include <onejit/stmt.hpp>

namespace onejit {

////////////////////////////////////////////////////////////////////////////////
class Stmt1 : public Stmt {
  using Base = Stmt;
  friend class Node;
  friend class Func;

public:
  /**
   * construct an invalid Stmt1
   * exists only to allow placing Stmt1 in containers
   * and similar uses that require a default constructor.
   *
   * to create a valid Stmt1, use Func::new_stmt1()
   */
  constexpr Stmt1() : Base{STMT_1, Bad, BAD_ST1} {
  }

  static constexpr uint32_t children() {
    return 1;
  }

  // shortcut for child(0)
  Node body() const {
    return child(0);
  }

protected:
  /* construct an invalid Stmt1 */
  constexpr explicit Stmt1(OpStmt1 op) : Base{STMT_1, Bad, op} {
  }

  // downcast Node to Stmt1
  constexpr explicit Stmt1(const Node &node) : Base{node} {
  }

  // downcast helper
  static constexpr bool is_allowed_type(Type t) {
    return t == STMT_1;
  }

  static Stmt1 create(OpStmt1 op, const Node &body, Code *holder);
};

std::ostream &operator<<(std::ostream &out, const Stmt1 &st);

////////////////////////////////////////////////////////////////////////////////
class DefaultStmt : public Stmt1 {
  using Base = Stmt1;
  friend class Node;
  friend class Func;

public:
  /**
   * construct an invalid DefaultStmt.
   * exists only to allow placing DefaultStmt in containers
   * and similar uses that require a default constructor.
   *
   * to create a valid DefaultStmt, use Func::new_default()
   */
  constexpr DefaultStmt() : Base{DEFAULT} {
  }

  static constexpr OpStmt1 op() {
    return DEFAULT;
  }

private:
  // downcast Node to DefaultStmt
  constexpr explicit DefaultStmt(const Node &node) : Base{node} {
  }

  // downcast helper
  static constexpr bool is_allowed_op(uint16_t op) {
    return op == DEFAULT;
  }

  static DefaultStmt create(const Node &body, Code *holder) {
    return DefaultStmt{Stmt1::create(DEFAULT, body, holder)};
  }
};

} // namespace onejit

#endif // ONEJIT_STMT1_HPP