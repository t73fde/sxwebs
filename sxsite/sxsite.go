//-----------------------------------------------------------------------------
// Copyright (c) 2024-present Detlef Stern
//
// This file is part of sxwebs.
//
// sxwebs is licensed under the latest version of the EUPL // (European Union
// Public License). Please see file LICENSE.txt for your rights and obligations
// under this license.
//
// SPDX-License-Identifier: EUPL-1.2
// SPDX-FileCopyrightText: 2024-present Detlef Stern
//-----------------------------------------------------------------------------

// Package sxsite allows to work with webs/site in a Sx environment.
package sxsite

import (
	"fmt"

	"t73f.de/r/sx"
	"t73f.de/r/sx/sxbuiltins"
	"t73f.de/r/sx/sxeval"
	"t73f.de/r/webs/site"
	"t73f.de/r/webs/urlbuilder"
)

// MakeURLForBuiltin returns a builtin that provides the (url-for node-id args...)
// function. It is specific to a webs/site.Site.
func MakeURLForBuiltin(st *site.Site) *sxeval.Builtin {
	return &sxeval.Builtin{
		Name:     "url-for",
		MinArity: 1,
		MaxArity: -1,
		TestPure: sxeval.AssertPure,
		Fn1: func(_ *sxeval.Environment, arg sx.Object, _ *sxeval.Frame) (sx.Object, error) {
			nodeID, errArg := sxbuiltins.GetString(arg, 0)
			if errArg != nil {
				return nil, errArg
			}
			return builderToSx(nodeID, st.BuilderFor(nodeID.GetValue()))
		},
		Fn: func(_ *sxeval.Environment, args sx.Vector, _ *sxeval.Frame) (sx.Object, error) {
			nodeID, err := sxbuiltins.GetString(args[0], 0)
			if err != nil {
				return nil, err
			}
			vals := make([]any, 0, len(args)-1)
			for i := 1; i < len(args); i++ {
				val, errArg := sxbuiltins.GetString(args[i], i) // TODO: more than just string?
				if errArg != nil {
					return nil, errArg
				}
				vals = append(vals, val.GetValue())
			}
			return builderToSx(nodeID, st.BuilderFor(nodeID.GetValue(), vals...))
		},
	}
}
func builderToSx(nodeID sx.Object, ub *urlbuilder.URLBuilder) (sx.Object, error) {
	if ub == nil {
		return nil, fmt.Errorf("node id not found: %v", nodeID)
	}
	return sx.MakeString(ub.String()), nil
}

// MakeMakeURLBuiltin returns a builtin that provides the (make-url path...)
// function. It is specific to a webs/site.Site.
func MakeMakeURLBuiltin(st *site.Site) *sxeval.Builtin {
	return &sxeval.Builtin{
		Name:     "make-url",
		MinArity: 0,
		MaxArity: -1,
		TestPure: sxeval.AssertPure,
		Fn0: func(*sxeval.Environment, *sxeval.Frame) (sx.Object, error) {
			return sx.MakeString(st.MakeURLBuilder().String()), nil
		},
		Fn1: func(_ *sxeval.Environment, arg sx.Object, _ *sxeval.Frame) (sx.Object, error) {
			s, err := sxbuiltins.GetString(arg, 0)
			if err != nil {
				return nil, err
			}
			ub := st.MakeURLBuilder()
			ub = ub.AddPath(s.GetValue())
			return sx.MakeString(ub.String()), nil
		},
		Fn: func(_ *sxeval.Environment, args sx.Vector, _ *sxeval.Frame) (sx.Object, error) {
			ub := st.MakeURLBuilder()
			for i, arg := range args {
				sVal, err := sxbuiltins.GetString(arg, i)
				if err != nil {
					return nil, err
				}
				ub = ub.AddPath(sVal.GetValue())
			}
			return sx.MakeString(ub.String()), nil
		},
	}
}
