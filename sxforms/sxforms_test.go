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

package sxforms_test

import (
	"maps"
	"net/url"
	"slices"
	"testing"

	"t73f.de/r/sxwebs/sxforms"
)

func TestSimpleRequiredForm(t *testing.T) {
	f := sxforms.Make(
		sxforms.MakeTextField("username", "User name", sxforms.InputRequired("username")),
		sxforms.MakePasswordField("password", "Password", sxforms.InputRequired("password")),
		sxforms.MakeSubmitField("submit", "Login"),
	)
	f.SetFormData(nil)
	if got := f.IsValid(); got {
		t.Error("empty form must not validate")
	}
	gotErrs := f.Errors()
	if len(gotErrs) == 0 {
		t.Error("form did not validate, but there are no errors")
	}
	expErrs := map[string][]error{
		"password": {sxforms.StopValidationError("password")},
		"username": {sxforms.StopValidationError("username")},
	}
	if !maps.EqualFunc(expErrs, gotErrs, slices.Equal) {
		t.Errorf("expected errors: %v, but got %v", expErrs, gotErrs)
	}

	f.SetFormData(url.Values{"username": nil, "password": nil})
	if got := f.IsValid(); got {
		t.Error("nil form must not validate")
	}

	f.SetFormData(url.Values{"username": {"user"}, "password": {"pass"}})
	if got := f.IsValid(); !got {
		t.Error("normal form must validate")
	}
	expData := map[string]string{"password": "pass", "username": "user"}
	if gotData := f.Data(); !maps.Equal(expData, gotData) {
		t.Errorf("expected data %v, but got %v", expData, gotData)
	}
}
