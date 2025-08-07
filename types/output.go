package types

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	ExitSuccess  = 0
	ExitBlocking = 2
)

type BaseOutput struct {
	Continue   *bool   `json:"continue,omitempty"`
	StopReason *string `json:"stopReason,omitempty"`
}

type PreToolUseOutput struct {
	BaseOutput
	AllowTool     *bool                  `json:"allowTool,omitempty"`
	ModifiedInput map[string]interface{} `json:"modifiedInput,omitempty"`
}

type PostToolUseOutput struct {
	BaseOutput
	ProcessResult *bool       `json:"processResult,omitempty"`
	Message       *string     `json:"message,omitempty"`
	Data          interface{} `json:"data,omitempty"`
}

type NotificationOutput struct {
	BaseOutput
	Acknowledged bool        `json:"acknowledged"`
	Response     *string     `json:"response,omitempty"`
	Data         interface{} `json:"data,omitempty"`
}

type UserPromptSubmitOutput struct {
	BaseOutput
	AllowSubmit    *bool   `json:"allowSubmit,omitempty"`
	ModifiedPrompt *string `json:"modifiedPrompt,omitempty"`
}

type StopOutput struct {
	BaseOutput
	AllowStop *bool   `json:"allowStop,omitempty"`
	Message   *string `json:"message,omitempty"`
}

type SubagentStopOutput struct {
	BaseOutput
	AllowStop *bool   `json:"allowStop,omitempty"`
	Message   *string `json:"message,omitempty"`
}

type PreCompactOutput struct {
	BaseOutput
	AllowCompact *bool   `json:"allowCompact,omitempty"`
	Message      *string `json:"message,omitempty"`
}

type SessionStartOutput struct {
	BaseOutput
	Message *string     `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type HookOutput interface {
	ToJSON() ([]byte, error)
	ExitWith() int
}

func (o BaseOutput) ToJSON() ([]byte, error) {
	return json.Marshal(o)
}

func (o BaseOutput) ExitWith() int {
	if o.Continue != nil && !*o.Continue {
		return ExitBlocking
	}
	return ExitSuccess
}

func (o PreToolUseOutput) ToJSON() ([]byte, error) {
	return json.Marshal(o)
}

func (o PreToolUseOutput) ExitWith() int {
	if o.Continue != nil && !*o.Continue {
		return ExitBlocking
	}
	if o.AllowTool != nil && !*o.AllowTool {
		return ExitBlocking
	}
	return ExitSuccess
}

func (o PostToolUseOutput) ToJSON() ([]byte, error) {
	return json.Marshal(o)
}

func (o PostToolUseOutput) ExitWith() int {
	if o.Continue != nil && !*o.Continue {
		return ExitBlocking
	}
	return ExitSuccess
}

func (o NotificationOutput) ToJSON() ([]byte, error) {
	return json.Marshal(o)
}

func (o NotificationOutput) ExitWith() int {
	if o.Continue != nil && !*o.Continue {
		return ExitBlocking
	}
	return ExitSuccess
}

func (o UserPromptSubmitOutput) ToJSON() ([]byte, error) {
	return json.Marshal(o)
}

func (o UserPromptSubmitOutput) ExitWith() int {
	if o.Continue != nil && !*o.Continue {
		return ExitBlocking
	}
	if o.AllowSubmit != nil && !*o.AllowSubmit {
		return ExitBlocking
	}
	return ExitSuccess
}

func (o StopOutput) ToJSON() ([]byte, error) {
	return json.Marshal(o)
}

func (o StopOutput) ExitWith() int {
	if o.Continue != nil && !*o.Continue {
		return ExitBlocking
	}
	if o.AllowStop != nil && !*o.AllowStop {
		return ExitBlocking
	}
	return ExitSuccess
}

func (o SubagentStopOutput) ToJSON() ([]byte, error) {
	return json.Marshal(o)
}

func (o SubagentStopOutput) ExitWith() int {
	if o.Continue != nil && !*o.Continue {
		return ExitBlocking
	}
	if o.AllowStop != nil && !*o.AllowStop {
		return ExitBlocking
	}
	return ExitSuccess
}

func (o PreCompactOutput) ToJSON() ([]byte, error) {
	return json.Marshal(o)
}

func (o PreCompactOutput) ExitWith() int {
	if o.Continue != nil && !*o.Continue {
		return ExitBlocking
	}
	if o.AllowCompact != nil && !*o.AllowCompact {
		return ExitBlocking
	}
	return ExitSuccess
}

func (o SessionStartOutput) ToJSON() ([]byte, error) {
	return json.Marshal(o)
}

func (o SessionStartOutput) ExitWith() int {
	if o.Continue != nil && !*o.Continue {
		return ExitBlocking
	}
	return ExitSuccess
}

func OutputAndExit(output HookOutput) {
	if jsonData, err := output.ToJSON(); err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling output: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Print(string(jsonData))
	}
	os.Exit(output.ExitWith())
}

func Success() HookOutput {
	return BaseOutput{}
}

func Block(reason string) HookOutput {
	continueVal := false
	return BaseOutput{
		Continue:   &continueVal,
		StopReason: &reason,
	}
}