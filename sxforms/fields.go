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
	"fmt"
	"time"

	"t73f.de/r/sx"
	"t73f.de/r/sxwebs/sxhtml"
)

// Field represents a HTTP form field.
type Field interface {
	Name() string
	Value() string
	Clear()
	SetValue(string) error
	Validators() []Validator
	Render(string, []string) sx.Object
}

// ----- <input ...> fields

// InputElement represents a HTTP <input> field.
type InputElement struct {
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

func (fd *InputElement) Name() string { return fd.name }
func (fd *InputElement) Value() string {
	if fd.itype == itypeSubmit {
		return ""
	}
	return fd.value
}

func (fd *InputElement) Clear() {
	if fd.itype != itypeSubmit {
		fd.value = ""
	}
}

// Time layouts of data coming from HTML forms
const (
	htmlDateLayout = "2006-01-02"
)

func (fd *InputElement) SetValue(value string) error {
	fd.value = value
	switch fd.itype {
	case itypeDate:
		if _, err := time.Parse(htmlDateLayout, value); err != nil {
			return err
		}
	}
	return nil
}

func (fd *InputElement) Validators() []Validator { return fd.validators }

func (fd *InputElement) Render(fieldID string, messages []string) sx.Object {
	var flb sx.ListBuilder
	flb.Add(sx.MakeSymbol("div"))
	if label := renderLabel(fieldID, fd.label); label != nil {
		flb.Add(label)
	}
	flb.ExtendBang(renderMessages(messages))

	var attrLb sx.ListBuilder
	attrLb.Add(sxhtml.SymAttr)
	attrLb.Add(sx.Cons(sx.MakeSymbol("id"), sx.MakeString(fieldID)))
	attrLb.Add(sx.Cons(sx.MakeSymbol("name"), sx.MakeString(fd.name)))
	attrLb.Add(sx.Cons(sx.MakeSymbol("type"), sx.MakeString(fd.itype)))
	attrLb.Add(sx.Cons(sx.MakeSymbol("value"), sx.MakeString(fd.value)))
	if fd.autofocus {
		attrLb.Add(sx.Cons(sx.MakeSymbol("autofocus"), sx.Nil()))
	}
	attrLb.ExtendBang(renderValidators(fd.validators))
	flb.Add(sx.MakeList(sx.MakeSymbol("input"), attrLb.List()))

	return flb.List()
}

// SetAutofocus for the field.
func (fd *InputElement) SetAutofocus() *InputElement { fd.autofocus = true; return fd }

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

// PasswordField builds a new password field.
func PasswordField(name, label string, validators ...Validator) *InputElement {
	return &InputElement{
		itype:      itypePassword,
		name:       name,
		label:      label,
		validators: validators,
	}
}

// SubmitField builds a new submit field.
func SubmitField(name, value string) *InputElement {
	return &InputElement{
		itype: itypeSubmit,
		name:  name,
		value: value,
	}
}

// ----- <textarea ...>...</textarea> field

// TextAreaElement represents the corresponding textarea form element.
type TextAreaElement struct {
	name       string
	label      string
	rows       uint32
	cols       uint32
	value      string
	validators []Validator
}

// TextAreaField creates a new text area element.
func TextAreaField(name, label string, validators ...Validator) *TextAreaElement {
	return &TextAreaElement{
		name:       name,
		label:      label,
		validators: validators,
	}
}
func (tae *TextAreaElement) Name() string  { return tae.name }
func (tae *TextAreaElement) Value() string { return tae.value }
func (tae *TextAreaElement) Clear()        { tae.value = "" }
func (tae *TextAreaElement) SetValue(value string) error {
	tae.value = value
	return nil
}
func (tae *TextAreaElement) Validators() []Validator { return tae.validators }
func (tae *TextAreaElement) Render(fieldID string, messages []string) sx.Object {
	var flb sx.ListBuilder
	flb.Add(sx.MakeSymbol("div"))
	if label := renderLabel(fieldID, tae.label); label != nil {
		flb.Add(label)
	}
	flb.ExtendBang(renderMessages(messages))
	var attrLb sx.ListBuilder
	attrLb.Add(sxhtml.SymAttr)
	attrLb.Add(sx.Cons(sx.MakeSymbol("id"), sx.MakeString(fieldID)))
	attrLb.Add(sx.Cons(sx.MakeSymbol("name"), sx.MakeString(tae.name)))
	if rows := tae.rows; rows > 0 {
		attrLb.Add(sx.Cons(sx.MakeSymbol("rows"), sx.MakeString(fmt.Sprint(rows))))
	}
	if cols := tae.cols; cols > 0 {
		attrLb.Add(sx.Cons(sx.MakeSymbol("cols"), sx.MakeString(fmt.Sprint(cols))))
	}
	attrLb.ExtendBang(renderValidators(tae.validators))

	flb.Add(sx.MakeList(sx.MakeSymbol("textarea"), attrLb.List(), sx.MakeString(tae.value)))
	return flb.List()
}

// ----- <select ...>...</select> field

// SelectElement represents the corresponding select form element.
type SelectElement struct {
	name       string
	label      string
	choices    []string
	value      string
	validators []Validator
}

// TextAreaField creates a new text area element.
func SelectField(name, label string, choices []string, validators ...Validator) *SelectElement {
	if len(choices)%2 != 0 {
		panic(fmt.Sprintf("choices must have even number of values: %v", choices))
	}
	return &SelectElement{
		name:       name,
		label:      label,
		choices:    choices,
		validators: validators,
	}
}
func (se *SelectElement) Name() string  { return se.name }
func (se *SelectElement) Value() string { return se.value }
func (se *SelectElement) Clear()        { se.value = "" }
func (se *SelectElement) SetValue(value string) error {
	se.value = value
	for i := 0; i < len(se.choices); i += 2 {
		if se.choices[i] == value {
			return nil
		}
	}
	return fmt.Errorf("no such choice: %q", value)
}
func (se *SelectElement) Validators() []Validator { return se.validators }
func (se *SelectElement) Render(fieldID string, messages []string) sx.Object {
	var flb sx.ListBuilder
	flb.Add(sx.MakeSymbol("div"))
	if label := renderLabel(fieldID, se.label); label != nil {
		flb.Add(label)
	}
	flb.ExtendBang(renderMessages(messages))
	var attrLb sx.ListBuilder
	attrLb.Add(sxhtml.SymAttr)
	attrLb.Add(sx.Cons(sx.MakeSymbol("id"), sx.MakeString(fieldID)))
	attrLb.Add(sx.Cons(sx.MakeSymbol("name"), sx.MakeString(se.name)))
	attrLb.ExtendBang(renderValidators(se.validators))

	var wlb sx.ListBuilder
	wlb.Add(sx.MakeSymbol("select"))
	wlb.Add(attrLb.List())
	for i := 0; i < len(se.choices); i += 2 {
		choice := se.choices[i]
		text := se.choices[i+1]
		var alb sx.ListBuilder
		alb.Add(sxhtml.SymAttr)
		alb.Add(sx.Cons(sx.MakeSymbol("value"), sx.MakeString(choice)))
		if se.value == choice {
			alb.Add(sx.Cons(sx.MakeSymbol("selected"), sx.Nil()))
		}
		wlb.Add(sx.MakeList(sx.MakeSymbol("option"), alb.List(), sx.MakeString(text)))
	}

	flb.Add(wlb.List())
	return flb.List()
}

// ----- General utility functions for rendering etc.

func renderLabel(fieldID, label string) *sx.Pair {
	if label == "" {
		return nil
	}
	return sx.MakeList(
		sx.MakeSymbol("label"),
		sx.MakeList(
			sxhtml.SymAttr,
			sx.Cons(sx.MakeSymbol("for"), sx.MakeString(fieldID)),
		),
		sx.MakeString(label),
	)
}

func renderMessages(messages []string) *sx.Pair {
	var lb sx.ListBuilder
	for _, msg := range messages {
		lb.Add(sx.MakeList(
			sx.MakeSymbol("span"),
			sx.MakeList(
				sxhtml.SymAttr,
				sx.Cons(sx.MakeSymbol("class"), sx.MakeString("message")),
			),
			sx.MakeString(msg),
		))
	}
	return lb.List()
}

func renderValidators(validators []Validator) *sx.Pair {
	var lb sx.ListBuilder
	for _, validator := range validators {
		lb.ExtendBang(validator.Attributes())
	}
	return lb.List()
}
