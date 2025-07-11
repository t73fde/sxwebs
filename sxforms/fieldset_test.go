// -----------------------------------------------------------------------------
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
// -----------------------------------------------------------------------------

package sxforms_test

import (
	"testing"

	"t73f.de/r/sxwebs/sxforms"
)

func TestBasicFieldset(t *testing.T) {
	cb1 := sxforms.CheckboxField("admin", "Admin")
	cb2 := sxforms.CheckboxField("user", "User")
	fs := sxforms.FieldsetField("fieldset", "I am legend", cb1, cb2)
	f := sxforms.Define(
		sxforms.DateField("begin", "Start"),
		fs,
		sxforms.DateField("end", "Stop"),
	)
	if got, err := f.Field("admin"); err != nil || got != cb1 {
		t.Error("unable to find admin field")
	}
	if got, err := f.Field("user"); err != nil || got != cb2 {
		t.Error("unable to find user field")
	}
	exp := "(form (@ (action . \"\") (method . \"POST\")) (div (label (@ (for . \"begin\")) \"Start\") (input (@ (id . \"begin\") (name . \"begin\") (type . \"date\") (value . \"\")))) (fieldset (@ (id . \"fieldset\") (name . \"fieldset\")) (legend \"I am legend\") (div (input (@ (id . \"admin\") (name . \"admin\") (type . \"checkbox\") (value . \"admin\"))) (label (@ (for . \"admin\")) \"Admin\")) (div (input (@ (id . \"user\") (name . \"user\") (type . \"checkbox\") (value . \"user\"))) (label (@ (for . \"user\")) \"User\"))) (div (label (@ (for . \"end\")) \"Stop\") (input (@ (id . \"end\") (name . \"end\") (type . \"date\") (value . \"\")))))"
	got := f.Render().String()
	if got != exp {
		t.Errorf("\nexpected %q\nbut got  %q", exp, got)
	}
}
