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

// Package forms handles HTML form data.
package sxforms

import (
	"net/http"
	"net/url"
	"strings"

	"t73f.de/r/sx"
	"t73f.de/r/sxwebs/sxhtml"
)

// Form represents a HTML form.
type Form struct {
	action     string
	method     string
	fields     []Field
	fieldnames map[string]Field
	messages   Messages
}

// Define builds a new form.
func Define(fields ...Field) *Form {
	fieldnames := make(map[string]Field, len(fields))
	for _, field := range fields {
		fieldnames[field.Name()] = field
	}
	return &Form{
		method:     http.MethodPost,
		fields:     fields,
		fieldnames: fieldnames,
	}
}

// SetAction updates the "action" URL attribute.
func (f *Form) SetAction(action string) *Form { f.action = action; return f }

// SetMethodGET updates the "method" attribute to the value "GET".
func (f *Form) SetMethodGET() *Form { f.method = http.MethodGet; return f }

// Messages contains all messages, as a map of field names to a list of string.
// Messages for the whole form will use the empty string as a field name.
type Messages map[string][]string

// Data contains all form data, as a map of field names to field values.
type Data map[string]string

// Fields return the sequence of fields.
func (f *Form) Fields() []Field { return f.fields }

// Data returns the map of field names to values.
func (f *Form) Data() Data {
	if len(f.fieldnames) == 0 {
		return nil
	}
	data := make(Data, len(f.fieldnames))
	for name, field := range f.fieldnames {
		if value := field.Value(); value != "" {
			data[name] = value
		}
	}
	return data
}

// SetFormValues populates the form with the given URL values.
func (f *Form) SetFormValues(data url.Values) {
	for name, values := range data {
		field, found := f.fieldnames[name]
		if !found {
			// Unknown field name --> ignore
			continue
		}
		value := ""
		if len(values) > 0 {
			value = values[0]
		}
		field.SetValue(strings.TrimSpace(value))
	}
	f.messages = nil
}

// ValidRequestForm populates the form with the values of the given HTTP request,
// and validates them.
func (f *Form) ValidRequestForm(r *http.Request) bool {
	if f.method == http.MethodGet {
		f.SetFormValues(r.URL.Query())
	} else {
		if err := r.ParseForm(); err != nil {
			f.messages = Messages{"": {err.Error()}}
			return false
		}
		f.SetFormValues(r.PostForm)
	}
	return f.IsValid()
}

// IsValid returns true if the form has been successfully validates.
func (f *Form) IsValid() bool {
	var messages Messages
	for _, field := range f.fields {
		fieldName := field.Name()
		for _, validator := range field.Validators() {
			if err := validator.Check(field); err != nil {
				if len(messages) == 0 {
					messages = Messages{fieldName: {err.Error()}}
				} else {
					messages[fieldName] = append(messages[fieldName], err.Error())
				}
				if _, isStop := err.(StopValidationError); isStop {
					break
				}
			}
		}
	}
	if len(messages) == 0 {
		return true
	}
	f.messages = messages
	return false
}

// Messages return the map of error messages, from an earlier validation.
func (f *Form) Messages() Messages { return f.messages }

// Render the form as an sx.Object.
func (f *Form) Render() sx.Object {
	var lb sx.ListBuilder
	lb.Add(sx.MakeSymbol("form"))
	lb.Add(sx.MakeList(
		sxhtml.SymAttr,
		sx.Cons(sx.MakeSymbol("action"), sx.MakeString(f.action)),
		sx.Cons(sx.MakeSymbol("method"), sx.MakeString(f.method)),
	))
	for _, field := range f.fields {
		fieldID := field.Name()

		var flb sx.ListBuilder
		flb.Add(sx.MakeSymbol("div"))
		if label := field.Label(); label != "" {
			flb.Add(sx.MakeList(
				sx.MakeSymbol("label"),
				sx.MakeList(
					sxhtml.SymAttr,
					sx.Cons(sx.MakeSymbol("for"), sx.MakeString(fieldID)),
				),
				sx.MakeString(label),
			))
		}

		for _, msg := range f.messages[field.Name()] {
			flb.Add(sx.MakeList(
				sx.MakeSymbol("span"),
				sx.MakeList(
					sxhtml.SymAttr,
					sx.Cons(sx.MakeSymbol("class"), sx.MakeString("message")),
				),
				sx.MakeString(msg),
			))
		}

		flb.Add(field.Render(fieldID))
		lb.Add(flb.List())
	}
	return lb.List()
}
