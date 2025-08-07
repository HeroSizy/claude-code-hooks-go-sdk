package types

import (
	"encoding/json"
)

type BaseInput struct {
	SessionID      string `json:"session_id"`
	TranscriptPath string `json:"transcript_path"`
	CWD            string `json:"cwd"`
	HookEventName  string `json:"hook_event_name"`
}

type PreToolUseInput struct {
	BaseInput
	ToolName  ToolName               `json:"tool_name"`
	ToolInput map[string]interface{} `json:"tool_input"`
}

type PostToolUseInput struct {
	BaseInput
	ToolName     ToolName               `json:"tool_name"`
	ToolInput    map[string]interface{} `json:"tool_input"`
	ToolResponse interface{}            `json:"tool_response"`
}

type NotificationInput struct {
	BaseInput
	Message string `json:"message"`
}

type UserPromptSubmitInput struct {
	BaseInput
	Prompt string `json:"prompt"`
}

type StopInput struct {
	BaseInput
	StopHookActive bool `json:"stop_hook_active"`
}

type SubagentStopInput struct {
	BaseInput
	StopHookActive bool `json:"stop_hook_active"`
}

type PreCompactInput struct {
	BaseInput
	Trigger            CompactTrigger `json:"trigger"`
	CustomInstructions string         `json:"custom_instructions"`
}

type SessionStartInput struct {
	BaseInput
	Source SessionSource `json:"source"`
}

type HookInput interface {
	GetSessionID() string
	GetTranscriptPath() string
	GetCWD() string
	GetEventName() string
}

func (b BaseInput) GetSessionID() string {
	return b.SessionID
}

func (b BaseInput) GetTranscriptPath() string {
	return b.TranscriptPath
}

func (b BaseInput) GetCWD() string {
	return b.CWD
}

func (b BaseInput) GetEventName() string {
	return b.HookEventName
}

func ParseInput(data []byte) (HookInput, EventName, error) {
	var base BaseInput
	if err := json.Unmarshal(data, &base); err != nil {
		return nil, "", err
	}

	eventName := EventName(base.HookEventName)
	if !eventName.IsValid() {
		return nil, eventName, &InvalidEventError{EventName: base.HookEventName}
	}

	switch eventName {
	case EventPreToolUse:
		var input PreToolUseInput
		err := json.Unmarshal(data, &input)
		return input, eventName, err
	case EventPostToolUse:
		var input PostToolUseInput
		err := json.Unmarshal(data, &input)
		return input, eventName, err
	case EventNotification:
		var input NotificationInput
		err := json.Unmarshal(data, &input)
		return input, eventName, err
	case EventUserPromptSubmit:
		var input UserPromptSubmitInput
		err := json.Unmarshal(data, &input)
		return input, eventName, err
	case EventStop:
		var input StopInput
		err := json.Unmarshal(data, &input)
		return input, eventName, err
	case EventSubagentStop:
		var input SubagentStopInput
		err := json.Unmarshal(data, &input)
		return input, eventName, err
	case EventPreCompact:
		var input PreCompactInput
		err := json.Unmarshal(data, &input)
		return input, eventName, err
	case EventSessionStart:
		var input SessionStartInput
		err := json.Unmarshal(data, &input)
		return input, eventName, err
	default:
		return nil, eventName, &InvalidEventError{EventName: base.HookEventName}
	}
}

type InvalidEventError struct {
	EventName string
}

func (e *InvalidEventError) Error() string {
	return "invalid event name: " + e.EventName
}