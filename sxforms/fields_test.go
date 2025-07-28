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

	"t73f.de/r/sx"
	"t73f.de/r/sxwebs/sxforms"
	"t73f.de/r/sxwebs/sxhtml"
)

func TestFlowContent(t *testing.T) {
	form := sxforms.Define(
		sxforms.FlowContentField("fce1", sx.MakeList(sxhtml.MakeSymbol("p"), sx.MakeString("Test"))),
	)

	exp := "(form (@ (action . \"\") (method . \"POST\")) (p \"Test\"))"
	if got := form.Render().String(); got != exp {
		t.Errorf("expected: %q, but got: %q", exp, got)
	}
}
