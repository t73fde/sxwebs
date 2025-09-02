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

// Package sxhtml represents HTML as s-expressions.
package sxhtml

import (
	"io"
	"sort"
	"strings"

	"t73f.de/r/sx"
	"t73f.de/r/webs/htmls/tags"
)

type attrType int

const (
	_         attrType = iota
	attrPlain          // No further escape needed
	attrURL            // Escape URL
	attrCSS            // Special CSS escaping
	attrJS             // Escape JavaScript
)

// MakeSymbol creates a symbol to be used for HTML purposes.
func MakeSymbol(s string) *sx.Symbol {
	return sx.MakeSymbol(s)
}

// Names for special symbols.
const (
	nameCDATA         = "@C"
	nameNoEscape      = "@H"
	nameListSplice    = "@L"
	nameInlineComment = "@@"
	nameBlockComment  = "@@@"
	nameDoctype       = "@@@@"
)

// Some often used symbols.
var (
	SymCDATA         = MakeSymbol(nameCDATA)
	SymNoEscape      = MakeSymbol(nameNoEscape)
	SymListSplice    = MakeSymbol(nameListSplice)
	SymInlineComment = MakeSymbol(nameInlineComment)
	SymBlockComment  = MakeSymbol(nameBlockComment)
	SymDoctype       = MakeSymbol(nameDoctype)
)

// Generator is the object that allows to generate HTML.
type Generator struct {
	withNewline bool
}

// SetNewline will add new-line characters before certain tags.
func (gen *Generator) SetNewline() *Generator { gen.withNewline = true; return gen }

// NewGenerator creates a new generator.
func NewGenerator() *Generator { return &Generator{} }

// WriteHTML emit HTML code for the s-expression to the given writer.
func (gen *Generator) WriteHTML(w io.Writer, obj sx.Object) error {
	enc := myEncoder{gen: gen, pr: printer{w: w}, lastWasTag: true}
	enc.generate(obj)
	return enc.pr.err
}

// WriteListHTML emits HTML code for a list of s-expressions to the given writer.
func (gen *Generator) WriteListHTML(w io.Writer, lst *sx.Pair) error {
	enc := myEncoder{gen: gen, pr: printer{w: w}, lastWasTag: true}
	for elem := range lst.Values() {
		enc.generate(elem)
	}
	return enc.pr.err
}

type myEncoder struct {
	gen        *Generator
	pr         printer
	lastWasTag bool
}

func (enc *myEncoder) generate(obj sx.Object) {
	switch o := obj.(type) {
	case sx.String:
		enc.pr.printHTML(o.GetValue())
		enc.lastWasTag = false
	case sx.Number:
		enc.pr.printString(o.String())
		enc.lastWasTag = false
	case *sx.Pair:
		if o.IsNil() {
			enc.lastWasTag = false
			return
		}
		if sym, isSymbol := sx.GetSymbol(o.Car()); isSymbol {
			tail := o.Tail()
			if s := sym.GetValue(); s[0] == '@' {
				switch s {
				case nameCDATA:
					enc.writeCDATA(tail)
				case nameNoEscape:
					enc.writeNoEscape(tail)
				case nameInlineComment:
					enc.writeComment(tail)
				case nameBlockComment:
					enc.writeCommentML(tail)
				case nameListSplice:
					enc.generateList(tail)
				case nameDoctype:
					enc.writeDoctype(tail)
				default:
					enc.writeTag(sym, tail)
					return
				}
				enc.lastWasTag = false
				return
			}
			enc.writeTag(sym, tail)
		}
	default:
		enc.lastWasTag = false
	}
}

func (enc *myEncoder) generateList(lst *sx.Pair) {
	for obj := range lst.Values() {
		enc.generate(obj)
	}
}

func (enc *myEncoder) writeCDATA(elems *sx.Pair) {
	enc.pr.printString("<![CDATA[")
	enc.writeNoEscape(elems)
	enc.pr.printString("]]>")
}

func (enc *myEncoder) writeNoEscape(elems *sx.Pair) {
	for obj := range elems.Values() {
		if s, isString := sx.GetString(obj); isString {
			enc.pr.printString(s.GetValue())
		}
	}
}

func (enc *myEncoder) writeComment(elems *sx.Pair) {
	enc.pr.printString("<!--")
	for obj := range elems.Values() {
		enc.pr.printString(" ")
		enc.printCommentObj(obj)
	}
	enc.pr.printString(" -->")
}
func (enc *myEncoder) writeCommentML(elems *sx.Pair) {
	enc.pr.printString("<!--")
	for obj := range elems.Values() {
		enc.pr.printString("\n")
		enc.printCommentObj(obj)
	}
	enc.pr.printString("\n-->\n")
}
func (enc *myEncoder) printCommentObj(obj sx.Object) {
	enc.pr.printComment(obj.GoString())
}

func (enc *myEncoder) writeDoctype(elems *sx.Pair) {
	// TODO: check for multiple doctypes, error on second
	enc.pr.printString("<!DOCTYPE html>\n")
	enc.generateList(elems)
}

