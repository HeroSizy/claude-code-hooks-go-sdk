package handler

import (
	"fmt"

	"github.com/anthropics/claude-code-hooks-go-sdk/types"
)

type ResolutionMode int

const (
	ResolutionModeBlockAny ResolutionMode = iota
	ResolutionModeFirstWin
	ResolutionModeMerge
)

type Resolver interface {
	Resolve(results []HandlerResult) (types.HookOutput, error)
}

type BlockAnyResolver struct{}

func (r *BlockAnyResolver) Resolve(results []HandlerResult) (types.HookOutput, error) {
	if len(results) == 0 {
		return types.Success(), nil
	}
	
	// Check for any blocking results first
	for _, result := range results {
		if result.Error != nil {
			return nil, result.Error
		}
		
		if result.Output != nil && result.Output.ExitWith() == types.ExitBlocking {
			return result.Output, nil
		}
	}
	
	// If no blocking results, return the last successful result
	for i := len(results) - 1; i >= 0; i-- {
		if results[i].Output != nil {
			return results[i].Output, nil
		}
	}
	
	return types.Success(), nil
}

type FirstWinResolver struct{}

func (r *FirstWinResolver) Resolve(results []HandlerResult) (types.HookOutput, error) {
	if len(results) == 0 {
		return types.Success(), nil
	}
	
	// Return first successful result
	for _, result := range results {
		if result.Error != nil {
			return nil, result.Error
		}
		
		if result.Output != nil {
			return result.Output, nil
		}
	}
	
	return types.Success(), nil
}

type MergeResolver struct{}

func (r *MergeResolver) Resolve(results []HandlerResult) (types.HookOutput, error) {
	if len(results) == 0 {
		return types.Success(), nil
	}
	
	// Check for errors first
	for _, result := range results {
		if result.Error != nil {
			return nil, result.Error
		}
	}
	
	// Check for any blocking result
	for _, result := range results {
		if result.Output != nil && result.Output.ExitWith() == types.ExitBlocking {
			return result.Output, nil
		}
	}
	
	// For merge, we'll implement type-specific merging
	// This is a simplified version - in a real implementation,
	// you'd need type-specific merge logic for each output type
	return r.mergeOutputs(results)
}

func (r *MergeResolver) mergeOutputs(results []HandlerResult) (types.HookOutput, error) {
	// Find the first non-nil output and use it as base
	// In a more sophisticated implementation, this would merge
	// compatible fields from multiple outputs
	for _, result := range results {
		if result.Output != nil {
			return result.Output, nil
		}
	}
	
	return types.Success(), nil
}

type CustomResolver struct {
	ResolveFunc func(results []HandlerResult) (types.HookOutput, error)
}

func (r *CustomResolver) Resolve(results []HandlerResult) (types.HookOutput, error) {
	if r.ResolveFunc == nil {
		return types.Success(), fmt.Errorf("custom resolver function not provided")
	}
	return r.ResolveFunc(results)
}

func GetResolver(mode ResolutionMode) Resolver {
	switch mode {
	case ResolutionModeBlockAny:
		return &BlockAnyResolver{}
	case ResolutionModeFirstWin:
		return &FirstWinResolver{}
	case ResolutionModeMerge:
		return &MergeResolver{}
	default:
		return &BlockAnyResolver{}
	}
}