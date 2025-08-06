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
// SPDX-FileCopyrightText: 2024-present Detlef Stern
// -----------------------------------------------------------------------------

package sxforms

// ----- <input ...> fields

import (
	"time"

	"t73f.de/r/sx"
	"t73f.de/r/sxwebs/sxhtml"
)

// InputElement represents a HTTP <input> field.
type InputElement struct {
	name       string
	label      string
	value      string
	validators Validators
	disabled   bool
	itype      inputType
}

type inputType uint

// Constants for inputType
const (
	_ inputType = iota
	itypeCheckbox
	itypeDate
	itypeDatetime
	itypeEmail
	itypeNumber
	itypePassword
	itypeText
)

// Name returns the name of this element.
func (fd *InputElement) Name() string { return fd.name }

// Value returns the value of the input element.
func (fd *InputElement) Value() string { return fd.value }

// Clear the input element.
func (fd *InputElement) Clear() { fd.value = "" }

// Time layouts of data coming from HTML forms.
const (
	htmlDateLayout     = "2006-01-02"
	htmlDatetimeLayout = "2006-01-02T15:04"
)

// SetValue sets the value of this input element.
func (fd *InputElement) SetValue(value string) (err error) {
	fd.value = value
	switch fd.itype {
	case itypeDate:
		if value != "" {
			_, err = time.Parse(htmlDateLayout, value)
		}
	case itypeDatetime:
		if value != "" {
			_, err = time.Parse(htmlDatetimeLayout, value)
		}
	}
	return err
}

// Validators returns all currently active Validators.
func (fd *InputElement) Validators() Validators {
	if fd.disabled {
		return nil
	}
	return fd.validators
}

// Disable the input element.
func (fd *InputElement) Disable() { fd.disabled = true }

// Render the form input element as SxHTML.
func (fd *InputElement) Render(fieldID string, messages []string) *sx.Pair {
	var flb sx.ListBuilder
	flb.Add(sxhtml.MakeSymbol("div"))
	if label := renderLabel(fd, fieldID, fd.label); label != nil {
		flb.Add(label)
	}
	flb.ExtendBang(renderMessages(messages))

	var attrLb sx.ListBuilder
	attrLb.AddN(
		sx.Cons(sxhtml.MakeSymbol("id"), sx.MakeString(fieldID)),
		sx.Cons(sxhtml.MakeSymbol("name"), sx.MakeString(fd.name)),
		sx.Cons(sxhtml.MakeSymbol("type"), inputTypeString[fd.itype]),
		sx.Cons(sxhtml.MakeSymbol("value"), sx.MakeString(fd.value)),
	)
	addEnablingAttributes(&attrLb, fd.disabled, fd.validators)

	flb.Add(sx.MakeList(sxhtml.MakeSymbol("input"), attrLb.List()))
	return flb.List()
}

var inputTypeString = map[inputType]sx.String{
	itypeCheckbox: sx.MakeString("checkbox"),
	itypeDate:     sx.MakeString("date"),
	itypeDatetime: sx.MakeString("datetime-local"),
	itypeEmail:    sx.MakeString("email"),
	itypeNumber:   sx.MakeString("number"),
	itypePassword: sx.MakeString("password"),
	itypeText:     sx.MakeString("text"),
}

// TextField builds a new text field.
func TextField(name, label string, validators ...Validator) *InputElement {
	return &InputElement{
		itype:      itypeText,
		name:       name,
		label:      label,
		validators: validators,
	}
}

// DateField builds a new field to enter dates.
func DateField(name, label string, validators ...Validator) *InputElement {
	return &InputElement{
		itype:      itypeDate,
		name:       name,
		label:      label,
		validators: validators,
	}
}

// DateValue returns the date as a string suitable for a HTML date field value.
func DateValue(t time.Time) string {
	if t.Equal(time.Time{}) {
		return ""
	}
	return t.Format(htmlDateLayout)
}

// DatetimeField builds a new field to enter a local date/time.
func DatetimeField(name, label string, validators ...Validator) *InputElement {
	return &InputElement{
		itype:      itypeDatetime,
		name:       name,
		label:      label,
		validators: validators,
	}
}

// DatetimeValue returns the time as a string suitable for a HTML datetime-local field value.
func DatetimeValue(t time.Time) string {
	if t.Equal(time.Time{}) {
		return ""
	}
	return t.Format(htmlDatetimeLayout)
}

// PasswordField builds a new password field.
func PasswordField(name, label string, validators ...Validator) *InputElement {
	return &InputElement{
		itype:      itypePassword,
		name:       name,
		label:      label,
		validators: validators,
	}
}

// EmailField builds a new e-mail field.
func EmailField(name, label string, validators ...Validator) *InputElement {
	return &InputElement{
		itype:      itypeEmail,
		name:       name,
		label:      label,
		validators: validators,
	}
}

// NumberField builds a new number field.
func NumberField(name, label string, validators ...Validator) *InputElement {
	return &InputElement{
		itype:      itypeNumber,
		name:       name,
		label:      label,
		validators: validators,
	}
}
