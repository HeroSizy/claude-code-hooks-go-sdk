package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/anthropics/claude-code-hooks-go-sdk/handler"
	"github.com/anthropics/claude-code-hooks-go-sdk/types"
)

type SecurityHandler struct{}

func (h *SecurityHandler) HandlePreToolUse(input types.PreToolUseInput) (types.PreToolUseOutput, error) {
	log.Printf("[SECURITY] Checking tool: %s", input.ToolName)
	
	// Block dangerous bash commands
	if input.ToolName == types.ToolBash {
		if cmd, ok := input.ToolInput["command"].(string); ok {
			dangerous := []string{"rm -rf", "sudo", "chmod 777", "format", "del"}
			for _, danger := range dangerous {
				if strings.Contains(strings.ToLower(cmd), danger) {
					continueVal := false
					allowTool := false
					return types.PreToolUseOutput{
						BaseOutput: types.BaseOutput{
							Continue:   &continueVal,
							StopReason: stringPtr(fmt.Sprintf("Dangerous command blocked: %s", danger)),
						},
						AllowTool: &allowTool,
					}, nil
				}
			}
		}
	}
	
	log.Printf("[SECURITY] Tool %s approved", input.ToolName)
	return types.PreToolUseOutput{}, nil
}

type AuditHandler struct{}

func (h *AuditHandler) HandlePreToolUse(input types.PreToolUseInput) (types.PreToolUseOutput, error) {
	log.Printf("[AUDIT] Tool execution: %s in session %s", input.ToolName, input.SessionID)
	return types.PreToolUseOutput{}, nil
}

func (h *AuditHandler) HandlePostToolUse(input types.PostToolUseInput) (types.PostToolUseOutput, error) {
	log.Printf("[AUDIT] Tool completed: %s", input.ToolName)
	return types.PostToolUseOutput{}, nil
}

type MetricsHandler struct{}

func (h *MetricsHandler) HandlePreToolUse(input types.PreToolUseInput) (types.PreToolUseOutput, error) {
	log.Printf("[METRICS] Incrementing counter for tool: %s", input.ToolName)
	return types.PreToolUseOutput{}, nil
}

func (h *MetricsHandler) HandlePostToolUse(input types.PostToolUseInput) (types.PostToolUseOutput, error) {
	log.Printf("[METRICS] Recording execution time for tool: %s", input.ToolName)
	return types.PostToolUseOutput{}, nil
}

type ContentFilterHandler struct{}

func (h *ContentFilterHandler) HandleUserPromptSubmit(input types.UserPromptSubmitInput) (types.UserPromptSubmitOutput, error) {
	log.Printf("[FILTER] Checking prompt content")
	
	// Block prompts with sensitive content
	blocked := []string{"password", "secret", "private key", "ssn", "credit card"}
	prompt := strings.ToLower(input.Prompt)
	
	for _, word := range blocked {
		if strings.Contains(prompt, word) {
			continueVal := false
			allowSubmit := false
			return types.UserPromptSubmitOutput{
				BaseOutput: types.BaseOutput{
					Continue:   &continueVal,
					StopReason: stringPtr(fmt.Sprintf("Blocked sensitive content: %s", word)),
				},
				AllowSubmit: &allowSubmit,
			}, nil
		}
	}
	
	log.Printf("[FILTER] Prompt approved")
	return types.UserPromptSubmitOutput{}, nil
}

func demonstrateSyncExecution() {
	log.Println("=== SYNC EXECUTION DEMO ===")
	
	router := handler.NewRouter().
		OnPreToolUse(
			&SecurityHandler{},    // Runs first
			&AuditHandler{},       // Runs second  
			&MetricsHandler{},     // Runs third
		).
		WithExecution(handler.ExecutionModeSync).
		WithResolution(handler.ResolutionModeBlockAny)
	
	log.Println("Configured sync execution with BlockAny resolution")
	handler.Execute(router)
}

func demonstrateAsyncExecution() {
	log.Println("=== ASYNC EXECUTION DEMO ===")
	
	router := handler.NewRouter().
		OnPreToolUse(
			&SecurityHandler{},    // Runs concurrently
			&AuditHandler{},       // Runs concurrently
			&MetricsHandler{},     // Runs concurrently
		).
		OnPostToolUse(
			&AuditHandler{},
			&MetricsHandler{},
		).
		WithExecution(handler.ExecutionModeAsync).
		WithResolution(handler.ResolutionModeBlockAny).
		WithTimeout(10 * time.Second)
	
	log.Println("Configured async execution with 10s timeout")
	handler.Execute(router)
}

func demonstrateMixedHandlers() {
	log.Println("=== MIXED HANDLERS DEMO ===")
	
	router := handler.NewRouter().
		OnPreToolUse(
			&SecurityHandler{},
			&AuditHandler{},
		).
		OnUserPromptSubmit(
			&ContentFilterHandler{},
		).
		OnPostToolUse(
			&MetricsHandler{},
		).
		WithExecution(handler.ExecutionModeSync).
		WithResolution(handler.ResolutionModeFirstWin)
	
	log.Println("Configured multiple event types with FirstWin resolution")
	handler.Execute(router)
}

func main() {
	// Choose demo based on environment variable or default to sync
	demo := "sync" // Could read from os.Getenv("DEMO_MODE")
	
	switch demo {
	case "async":
		demonstrateAsyncExecution()
	case "mixed":
		demonstrateMixedHandlers()
	default:
		demonstrateSyncExecution()
	}
}

func stringPtr(s string) *string {
	return &s
}