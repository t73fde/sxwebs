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

import (
	"fmt"
	"strings"
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
	Validators() Validators
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
	validators Validators
	disabled   bool
}

// Constants for InputField.itype
const (
	itypeCheckbox = "checkbox"
	itypeDate     = "date"
	itypeDatetime = "datetime-local"
	itypeEmail    = "email"
	itypeNumber   = "number"
	itypePassword = "password"
	itypeText     = "text"
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
func (fd *InputElement) SetValue(value string) error {
	fd.value = value
	switch fd.itype {
	case itypeDate:
		if value != "" {
			if _, err := time.Parse(htmlDateLayout, value); err != nil {
				return err
			}
		}
	}
	return nil
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
	flb.Add(sx.MakeSymbol("div"))
	if label := renderLabel(fd, fieldID, fd.label); label != nil {
		flb.Add(label)
	}
	flb.ExtendBang(renderMessages(messages))

	var attrLb sx.ListBuilder
	attrLb.AddN(
		sxhtml.SymAttr,
		sx.Cons(sx.MakeSymbol("id"), sx.MakeString(fieldID)),
		sx.Cons(sx.MakeSymbol("name"), sx.MakeString(fd.name)),
		sx.Cons(sx.MakeSymbol("type"), sx.MakeString(fd.itype)),
		sx.Cons(sx.MakeSymbol("value"), sx.MakeString(fd.value)),
	)
	addEnablingAttributes(&attrLb, fd.disabled, fd.validators)

	flb.Add(sx.MakeList(sx.MakeSymbol("input"), attrLb.List()))
	return flb.List()
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
func DateValue(t time.Time) string { return t.Format(htmlDateLayout) }

// DatetimeField builds a new field to enter dates.
func DatetimeField(name, label string, validators ...Validator) *InputElement {
	return &InputElement{
		itype:      itypeDatetime,
		name:       name,
		label:      label,
		validators: validators,
	}
}

// DatetimeValue returns the time as a string suitable for a HTML datetime-local field value.
func DatetimeValue(t time.Time) string { return t.Format(htmlDatetimeLayout) }

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

// ----- Submit input element

// SubmitElement represents an element <input type="submit" ...>
type SubmitElement struct {
	name     string
	label    string
	value    string
	prio     uint8
	disabled bool
	isCancel bool
}

// SubmitField builds a new submit field.
func SubmitField(name, label string) *SubmitElement {
	return &SubmitElement{
		name:  name,
		label: label,
	}
}

// SetPriority sets the importance of the field. Only the values 0, 1, 2, and 3
// are allowed, with 0 being the highest priority.
func (se *SubmitElement) SetPriority(prio uint8) *SubmitElement {
	se.prio = min(prio, uint8(len(submitPrioClass)-1))
	return se
}

var submitPrioClass = map[uint8]string{
	0: "primary",
	1: "secondary",
	2: "tertiary",
	3: "cancel",
}

// MarkCancel marks the submit field as an action that disables form
// validation, if this field causes the form to be sent.
func (se *SubmitElement) MarkCancel() *SubmitElement {
	se.isCancel = true
	return se
}

// Name returns the name of this element.
func (se *SubmitElement) Name() string { return se.name }

// Value returns the value of this element.
func (se *SubmitElement) Value() string { return se.value }

// Clear the element.
func (se *SubmitElement) Clear() { se.value = "" }

// SetValue sets the value of this element.
func (se *SubmitElement) SetValue(value string) error { se.value = value; return nil }

// Validators return the currently active validators.
func (se *SubmitElement) Validators() Validators { return nil }

// Disable the submit element.
func (se *SubmitElement) Disable() { se.disabled = true }

// Render the submit element as SxHTML.
func (se *SubmitElement) Render(fieldID string, _ []string) *sx.Pair {
	var attrLb sx.ListBuilder
	attrLb.AddN(
		sxhtml.SymAttr,
		sx.Cons(sx.MakeSymbol("id"), sx.MakeString(fieldID)),
		sx.Cons(sx.MakeSymbol("name"), sx.MakeString(se.name)),
		sx.Cons(sx.MakeSymbol("type"), sx.MakeString("submit")),
		sx.Cons(sx.MakeSymbol("value"), sx.MakeString(se.label)),
		sx.Cons(sx.MakeSymbol("class"), sx.MakeString(submitPrioClass[se.prio])),
	)
	addBoolAttribute(&attrLb, sx.MakeSymbol("disabled"), se.disabled)

	return sx.MakeList(sx.MakeSymbol("input"), attrLb.List())
}

// ----- Checkbox field

// CheckboxElement represents a checkbox.
type CheckboxElement struct {
	name     string
	label    string
	value    string
	disabled bool
}

// CheckboxField provides a checkbox.
func CheckboxField(name, label string) *CheckboxElement {
	return &CheckboxElement{
		name:  name,
		label: label,
	}
}

// Name returns the name of this element.
func (cbe *CheckboxElement) Name() string { return cbe.name }

// Value returns the value of this element.
func (cbe *CheckboxElement) Value() string { return cbe.value }

// Clear the element.
func (cbe *CheckboxElement) Clear() { cbe.value = "" }

// SetValue sets the value of this element.
func (cbe *CheckboxElement) SetValue(value string) error { cbe.value = value; return nil }

// Validators return the currently active validators.
func (cbe *CheckboxElement) Validators() Validators { return nil }

// Disable the checkbox element.
func (cbe *CheckboxElement) Disable() { cbe.disabled = true }

// Render the checkbox element as SxHTML.
func (cbe *CheckboxElement) Render(fieldID string, _ []string) *sx.Pair {

	var attrLb sx.ListBuilder
	attrLb.AddN(
		sxhtml.SymAttr,
		sx.Cons(sx.MakeSymbol("id"), sx.MakeString(fieldID)),
		sx.Cons(sx.MakeSymbol("name"), sx.MakeString(cbe.name)),
		sx.Cons(sx.MakeSymbol("type"), sx.MakeString("checkbox")),
		sx.Cons(sx.MakeSymbol("value"), sx.MakeString(cbe.value)),
	)
	addEnablingAttributes(&attrLb, cbe.disabled, nil)

	var flb sx.ListBuilder
	flb.Add(sx.MakeSymbol("div"))
	flb.Add(sx.MakeList(sx.MakeSymbol("input"), attrLb.List()))
	if label := renderLabel(cbe, fieldID, cbe.label); label != nil {
		flb.Add(label)
	}
	return flb.List()
}

// ----- <textarea ...>...</textarea> field

// TextAreaElement represents the corresponding textarea form element.
type TextAreaElement struct {
	name       string
	label      string
	rows       uint32
	cols       uint32
	value      string
	validators Validators
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

// SetRows sets the number of rows for the text area element.
func (tae *TextAreaElement) SetRows(rows uint32) *TextAreaElement {
	tae.rows = rows
	return tae
}

// SetCols sets the number of columns for the text area, i.e. the number of
// possibly visible lines.
func (tae *TextAreaElement) SetCols(cols uint32) *TextAreaElement {
	tae.cols = cols
	return tae
}

// Name returns the name of the text area element.
func (tae *TextAreaElement) Name() string { return tae.name }

// Value returns the value of the text area.
func (tae *TextAreaElement) Value() string { return tae.value }

// Clear the text area.
func (tae *TextAreaElement) Clear() { tae.value = "" }

// SetValue sets the value of the text area. Sequences of '\r\n' will be replaced by '\n'.
func (tae *TextAreaElement) SetValue(value string) error {
	tae.value = strings.ReplaceAll(value, "\r\n", "\n") // Unify Windows/Unix EOL handling
	return nil
}

// Validators returns the currently active validators for this text area.
func (tae *TextAreaElement) Validators() Validators {
	if tae.disabled {
		return nil
	}
	return tae.validators
}

// Disable the text area element.
func (tae *TextAreaElement) Disable() { tae.disabled = true }

// Render the text area as SxHTML.
func (tae *TextAreaElement) Render(fieldID string, messages []string) *sx.Pair {
	var flb sx.ListBuilder
	flb.Add(sx.MakeSymbol("div"))
	if label := renderLabel(tae, fieldID, tae.label); label != nil {
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
	addEnablingAttributes(&attrLb, tae.disabled, tae.validators)

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
	validators Validators
	disabled   bool
}

// SelectField creates a new select element.
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

// Name returns the element name.
func (se *SelectElement) Name() string { return se.name }

// Value returns the value of the select element.
func (se *SelectElement) Value() string { return se.value }

// Clear the select element.
func (se *SelectElement) Clear() { se.value = "" }

// SetValue sets the value of the select element.
func (se *SelectElement) SetValue(value string) error {
	se.value = value
	for i := 0; i < len(se.choices); i += 2 {
		if se.choices[i] == value {
			return nil
		}
	}
	return fmt.Errorf("no such choice: %q", value)
}

// Validators return the active validators for the select element.
func (se *SelectElement) Validators() Validators {
	if se.disabled {
		return nil
	}
	return se.validators
}

// Disable the field.
func (se *SelectElement) Disable() { se.disabled = true }

// Render the select element as SxHTML.
func (se *SelectElement) Render(fieldID string, messages []string) *sx.Pair {
	var flb sx.ListBuilder
	flb.Add(sx.MakeSymbol("div"))
	if label := renderLabel(se, fieldID, se.label); label != nil {
		flb.Add(label)
	}
	flb.ExtendBang(renderMessages(messages))
	var attrLb sx.ListBuilder
	attrLb.AddN(
		sxhtml.SymAttr,
		sx.Cons(sx.MakeSymbol("id"), sx.MakeString(fieldID)),
		sx.Cons(sx.MakeSymbol("name"), sx.MakeString(se.name)),
	)
	addEnablingAttributes(&attrLb, se.disabled, se.validators)

	var wlb sx.ListBuilder
	wlb.AddN(sx.MakeSymbol("select"), attrLb.List())
	for i := 0; i < len(se.choices); i += 2 {
		choice := se.choices[i]
		text := se.choices[i+1]
		var alb sx.ListBuilder
		alb.AddN(sxhtml.SymAttr, sx.Cons(sx.MakeSymbol("value"), sx.MakeString(choice)))
		addBoolAttribute(&alb, sx.MakeSymbol("disabled"), choice == "")
		addBoolAttribute(&alb, sx.MakeSymbol("selected"), se.value == choice)
		wlb.Add(sx.MakeList(sx.MakeSymbol("option"), alb.List(), sx.MakeString(text)))
	}

	flb.Add(wlb.List())
	return flb.List()
}

// ----- General utility functions for rendering etc.

func renderLabel(field Field, fieldID, label string) *sx.Pair {
	if label == "" {
		return nil
	}
	var lb sx.ListBuilder
	lb.AddN(
		sx.MakeSymbol("label"),
		sx.MakeList(
			sxhtml.SymAttr,
			sx.Cons(sx.MakeSymbol("for"), sx.MakeString(fieldID)),
		),
		sx.MakeString(label),
	)
	if field.Validators().HasRequired() {
		lb.Add(sx.MakeString("*"))
	}
	return lb.List()
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

// addEnablingAttributes adds some attributes, depending whether the field is
// disabled or not. If it is disabled, the "disabled" attribute will be added,
// and no validator attributes are added too.
// Otherwise, the field is enable and therefore the attributes of an validator
// will be added.
func addEnablingAttributes(lb *sx.ListBuilder, disabled bool, validators []Validator) {
	if disabled {
		lb.Add(sx.Cons(sx.MakeSymbol("disabled"), sx.Nil()))
	} else {
		lb.ExtendBang(renderValidators(validators))
	}
}
