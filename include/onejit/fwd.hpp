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
 * fwd.hpp
 *
 *  Created on Jan 09, 2021
 *      Author Massimiliano Ghilardi
 */

#ifndef ONEJIT_FWD_HPP
#define ONEJIT_FWD_HPP

#include <onestl/fwd.hpp>

#include <cstdint> // uint*_t

#define ONEJIT_NOINLINE ONESTL_NOINLINE
#define ONEJIT_NORETURN ONESTL_NORETURN

#define ONEJIT_N_OF(array) (sizeof(array) / sizeof(array[0]))

namespace onejit {

template <class T> using Buffer = ::onestl::Buffer<T>;
using Chars = ::onestl::Chars;
template <class T> using Span = ::onestl::Span<T>;
using String = ::onestl::String;
template <class T> using Vector = ::onestl::Vector<T>;
template <class T> using View = ::onestl::View<T>;

class ArchId;
class AssignStmt;
class AssignTupleStmt;
class LogSize;
class BlockStmt;
class BinaryExpr;
class BreakStmt;
class Compiler;
class CallExpr;
class CaseStmt;
class Code;
class CodeParser;
class CondStmt;
class Const;
class ConstExpr;
union ConstFloat32;
union ConstFloat64;
class ContinueStmt;
class DefaultStmt;
class Error;
class Expr;
class FallthroughStmt;
class ForStmt;
class Func;
class FuncType;
class GotoStmt;
class IfStmt;
class JumpIfStmt;
class Label;
class Kind;
class MemExpr;
class Node;
union NodeHeader;
class ReturnStmt;
class Stmt0;
class Stmt1;
class Stmt2;
class Stmt3;
class Stmt4;
class StmtN;
class SwitchStmt;
class Test;
class UnaryExpr;
union Var;
class VarId;
class VarExpr;

using Bits = LogSize;
using CodeItem = uint32_t;
using Offset = uint32_t;
using SimdN = LogSize;

using CaseStmts = View<CaseStmt>;
using CodeItems = View<CodeItem>;
using Exprs = View<Expr>;
using Kinds = View<Kind>;
using Nodes = View<Node>;

} // namespace onejit

#endif // ONEJIT_FWD_HPP
