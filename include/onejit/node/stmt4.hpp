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
 * stmt4.hpp
 *
 *  Created on Jan 19, 2021
 *      Author Massimiliano Ghilardi
 */

#ifndef ONEJIT_NODE_STMT4_HPP
#define ONEJIT_NODE_STMT4_HPP

#include <onejit/fmt.hpp>
#include <onejit/node/expr.hpp>
#include <onejit/node/stmt.hpp>
#include <onejit/opstmt.hpp>

namespace onejit {

////////////////////////////////////////////////////////////////////////////////
class Stmt4 : public Stmt {
  using Base = Stmt;
  friend class Node;
  friend class Func;

public:
  /**
   * construct an invalid Stmt4.
   * exists only to allow placing Stmt4 in containers
   * and similar uses that require a default constructor.
   *
   * to create a valid Stmt4, use Func::new_stmt4()
   */
  constexpr Stmt4() noexcept : Base{} {
  }

  static constexpr Type type() noexcept {
    return STMT_4;
  }

  constexpr OpStmt4 op() const noexcept {
    return OpStmt4(Base::op());
  }

  static constexpr uint32_t children() noexcept {
    return 4;
  }

  const Fmt &format(const Fmt &out, size_t depth = 0) const;

protected:
  // downcast Node to Stmt4
  constexpr explicit Stmt4(const Node &node) noexcept : Base{node} {
  }

  // downcast helper
  static constexpr bool is_allowed_type(Type t) noexcept {
    return t == STMT_4;
  }

  // used by subclasses
  Stmt4(Func &func, const Node &child0, const Node &child1, //
        const Node &child2, const Node &child3, OpStmt4 op) noexcept
      : Base{create(func, child0, child1, child2, child3, op)} {
  }

private:
  constexpr bool child_result_is_used(uint32_t i) const noexcept {
    return op() == FOR ? i == 1 : true;
  }

  static Node create(Func &func, const Node &child0, const Node &child1, //
                     const Node &child2, const Node &child3, OpStmt4 op) noexcept;
};

////////////////////////////////////////////////////////////////////////////////
class For : public Stmt4 {
  using Base = Stmt4;
  friend class Node;
  friend class Func;

public:
  /**
   * construct an invalid For.
   * exists only to allow placing For in containers
   * and similar uses that require a default constructor.
   *
   * to create a valid For, use one of the other constructors
   */
  constexpr For() noexcept : Base{} {
  }

  For(Func &func, const Node &init, const Expr &test, const Node &post, const Node &body) noexcept
      : Base{func, init, test, post, body, FOR} {
  }

  static constexpr OpStmt4 op() noexcept {
    return FOR;
  }

  // shortcut for child(0)
  Node init() const noexcept {
    return child(0);
  }

  // shortcut for child_is<Expr>(1)
  Expr test() const noexcept {
    return child_is<Expr>(1);
  }

  // shortcut for child(2)
  Node post() const noexcept {
    return child(2);
  }

  // shortcut for child(3)
  Node body() const noexcept {
    return child(3);
  }

private:
  // downcast Node to For
  constexpr explicit For(const Node &node) noexcept : Base{node} {
  }

  // downcast helper
  static constexpr bool is_allowed_op(uint16_t op) noexcept {
    return op == FOR;
  }

  static constexpr bool child_result_is_used(uint32_t i) noexcept {
    return i == 1;
  }
};

} // namespace onejit

#endif // ONEJIT_NODE_STMT4_HPP