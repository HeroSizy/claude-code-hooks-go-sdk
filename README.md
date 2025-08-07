# Claude Code Hooks Go SDK

A production-ready Go SDK for building [Claude Code hooks](https://docs.anthropic.com/en/docs/claude-code/hooks) with type-safe, idiomatic Go interfaces.

## Installation

```bash
go get github.com/anthropics/claude-code-hooks-go-sdk
```

## Quick Start

Create a simple hook that logs tool usage:

```go
package main

import (
    "log"
    "github.com/anthropics/claude-code-hooks-go-sdk/handler"
    "github.com/anthropics/claude-code-hooks-go-sdk/types"
)

type MyHandler struct{}

func (h *MyHandler) HandlePreToolUse(input types.PreToolUseInput) (types.PreToolUseOutput, error) {
    log.Printf("About to use tool: %s", input.ToolName)
    return types.PreToolUseOutput{}, nil
}

func main() {
    h := &MyHandler{}
    multiHandler := &handler.MultiHandler{
        PreToolUseHandler: h,
    }
    handler.Execute(multiHandler)
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
    if input.ToolName == "Bash" {
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

## Handler Patterns

### Multi-Handler Pattern
Handle multiple event types with one struct:

```go
type MyHandler struct{}

func (h *MyHandler) HandlePreToolUse(input types.PreToolUseInput) (types.PreToolUseOutput, error) {
    // Handle PreToolUse
    return types.PreToolUseOutput{}, nil
}

func (h *MyHandler) HandlePostToolUse(input types.PostToolUseInput) (types.PostToolUseOutput, error) {
    // Handle PostToolUse  
    return types.PostToolUseOutput{}, nil
}

func main() {
    h := &MyHandler{}
    multiHandler := &handler.MultiHandler{
        PreToolUseHandler:  h,
        PostToolUseHandler: h,
    }
    handler.Execute(multiHandler)
}
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
    case types.EventPostToolUse:
        postInput := input.(types.PostToolUseInput)
        log.Printf("PostToolUse: %s", postInput.ToolName)
        return types.PostToolUseOutput{}, nil
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
- **PreToolUse**: `ToolName`, `ToolInput`
- **PostToolUse**: `ToolName`, `ToolInput`, `ToolResponse`
- **UserPromptSubmit**: `Prompt`
- **Stop/SubagentStop**: `StopHookActive`
- **Notification**: `Message`
- **PreCompact**: `Trigger`, `CustomInstructions`
- **SessionStart**: `Source`

## Testing Hooks

Use the router's `Process` method for testing:

```go
func TestMyHook(t *testing.T) {
    h := &MyHandler{}
    multiHandler := &handler.MultiHandler{PreToolUseHandler: h}
    router := handler.NewRouter(multiHandler)
    
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