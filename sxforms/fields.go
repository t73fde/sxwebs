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
	Render(string) sx.Object
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

func (fd *InputField) Name() string  { return fd.name }
func (fd *InputField) Label() string { return fd.label }
func (fd *InputField) Value() string {
	if fd.itype == "submit" {
		return ""
	}
	return fd.value
}

func (fd *InputField) Clear() {
	if fd.itype != "submit" {
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
	case "date":
		if _, err := time.Parse(htmlDateLayout, value); err != nil {
			return err
		}
	}
	return nil
}

func (fd *InputField) Validators() []Validator { return fd.validators }
func (fd *InputField) Render(fieldID string) sx.Object {
	var lb sx.ListBuilder
	lb.Add(sxhtml.SymAttr)
	lb.Add(sx.Cons(sx.MakeSymbol("id"), sx.MakeString(fieldID)))
	lb.Add(sx.Cons(sx.MakeSymbol("name"), sx.MakeString(fd.name)))
	lb.Add(sx.Cons(sx.MakeSymbol("type"), sx.MakeString(fd.itype)))
	lb.Add(sx.Cons(sx.MakeSymbol("value"), sx.MakeString(fd.value)))
	if fd.autofocus {
		lb.Add(sx.Cons(sx.MakeSymbol("autofocus"), sx.Nil()))
	}
	for _, validator := range fd.validators {
		if valAttrs := validator.Attributes(); valAttrs != nil {
			lb.ExtendBang(valAttrs)
		}
	}
	return sx.MakeList(sx.MakeSymbol("input"), lb.List())
}

// SetAutofocus for the field.
func (fd *InputField) SetAutofocus() *InputField { fd.autofocus = true; return fd }

// TextField builds a new text field.
func TextField(name, label string, validators ...Validator) *InputField {
	return &InputField{
		itype:      "text",
		name:       name,
		label:      label,
		validators: validators,
	}
}

// DateField builds a new field to enter dates.
func DateField(name, label string, validators ...Validator) *InputField {
	return &InputField{
		itype:      "date",
		name:       name,
		label:      label,
		validators: validators,
	}
}

// PasswordField builds a new password field.
func PasswordField(name, label string, validators ...Validator) *InputField {
	return &InputField{
		itype:      "password",
		name:       name,
		label:      label,
		validators: validators,
	}
}

// SubmitField builds a new submit field.
func SubmitField(name, value string) *InputField {
	return &InputField{
		itype: "submit",
		name:  name,
		value: value,
	}
}
