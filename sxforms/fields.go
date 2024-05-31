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
	Disable()
	Render(string, []string) *sx.Pair
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
	disabled   bool
}

// Constants for InputField.itype
const (
	itypeDate     = "date"
	itypePassword = "password"
	itypeText     = "text"
)

func (fd *InputElement) Name() string  { return fd.name }
func (fd *InputElement) Value() string { return fd.value }
func (fd *InputElement) Clear()        { fd.value = "" }

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

func (fd *InputElement) Disable() { fd.disabled = true }

func (fd *InputElement) Render(fieldID string, messages []string) *sx.Pair {
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
	addBoolAttribute(&attrLb, sx.MakeSymbol("autofocus"), fd.autofocus)
	addBoolAttribute(&attrLb, sx.MakeSymbol("disabled"), fd.disabled)
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

// DateValue returns the date as a string suitable for a HTML date field value.
func DateValue(t time.Time) string { return t.Format(htmlDateLayout) }

// PasswordField builds a new password field.
func PasswordField(name, label string, validators ...Validator) *InputElement {
	return &InputElement{
		itype:      itypePassword,
		name:       name,
		label:      label,
		validators: validators,
	}
}

// ----- Submit input element

// SubmitElement represents an element <input type="submit" ...>
type SubmitElement struct {
	name     string
	label    string
	value    string
	prio     uint8
	disabled bool
}

// SubmitField builds a new submit field.
func SubmitField(name, label string) *SubmitElement {
	return &SubmitElement{
		name:  name,
		label: label,
	}
}

// SetPriority sets the importance of the field. Only the values 0, 1, and 2
// are allowed, with 0 being the highest priority.
func (se *SubmitElement) SetPriority(prio uint8) *SubmitElement {
	se.prio = min(prio, uint8(len(submitPrioClass)-1))
	return se
}

var submitPrioClass = map[uint8]string{
	0: "primary",
	1: "secondary",
	2: "tertiary",
}

func (se *SubmitElement) Name() string                { return se.name }
func (se *SubmitElement) Value() string               { return se.value }
func (se *SubmitElement) Clear()                      { se.value = "" }
func (se *SubmitElement) SetValue(value string) error { se.value = value; return nil }
func (se *SubmitElement) Validators() []Validator     { return nil }
func (se *SubmitElement) Disable()                    { se.disabled = true }
func (se *SubmitElement) Render(fieldID string, messages []string) *sx.Pair {
	var attrLb sx.ListBuilder
	attrLb.Add(sxhtml.SymAttr)
	attrLb.Add(sx.Cons(sx.MakeSymbol("id"), sx.MakeString(fieldID)))
	attrLb.Add(sx.Cons(sx.MakeSymbol("name"), sx.MakeString(se.name)))
	attrLb.Add(sx.Cons(sx.MakeSymbol("type"), sx.MakeString("submit")))
	attrLb.Add(sx.Cons(sx.MakeSymbol("value"), sx.MakeString(se.label)))
	attrLb.Add(sx.Cons(sx.MakeSymbol("class"), sx.MakeString(submitPrioClass[se.prio])))
	addBoolAttribute(&attrLb, sx.MakeSymbol("disabled"), se.disabled)

	return sx.MakeList(sx.MakeSymbol("input"), attrLb.List())
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
	disabled   bool
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
func (tae *TextAreaElement) Disable()                { tae.disabled = true }
func (tae *TextAreaElement) Render(fieldID string, messages []string) *sx.Pair {
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
	addBoolAttribute(&attrLb, sx.MakeSymbol("disabled"), tae.disabled)
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
	disabled   bool
}

// TextAreaField creates a new text area element.
func SelectField(name, label string, choices []string, validators ...Validator) *SelectElement {
	se := &SelectElement{
		name:       name,
		label:      label,
		validators: validators,
	}
	se.SetChoices(choices)
	return se
}

// SetChoices allows to update the choices after field creation, e.g. for
// dynamically generated choices.
func (se *SelectElement) SetChoices(choices []string) {
	if len(choices) == 0 || len(choices) == 1 {
		se.choices = nil
	} else if len(choices)%2 != 0 {
		se.choices = choices[0 : len(choices)-2]
	} else {
		se.choices = choices
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
func (se *SelectElement) Disable()                { se.disabled = true }
func (se *SelectElement) Render(fieldID string, messages []string) *sx.Pair {
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
	addBoolAttribute(&attrLb, sx.MakeSymbol("disabled"), se.disabled)
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
		addBoolAttribute(&alb, sx.MakeSymbol("selected"), se.value == choice)
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

func addBoolAttribute(lb *sx.ListBuilder, sym *sx.Symbol, val bool) {
	if val {
		lb.Add(sx.Cons(sym, sx.Nil()))
	}
}

func renderValidators(validators []Validator) *sx.Pair {
	var lb sx.ListBuilder
	for _, validator := range validators {
		lb.ExtendBang(validator.Attributes())
	}
	return lb.List()
}
