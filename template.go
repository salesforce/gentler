// Copyright (c) 2023, Salesforce, Inc.
// All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/BSD-3-Clause

package gentler

import (
	"bytes"
	"text/template"
)

type PromptTemplate[I any, O any] struct {
	tmpl   *template.Template
	parser func(string) (O, error)
}

func NewPromptTemplate[I any, O any](t string, parser func(string) (O, error)) (PromptTemplate[I, O], error) {
	pt := PromptTemplate[I, O]{}

	tt, err := template.New("").Parse(t)
	if err != nil {
		return pt, err
	}

	pt.tmpl = tt
	pt.parser = parser
	return pt, nil
}

func (t PromptTemplate[I, O]) ToPrompt(data I) (string, error) {
	buf := &bytes.Buffer{}
	if err := t.tmpl.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (t PromptTemplate[I, O]) ParseResult(resp string) (O, error) {
	return t.parser(resp)
}
