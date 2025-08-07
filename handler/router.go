package handler

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/anthropics/claude-code-hooks-go-sdk/types"
)

type RouterConfig struct {
	handlers       map[types.EventName][]Handler
	executionMode  ExecutionMode
	resolutionMode ResolutionMode
	timeout        time.Duration
}

type Router struct {
	config RouterConfig
}

func NewRouter() *Router {
	return &Router{
		config: RouterConfig{
			handlers:       make(map[types.EventName][]Handler),
			executionMode:  ExecutionModeSync,
			resolutionMode: ResolutionModeBlockAny,
			timeout:        30 * time.Second,
		},
	}
}

func (r *Router) OnPreToolUse(handlers ...PreToolUseHandler) *Router {
	eventHandlers := make([]Handler, len(handlers))
	for i, h := range handlers {
		eventHandlers[i] = AdaptPreToolUse(h)
	}
	r.config.handlers[types.EventPreToolUse] = eventHandlers
	return r
}

func (r *Router) OnPostToolUse(handlers ...PostToolUseHandler) *Router {
	eventHandlers := make([]Handler, len(handlers))
	for i, h := range handlers {
		eventHandlers[i] = AdaptPostToolUse(h)
	}
	r.config.handlers[types.EventPostToolUse] = eventHandlers
	return r
}

func (r *Router) OnNotification(handlers ...NotificationHandler) *Router {
	eventHandlers := make([]Handler, len(handlers))
	for i, h := range handlers {
		eventHandlers[i] = AdaptNotification(h)
	}
	r.config.handlers[types.EventNotification] = eventHandlers
	return r
}

func (r *Router) OnUserPromptSubmit(handlers ...UserPromptSubmitHandler) *Router {
	eventHandlers := make([]Handler, len(handlers))
	for i, h := range handlers {
		eventHandlers[i] = AdaptUserPromptSubmit(h)
	}
	r.config.handlers[types.EventUserPromptSubmit] = eventHandlers
	return r
}

func (r *Router) OnStop(handlers ...StopHandler) *Router {
	eventHandlers := make([]Handler, len(handlers))
	for i, h := range handlers {
		eventHandlers[i] = AdaptStop(h)
	}
	r.config.handlers[types.EventStop] = eventHandlers
	return r
}

func (r *Router) OnSubagentStop(handlers ...SubagentStopHandler) *Router {
	eventHandlers := make([]Handler, len(handlers))
	for i, h := range handlers {
		eventHandlers[i] = AdaptSubagentStop(h)
	}
	r.config.handlers[types.EventSubagentStop] = eventHandlers
	return r
}

func (r *Router) OnPreCompact(handlers ...PreCompactHandler) *Router {
	eventHandlers := make([]Handler, len(handlers))
	for i, h := range handlers {
		eventHandlers[i] = AdaptPreCompact(h)
	}
	r.config.handlers[types.EventPreCompact] = eventHandlers
	return r
}

func (r *Router) OnSessionStart(handlers ...SessionStartHandler) *Router {
	eventHandlers := make([]Handler, len(handlers))
	for i, h := range handlers {
		eventHandlers[i] = AdaptSessionStart(h)
	}
	r.config.handlers[types.EventSessionStart] = eventHandlers
	return r
}

func (r *Router) WithExecution(mode ExecutionMode) *Router {
	r.config.executionMode = mode
	return r
}

func (r *Router) WithResolution(mode ResolutionMode) *Router {
	r.config.resolutionMode = mode
	return r
}

func (r *Router) WithTimeout(timeout time.Duration) *Router {
	r.config.timeout = timeout
	return r
}

func (r *Router) Run() error {
	return r.RunWithReader(os.Stdin)
}

func (r *Router) RunWithReader(reader io.Reader) error {
	input, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	output, err := r.Process(input)
	if err != nil {
		return err
	}

	types.OutputAndExit(output)
	return nil
}

func (r *Router) Process(input []byte) (types.HookOutput, error) {
	hookInput, eventName, err := types.ParseInput(input)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}

	return r.HandleEvent(hookInput, eventName)
}

func (r *Router) HandleEvent(input types.HookInput, eventName types.EventName) (types.HookOutput, error) {
	handlers, exists := r.config.handlers[eventName]
	if !exists || len(handlers) == 0 {
		return types.Success(), nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.config.timeout)
	defer cancel()

	executor := GetExecutor(r.config.executionMode)
	results, err := executor.Execute(ctx, input, eventName, handlers)
	if err != nil {
		return nil, err
	}

	resolver := GetResolver(r.config.resolutionMode)
	return resolver.Resolve(results)
}

// Convenience functions for simple use cases
func Execute(router *Router) {
	if err := router.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Hook execution failed: %v\n", err)
		os.Exit(1)
	}
}

// Simple single handler execution (backward compatibility helper)
func ExecuteSingle(handler Handler) {
	router := NewRouter()
	
	// Register the single handler for all events it supports
	if h, ok := handler.(PreToolUseHandler); ok {
		router.OnPreToolUse(h)
	}
	if h, ok := handler.(PostToolUseHandler); ok {
		router.OnPostToolUse(h)
	}
	if h, ok := handler.(NotificationHandler); ok {
		router.OnNotification(h)
	}
	if h, ok := handler.(UserPromptSubmitHandler); ok {
		router.OnUserPromptSubmit(h)
	}
	if h, ok := handler.(StopHandler); ok {
		router.OnStop(h)
	}
	if h, ok := handler.(SubagentStopHandler); ok {
		router.OnSubagentStop(h)
	}
	if h, ok := handler.(PreCompactHandler); ok {
		router.OnPreCompact(h)
	}
	if h, ok := handler.(SessionStartHandler); ok {
		router.OnSessionStart(h)
	}
	
	Execute(router)
}

// Function handler support
type FuncHandler func(input types.HookInput, eventName types.EventName) (types.HookOutput, error)

func (f FuncHandler) HandleEvent(input types.HookInput, eventName types.EventName) (types.HookOutput, error) {
	return f(input, eventName)
}

func ExecuteFunc(handlerFunc FuncHandler) {
	ExecuteSingle(handlerFunc)
}