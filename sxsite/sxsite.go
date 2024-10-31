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
		Fn1: func(_ *sxeval.Environment, arg sx.Object) (sx.Object, error) {
			nodeID, errArg := sxbuiltins.GetString(arg, 0)
			if errArg != nil {
				return nil, errArg
			}
			return builderToSx(nodeID, st.BuilderFor(nodeID.GetValue()))
		},
		Fn2: func(_ *sxeval.Environment, arg0, arg1 sx.Object) (sx.Object, error) {
			nodeID, errArg := sxbuiltins.GetString(arg0, 0)
			if errArg != nil {
				return nil, errArg
			}
			key, errArg := sxbuiltins.GetString(arg1, 1)
			if errArg != nil {
				return nil, errArg
			}
			return builderToSx(nodeID, st.BuilderFor(nodeID.GetValue(), key.GetValue()))
		},
		Fn: func(_ *sxeval.Environment, args sx.Vector) (sx.Object, error) {
			nodeID, err := sxbuiltins.GetString(args[0], 0)
			if err != nil {
				return nil, err
			}
			vals := make([]string, 0, len(args)-1)
			for i := 1; i < len(args); i++ {
				val, errArg := sxbuiltins.GetString(args[i], i)
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
		Fn0: func(_ *sxeval.Environment) (sx.Object, error) {
			return sx.MakeString(st.MakeURLBuilder().String()), nil
		},
		Fn1: func(_ *sxeval.Environment, arg sx.Object) (sx.Object, error) {
			s, err := sxbuiltins.GetString(arg, 0)
			if err != nil {
				return nil, err
			}
			ub := st.MakeURLBuilder()
			ub = ub.AddPath(s.GetValue())
			return sx.MakeString(ub.String()), nil
		},
		Fn2: func(_ *sxeval.Environment, arg0, arg1 sx.Object) (sx.Object, error) {
			ub := st.MakeURLBuilder()
			s, err := sxbuiltins.GetString(arg0, 0)
			if err != nil {
				return nil, err
			}
			ub = ub.AddPath(s.GetValue())
			s, err = sxbuiltins.GetString(arg1, 1)
			if err != nil {
				return nil, err
			}
			ub = ub.AddPath(s.GetValue())
			return sx.MakeString(ub.String()), nil
		},
		Fn: func(_ *sxeval.Environment, args sx.Vector) (sx.Object, error) {
			ub := st.MakeURLBuilder()
			for i := 0; i < len(args); i++ {
				sVal, err := sxbuiltins.GetString(args[i], i)
				if err != nil {
					return nil, err
				}
				ub = ub.AddPath(sVal.GetValue())
			}
			return sx.MakeString(ub.String()), nil
		},
	}
}
