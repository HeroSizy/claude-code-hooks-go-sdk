package main

import (
	"log"
	"strings"

	"github.com/anthropics/claude-code-hooks-go-sdk/handler"
	"github.com/anthropics/claude-code-hooks-go-sdk/types"
)

type ExampleHandler struct{}

func (h *ExampleHandler) HandlePreToolUse(input types.PreToolUseInput) (types.PreToolUseOutput, error) {
	log.Printf("PreToolUse: Tool=%s", input.ToolName)
	
	if input.ToolName == types.ToolBash {
		log.Printf("Bash command detected in session %s", input.SessionID)
	}
	
	return types.PreToolUseOutput{}, nil
}

func (h *ExampleHandler) HandlePostToolUse(input types.PostToolUseInput) (types.PostToolUseOutput, error) {
	log.Printf("PostToolUse: Tool=%s", input.ToolName)
	log.Printf("Tool response: %v", input.ToolResponse)
	
	return types.PostToolUseOutput{}, nil
}

func (h *ExampleHandler) HandleUserPromptSubmit(input types.UserPromptSubmitInput) (types.UserPromptSubmitOutput, error) {
	log.Printf("UserPromptSubmit: %s", input.Prompt)
	
	if strings.Contains(strings.ToLower(input.Prompt), "dangerous") {
		continueVal := false
		allowSubmit := false
		return types.UserPromptSubmitOutput{
			BaseOutput: types.BaseOutput{
				Continue:   &continueVal,
				StopReason: stringPtr("Prompt contains potentially dangerous content"),
			},
			AllowSubmit: &allowSubmit,
		}, nil
	}
	
	return types.UserPromptSubmitOutput{}, nil
}

func (h *ExampleHandler) HandleStop(input types.StopInput) (types.StopOutput, error) {
	log.Printf("Stop event: StopHookActive=%t", input.StopHookActive)
	return types.StopOutput{}, nil
}

func main() {
	h := &ExampleHandler{}
	
	// New multi-handler API - single handler for multiple events
	router := handler.NewRouter().
		OnPreToolUse(h).
		OnPostToolUse(h).
		OnUserPromptSubmit(h).
		OnStop(h)
	
	handler.Execute(router)
}

func stringPtr(s string) *string {
	return &s
}