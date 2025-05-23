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
	"strconv"
	"unicode/utf8"

	"t73f.de/r/sx"
)

// Validator is used to check if a field value is valid.
// In addition, it supports field rendering by adding HTML form field attributes.
type Validator interface {
	Check(*Form, Field) error

	// Attributes contain additional HTML attributes for a field.
	Attributes() *sx.Pair
}

// Validators is a sequence of Validator.
type Validators []Validator

// HasRequired returns true, if there is at least the Required validator.
func (vs Validators) HasRequired() bool {
	for _, v := range vs {
		if _, ok := v.(Required); ok {
			return true
		}
	}
	return false
}

// ValidationError is an error that wraps a validator error message that should
// allow further validation of the field.
type ValidationError string

func (ve ValidationError) Error() string { return string(ve) }

// StopValidationError is a validation error that stops further validation of the field.
type StopValidationError string

func (sve StopValidationError) Error() string { return string(sve) }

// ----- Required: field must have a value.

// Required is a validator that checks if data is available.
type Required struct{ Message string }

// Check the given field w.r.t. to this validator.
func (ir Required) Check(_ *Form, field Field) error {
	if field.Value() != "" {
		return nil
	}
	if ir.Message == "" {
		return StopValidationError("Required")
	}
	return StopValidationError(ir.Message)
}

// Attributes returns HTML attributes as a Sx cons list.
func (Required) Attributes() *sx.Pair {
	return sx.MakeList(sx.Cons(sx.MakeSymbol("required"), sx.Nil()))
}

// ----- Optional: field must not have a value, could be missing.

// Optional is a validator that stops all further validation, if field has no value.
type Optional struct{}

// Check the given field w.r.t. to this validator.
func (Optional) Check(_ *Form, field Field) error {
	if field.Value() != "" {
		return nil
	}
	return StopValidationError("")
}

// Attributes returns HTML attributes as a Sx cons list.
func (Optional) Attributes() *sx.Pair { return nil }

// ----- MinMaxLength: field must have a value of a specific length.

// MinMaxLength is a validator that checks for a length.
type MinMaxLength struct {
	MinLength int
	MaxLength int
}

// Check the given field w.r.t. to this validator.
func (mml *MinMaxLength) Check(_ *Form, field Field) error {
	if minl, curl := mml.MinLength, utf8.RuneCountInString(field.Value()); minl > 0 && curl < minl {
		return ValidationError(fmt.Sprintf("minimum length of %s is %d, but got %d", field.Name(), minl, curl))
	}
	if maxl, curl := mml.MaxLength, utf8.RuneCountInString(field.Value()); maxl > 0 && curl > maxl {
		return ValidationError(fmt.Sprintf("maximum length of %s is %d, but got %d", field.Name(), maxl, curl))
	}
	return nil
}

// Attributes returns HTML attributes as a Sx cons list.
func (mml *MinMaxLength) Attributes() (result *sx.Pair) {
	if minl := mml.MinLength; minl > 0 {
		result = result.Cons(sx.Cons(sx.MakeSymbol("minlength"), sx.MakeString(strconv.Itoa(minl))))
	}
	if maxl := mml.MaxLength; maxl > 0 {
		result = result.Cons(sx.Cons(sx.MakeSymbol("maxlength"), sx.MakeString(strconv.Itoa(maxl))))
	}
	return result
}

// ----- MinValue: field must have a minimum value.

// MinValue is a validator that checks for a minimum value.
type MinValue struct {
	Value string
}

// Check the given field w.r.t. to this validator.
func (mv *MinValue) Check(_ *Form, field Field) error {
	val := field.Value()
	switch f := field.(type) {
	case *InputElement:
		switch f.itype {
		case itypeNumber:
			fvalue, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return ValidationError(fmt.Sprintf("%s does not contain a number: %v", field.Name(), val))
			}
			mvalue, err := strconv.ParseFloat(mv.Value, 64)
			if err == nil && fvalue < mvalue {
				return ValidationError(fmt.Sprintf(
					"minimum value of %s is %v, but got %v", field.Name(), mv.Value, val))
			}
		case itypeDate: // TODO
		case itypeDatetime: // TODO
		}
	}
	return nil
}

// Attributes returns HTML attributes as a Sx cons list.
func (mv *MinValue) Attributes() (result *sx.Pair) {
	return sx.MakeList(sx.Cons(sx.MakeSymbol("min"), sx.MakeString(mv.Value)))
}

// ----- MaxValue: field must have a maximum value.

// MaxValue is a validator that checks for a maximum value.
type MaxValue struct {
	Value string
}

