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
	f := sxforms.Define(
		sxforms.TextField("username", "User name", sxforms.Required("username")),
		sxforms.PasswordField("password", "Password", sxforms.Required("password")),
		sxforms.SubmitField("submit", "Login"),
	)
	f.SetFormValues(nil)
	if got := f.IsValid(); got {
		t.Error("empty form must not validate")
	}
	gotMsgs := f.Messages()
	if len(gotMsgs) == 0 {
		t.Error("form did not validate, but there are no messages")
	}
	expMsgs := sxforms.Messages{
		"password": {"password"},
		"username": {"username"},
	}
	if !maps.EqualFunc(expMsgs, gotMsgs, slices.Equal) {
		t.Errorf("expected errors: %v, but got %v", expMsgs, gotMsgs)
	}

	f.SetFormValues(url.Values{"username": nil, "password": nil})
	if got := f.IsValid(); got {
		t.Error("nil form must not validate")
	}

	f.SetFormValues(url.Values{"username": {"user"}, "password": {"pass"}})
	if got := f.IsValid(); !got {
		t.Error("normal form must validate")
	}
	expData := sxforms.Data{"password": "pass", "username": "user"}
	if gotData := f.Data(); !maps.Equal(expData, gotData) {
		t.Errorf("expected data %v, but got %v", expData, gotData)
	}
}
