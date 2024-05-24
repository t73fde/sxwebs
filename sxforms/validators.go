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

import "t73f.de/r/sx"

// Validator is used to check if a field value is valid.
// In addition, it supports field rendering by adding HTML form field attributes.
type Validator interface {
	Check(Field) error

	// Attributes contain additional HTML attributes for a field.
	Attributes() *sx.Pair
}

// ValidationError is an error that wraps a validator error message that should
// allow further validation of the field.
type ValidationError string

func (ve ValidationError) Error() string { return string(ve) }

// StopValidationError is a validation error that stops further validation of the field.
type StopValidationError string

func (sve StopValidationError) Error() string { return string(sve) }

// Required is a validator that checks if data is available.
type Required struct{ Message string }

func (ir Required) Check(field Field) error {
	if field.Value() != "" {
		return nil
	}
	if ir.Message == "" {
		return StopValidationError("Required")
	}
	return StopValidationError(ir.Message)
}

func (Required) Attributes() *sx.Pair {
	return sx.MakeList(sx.Cons(sx.MakeSymbol("required"), sx.Nil()))
}
