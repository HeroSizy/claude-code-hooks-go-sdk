package handler

import (
	"fmt"
	"io"
	"os"

	"github.com/anthropics/claude-code-hooks-go-sdk/types"
)

type Router struct {
	handler Handler
}

func NewRouter(handler Handler) *Router {
	return &Router{handler: handler}
}

func (r *Router) Run() error {
	return r.RunWithReader(os.Stdin)
}

func (r *Router) RunWithReader(reader io.Reader) error {
	input, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	hookInput, eventName, err := types.ParseInput(input)
	if err != nil {
		return fmt.Errorf("failed to parse input: %w", err)
	}

	output, err := r.handler.HandleEvent(hookInput, eventName)
	if err != nil {
		return fmt.Errorf("handler error: %w", err)
	}

	types.OutputAndExit(output)
	return nil
}

func (r *Router) Process(input []byte) (types.HookOutput, error) {
	hookInput, eventName, err := types.ParseInput(input)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}

	return r.handler.HandleEvent(hookInput, eventName)
}

func Execute(handler Handler) {
	router := NewRouter(handler)
	if err := router.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Hook execution failed: %v\n", err)
		os.Exit(1)
	}
}

func ExecuteFunc(handlerFunc FuncHandler) {
	Execute(handlerFunc)
}