func (enc *myEncoder) writeTag(sym *sx.Symbol, elems *sx.Pair) {
	tag := sym.GetValue()
	if isIgnorableEmptyTag(tag) && ignoreEmptyStrings(elems) == nil {
		return
	}
	withNewline := enc.gen.withNewline && isNewLineTag(tag)
	tagName := sym.String()
	if withNewline && (!enc.lastWasTag || isAlwaysNewLineTag(tag)) {
		enc.pr.printStrings("\n<", tagName)
	} else {
		enc.pr.printStrings("<", tagName)
	}
	if attrs := getAttributes(elems); attrs != nil {
		enc.writeAttributes(attrs)
		elems = elems.Tail()
	}
	enc.pr.printString(">")
	if tags.IsVoid(tag) {
		enc.lastWasTag = withNewline
		return
	}

	enc.generateList(elems)
	if withNewline {
		enc.pr.printStrings("</", tagName, ">\n")
	} else {
		enc.pr.printStrings("</", tagName, ">")
	}
	enc.lastWasTag = withNewline
}

func isIgnorableEmptyTag(tag string) bool {
	// tags that can be ignored if empty
	switch tag {
	case "div", "span", "code", "kbd", "p", "samp":
		return true
	}
	return false
}

func ignoreEmptyStrings(elem *sx.Pair) *sx.Pair {
	for node := range elem.Pairs() {
		if s, isString := sx.GetString(node.Car()); !isString || s.GetValue() != "" {
			return node
		}
	}
	return nil
}

func isNewLineTag(tag string) bool {
	switch tag {
	case nameCDATA,
		"head", "link", "meta", "title", "script", "body",
		"article", "details", "div", "header", "footer", "form",
		"main", "summary",
		"h1", "h2", "h3", "h4", "h5", "h6",
		"li", "ol", "ul", "dd", "dt", "dl",
		"table", "thead", "tbody", "tr",
		"section", "input":
		return true
	}
	return false
}
func isAlwaysNewLineTag(tag string) bool {
	switch tag {
	case "head", "link", "meta", "title", "div":
		return true
	}
	return false
}

func getAttributes(lst *sx.Pair) *sx.Pair {
	if pair, isPair := sx.GetPair(lst.Car()); isPair && pair != nil {
		if _, isAttr := sx.GetPair(pair.Car()); isAttr {
			return pair
		}
	}
	return nil
}

func (enc *myEncoder) writeAttributes(attrs *sx.Pair) {
	length := attrs.Length()
	found := make(map[string]struct{}, length)
	empty := make(map[string]struct{}, length)
	a := make(map[string]string, length)
	for val := range attrs.Values() {
		pair, isPair := sx.GetPair(val)
		if !isPair {
			continue
		}
		sym, isSymbol := sx.GetSymbol(pair.Car())
		if !isSymbol {
			continue
		}
		key := sym.String()
		if _, found := found[key]; found {
			continue
		}
		found[key] = struct{}{}
		if cdr := pair.Cdr(); !sx.IsNil(cdr) {
			var obj sx.Object
			if tail, isTail := sx.GetPair(cdr); isTail {
				obj = tail.Car()
			} else {
				obj = cdr
			}
			var s string
			switch o := obj.(type) {
			case sx.String:
				s = o.GetValue()
			case *sx.Symbol:
				s = o.GetValue()
			case sx.Number:
				s = o.GoString()
			default:
				continue
			}
			a[key] = strings.TrimSpace(s)
		} else {
			a[key] = ""
			empty[key] = struct{}{}
		}
	}

	keys := make([]string, 0, len(a))
	for key := range a {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		enc.pr.printStrings(" ", key)
		if _, isEmpty := empty[key]; !isEmpty {
			enc.pr.printString(`=`)
			enc.pr.printAttributeValue(getAttributeType(key), a[key])
		}
	}
}

func getAttributeType(key string) attrType {
	if dataName, isData := strings.CutPrefix(key, "data-"); isData {
		key = dataName
	} else if prefix, rest, hasPrefix := strings.Cut(key, ":"); hasPrefix {
		if prefix == "xmlns" {
			return attrURL
		}
		key = rest
	}

	// Attributes with URL values: https://html.spec.whatwg.org/multipage/indices.html#attributes-1
	switch key {
	case "action", "cite", "data", "formaction", "href", "itemid", "itemprop",
		"itemtype", "ping", "poster", "src":

		return attrURL
	}

	// Names that contain something similar to URL are treated as URLs
	if strings.HasSuffix(key, "uri") || strings.HasSuffix(key, "url") || strings.HasSuffix(key, "doi") {
		return attrURL
	}

	if key == "style" {
		return attrCSS
	}

	// Attribute names starting with "on" (e.g. "onload") are treated as JavaScript values.
	if strings.HasPrefix(key, "on") {
		return attrJS
	}

	return attrPlain
}
