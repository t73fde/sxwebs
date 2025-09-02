//-----------------------------------------------------------------------------
// Copyright (c) 2025-present Detlef Stern
//
// This file is part of sxwebs.
//
// sxwebs is licensed under the latest version of the EUPL // (European Union
// Public License). Please see file LICENSE.txt for your rights and obligations
// under this license.
//
// SPDX-License-Identifier: EUPL-1.2
// SPDX-FileCopyrightText: 2025-present Detlef Stern
//-----------------------------------------------------------------------------

// Package sxhtmls allows to convert HTML representations: webs/sxhtml and
// t73f.de/r/webs/htmls.
package sxhtmls

import (
	"errors"
	"fmt"
	"strings"

	"t73f.de/r/sx"
	"t73f.de/r/sxwebs/sxhtml"
	"t73f.de/r/webs/htmls"
)

// ToSxHTML transforms an htmls.Node into a SxHTML object.
func ToSxHTML(n *htmls.Node) (sx.Object, error) {
	if n == nil {
		return sx.Nil(), nil
	}
	switch n.Type {
	case htmls.TextNode:
		return sx.MakeString(n.Data), nil
	case htmls.ElementNode:
		// no-op, fall through switch
	case htmls.RawNode:
		return sx.MakeList(sxhtml.SymNoEscape, sx.MakeString(n.Data)), nil
	case htmls.CommentNode:
		return sx.MakeList(sxhtml.SymBlockComment, sx.MakeString(n.Data)), nil
	default:
		return nil, fmt.Errorf("unknown node type: %v", n.Type)
	}

	tag, err := makeSymbol(n.Data)
	if err != nil {
		return nil, err
	}

	var lb sx.ListBuilder
	lb.Add(tag)
	attrs, err := toSxAttrs(n.Attributes)
	if err != nil {
		return nil, err
	}
	if attrs != nil {
		lb.Add(attrs)
	}
	for _, child := range n.Children {
		obj, errChild := ToSxHTML(child)
		if errChild != nil {
			return nil, errChild
		}
		lb.Add(obj)
	}
	return lb.List(), nil
}

func toSxAttrs(attrs []htmls.Attribute) (*sx.Pair, error) {
	if len(attrs) == 0 {
		return nil, nil
	}
	var lb sx.ListBuilder
	for _, attr := range attrs {
		sym, err := makeSymbol(attr.Key)
		if err != nil {
			return nil, err
		}
		lb.Add(sx.Cons(sym, sx.MakeString(attr.Value)))
	}
	return lb.List(), nil
}

func makeSymbol(s string) (*sx.Symbol, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, errEmptySymbol
	}
	return sxhtml.MakeSymbol(s), nil
}

var errEmptySymbol = errors.New("empty symbol string")
