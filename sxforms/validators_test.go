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

func TestValidatorAnyNoneOf(t *testing.T) {
	valid, invalid := []string{"a", "c"}, []string{"b", "d"}
	anyOf := sxforms.AnyOf(valid...)
	noneOf := sxforms.NoneOf(valid...)

	if got := anyOf.Attributes(); got != nil {
		t.Error("non nil attributes for anyOf", got)
	}
	if got := noneOf.Attributes(); got != nil {
		t.Error("non nil attributes for noneOf", got)
	}

	fieldAny := sxforms.TextField("any", "any", anyOf)
	fieldNone := sxforms.TextField("none", "none", noneOf)
	form := sxforms.Define(fieldAny, fieldNone)

	for _, v := range valid {
		if err := fieldAny.SetValue(v); err != nil {
			t.Errorf("fieldAny.SetValue(%q) failed: %v", v, err)
		}
		if err := anyOf.Check(form, fieldAny); err != nil {
			t.Error("unexpected error for any:", err)
		}

		if err := fieldNone.SetValue(v); err != nil {
			t.Errorf("fieldNone.SetValue(%q) failed: %v", v, err)
		}
		if err := noneOf.Check(form, fieldNone); err == nil {
			t.Error("expected error for none, but got nil")
		}
	}

	for _, v := range invalid {
		if err := fieldAny.SetValue(v); err != nil {
			t.Errorf("fieldAny.SetValue(%q) failed: %v", v, err)
		}
		if err := anyOf.Check(form, fieldAny); err == nil {
			t.Error("expected error for any, but got nil")
		}

		if err := fieldNone.SetValue(v); err != nil {
			t.Errorf("fieldNone.SetValue(%q) failed: %v", v, err)
		}
		if err := noneOf.Check(form, fieldNone); err != nil {
			t.Error("unexpected error for none:", err)
		}
	}
}
