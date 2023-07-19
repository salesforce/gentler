// Copyright (c) 2023, Salesforce, Inc.
// All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/BSD-3-Clause

package gentler

import (
	"context"

	"github.com/salesforce/gentler/future"
)

type LLMProvider interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

type Generator[I any, O any] struct {
	llm      LLMProvider
	template PromptTemplate[I, O]
}

func NewGenerator[I any, O any](llm LLMProvider, template PromptTemplate[I, O]) Generator[I, O] {
	return Generator[I, O]{
		llm,
		template,
	}
}

func (g Generator[I, O]) Execute(ctx context.Context, data I) (O, error) {
	var output O
	prompt, err := g.template.ToPrompt(data)

	if err != nil {
		return output, err
	}

	rawOutput, err := g.llm.Generate(ctx, prompt)
	if err != nil {
		return output, err
	}

	return g.template.parser(rawOutput)

}

func (g Generator[I, O]) ExecuteAsync(ctx context.Context, data I) *future.Future[O] {
	return future.New(func() (O, error) {
		return g.Execute(ctx, data)
	})
}
