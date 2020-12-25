/*
 * gomacrojit - JIT compiler in Go
 *
 * Copyright (C) 2019-2020 Massimiliano Ghilardi
 *
 *     This Source Code Form is subject to the terms of the Mozilla Public
 *     License, v. 2.0. If a copy of the MPL was not distributed with this
 *     file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 *
 * stmt.go
 *
 *  Created on Dec 24, 2020
 *      Author Massimiliano Ghilardi
 */

package jit

type Stmt interface {
	stmt() // private marker
}

// ============================ ExprStmt =========================================

type ExprStmt struct {
	expr Expr
}

func ToStmt(expr Expr) *ExprStmt {
	return &ExprStmt{
		expr: expr,
	}
}

func (s *ExprStmt) Expr() Expr {
	return s.expr
}

// implement Stmt interface
func (s *ExprStmt) stmt() {
}

// ============================ IfStmt ===========================================

type IfStmt struct {
	cond  Expr
	then_ Stmt
	else_ Stmt
}

func If(cond Expr, then_ Stmt, else_ Stmt) *IfStmt {
	if kind := cond.Kind(); kind != Bool {
		badOpKind(JUMP_IF, kind)
	}
	return &IfStmt{
		cond:  cond,
		then_: then_,
		else_: else_,
	}
}

func (s *IfStmt) Cond() Expr {
	return s.cond
}

func (s *IfStmt) Then() Stmt {
	return s.then_
}

func (s *IfStmt) Else() Stmt {
	return s.else_
}

// implement Stmt interface
func (s *IfStmt) stmt() {
}

// ============================ BlockStmt ========================================

type BlockStmt struct {
	list []Stmt
}

func (s *BlockStmt) Len() int {
	return len(s.list)
}

func (s *BlockStmt) At(i int) Stmt {
	return s.list[i]
}

func Block(list ...Stmt) *BlockStmt {
	return &BlockStmt{
		list: list,
	}
}

func BlockSlice(list []Stmt) *BlockStmt {
	return &BlockStmt{
		list: append([]Stmt(nil), list...),
	}
}

func (s *BlockStmt) stmt() {
}

// ============================ ForStmt ========================================

type ForStmt struct {
	init Stmt
	cond Expr
	post Stmt
	body Stmt
}

func (s *ForStmt) Init() Stmt {
	return s.init
}

func (s *ForStmt) Cond() Expr {
	return s.cond
}

func (s *ForStmt) Post() Stmt {
	return s.post
}

func (s *ForStmt) Body() Stmt {
	return s.body
}

func For(init Stmt, cond Expr, post Stmt, body Stmt) *ForStmt {
	if cond != nil && cond.Kind() != Bool {
		badOpKind(JUMP_IF, cond.Kind())
	}
	return &ForStmt{
		init: init,
		cond: cond,
		post: post,
		body: body,
	}
}

func (s *ForStmt) stmt() {
}