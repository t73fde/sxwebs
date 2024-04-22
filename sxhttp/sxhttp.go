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
// SPDX-FileCopyrightText: 2023-present Detlef Stern
//-----------------------------------------------------------------------------

// Package sxhttp encapsulates net/http definitions as Sx objects.
package sxhttp

import (
	"fmt"
	"net/http"

	"t73f.de/r/sx"
	"t73f.de/r/sx/sxeval"
)

// SxRequest is a http.Request, seen as a Sx object.
type SxRequest http.Request

func MakeRequest(r *http.Request) *SxRequest { return (*SxRequest)(r) }
func (r *SxRequest) GetValue() *http.Request { return (*http.Request)(r) }

func (r *SxRequest) IsNil() bool { return r == nil }
func (*SxRequest) IsAtom() bool  { return true }
func (r *SxRequest) IsEqual(other sx.Object) bool {
	if r == nil {
		return sx.IsNil(other)
	}
	if sx.IsNil(other) {
		return false
	}
	otherReq, isReq := other.(*SxRequest)
	return isReq && r == otherReq
}
func (r *SxRequest) String() string {
	return fmt.Sprintf("#<SxRequest:%v>", r.GetValue())
}
func (r *SxRequest) GoString() string { return r.String() }

// GetRequest returns the given sx.Object as a SxRequest, if possible.
func GetRequest(obj sx.Object) (*SxRequest, bool) {
	if sx.IsNil(obj) {
		return nil, false
	}
	sym, ok := obj.(*SxRequest)
	return sym, ok
}

// GetBuiltinRequest returns the given sx.Object as a SxRequest. If this is not
// possible, an error is returned.
//
// This function can be used as a helper function to implement sxeval.Builtin.
func GetBuiltinRequest(arg sx.Object, pos int) (*SxRequest, error) {
	if r, isRequest := GetRequest(arg); isRequest {
		return r, nil
	}
	return nil, fmt.Errorf("argument %d is not a http request, but %T/%v", pos+1, arg, arg)
}

var URLPath = sxeval.Builtin{
	Name:     "url-path",
	MinArity: 1,
	MaxArity: 1,
	TestPure: sxeval.AssertPure,
	Fn1: func(_ *sxeval.Environment, arg sx.Object) (sx.Object, error) {
		r, err := GetBuiltinRequest(arg, 0)
		if err != nil {
			return sx.Nil(), err
		}
		return sx.MakeString(r.GetValue().URL.Path), nil
	},
}

// ----- ResponseWriter ------------------------------------------------------

// SxResponseWriter is a http.ResponseWriter, seen as a Sx object.
type SxResponseWriter struct{ val http.ResponseWriter }

func MakeResponseWriter(w http.ResponseWriter) *SxResponseWriter { return &SxResponseWriter{w} }
func (w *SxResponseWriter) GetValue() http.ResponseWriter        { return w.val }

func (w *SxResponseWriter) IsNil() bool { return w == nil }
func (*SxResponseWriter) IsAtom() bool  { return true }
func (w *SxResponseWriter) IsEqual(other sx.Object) bool {
	if w == nil {
		return sx.IsNil(other)
	}
	if sx.IsNil(other) {
		return false
	}
	otherResp, isResp := other.(*SxResponseWriter)
	return isResp && w.val == otherResp.val
}
func (w *SxResponseWriter) String() string {
	return fmt.Sprintf("#<SxResponseWriter:%v>", w.GetValue())
}
func (w *SxResponseWriter) GoString() string { return w.String() }
