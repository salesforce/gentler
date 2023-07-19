// Copyright (c) 2023, Salesforce, Inc.
// All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/BSD-3-Clause

package gentler

type task[T any] struct {
	fn     func() T
	output chan T
}

type TaskPool[T any] struct {
	concurrency int
	taskQueue   chan task[T]
}

func NewTaskPool[T any](concurrency int) *TaskPool[T] {
	s := TaskPool[T]{
		concurrency: concurrency,
		taskQueue:   make(chan task[T]),
	}

	for i := 0; i < s.concurrency; i++ {
		go func() {
			for task := range s.taskQueue {
				task.output <- task.fn()
			}
		}()
	}

	return &s
}

func (s *TaskPool[T]) Run(fn func() T) T {
	output := make(chan T)
	defer close(output)
	s.taskQueue <- task[T]{
		fn:     fn,
		output: output,
	}
	return <-output
}
