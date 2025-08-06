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

// ----- Submit input element

// SubmitElement represents an element <input type="submit" ...>
type SubmitElement struct {
	name           string
	label          string
	value          string
	prio           uint8
	disabled       bool
	noFormValidate bool
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
	3: "cancel", // must always be the last, see se.SetCancel()
}

// NoFormValidate marks the submit field as an action that disables form
// validation, if this field causes the form to be sent.
func (se *SubmitElement) NoFormValidate() *SubmitElement {
	se.noFormValidate = true
	return se
}

// SetCancel marks the submit field to work as as a cancel button.
func (se *SubmitElement) SetCancel() *SubmitElement {
	se.prio = uint8(len(submitPrioClass) - 1)
	se.noFormValidate = true
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
		sx.Cons(sxhtml.MakeSymbol("id"), sx.MakeString(fieldID)),
		sx.Cons(sxhtml.MakeSymbol("name"), sx.MakeString(se.name)),
		sx.Cons(sxhtml.MakeSymbol("type"), sx.MakeString("submit")),
		sx.Cons(sxhtml.MakeSymbol("value"), sx.MakeString(se.label)),
		sx.Cons(sxhtml.MakeSymbol("class"), sx.MakeString(submitPrioClass[se.prio])),
	)
	addBoolAttribute(&attrLb, sxhtml.MakeSymbol("disabled"), se.disabled)
	addBoolAttribute(&attrLb, sxhtml.MakeSymbol("formnovalidate"), se.noFormValidate)

	return sx.MakeList(sxhtml.MakeSymbol("input"), attrLb.List())
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

// SetChecked sets the value of the checkbox element.
func (cbe *CheckboxElement) SetChecked(val bool) {
	if val {
		cbe.value = "on"
	} else {
		cbe.value = ""
	}
}

// CheckedValue returns the date as a string suitable for a HTML checked field value.
func CheckedValue(b bool) string {
	if b {
		return "on"
	}
	return ""
}

// Validators return the currently active validators.
func (cbe *CheckboxElement) Validators() Validators { return nil }

// Disable the checkbox element.
func (cbe *CheckboxElement) Disable() { cbe.disabled = true }

// Render the checkbox element as SxHTML.
func (cbe *CheckboxElement) Render(fieldID string, _ []string) *sx.Pair {

	var attrLb sx.ListBuilder
	attrLb.AddN(
		sx.Cons(sxhtml.MakeSymbol("id"), sx.MakeString(fieldID)),
		sx.Cons(sxhtml.MakeSymbol("name"), sx.MakeString(cbe.name)),
		sx.Cons(sxhtml.MakeSymbol("type"), sx.MakeString("checkbox")),
		sx.Cons(sxhtml.MakeSymbol("value"), sx.MakeString(cbe.name)),
	)
	if cbe.value != "" {
		attrLb.Add(sx.Cons(sxhtml.MakeSymbol("checked"), sx.Nil()))
	}
	addEnablingAttributes(&attrLb, cbe.disabled, nil)

	var flb sx.ListBuilder
	flb.Add(sxhtml.MakeSymbol("div"))
	flb.Add(sx.MakeList(sxhtml.MakeSymbol("input"), attrLb.List()))
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
	flb.Add(sxhtml.MakeSymbol("div"))
	if label := renderLabel(tae, fieldID, tae.label); label != nil {
		flb.Add(label)
	}
	flb.ExtendBang(renderMessages(messages))
	var attrLb sx.ListBuilder
	attrLb.AddN(
		sx.Cons(sxhtml.MakeSymbol("id"), sx.MakeString(fieldID)),
		sx.Cons(sxhtml.MakeSymbol("name"), sx.MakeString(tae.name)),
	)
	if rows := tae.rows; rows > 0 {
		attrLb.Add(sx.Cons(sxhtml.MakeSymbol("rows"), sx.MakeString(fmt.Sprint(rows))))
	}
	if cols := tae.cols; cols > 0 {
		attrLb.Add(sx.Cons(sxhtml.MakeSymbol("cols"), sx.MakeString(fmt.Sprint(cols))))
	}
	addEnablingAttributes(&attrLb, tae.disabled, tae.validators)

	flb.Add(sx.MakeList(sxhtml.MakeSymbol("textarea"), attrLb.List(), sx.MakeString(tae.value)))
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
	flb.Add(sxhtml.MakeSymbol("div"))
	if label := renderLabel(se, fieldID, se.label); label != nil {
		flb.Add(label)
	}
	flb.ExtendBang(renderMessages(messages))
	var attrLb sx.ListBuilder
	attrLb.AddN(
		sx.Cons(sxhtml.MakeSymbol("id"), sx.MakeString(fieldID)),
		sx.Cons(sxhtml.MakeSymbol("name"), sx.MakeString(se.name)),
	)
	addEnablingAttributes(&attrLb, se.disabled, se.validators)

	var wlb sx.ListBuilder
	wlb.AddN(sxhtml.MakeSymbol("select"), attrLb.List())
	for i := 0; i < len(se.choices); i += 2 {
		choice := se.choices[i]
		text := se.choices[i+1]
		var alb sx.ListBuilder
		alb.Add(sx.Cons(sxhtml.MakeSymbol("value"), sx.MakeString(choice)))
		addBoolAttribute(&alb, sxhtml.MakeSymbol("disabled"), choice == "")
		addBoolAttribute(&alb, sxhtml.MakeSymbol("selected"), se.value == choice)
		wlb.Add(sx.MakeList(sxhtml.MakeSymbol("option"), alb.List(), sx.MakeString(text)))
	}

	flb.Add(wlb.List())
	return flb.List()
}

// ----- Flow Content -----

// FlowContentElement adds some flow content to the form.
type FlowContentElement struct {
	name    string
	content *sx.Pair
}

// FlowContentField allows to add some text (aka flow content) to the form.
func FlowContentField(name string, content *sx.Pair) *FlowContentElement {
	return &FlowContentElement{name: name, content: content}
}

// Name returns the element name.
func (fce *FlowContentElement) Name() string { return fce.name }

// Value returns the value of the select element.
func (*FlowContentElement) Value() string { return "" }

// Clear the select element.
func (*FlowContentElement) Clear() {}

// SetValue sets the value of the select element.
func (*FlowContentElement) SetValue(value string) error {
	return fmt.Errorf("flow content has no specific value")
}

// Validators return the active validators for the select element.
func (*FlowContentElement) Validators() Validators { return nil }

// Disable the field.
func (*FlowContentElement) Disable() {}

// Render the select element as SxHTML.
func (fce *FlowContentElement) Render(fieldID string, messages []string) *sx.Pair {
	return fce.content
}

// ----- General utility functions for rendering etc.

func renderLabel(field Field, fieldID, label string) *sx.Pair {
	if label == "" {
		return nil
	}
	var lb sx.ListBuilder
	lb.AddN(
		sxhtml.MakeSymbol("label"),
		sx.MakeList(sx.Cons(sxhtml.MakeSymbol("for"), sx.MakeString(fieldID))),
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
			sxhtml.MakeSymbol("span"),
			sx.MakeList(sx.Cons(sxhtml.MakeSymbol("class"), sx.MakeString("message"))),
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
		lb.Add(sx.Cons(sxhtml.MakeSymbol("disabled"), sx.Nil()))
	} else {
		lb.ExtendBang(renderValidators(validators))
	}
}
