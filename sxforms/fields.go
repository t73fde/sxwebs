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

import (
	"time"

	"t73f.de/r/sx"
	"t73f.de/r/sxwebs/sxhtml"
)

// Field represents a HTTP form field.
type Field interface {
	Name() string
	Label() string
	Value() string
	Clear()
	SetValue(string) error
	Validators() []Validator
	Render(string, []string) sx.Object
}

// InputField represents a HTTP <input> field.
type InputField struct {
	itype      string
	name       string
	label      string
	value      string
	validators []Validator
	autofocus  bool
}

// Constants for InputField.itype
const (
	itypeDate     = "date"
	itypePassword = "password"
	itypeSubmit   = "submit"
	itypeText     = "text"
)

func (fd *InputField) Name() string  { return fd.name }
func (fd *InputField) Label() string { return fd.label }
func (fd *InputField) Value() string {
	if fd.itype == itypeSubmit {
		return ""
	}
	return fd.value
}

func (fd *InputField) Clear() {
	if fd.itype != itypeSubmit {
		fd.value = ""
	}
}

// Time layouts of data coming from HTML forms
const (
	htmlDateLayout = "2006-01-02"
)

func (fd *InputField) SetValue(value string) error {
	fd.value = value
	switch fd.itype {
	case itypeDate:
		if _, err := time.Parse(htmlDateLayout, value); err != nil {
			return err
		}
	}
	return nil
}

func (fd *InputField) Validators() []Validator { return fd.validators }

func (fd *InputField) Render(fieldID string, messages []string) sx.Object {
	var flb sx.ListBuilder
	flb.Add(sx.MakeSymbol("div"))
	if label := fd.label; label != "" {
		flb.Add(sx.MakeList(
			sx.MakeSymbol("label"),
			sx.MakeList(
				sxhtml.SymAttr,
				sx.Cons(sx.MakeSymbol("for"), sx.MakeString(fieldID)),
			),
			sx.MakeString(label),
		))
	}

	for _, msg := range messages {
		flb.Add(sx.MakeList(
			sx.MakeSymbol("span"),
			sx.MakeList(
				sxhtml.SymAttr,
				sx.Cons(sx.MakeSymbol("class"), sx.MakeString("message")),
			),
			sx.MakeString(msg),
		))
	}

	var wlb sx.ListBuilder
	wlb.Add(sxhtml.SymAttr)
	wlb.Add(sx.Cons(sx.MakeSymbol("id"), sx.MakeString(fieldID)))
	wlb.Add(sx.Cons(sx.MakeSymbol("name"), sx.MakeString(fd.name)))
	wlb.Add(sx.Cons(sx.MakeSymbol("type"), sx.MakeString(fd.itype)))
	wlb.Add(sx.Cons(sx.MakeSymbol("value"), sx.MakeString(fd.value)))
	if fd.autofocus {
		wlb.Add(sx.Cons(sx.MakeSymbol("autofocus"), sx.Nil()))
	}
	for _, validator := range fd.validators {
		if valAttrs := validator.Attributes(); valAttrs != nil {
			wlb.ExtendBang(valAttrs)
		}
	}
	flb.Add(sx.MakeList(sx.MakeSymbol("input"), wlb.List()))

	return flb.List()
}

// SetAutofocus for the field.
func (fd *InputField) SetAutofocus() *InputField { fd.autofocus = true; return fd }

// TextField builds a new text field.
func TextField(name, label string, validators ...Validator) *InputField {
	return &InputField{
		itype:      itypeText,
		name:       name,
		label:      label,
		validators: validators,
	}
}

// DateField builds a new field to enter dates.
func DateField(name, label string, validators ...Validator) *InputField {
	return &InputField{
		itype:      itypeDate,
		name:       name,
		label:      label,
		validators: validators,
	}
}

// PasswordField builds a new password field.
func PasswordField(name, label string, validators ...Validator) *InputField {
	return &InputField{
		itype:      itypePassword,
		name:       name,
		label:      label,
		validators: validators,
	}
}

// SubmitField builds a new submit field.
func SubmitField(name, value string) *InputField {
	return &InputField{
		itype: itypeSubmit,
		name:  name,
		value: value,
	}
}
