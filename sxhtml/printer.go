//-----------------------------------------------------------------------------
// Copyright (c) 2023-present Detlef Stern
//
// This file is part of sxwebs.
//
// sxwebs is licensed under the latest version of the EUPL // (European Union
// Public License). Please see file LICENSE.txt for your rights and obligations
// under this license.
//
// SPDX-License-Identifier: EUPL-1.2
// SPDX-FileCopyrightText: 2023-present Detlef Stern
//-----------------------------------------------------------------------------

package sxhtml

import (
	"fmt"
	"io"
	"strings"

	"t73f.de/r/webs/htmls/render"
)

type printer struct {
	w   io.Writer
	err error
}

func (pr *printer) printString(s string) {
	if pr.err == nil {
		_, pr.err = io.WriteString(pr.w, s)
	}
}
func (pr *printer) printStrings(sl ...string) {
	if pr.err == nil {
		for _, s := range sl {
			_, err := io.WriteString(pr.w, s)
			if err != nil {
				pr.err = err
				return
			}
		}
	}
}

func (pr *printer) printHTML(s string) {
	if pr.err == nil {
		pr.err = render.Escape(pr.w, s)
	}
}

func (pr *printer) printComment(s string) {
	if pr.err == nil {
		pr.err = render.EscapeComment(pr.w, s)
	}
}

func (pr *printer) printAttributeValue(t attrType, s string) {
	if pr.err == nil {
		switch t {
		case attrPlain, attrCSS, attrJS:
			pr.err = render.EscapeAttrValue(pr.w, s)
		case attrURL:
			var sb strings.Builder
			sb.Grow(len(s) * 2)
			if pr.err = render.EscapeURL(&sb, s); pr.err == nil {
				pr.err = render.EscapeAttrValue(pr.w, sb.String())
			}
		default:
			pr.err = fmt.Errorf("unknown attribute type: %v", t)
		}
	}
}
