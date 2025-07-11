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

package sxforms

import (
	"t73f.de/r/sx"
	"t73f.de/r/sxwebs/sxhtml"
)

// Fieldset represents an HTML <fieldset>
type Fieldset struct {
	form     *Form
	name     string
	legend   string
	fields   []Field
	disabled bool
}

func (fs *Fieldset) setForm(f *Form) {
	for _, fd := range fs.fields {
		f.addName(fd)
	}
	fs.form = f
}

// FieldsetField builds a Fieldset.
func FieldsetField(name, legend string, fields ...Field) *Fieldset {
	return &Fieldset{
		form:     nil,
		name:     name,
		legend:   legend,
		fields:   fields,
		disabled: false,
	}
}

// Name the Fieldset.
func (fs *Fieldset) Name() string { return fs.name }

// Value returns the value of the Fieldset: there is no value.
func (Fieldset) Value() string { return "" }

// Clear the Fieldset.
func (fs *Fieldset) Clear() {
	for _, f := range fs.fields {
		f.Clear()
	}
}

// SetValue resetturns the value of the Fieldset: there is no value -> ignore
func (Fieldset) SetValue(string) error { return nil }

// Validators returns the validators for this Fieldset: there are no validators.
func (Fieldset) Validators() Validators { return nil }

// Disable the Fieldset.
func (fs *Fieldset) Disable() {
	for _, f := range fs.fields {
		f.Disable()
	}
}

// Render the Fieldset.
func (fs *Fieldset) Render(fieldID string, messages []string) *sx.Pair {
	var attrLb sx.ListBuilder
	attrLb.AddN(
		sxhtml.SymAttr,
		sx.Cons(sx.MakeSymbol("id"), sx.MakeString(fieldID)),
		sx.Cons(sx.MakeSymbol("name"), sx.MakeString(fs.name)),
	)
	addBoolAttribute(&attrLb, sx.MakeSymbol("disabled"), fs.disabled)

	form := fs.form
	var lb sx.ListBuilder
	lb.AddN(
		sx.MakeSymbol("fieldset"),
		attrLb.List(),
	)
	if legend := fs.legend; legend != "" {
		lb.Add(sx.MakeList(sx.MakeSymbol("legend"), sx.MakeString(legend)))
	}
	lb.ExtendBang(renderMessages(messages))
	for _, field := range fs.fields {
		lb.Add(field.Render(form.calcFieldID(field), form.messages[field.Name()]))
	}
	return lb.List()
}
