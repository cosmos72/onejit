/*
 * Copyright (C) 2021 Massimiliano Ghilardi
 *
 *     This Source Code Form is subject to the terms of the Mozilla Public
 *     License, v. 2.0. If a copy of the MPL was not distributed with this
 *     file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 *
 * resolver_const.go
 *
 *  Created on: Apr 10, 2021
 *      Author: Massimiliano Ghilardi
 */

package typecheck

import (
	"github.com/cosmos72/onejit/go/ast"
	"github.com/cosmos72/onejit/go/constant"
	"github.com/cosmos72/onejit/go/token"
	"github.com/cosmos72/onejit/go/types"
)

func (r *Resolver) recoverFromPanic(node *ast.Node) {
	switch fail := recover().(type) {
	case nil:
	case *token.Error:
		r.errors.errs = append(r.errors.errs, fail)
	case string:
		r.error(*node, fail)
	default:
		panic(fail)
	}
}

func (r *Resolver) declareObjConst(obj *Object) {
	if obj.Value() != nil && obj.Type() != nil {
		return
	}
	decl := obj.Decl()
	// TODO support iota
	if decl.init == nil && decl.typ == nil {
		return
	}
	// FIXME: remove this hack when makeType() is finished
	defer r.recoverFromPanic(&decl.node)

	var v constant.Value
	var t *types.Complete
	if decl.init != nil {
		v = r.makeConst(decl.init)
		t = v.Type()
	}
	if decl.typ != nil {
		t = completeType(r.makeType(decl.typ))
		if decl.init == nil {
			v = constant.MakeZero(t)
		} else {
			v = v.To(t)
		}
	}
	obj.SetType(t)
	obj.SetValue(v)
}

func (r *Resolver) makeConst(node ast.Node) constant.Value {
	if obj := r.resolved[node]; obj != nil && obj.Class() == types.ConstObj {
		if v, ok := obj.Value().(constant.Value); ok {
			return v
		}
	}
	var v constant.Value // the zero value is Invalid
	if op := node.Op(); op.IsLiteral() && op != token.IDENT {
		atom := node.(*ast.Atom)
		v = constant.MakeFromLiteral(atom.Lit, op)
	} else if op.IsOperator() {
		if node.Len() == 1 {
			v = r.makeConst(node.At(0))
			v = constant.UnaryOp(op, v)
		} else if node.Len() == 2 {
			xv := r.makeConst(node.At(0))
			yv := r.makeConst(node.At(1))
			v = constant.BinaryOp(xv, op, yv)
		}
	} else if op == token.CALL {
		// TODO
		r.error(node, "unimplemented: type conversion on constant")
	}
	if v.Kind() == constant.Invalid {
		r.error(node, "const initializer "+node.String()+" is not a constant")
	}
	return v
}
