// Copyright (c) 2023, Salesforce, Inc.
// All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/BSD-3-Clause

package gentler

import (
	"bytes"
	"text/template"
)

type PromptIO[I any, O any] struct {
	tmpl   *template.Template
	parser func(string) (O, error)
}

func NewPromptIO[I any, O any](t string, parser func(string) (O, error)) (PromptIO[I, O], error) {
	pt := PromptIO[I, O]{}

	tt, err := template.New("").Parse(t)
	if err != nil {
		return pt, err
	}

	pt.tmpl = tt
	pt.parser = parser
	return pt, nil
}

func (t PromptIO[I, O]) ToPrompt(data I) (string, error) {
	buf := &bytes.Buffer{}
	if err := t.tmpl.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (t PromptIO[I, O]) ParseResult(resp string) (O, error) {
	return t.parser(resp)
}
