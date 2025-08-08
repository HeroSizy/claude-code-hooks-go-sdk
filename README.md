# Claude Code Hooks Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/HeroSizy/claude-code-hooks-go-sdk.svg)](https://pkg.go.dev/github.com/HeroSizy/claude-code-hooks-go-sdk)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/HeroSizy/claude-code-hooks-go-sdk)](https://goreportcard.com/report/github.com/HeroSizy/claude-code-hooks-go-sdk)
[![CI](https://github.com/HeroSizy/claude-code-hooks-go-sdk/actions/workflows/ci.yaml/badge.svg)](https://github.com/HeroSizy/claude-code-hooks-go-sdk/actions/workflows/ci.yaml)
[![Release](https://img.shields.io/github/v/release/HeroSizy/claude-code-hooks-go-sdk?sort=semver&color=blue)](https://github.com/HeroSizy/claude-code-hooks-go-sdk/releases/latest)

A production-ready Go SDK for building [Claude Code hooks](https://docs.anthropic.com/en/docs/claude-code/hooks) with type-safe, idiomatic Go interfaces.

## Installation

```bash
go get github.com/HeroSizy/claude-code-hooks-go-sdk
```

## Quick Start

Create a simple hook that logs tool usage:

```go
package main

import (
    "log"
    "github.com/HeroSizy/claude-code-hooks-go-sdk/handler"
    "github.com/HeroSizy/claude-code-hooks-go-sdk/types"
)

type MyHandler struct{}

func (h *MyHandler) HandlePreToolUse(input types.PreToolUseInput) (types.PreToolUseOutput, error) {
    log.Printf("About to use tool: %s", input.ToolName)
    return types.PreToolUseOutput{}, nil
}

func main() {
    h := &MyHandler{}
    
    router := handler.NewRouter().
        OnPreToolUse(h)
    
    handler.Execute(router)
}
```

Build and configure in your Claude Code settings:

```bash
go build -o my-hook main.go

# In your Claude Code settings.json:
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "*",
        "hooks": [
          {
            "type": "command",
            "command": "/path/to/my-hook"
          }
        ]
      }
    ]
  }
}
```

## Hook Events

The SDK supports all 8 Claude Code hook events:

### PreToolUse
Triggered before Claude Code executes a tool.

```go
func (h *MyHandler) HandlePreToolUse(input types.PreToolUseInput) (types.PreToolUseOutput, error) {
    // Block dangerous commands
    if input.ToolName == types.ToolBash {
        if cmd, ok := input.ToolInput["command"].(string); ok {
            if strings.Contains(cmd, "rm -rf") {
                allowTool := false
                return types.PreToolUseOutput{
                    BaseOutput: types.BaseOutput{
                        Continue:   &[]bool{false}[0],
                        StopReason: &[]string{"Dangerous command blocked"}[0],
                    },
                    AllowTool: &allowTool,
                }, nil
            }
        }
    }
    return types.PreToolUseOutput{}, nil
}
```

### PostToolUse
Triggered after Claude Code executes a tool.

```go
func (h *MyHandler) HandlePostToolUse(input types.PostToolUseInput) (types.PostToolUseOutput, error) {
    if !input.Success {
        log.Printf("Tool %s failed: %v", input.ToolName, input.Error)
    }
    return types.PostToolUseOutput{}, nil
}
```

### UserPromptSubmit
Triggered when user submits a prompt.

```go
func (h *MyHandler) HandleUserPromptSubmit(input types.UserPromptSubmitInput) (types.UserPromptSubmitOutput, error) {
    // Block prompts containing sensitive keywords
    if containsSensitive(input.Prompt) {
        allowSubmit := false
        return types.UserPromptSubmitOutput{
            BaseOutput: types.BaseOutput{
                Continue:   &[]bool{false}[0],
                StopReason: &[]string{"Sensitive content detected"}[0],
            },
            AllowSubmit: &allowSubmit,
        }, nil
    }
    return types.UserPromptSubmitOutput{}, nil
}
```

### Other Events
- **Notification**: System notifications from Claude Code
- **Stop**: Claude Code session stop events  
- **SubagentStop**: Subagent termination events
- **PreCompact**: Before transcript compaction
- **SessionStart**: New session initialization

## Multi-Handler Patterns

### Single Handler, Multiple Events
Handle multiple event types with one struct:

```go
type MyHandler struct{}

func (h *MyHandler) HandlePreToolUse(input types.PreToolUseInput) (types.PreToolUseOutput, error) {
    log.Printf("PreToolUse: %s", input.ToolName)
    return types.PreToolUseOutput{}, nil
}

func (h *MyHandler) HandlePostToolUse(input types.PostToolUseInput) (types.PostToolUseOutput, error) {
    log.Printf("PostToolUse: %s", input.ToolName)
    return types.PostToolUseOutput{}, nil
}

func main() {
    h := &MyHandler{}
    
    router := handler.NewRouter().
        OnPreToolUse(h).
        OnPostToolUse(h)
    
    handler.Execute(router)
}
```

### Multiple Handlers, Single Event
Run multiple handlers for the same event:

```go
type SecurityHandler struct{}
func (h *SecurityHandler) HandlePreToolUse(input types.PreToolUseInput) (types.PreToolUseOutput, error) {
    // Security checks
    return types.PreToolUseOutput{}, nil
}

type AuditHandler struct{}  
func (h *AuditHandler) HandlePreToolUse(input types.PreToolUseInput) (types.PreToolUseOutput, error) {
    // Audit logging
    return types.PreToolUseOutput{}, nil
}

func main() {
    router := handler.NewRouter().
        OnPreToolUse(
            &SecurityHandler{},    // Runs first (sync mode)
            &AuditHandler{},       // Runs second
        ).
        WithExecution(handler.ExecutionModeSync).
        WithResolution(handler.ResolutionModeBlockAny)
    
    handler.Execute(router)
}
```

### Execution Modes

#### Synchronous Execution (Default)
Handlers run sequentially, stopping on first block/error:

```go
router := handler.NewRouter().
    OnPreToolUse(handler1, handler2, handler3).
    WithExecution(handler.ExecutionModeSync)     // Default
```

#### Asynchronous Execution  
Handlers run concurrently, results combined at the end:

```go
router := handler.NewRouter().
    OnPreToolUse(handler1, handler2, handler3).
    WithExecution(handler.ExecutionModeAsync).
    WithTimeout(10 * time.Second)
```

#### Pipeline Execution
Handlers run sequentially (useful for middleware-style processing):

```go
router := handler.NewRouter().
    OnPreToolUse(
        inputValidator,    // Validates input
        policyChecker,     // Checks policies  
        auditLogger,       // Logs decision
    ).
    WithExecution(handler.ExecutionModePipeline)
```

### Result Resolution Strategies

#### BlockAny (Default)
Block operation if any handler blocks:

```go
router.WithResolution(handler.ResolutionModeBlockAny)  // Default
```

#### FirstWin  
Use result from first successful handler:

```go
router.WithResolution(handler.ResolutionModeFirstWin)
```

#### Merge
Combine compatible results from all handlers:

```go
router.WithResolution(handler.ResolutionModeMerge)
```

### Function Handler Pattern
Use a single function for all events:

```go
func myHookHandler(input types.HookInput, eventName types.EventName) (types.HookOutput, error) {
    switch eventName {
    case types.EventPreToolUse:
        preInput := input.(types.PreToolUseInput)
        log.Printf("PreToolUse: %s", preInput.ToolName)
        return types.PreToolUseOutput{}, nil
    default:
        return types.Success(), nil
    }
}

func main() {
    handler.ExecuteFunc(myHookHandler)
}
```

## Output Control

### Allow/Block Operations
Control whether operations should proceed:

```go
// Block the operation
allowTool := false
continueVal := false
return types.PreToolUseOutput{
    BaseOutput: types.BaseOutput{
        Continue:   &continueVal,
        StopReason: &[]string{"Operation blocked"}[0],
    },
    AllowTool: &allowTool,
}, nil

// Allow the operation (default)
return types.PreToolUseOutput{}, nil
```

### Convenience Functions
```go
// Simple success
return types.Success(), nil

// Block with reason  
return types.Block("Operation not allowed"), nil
```

## Input Types

All hook inputs contain common fields:
- `SessionID`: Unique session identifier
- `TranscriptPath`: Path to session transcript  
- `CWD`: Current working directory
- `HookEventName`: The event type name

Event-specific fields:
- **PreToolUse**: `ToolName` (ToolName enum), `ToolInput`
- **PostToolUse**: `ToolName` (ToolName enum), `ToolInput`, `ToolResponse`
- **UserPromptSubmit**: `Prompt`
- **Stop/SubagentStop**: `StopHookActive`
- **Notification**: `Message`
- **PreCompact**: `Trigger` (CompactTrigger enum), `CustomInstructions`
- **SessionStart**: `Source` (SessionSource enum)

## Enums

The SDK provides typed enums for predefined values:

### Tool Names
```go
types.ToolTask, types.ToolBash, types.ToolGlob, types.ToolGrep, 
types.ToolRead, types.ToolEdit, types.ToolMultiEdit, types.ToolWrite, 
types.ToolWebFetch, types.ToolWebSearch
```

### Compact Triggers
```go
types.CompactTriggerManual, types.CompactTriggerAuto
```

### Session Sources
```go
types.SessionSourceStartup, types.SessionSourceResume, types.SessionSourceClear
```

## Testing Hooks

Use the router's `Process` method for testing:

```go
func TestMyHook(t *testing.T) {
    h := &MyHandler{}
    router := handler.NewRouter().OnPreToolUse(h)
    
    input := `{
        "session_id": "test-session",
        "transcript_path": "/tmp/transcript",
        "cwd": "/tmp",
        "hook_event_name": "PreToolUse",
        "tool_name": "Bash",
        "tool_input": {"command": "ls"}
    }`
    
    output, err := router.Process([]byte(input))
    assert.NoError(t, err)
    assert.NotNil(t, output)
}

func TestMultipleHandlers(t *testing.T) {
    securityHandler := &SecurityHandler{}
    auditHandler := &AuditHandler{}
    
    router := handler.NewRouter().
        OnPreToolUse(securityHandler, auditHandler).
        WithExecution(handler.ExecutionModeAsync).
        WithTimeout(5 * time.Second)
    
    // Test that both handlers are called
    output, err := router.Process([]byte(input))
    assert.NoError(t, err)
    assert.Equal(t, types.ExitSuccess, output.ExitWith())
}
```

## Building Production Hooks

1. **Build static binaries**:
```bash
CGO_ENABLED=0 go build -ldflags="-w -s" -o my-hook main.go
```

2. **Handle errors gracefully**:
```go
func (h *MyHandler) HandlePreToolUse(input types.PreToolUseInput) (types.PreToolUseOutput, error) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Hook panic recovered: %v", r)
        }
    }()
    
    // Your hook logic
    return types.PreToolUseOutput{}, nil
}
```

3. **Use structured logging**:
```go
import "log/slog"

func (h *MyHandler) HandlePreToolUse(input types.PreToolUseInput) (types.PreToolUseOutput, error) {
    slog.Info("PreToolUse hook triggered", 
        "tool", input.ToolName,
        "session", input.SessionID,
    )
    return types.PreToolUseOutput{}, nil
}
```

## Exit Codes

The SDK automatically handles exit codes:
- `0`: Success (default)
- `2`: Blocking error (when `Continue: false` or operation blocked)

## License

This SDK follows the same license as Claude Code. See [Claude Code documentation](https://docs.anthropic.com/en/docs/claude-code) for details.

## Contributing

This is an unofficial SDK. For issues with Claude Code itself, please refer to the [official documentation](https://docs.anthropic.com/en/docs/claude-code/hooks).