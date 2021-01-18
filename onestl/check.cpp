/*
 * onestl - Tiny STL C++ library
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
 * check.cpp
 *
 *  Created on Jan 14, 2020
 *      Author Massimiliano Ghilardi
 */

#include "onestl/check.hpp"
#include "onestl/chars.hpp"

#include <stdexcept>
#include <string>

namespace onestl {

void ONESTL_NORETURN CheckFailed::throw_error() const {
  std::stringstream buf;
  buf << Chars("check failed at ") << file_ << ':' << line_ //
      << '\t' << lhs_ << ' ' << op_ << ' ' << rhs_          //
      << Chars("\n\t") << lhs_ << Chars("\t= ") << lval_    //
      << Chars("\n\t") << rhs_ << Chars("\t= ") << rval_;

  throw std::range_error(buf.str());
}

} // namespace onestl