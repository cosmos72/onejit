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
 * compiler.cpp
 *
 *  Created on Feb 19, 2021
 *      Author Massimiliano Ghilardi
 */

#include <onejit/compiler.hpp>
#include <onejit/func.hpp>
#include <onejit/ir.hpp>
#include <onejit/x64/compiler.hpp>
#include <onejit/x64/mem.hpp>

namespace onejit {

Compiler &Compiler::compile_x64(Func &func, Opt flags) noexcept {
  compile(func, flags);
  if (*this && error_.empty()) {
    // pass our internal buffers node_ and error_ to x64::Compiler
    onejit::x64::Compiler{}.compile(func, allocator_, node_, flowgraph_, error_, //
                                    flags, abi_autodetect(abi_));
  }
  return *this;
}

// ===============================  x64::Compiler  =============================

namespace x64 {

Compiler::operator bool() const noexcept {
  return good_ && func_ && *func_;
}

Compiler &Compiler::compile(Func &func, reg::Allocator &allocator, Array<Node> &node_vec,
                            FlowGraph &flowgraph, Array<Error> &error_vec, Opt flags,
                            Abi abi) noexcept {
  if (func.get_compiled(X64)) {
    // already compiled for x86_64
    return *this;
  }
  Node node = func.get_compiled(NOARCH);
  if (!node) {
    error_vec.append(Error{node, "function not compiled yet, cannot materialize it for x64 arch"});
    return *this;
  }

  node_vec.clear();
  func_ = &func;
  allocator_ = &allocator;
  node_ = &node_vec;
  flowgraph_ = &flowgraph;
  error_ = &error_vec;
  flags_ = flags;
  good_ = bool(func);

  return compile(node).allocate_regs(abi).finish();
}

Compiler &Compiler::allocate_regs(Abi abi) noexcept {
  Vars vars = func_->vars();
  if (allocator_->reset(vars.size())) {
    fill_interference_graph();
    set_reg_hints(abi);
    // x86_64 has 16 general registers, we reserve RSP and RBX
    allocator_->allocate_regs(14);
  }
  return *this;
}

Compiler &Compiler::fill_interference_graph() noexcept {
  if (!flowgraph_->build(*node_, *error_)) {
    good_ = false;
    return *this;
  }
  BitSet &live = allocator_->get_bitset();
  live.fill(false);
  // TODO propagate liveness analysis across jumps
  for (size_t i = node_->size(); i > 0; i--) {
    Node node = (*node_)[i - 1];
    update_live_regs(live, node);
  }
  return *this;
}

void Compiler::update_live_regs(BitSet & /*live*/, Node /*node*/) noexcept {
}

Compiler &Compiler::set_reg_hints(Abi abi) noexcept {
  (void)abi; // TODO implement
  return *this;
}

Compiler &Compiler::finish() noexcept {
  if (*this && node_) {
    Node compiled;
    switch (node_->size()) {
    case 0:
      compiled = VoidConst;
      break;
    case 1:
      compiled = (*node_)[0];
      break;
    default:
      compiled = Block{*func_, *node_};
      break;
    }
    func_->set_compiled(X64, compiled);
  }
  return *this;
}

// ===============================  compile(Node)  =============================

Compiler &Compiler::compile(Node node) noexcept {
  const Type t = node.type();
  switch (t) {
  case STMT_0:
    return error(node, "unexpected Stmt0");
  case STMT_1:
    return compile(node.is<Stmt1>());
  case STMT_2:
    return compile(node.is<Stmt2>());
  case STMT_3:
    return error(node, "unexpected Stmt3");
  case STMT_4:
    return error(node, "unexpected Stmt4");
  case STMT_N:
    return compile(node.is<StmtN>());
  default:
    if (Expr expr = node.is<Expr>()) {
      return compile(expr);
    }
    return error(node, "unexpected Node");
  }
}

// ===============================  compile(Expr)  =============================

Compiler &Compiler::compile(Expr expr) noexcept {
  return add(simplify(expr));
}

Expr Compiler::simplify(Expr expr) noexcept {
  switch (expr.type()) {
  case VAR:
  case LABEL:
  case CONST:
    return expr;
  case MEM:
    return simplify(expr.is<onejit::Mem>());
  case UNARY:
    return simplify(expr.is<Unary>());
  case BINARY:
    return simplify(expr.is<Binary>());
  case TUPLE:
    return simplify(expr.is<Tuple>());
  default:
    error(expr, "unexpected Expr");
    return expr;
  }
}

Expr Compiler::simplify(onejit::Mem expr) noexcept {
  return expr; // TODO
}

Expr Compiler::simplify(Unary expr) noexcept {
  Expr x = expr.x();
  Expr simpl_x = to_var_mem_const(simplify(x));
  if (simpl_x != x) {
    return Unary{*func_, expr.kind(), expr.op(), simpl_x};
  }
  return expr;
}

Expr Compiler::simplify(Binary expr) noexcept {
  Expr x = expr.x(), y = expr.y();
  Expr simpl_x = x, simpl_y = y;
  simplify_binary(simpl_x, simpl_y);
  if (simpl_x != x || simpl_y != y) {
    return Binary{*func_, expr.op(), simpl_x, simpl_y};
  }
  return expr;
}

void Compiler::simplify_binary(Expr &x, Expr &y) noexcept {
  Expr simpl_x = to_var_mem_const(simplify(x));
  Expr simpl_y = to_var_mem_const(simplify(y));
  if (simpl_x.type() == MEM && simpl_y.type() == MEM) {
    // both arguments are memory.
    // not supported by x86_64 assembly, force one to register
    simpl_y = to_var(simpl_y);
  }
  x = simpl_x;
  y = simpl_y;
}

Expr Compiler::simplify(Tuple expr) noexcept {
  return expr; // TODO
}

// ===============================  compile(Stmt1)  ============================

Compiler &Compiler::compile(Stmt1 st) noexcept {
  static const OpStmt1 gen_st1[] = {X86_JMP, X86_INC, X86_DEC};
  static const OpStmt1 cond_jump[] = {X86_JA, X86_JAE, X86_JB, X86_JBE, X86_JE,
                                      X86_JG, X86_JGE, X86_JL, X86_JLE, X86_JNE};
  OpStmt1 op = st.op();
  if (op >= GOTO && op <= DEC) {
    op = gen_st1[op - GOTO];
  } else if (op >= ASM_JA && op <= ASM_JNE) {
    op = cond_jump[op - ASM_JA];
  } else {
    return error(st, "unexpected Stmt1 operation");
  }
  Expr expr = simplify(st.arg());
  return add(Stmt1{*func_, expr, op});
}

// ===============================  compile(Stmt2)  ============================

Compiler &Compiler::compile(Stmt2 st) noexcept {
  switch (st.op()) {
  case JUMP_IF:
    return compile(st.is<JumpIf>());
  case ASM_CMP: {
    Expr x = st.child_is<Expr>(0);
    Expr y = st.child_is<Expr>(1);
    simplify_binary(x, y);
    return add(Stmt2{*func_, x, y, X86_CMP});
  }
  default:
    if (Assign assign = st.is<Assign>()) {
      return compile(assign);
    }
    return error(st, "unexpected Stmt2");
  }
}

Compiler &Compiler::compile(Assign st) noexcept {
  Expr src = st.src(), dst = st.dst();
  // simplify src first: its side effects, if any, must be applied before dst
  //
  // do not call to_var_mem_const(simplify(src)): it creates redundant Vars
  src = simplify(src);
  dst = to_var_mem_const(simplify(dst));
  if (src.type() == MEM && dst.type() == MEM) {
    // both arguments are memory.
    // not supported by x86_64 assembly, force one to register
    src = to_var(src);
  }
  return add(simplify_assign(st, dst, src));
}

Node Compiler::simplify_assign(Assign st, Expr dst, Expr src) noexcept {
  // FIXME use st.kind() to select appropriate x86_64 instruction
  OpStmt2 op = st.op();
  if (op >= ADD_ASSIGN && op <= SHR_ASSIGN) {
    static const OpStmt2 xop[] =
        // FIXME: X86_DIV computes both quotient and remainder
        {X86_ADD, X86_SUB, X86_MUL, X86_DIV, REM_ASSIGN, //
         X86_AND, X86_OR,  X86_XOR, X86_SHL, X86_SHR};
    op = xop[op - ADD_ASSIGN];
    src = to_var_mem_const(src);
  } else if (op == ASSIGN) {
    switch (src.type()) {
    case VAR:
    case MEM:
    case CONST:
    case LABEL:
      op = X86_MOV;
      break;
    case UNARY:
      return simplify_assign(st, dst, src.is<Unary>());
    case BINARY:
      return simplify_assign(st, dst, src.is<Binary>());
    case TUPLE:
      return simplify_assign(st, dst, src.is<Tuple>());
    default:
      error(st, "unexpected Assign right argument");
      break;
    }
  } else {
    error(st, "unexpected Assign operation");
  }
  return Stmt2{*func_, dst, src, op};
}

Node Compiler::simplify_assign(Assign st, Expr dst, Unary src) noexcept {
  // TODO
  (void)dst;
  (void)src;
  return st;
}

Node Compiler::simplify_assign(Assign st, Expr dst, Tuple src) noexcept {
  switch (src.op()) {
  case ADD:
    if (dst.type() == VAR) {
      if (Mem mem{*this, Ptr, ChildRange{src}}) {
        return Stmt2{*func_, dst, mem, X86_LEA};
      }
    }
    break;
  case CALL: {
    StmtN set{*func_, Nodes{&dst, 1}, SET_};
    ChildRange ranges[] = {ChildRange{src, 0, 2},          // ftype, label
                           ChildRange{Block{*func_, set}}, // (set_ dst)
                           ChildRange{src, 2, sub_uint32(src.children(), 2)}};
    return StmtN{*func_, ChildRanges{ranges, 3}, X86_CALL_};
  }
  default:
    break;
  } // switch
  return st;
}

Node Compiler::simplify_assign(Assign st, Expr dst, Binary src) noexcept {
  // TODO
  (void)dst;
  (void)src;
  return st;
}

// ===============================  compile(StmtN)  ============================

Compiler &Compiler::compile(StmtN st) noexcept {
  switch (st.op()) {
  case ASSIGN_CALL:
    return compile(st.is<AssignCall>());
  case BLOCK:
    return compile(st.is<Block>());
  case RETURN:
    return compile(st.is<Return>());
  case SET_: // used in function prologue
    return add(st);
  default:
    return error(st, "unexpected StmtN");
  }
}

Compiler &Compiler::compile(Block st) noexcept {
  for (size_t i = 0, n = st.children(); i < n; i++) {
    compile(st.child(i));
  }
  return *this;
}

Compiler &Compiler::compile(AssignCall st) noexcept {
  return add(st); // TODO
}

Compiler &Compiler::compile(Return st) noexcept {
  ChildRange children{st, 0, st.children()};
  return add(StmtN{*func_, ChildRanges{&children, 1}, X86_RET});
}

////////////////////////////////////////////////////////////////////////////////

Var Compiler::to_var(Expr expr) noexcept {
  Var v = expr.is<Var>();
  if (expr && !v) {
    // copy Expr result to a Var
    v = Var{*func_, expr.kind()};
    // compile(Assign{...}) would cause infinite recursion
    add(Assign{*func_, ASSIGN, v, expr});
  }
  return v;
}

Expr Compiler::to_var_const(Expr expr) noexcept {
  switch (expr.type()) {
  case VAR:
  case CONST:
  case LABEL:
    return expr;
  default:
    return to_var(expr);
  }
}

Expr Compiler::to_var_mem_const(Expr expr) noexcept {
  switch (expr.type()) {
  case VAR:
  case MEM:
  case CONST:
  case LABEL:
    return expr;
  default:
    return to_var(expr);
  }
}

Compiler &Compiler::add(Node node) noexcept {
  if (node != VoidConst) {
    good_ = good_ && node_ && node_->append(node);
  }
  return *this;
}

Compiler &Compiler::error(Node where, Chars msg) noexcept {
  good_ = good_ && error_ && error_->append(Error{where, msg});
  return *this;
}

Compiler &Compiler::out_of_memory(Node where) noexcept {
  // always set good_ to false
  good_ = good_ && error_ && error_->append(Error{where, "out of memory"}) && false;
  return *this;
}

} // namespace x64
} // namespace onejit