// Check the given field w.r.t. to this validator.
func (mv *MaxValue) Check(_ *Form, field Field) error {
	val := field.Value()
	switch f := field.(type) {
	case *InputElement:
		switch f.itype {
		case itypeNumber:
			fvalue, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return ValidationError(fmt.Sprintf("%s does not contain a number: %v", field.Name(), val))
			}
			mvalue, err := strconv.ParseFloat(mv.Value, 64)
			if err == nil && fvalue > mvalue {
				return ValidationError(fmt.Sprintf(
					"minimum value of %s is %v, but got %v", field.Name(), mv.Value, val))
			}
		case itypeDate: // TODO
		case itypeDatetime: // TODO
		}
	}
	return nil
}

// Attributes returns HTML attributes as a Sx cons list.
func (mv *MaxValue) Attributes() *sx.Pair {
	return sx.MakeList(sx.Cons(sx.MakeSymbol("max"), sx.MakeString(mv.Value)))
}

// ----- UInt: field must have an unsigned integer value.

// UInt is a validator that checks for an unsigned integer value.
type UInt struct{}

// Check the given field w.r.t. to this validator.
func (u UInt) Check(_ *Form, field Field) error {
	val := field.Value()
	if _, err := strconv.ParseUint(val, 10, 64); err != nil {
		return ValidationError(fmt.Sprintf("%s does not contain an unsigned integer value: %v", field.Name(), val))
	}
	return nil
}

// Attributes returns HTML attributes as a Sx cons list.
func (UInt) Attributes() *sx.Pair { return sx.Nil() }

// ----- StringXXX: field must have a value that compares to a specific constant.

// StringLess performs a string comparison with the given field.
func StringLess(value string, msg string) Validator {
	return &stringCompare{value: value, op: -2, message: msg}
}

// StringLessEqual performs a string comparison with the given field.
func StringLessEqual(value string, msg string) Validator {
	return &stringCompare{value: value, op: -1, message: msg}
}

// StringEqual performs a string comparison with the given field.
func StringEqual(value string, msg string) Validator {
	return &stringCompare{value: value, op: 0, message: msg}
}

// StringGreaterEqual performs a string comparison with the given field.
func StringGreaterEqual(value string, msg string) Validator {
	return &stringCompare{value: value, op: 1, message: msg}
}

// StringGreater performs a string comparison with the given field.
func StringGreater(value string, msg string) Validator {
	return &stringCompare{value: value, op: 2, message: msg}
}

// stringCompare validates that the current field by comparing with the given one.
// Comparison is done via string comparison.
type stringCompare struct {
	value   string
	op      int
	message string
}

func (fsc *stringCompare) Check(_ *Form, field Field) error {
	return compareStringValues(fsc.op, field.Value(), fsc.value, fsc.message)
}

func (*stringCompare) Attributes() *sx.Pair { return nil }

func compareStringValues(op int, value, other string, msg string) error {
	var msgOp string
	switch op {
	case -2:
		if value < other {
			return nil
		}
		msgOp = "≥"
	case -1:
		if value <= other {
			return nil
		}
		msgOp = ">"
	case 0:
		if value == other {
			return nil
		}
		msgOp = "≠"
	case 1:
		if value >= other {
			return nil
		}
		msgOp = "<"
	case 2:
		if value > other {
			return nil
		}
		msgOp = "≤"
	default:
		return fmt.Errorf("comparison value not expected: %d", op)
	}
	if msg != "" {
		return ValidationError(msg)
	}
	return ValidationError(fmt.Sprintf("%v %s %v", value, msgOp, other))

}

// ----- FieldStringXXX: field must have a value that is compared to another field.

// FieldStringLess performs a string comparison with the given field.
func FieldStringLess(name string, msg string) Validator {
	return &fieldStringCompare{fieldname: name, op: -2, message: msg}
}

// FieldStringLessEqual performs a string comparison with the given field.
func FieldStringLessEqual(name string, msg string) Validator {
	return &fieldStringCompare{fieldname: name, op: -1, message: msg}
}

// FieldStringEqual performs a string comparison with the given field.
func FieldStringEqual(name string, msg string) Validator {
	return &fieldStringCompare{fieldname: name, op: 0, message: msg}
}

// FieldStringGreaterEqual performs a string comparison with the given field.
func FieldStringGreaterEqual(name string, msg string) Validator {
	return &fieldStringCompare{fieldname: name, op: 1, message: msg}
}

// FieldStringGreater performs a string comparison with the given field.
func FieldStringGreater(name string, msg string) Validator {
	return &fieldStringCompare{fieldname: name, op: 2, message: msg}
}

// fieldStringCompare validates that the current field by comparing with the given one.
// Comparison is done via string comparison.
type fieldStringCompare struct {
	fieldname string
	op        int
	message   string
}

func (fsc *fieldStringCompare) Check(f *Form, field Field) error {
	other, err := f.Field(fsc.fieldname)
	if err != nil {
		return err
	}
	return compareStringValues(fsc.op, field.Value(), other.Value(), fsc.message)
}

func (*fieldStringCompare) Attributes() *sx.Pair { return nil }
