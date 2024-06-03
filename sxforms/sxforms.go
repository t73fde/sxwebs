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
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"t73f.de/r/sx"
	"t73f.de/r/sxwebs/sxhtml"
)

// Form represents a HTML form.
type Form struct {
	action      string
	method      string
	maxFormSize int64
	fields      []Field
	fieldnames  map[string]Field
	messages    Messages
}

// Define builds a new form.
func Define(fields ...Field) *Form {
	fieldnames := make(map[string]Field, len(fields))
	for _, field := range fields {
		fieldnames[field.Name()] = field
	}
	return &Form{
		method:      http.MethodPost,
		maxFormSize: (10 << 20), // 10 MB
		fields:      fields,
		fieldnames:  fieldnames,
	}
}

// Add a field.
func (f *Form) Add(field Field) *Form {
	f.fields = append(f.fields, field)
	f.fieldnames[field.Name()] = field
	return f
}

// SetAction updates the "action" URL attribute.
func (f *Form) SetAction(action string) *Form { f.action = action; return f }

// SetMethodGET updates the "method" attribute to the value "GET".
func (f *Form) SetMethodGET() *Form { f.method = http.MethodGet; return f }

// Clear all field data and messages.
func (f *Form) Clear() {
	for _, field := range f.fields {
		field.Clear()
	}
	f.messages = nil
}

// Disable the form.
func (f *Form) Disable() *Form {
	for _, field := range f.fields {
		field.Disable()
	}
	return f
}

// Messages contains all messages, as a map of field names to a list of string.
// Messages for the whole form will use the empty string as a field name.
type Messages map[string][]string

// Add a new message for the given field.
func (m Messages) Add(fieldName, message string) Messages {
	if len(m) == 0 {
		return Messages{fieldName: {message}}
	}
	m[fieldName] = append(m[fieldName], message)
	return m
}

// Data contains all form data, as a map of field names to field values.
type Data map[string]string

// Get string data of a field. Return empty string for unknwon field.
func (d Data) Get(fieldName string) string {
	if len(d) == 0 {
		return ""
	}
	if value, found := d[fieldName]; found {
		return value
	}
	return ""
}

// GetDate returns the value of the given field as a time.Time, but only
// as a real date, with time 00:00:00.
func (d Data) GetDate(fieldName string) time.Time {
	if len(d) > 0 {
		if value, found := d[fieldName]; found {
			if result, err := time.Parse(htmlDateLayout, value); err == nil {
				return result
			}
		}
	}
	return time.Time{}
}

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

// SetData set field values according to the given data.
func (f *Form) SetData(data Data) bool {
	ok := true
	for name, value := range data {
		field, found := f.fieldnames[name]
		if !found {
			// Unknown field name --> ignore
			continue
		}
		err := field.SetValue(strings.TrimSpace(value))
		if err != nil {
			f.messages = f.messages.Add(name, err.Error())
			ok = false
		}
	}
	return ok
}

// SetFormValues populates the form with the given URL values.
func (f *Form) SetFormValues(vals url.Values, _ *multipart.Form) bool {
	if len(vals) == 0 {
		return true
	}
	data := make(Data, len(vals))
	for name, values := range vals {
		value := ""
		if len(values) > 0 {
			value = values[0]
		}
		data[name] = value
	}
	return f.SetData(data)
}

// ValidRequestForm populates the form with the values of the given HTTP request,
// and validates them.
func (f *Form) ValidRequestForm(r *http.Request) bool {
	if f.method == http.MethodPost {
		return f.ValidOnSubmit(r)
	}
	return f.SetFormValues(r.URL.Query(), nil) && f.IsValid()
}

// ValidOnSubmit return true, if the request method was "POST" and the
// populated values of the request were successully validated.
func (f *Form) ValidOnSubmit(r *http.Request) bool {
	if r.Method != http.MethodPost {
		return false
	}
	if err := f.parseForm(r); err != nil {
		f.messages = Messages{"": {err.Error()}}
		return false
	}
	return f.SetFormValues(r.PostForm, r.MultipartForm) && f.IsValid()
}

// parseForm uses the approriate form parser, depending on the request.
//
// Until there is no FileElement, an ordinary ParseForm is suffcient.
// When a FileElement is added, the form must use a different encoding
// "multipart/form-data", instead of the default value
// "application/x-www-form-urlencoded".
func (f *Form) parseForm(r *http.Request) (err error) {
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		ct, _, err = mime.ParseMediaType(ct)
		if err != nil {
			return err
		}
	}
	if ct == "multipart/form-data" {
		return r.ParseMultipartForm(f.maxFormSize)
	}
	return r.ParseForm()
}

// IsValid returns true if the form has been successfully validates.
func (f *Form) IsValid() bool {
	var messages Messages
	for _, field := range f.fields {
		fieldName := field.Name()
		for _, validator := range field.Validators() {
			if err := validator.Check(field); err != nil {
				messages = messages.Add(fieldName, err.Error())
				if _, isStop := err.(StopValidationError); isStop {
					break
				}
			}
		}
	}
	f.messages = messages
	return len(messages) == 0
}

// Messages return the map of error messages, from an earlier validation.
func (f *Form) Messages() Messages { return f.messages }

// Render the form as an sx list.
func (f *Form) Render() *sx.Pair {
	var lb sx.ListBuilder
	lb.Add(sx.MakeSymbol("form"))
	lb.Add(sx.MakeList(
		sxhtml.SymAttr,
		sx.Cons(sx.MakeSymbol("action"), sx.MakeString(f.action)),
		sx.Cons(sx.MakeSymbol("method"), sx.MakeString(f.method)),
	))
	var submitLb sx.ListBuilder
	for _, field := range f.fields {
		fieldID := f.calcFieldID(field)
		if submitField, isSubmit := field.(*SubmitElement); isSubmit {
			if submitLb.IsEmpty() {
				submitLb.Add(sx.MakeSymbol("div"))
			}
			submitLb.Add(submitField.Render(fieldID, nil))
			continue
		}
		if submitList := submitLb.List(); submitList != nil {
			lb.Add(submitList)
		}
		lb.Add(field.Render(fieldID, f.messages[field.Name()]))
	}
	if submitList := submitLb.List(); submitList != nil {
		lb.Add(submitList)
	}
	return lb.List()
}

func (*Form) calcFieldID(field Field) string { return field.Name() }
