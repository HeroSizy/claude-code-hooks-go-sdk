package handler

import (
	"context"
	"sync"

	"github.com/HeroSizy/claude-code-hooks-go-sdk/types"
)

type ExecutionMode int

const (
	ExecutionModeSync ExecutionMode = iota
	ExecutionModeAsync
	ExecutionModePipeline
)

type HandlerResult struct {
	Output types.HookOutput
	Error  error
	Index  int
}

type Executor interface {
	Execute(ctx context.Context, input types.HookInput, eventName types.EventName, handlers []Handler) ([]HandlerResult, error)
}

type SyncExecutor struct{}

func (e *SyncExecutor) Execute(ctx context.Context, input types.HookInput, eventName types.EventName, handlers []Handler) ([]HandlerResult, error) {
	results := make([]HandlerResult, 0, len(handlers))

	for i, handler := range handlers {
		select {
		case <-ctx.Done():
			return results, ctx.Err()
		default:
		}

		output, err := handler.HandleEvent(input, eventName)
		result := HandlerResult{
			Output: output,
			Error:  err,
			Index:  i,
		}
		results = append(results, result)

		// Stop on first error in sync mode
		if err != nil {
			return results, err
		}

		// Check if this handler blocked the operation
		if output != nil && isBlocking(output) {
			return results, nil
		}
	}

	return results, nil
}

type AsyncExecutor struct{}

func (e *AsyncExecutor) Execute(ctx context.Context, input types.HookInput, eventName types.EventName, handlers []Handler) ([]HandlerResult, error) {
	if len(handlers) == 0 {
		return nil, nil
	}

	results := make([]HandlerResult, len(handlers))
	var wg sync.WaitGroup

	for i, handler := range handlers {
		wg.Add(1)
		go func(index int, h Handler) {
			defer wg.Done()

			output, err := h.HandleEvent(input, eventName)
			results[index] = HandlerResult{
				Output: output,
				Error:  err,
				Index:  index,
			}
		}(i, handler)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return results, ctx.Err()
	case <-done:
		return results, nil
	}
}

type PipelineExecutor struct{}

func (e *PipelineExecutor) Execute(ctx context.Context, input types.HookInput, eventName types.EventName, handlers []Handler) ([]HandlerResult, error) {
	if len(handlers) == 0 {
		return nil, nil
	}

	results := make([]HandlerResult, 0, len(handlers))
	currentInput := input

	for i, handler := range handlers {
		select {
		case <-ctx.Done():
			return results, ctx.Err()
		default:
		}

		output, err := handler.HandleEvent(currentInput, eventName)
		result := HandlerResult{
			Output: output,
			Error:  err,
			Index:  i,
		}
		results = append(results, result)

		if err != nil {
			return results, err
		}

		// In pipeline mode, we don't chain outputs to inputs since
		// each event type has specific input/output types
		// Pipeline is mainly useful for middleware-style processing
		if output != nil && isBlocking(output) {
			return results, nil
		}
	}

	return results, nil
}

func isBlocking(output types.HookOutput) bool {
	return output.ExitWith() == types.ExitBlocking
}

func GetExecutor(mode ExecutionMode) Executor {
	switch mode {
	case ExecutionModeSync:
		return &SyncExecutor{}
	case ExecutionModeAsync:
		return &AsyncExecutor{}
	case ExecutionModePipeline:
		return &PipelineExecutor{}
	default:
		return &SyncExecutor{}
	}
}
