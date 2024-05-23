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
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"t73f.de/r/sx"
	"t73f.de/r/sxwebs/sxhtml"
)

// Form represents a HTML form.
type Form struct {
	fields     []Field
	fieldnames map[string]Field
	errors     map[string][]error
}

// Make builds a new form.
func Make(fields ...Field) *Form {
	fieldnames := make(map[string]Field, len(fields))
	for _, field := range fields {
		fieldnames[field.Name()] = field
	}
	return &Form{
		fields:     fields,
		fieldnames: fieldnames,
		errors:     nil,
	}
}

// Fields return the sequence of fields.
func (f *Form) Fields() []Field { return f.fields }

// Data returns the map of field names to values.
func (f *Form) Data() map[string]string {
	if len(f.fieldnames) == 0 {
		return nil
	}
	data := make(map[string]string, len(f.fieldnames))
	for name, field := range f.fieldnames {
		if value := field.Value(); value != "" {
			data[name] = value
		}
	}
	return data
}

// SetFormData populates the form with the given URL values.
func (f *Form) SetFormData(data url.Values) {
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
		field.SetValue(value)
	}
	f.errors = nil
}

// ValidateRequestForm populates the form with the values of the given HTTP request,
// and validates them.
func (f *Form) ValidateRequestForm(r *http.Request) bool {
	if err := r.ParseForm(); err != nil {
		f.errors = map[string][]error{"": {err}}
		return false
	}
	f.SetFormData(r.PostForm)
	return f.IsValid()
}

// IsValid returns true if the form has been successfully validates.
func (f *Form) IsValid() bool {
	var errors map[string][]error
	for _, field := range f.fields {
		fieldName := field.Name()
		value := field.Value()
		for _, validator := range field.Validators() {
			if err := validator.Check(value); err != nil {
				if len(errors) == 0 {
					errors = map[string][]error{fieldName: {err}}
				} else {
					errors[fieldName] = append(errors[fieldName], err)
				}
				if _, isStop := err.(StopValidationError); isStop {
					break
				}
			}
		}
	}
	if len(errors) == 0 {
		return true
	}
	f.errors = errors
	return false
}

// Errors return the map of errors from an earlier validation.
func (f *Form) Errors() map[string][]error { return f.errors }

// Messages return the map of error messages, from an earlier validation.
func (f *Form) Messages() map[string][]string {
	errors := f.errors
	if len(errors) == 0 {
		return nil
	}
	messages := make(map[string][]string, len(errors))
	for fieldname, fielderrors := range errors {
		msgs := make([]string, len(fielderrors))
		for errnum, err := range fielderrors {
			msgs[errnum] = err.Error()
		}
		messages[fieldname] = msgs
	}
	return messages
}

// Render the form as an sx.Object.
func (f *Form) Render() sx.Object {
	var lb sx.ListBuilder
	lb.Add(sx.MakeSymbol("form"))
	lb.Add(sx.MakeList(
		sxhtml.SymAttr,
		sx.Cons(sx.MakeSymbol("action"), sx.MakeString("")),
		sx.Cons(sx.MakeSymbol("method"), sx.MakeString("POST")),
	))
	for i, field := range f.fields {
		fieldID := fmt.Sprintf("field-%d", i)

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
		flb.Add(field.Render(fieldID))
		lb.Add(flb.List())
	}
	return lb.List()
}

// Field represents a HTTP form field.
type Field interface {
	Name() string
	Label() string
	Value() string
	SetValue(string)
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
func (fd *InputField) SetValue(value string)   { fd.value = value }
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
	return sx.MakeList(sx.MakeSymbol("input"), lb.List())
}

// SetAutofocus for the field.
func (fd *InputField) SetAutofocus() *InputField { fd.autofocus = true; return fd }

// MakeTextField builds a new text field.
func MakeTextField(name, label string, validators ...Validator) *InputField {
	return &InputField{
		itype:      "text",
		name:       name,
		label:      label,
		validators: validators,
	}
}

// MakePasswordField builds a new password field.
func MakePasswordField(name, label string, validators ...Validator) *InputField {
	return &InputField{
		itype:      "password",
		name:       name,
		label:      label,
		validators: validators,
	}
}

// MakeSubmitField builds a new submit field.
func MakeSubmitField(name, value string) *InputField {
	return &InputField{
		itype: "submit",
		name:  name,
		value: value,
	}
}

// Validator is use to check if a field value is valid.
type Validator interface {
	Check(string) error
}

// InputRequired is a validator that checks if data is available.
type InputRequired string

func (ir InputRequired) Check(value string) error {
	if strings.TrimSpace(value) != "" {
		return nil
	}
	if string(ir) == "" {
		return StopValidationError("Input required")
	}
	return StopValidationError(string(ir))
}

// StopValidationError is a validation error that stops further validation of the field.
type StopValidationError string

func (sve StopValidationError) Error() string { return string(sve) }
