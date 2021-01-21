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
 * functype.cpp
 *
 *  Created on Jan 21, 2020
 *      Author Massimiliano Ghilardi
 */

#include <onejit/code.hpp>
#include <onejit/functype.hpp>
#include <onestl/chars.hpp>

namespace onejit {

/// \return kind of i-th parameter, or Bad if out-of-bounds
Kind FuncType::param(uint16_t i) const noexcept {
  return child(i + result_n()).kind();
}

/// \return kind of i-th result, or Bad if out-of-bounds
Kind FuncType::result(uint16_t i) const noexcept {
  return i >= result_n() ? Bad : child(i).kind();
}

// param_n must fit uint16_t, and result_n must fit uint16_t
static constexpr bool functype_can_represent(size_t param_n, size_t result_n) {
  return (param_n | result_n) <= 0xFFFF;
}

FuncType FuncType::create(Kinds params, Kinds results, Code *holder) noexcept {
  const size_t param_n = params.size();
  const size_t result_n = results.size();
  while (holder) {
    const NodeHeader header{FTYPE, results.empty() ? Void : results[0], uint16_t(result_n)};
    CodeItem offset = holder->length();

    // put results first, as expected by param(uint32_t) and result(uint32_t)
    if (functype_can_represent(param_n, result_n)                           //
        && holder->add(header) && holder->add(uint32_t(param_n + result_n)) //
        && holder->add(results) && holder->add(params)) {

      return FuncType{Node{header, offset, holder}};
    }
    holder->truncate(offset);
    break;
  }
  return FuncType{};
}

std::ostream &operator<<(std::ostream &out, const FuncType &ftype) {
  out << '(' << ftype.type() << Chars(" (");
  for (size_t i = 0, n = ftype.param_n(); i < n; i++) {
    if (i) {
      out << ' ';
    }
    out << ftype.param(i);
  }
  out << Chars(") -> (");
  for (size_t i = 0, n = ftype.result_n(); i < n; i++) {
    if (i) {
      out << ' ';
    }
    out << ftype.result(i);
  }
  return out << Chars("))");
}

} // namespace onejit