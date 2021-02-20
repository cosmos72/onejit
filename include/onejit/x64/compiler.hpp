/*
 * onejit - in-memory assembler
 *
 * Copyright (C) 2021 Massimiliano Ghilardi
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
 * compiler.hpp
 *
 *  Created on Feb 19, 2021
 *      Author Massimiliano Ghilardi
 */
#ifndef ONEJIT_X64_COMPILER_HPP
#define ONEJIT_X64_COMPILER_HPP

#include <onejit/error.hpp>
#include <onestl/vector.hpp>

namespace onejit {
namespace x64 {

////////////////////////////////////////////////////////////////////////////////
// compiles code from portable intermediate representation
// (produced by onejit::Compiler::compile()) to x86_64 assembler
//
// usually not invoked directly - the public API is onejit::Compiler::x64()
class Compiler {
  friend class ::onejit::Compiler;

public:
  constexpr Compiler() noexcept //
      : func_{}, node_{}, error_{}, flags_{}, good_{true} {
  }

  Compiler(Compiler &&other) noexcept = default;
  Compiler &operator=(Compiler &&other) noexcept = default;

  ~Compiler() noexcept = default;

  // return false if out of memory
  explicit operator bool() const noexcept;

private:
  // private, use onejit::Compiler::x64() instead
  Compiler &compile(Func &func, Vector<Node> &node, Vector<Error> &error, Opt flags) noexcept;

  Compiler &compile(Assign stmt) noexcept;
  Compiler &compile(AssignCall stmt) noexcept;
  Compiler &compile(Block stmt) noexcept;
  Compiler &compile(Expr expr) noexcept;
  Compiler &compile(Node node) noexcept;
  Compiler &compile(Return stmt) noexcept;
  Compiler &compile(Stmt1 stmt) noexcept;
  Compiler &compile(Stmt2 stmt) noexcept;
  Compiler &compile(StmtN stmt) noexcept;

  Expr simplify(Binary expr) noexcept;
  Expr simplify(Expr expr) noexcept;
  Expr simplify(onejit::Mem expr) noexcept;
  Expr simplify(Tuple expr) noexcept;
  Expr simplify(Unary expr) noexcept;
  void simplify_binary(Expr &x, Expr &y) noexcept;

  // if expr is a Var, does nothing and returns it.
  // otherwise copies expression result to a new local variable and returns it.
  Var to_var(Expr expr) noexcept;

  // if expr is a Var or Const, does nothing and returns it
  // otherwise copies expression result to a new local variable and returns it.
  Expr to_var_const(Expr expr) noexcept;

  // if expr is a Var, Mem or Const, does nothing and returns it
  // otherwise copies expression result to a new local variable and returns it.
  Expr to_var_mem_const(Expr expr) noexcept;

  // add an already compiled node to compiled list
  Compiler &add(const Node &node) noexcept;

  // store compiled code into function.set_compiled(X64)
  // invoked by compile(Func)
  Compiler &finish() noexcept;

  Compiler &error(const Node &where, Chars msg) noexcept;

  Func *func_;
  Vector<Node> *node_;
  Vector<Error> *error_;
  Opt flags_;
  bool good_; // !good_ means out of memory
};

} // namespace x64
} // namespace onejit

#endif // ONEJIT_X64_COMPILER_HPP