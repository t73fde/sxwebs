//-----------------------------------------------------------------------------
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
//-----------------------------------------------------------------------------

package sxhtmls_test

import (
	"testing"

	"t73f.de/r/sxwebs/sxhtmls"
	"t73f.de/r/webs/htmls"
)

func TestToSxHTML(t *testing.T) {
	var testcases = []struct {
		name string
		node *htmls.Node
		exp  string
	}{
		{"nil", nil, "()"},
		{"text", htmls.Text("abc"), "\"abc\""},
		{"br", htmls.Elem("br", nil), "(br)"},
		{"ahref",
			htmls.Elem("a", htmls.Attrs("href", "https://t73f.de"),
				htmls.Text("Detlef Stern"),
			),
			"(a ((href . \"https://t73f.de\")) \"Detlef Stern\")",
		},
		{"raw", &htmls.Node{Type: htmls.RawNode, Data: "very raw"}, "(@H \"very raw\")"},
		{"comment",
			&htmls.Node{Type: htmls.CommentNode, Data: "just a comment"},
			"(@@@ \"just a comment\")"},
		{"list",
			htmls.Elem("ol", htmls.Attrs("start", "17", "reversed"),
				htmls.Elem("li", htmls.Attrs("value", "one"), htmls.Text("1")),
				htmls.Elem("li", htmls.Attrs("value", "three"), htmls.Text("3")),
			),
			"(ol ((start . \"17\") (reversed . \"\")) (li ((value . \"one\")) \"1\") (li ((value . \"three\")) \"3\"))"},
		{"err-type",
			htmls.Elem("em", nil, &htmls.Node{Type: htmls.NodeType(255)}),
			"{[{unknown node type: 255}]}"},
		{"no-tag",
			htmls.Elem("span", nil, htmls.Elem("", nil)),
			"{[{empty symbol string}]}"},
		{"no-attr",
			htmls.Elem("span", htmls.Attrs("", "val")),
			"{[{empty symbol string}]}"},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			obj, err := sxhtmls.ToSxHTML(tc.node)
			var got string
			if err == nil {
				got = obj.String()
			} else {
				got = "{[{" + err.Error() + "}]}"
			}
			if tc.exp != got {
				t.Errorf("\nexpected: %q\nbut got : %q", tc.exp, got)
			}
		})
	}
}
