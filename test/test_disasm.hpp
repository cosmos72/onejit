/*
 * onejit - JIT compiler in C++
 *
 * Copyright (C) 2018-2021 Massimiliano Ghilardi
 *
 *     This Source Code Form is subject to the terms of the Mozilla Public
 *     License, v. 2.0. If a copy of the MPL was not distributed with this
 *     file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * test_disasm.hpp
 *
 *  Created on Feb 04, 2021
 *      Author Massimiliano Ghilardi
 */

#include <onejit/fwd.hpp>
#include <onejit/test.hpp>
#include <onejit_config.h> // HAVE_LIBCAPSTONE

#ifdef HAVE_LIBCAPSTONE
#include <capstone.h>
#endif

#define TEST(lhs, op, rhs) (ONEJIT_TEST(lhs, op, rhs), ++testcount(), (void)0)

namespace onejit {

class TestDisasm {
public:
  TestDisasm() noexcept;
  ~TestDisasm() noexcept;

  constexpr explicit operator bool() const noexcept {
    return err_ == 0;
  }

  size_t &testcount() const noexcept {
    return testcount_;
  }

  const Fmt &disasm(const Fmt &fmt, Bytes bytes);

  // assemble instruction, then disassemble it with capstone
  // and compare its text representation
  void test_asm_disasm_x64(const Node &node, Assembler &assembler);

private:
  mutable size_t testcount_;

#ifdef HAVE_LIBCAPSTONE
  const Fmt &format(const Fmt &fmt, const cs_insn *insn);
  static int64_t fix_immediate(const cs_insn *insn, int64_t imm) noexcept;

  csh handle_;
  cs_err err_;
#else
  int err_;
#endif
};

} // namespace onejit
