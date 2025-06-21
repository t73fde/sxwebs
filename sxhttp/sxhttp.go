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
	"context"
	"fmt"
	"net/http"

	"t73f.de/r/sx"
	"t73f.de/r/sx/sxeval"
)

// ----- SxContext -----------------------------------------------------------

// SxContext is a context.Context, seen as a Sx object.
type SxContext struct{ val context.Context }

// MakeContext creates a SxContext from a context.Context.
func MakeContext(ctx context.Context) SxContext { return SxContext{ctx} }

// GetValue returns the context.Context value.
func (ctx SxContext) GetValue() context.Context { return ctx.val }

// IsNil returns true for a nil value.
func (SxContext) IsNil() bool { return false }

// IsAtom returns true for an atomic value.
func (SxContext) IsAtom() bool { return true }

// IsEqual returns true if the sx content is equal to the given object.
func (ctx SxContext) IsEqual(other sx.Object) bool {
	if other.IsNil() {
		return false
	}
	otherCtx, isCtx := other.(SxContext)
	return isCtx && ctx.val == otherCtx.val
}
func (ctx SxContext) String() string {
	return fmt.Sprintf("#<SxContext:%v>", ctx.val)
}

// GoString return the Go representation of the context.
func (ctx SxContext) GoString() string { return ctx.String() }

// GetContext returns the given sx.Object as a SxContext, if possible.
func GetContext(obj sx.Object) (SxContext, bool) {
	if obj.IsNil() {
		return SxContext{}, false
	}
	ctx, ok := obj.(SxContext)
	return ctx, ok
}

// GetBuiltinContext returns the given sx.Object as a SxContext. If this is not
// possible, an error is returned.
//
// This function can be used as a helper function to implement sxeval.Builtin.
func GetBuiltinContext(arg sx.Object, pos int) (SxContext, error) {
	if ctx, isCtx := GetContext(arg); isCtx {
		return ctx, nil
	}
	return SxContext{}, fmt.Errorf("argument %d is not a context, but %T/%v", pos+1, arg, arg)
}

// ----- SxRequest -----------------------------------------------------------

// SxRequest is a http.Request, seen as a Sx object.
type SxRequest http.Request

// MakeRequest creates a Sx object from a htt.Request.
func MakeRequest(r *http.Request) *SxRequest { return (*SxRequest)(r) }

// GetValue returns the underlying request object.
func (r *SxRequest) GetValue() *http.Request { return (*http.Request)(r) }

// IsNil returns true of the object id a nil value.
func (r *SxRequest) IsNil() bool { return r == nil }

// IsAtom returns true for an atomic value.
func (*SxRequest) IsAtom() bool { return true }

// IsEqual returns true if the other object is equal to this request object.
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

// GoString returns the Go representation.
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

// URLPath is a builtin that returns the URL path ob an request object.
var URLPath = sxeval.Builtin{
	Name:     "request-url-path",
	MinArity: 1,
	MaxArity: 1,
	TestPure: sxeval.AssertPure,
	Fn1: func(_ *sxeval.Environment, arg sx.Object, _ *sxeval.Frame) (sx.Object, error) {
		r, err := GetBuiltinRequest(arg, 0)
		if err != nil {
			return sx.Nil(), err
		}
		return sx.MakeString(r.GetValue().URL.Path), nil
	},
}

// Context is a builtin the returns the context object of an request object.
var Context = sxeval.Builtin{
	Name:     "request-context",
	MinArity: 1,
	MaxArity: 1,
	TestPure: sxeval.AssertPure,
	Fn1: func(_ *sxeval.Environment, arg sx.Object, _ *sxeval.Frame) (sx.Object, error) {
		r, err := GetBuiltinRequest(arg, 0)
		if err != nil {
			return sx.Nil(), err
		}
		return MakeContext(r.GetValue().Context()), nil
	},
}

// ----- SxResponseWriter ----------------------------------------------------

// SxResponseWriter is a http.ResponseWriter, seen as a Sx object.
type SxResponseWriter struct{ val http.ResponseWriter }

// MakeResponseWriter creates an object based on a response writer.
func MakeResponseWriter(w http.ResponseWriter) SxResponseWriter { return SxResponseWriter{w} }

// GetValue returns the underlying response writer value.
func (w SxResponseWriter) GetValue() http.ResponseWriter { return w.val }

// IsNil returns true, if the object is a nil value.
func (SxResponseWriter) IsNil() bool { return false }

// IsAtom returns true if this object is an atomic object.
func (SxResponseWriter) IsAtom() bool { return true }

// IsEqual returns true, if this response writer is equal to the given object.
func (w SxResponseWriter) IsEqual(other sx.Object) bool {
	if sx.IsNil(other) {
		return false
	}
	otherResp, isResp := other.(*SxResponseWriter)
	return isResp && w.val == otherResp.val
}
func (w SxResponseWriter) String() string {
	return fmt.Sprintf("#<SxResponseWriter:%v>", w.GetValue())
}

// GoString returns the Go representation.
func (w SxResponseWriter) GoString() string { return w.String() }
