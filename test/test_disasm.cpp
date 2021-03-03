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
 * test_disasm.cpp
 *
 *  Created on Feb 04, 2021
 *      Author Massimiliano Ghilardi
 */

#include "test_disasm.hpp"

#include <onejit/assembler.hpp>
#include <onejit/fmt.hpp>

#ifdef HAVE_LIBCAPSTONE
#include <arm64.h>
#include <capstone.h>
#include <x86.h>
#endif

namespace onejit {

#ifdef HAVE_LIBCAPSTONE

TestDisasm::TestDisasm() noexcept
    : testcount_{}, handle_{}, err_{cs_open(CS_ARCH_X86, CS_MODE_64, &handle_)} {
  if (err_ != CS_ERR_OK) {
    return;
  }
  cs_option(handle_, CS_OPT_DETAIL, CS_OPT_ON);
}

TestDisasm::~TestDisasm() noexcept {
  if (err_ != CS_ERR_OK) {
    return;
  }
  cs_close(&handle_);
}

const Fmt &TestDisasm::disasm(const Fmt &fmt, Bytes bytes) {
  if (err_ != CS_ERR_OK) {
    return fmt << "capstone initialization error";
  }
  cs_insn *insn = nullptr;
  size_t count = cs_disasm(handle_, bytes.data(), bytes.size(), 0x0, 0, &insn);
  if (count <= 0) {
    return fmt << "capstone disassemble error";
  }
  for (size_t i = 0; i < count; i++) {
    if (i) {
      fmt << '\n';
    }
    format(fmt, insn + i);
  }
  cs_free(insn, count);
  return fmt;
}

const Fmt &TestDisasm::format(const Fmt &fmt, const cs_insn *insn) {
  fmt << "(x86_" << Chars(cs_insn_name(handle_, insn->id));

  cs_detail *detail = insn->detail;
  for (size_t i = 0, n = detail->x86.op_count; i < n; i++) {
    cs_x86_op *op = detail->x86.operands + i;
    switch (op->type) {
    case X86_OP_REG:
      fmt << ' ' << Chars(cs_reg_name(handle_, op->reg));
      break;
    case X86_OP_IMM:
      fmt << ' ' << fix_immediate(insn, op->imm);
      break;
    case X86_OP_MEM:
      fmt << " (x86_mem" << (op->size * 8);
      if (op->mem.disp != 0) {
        fmt << ' ' << op->mem.disp;
      }
      if (op->mem.base != X86_REG_INVALID) {
        fmt << ' ' << Chars(cs_reg_name(handle_, op->mem.base));
      }
      if (op->mem.index != X86_REG_INVALID) {
        fmt << ' ' << Chars(cs_reg_name(handle_, op->mem.index));
        fmt << ' ' << op->mem.scale;
      }
      fmt << ')';
      break;
    default:
      break;
    }
  }
  return fmt << ')';
}

static bool is_jump(cs_group_type group) noexcept {
  switch (group) {
  case CS_GRP_JUMP:
  case CS_GRP_CALL:
  case CS_GRP_BRANCH_RELATIVE:
    return true;
  default:
    return false;
  }
}

static bool is_jump(const cs_detail *detail) noexcept {
  for (size_t i = 0, n = detail->groups_count; i < n; i++) {
    if (is_jump(cs_group_type(detail->groups[i]))) {
      return true;
    }
  }
  return false;
}

/**
 * capstone computes absolute destination address of jumps,
 * while we want relative address from end of disassembled instruction.
 *
 * fix it.
 */
int64_t TestDisasm::fix_immediate(const cs_insn *insn, int64_t imm) noexcept {
  if (is_jump(insn->detail)) {
    return int32_t(imm - insn->size);
  }
  return imm;
}

void TestDisasm::test_asm_disasm_x64(const Node &node, Assembler &assembler) {
  assembler.clear();
  assembler.x64(node);

  if (!assembler) {
    TEST(bool(assembler), ==, true);
  } else {
    CRange<Error> errors = assembler.errors();
    if (errors) {
      String errmsg;
      Fmt{&errmsg} << errors;
      TEST(errmsg, ==, Chars{});
    }
  }

  String disassembled;
  disasm(Fmt{&disassembled}, assembler.bytes());

  TEST(to_string(node, Syntax::CapstoneCompat), ==, disassembled);
}

#else // !HAVE_LIBCAPSTONE

TestDisasm::TestDisasm() noexcept : testcount_{}, err_{1} {
}

TestDisasm::~TestDisasm() noexcept {
}

const Fmt &TestDisasm::disasm(const Fmt &fmt, Bytes bytes) {
  (void)bytes;
  return fmt << "capstone support non compiled";
}

void TestDisasm::test_asm_disasm_x64(const Node &node, Assembler &assembler) {
  assembler.clear();
  assembler.x64(node);
}

#endif

} // namespace onejit
