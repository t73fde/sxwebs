// -----------------------------------------------------------------------------
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
// -----------------------------------------------------------------------------

package sxforms

// Attributes are technically validators that do not validate, but return HTML
// attributes.

import "t73f.de/r/sx"

// ----- Autofocus: where the input starts.

// AttrAutofocus sets the "autofocus" attribute.
type AttrAutofocus struct{}

// Check nothing.
func (AttrAutofocus) Check(*Form, Field) error { return nil }

// Attributes return the HTTP attributes.
func (AttrAutofocus) Attributes() *sx.Pair {
	return sx.MakeList(sx.Cons(sx.MakeSymbol("autofocus"), sx.Nil()))
}

// ----- Step: allow to increment / decrement value in HTML client.

// AttrStep is a non-validator that instructs the HTML client to increment / decrement
// the value in its user interface. It does not check anything.
type AttrStep struct {
	Value string
}

// Check the given field w.r.t. to this validator.
func (AttrStep) Check(*Form, Field) error { return nil }

// Attributes returns HTML attributes as a Sx cons list.
func (s AttrStep) Attributes() *sx.Pair {
	return sx.MakeList(sx.Cons(sx.MakeSymbol("step"), sx.MakeString(s.Value)))
}
