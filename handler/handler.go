package handler

import (
	"github.com/anthropics/claude-code-hooks-go-sdk/types"
)

type Handler interface {
	HandleEvent(input types.HookInput, eventName types.EventName) (types.HookOutput, error)
}

type PreToolUseHandler interface {
	HandlePreToolUse(input types.PreToolUseInput) (types.PreToolUseOutput, error)
}

type PostToolUseHandler interface {
	HandlePostToolUse(input types.PostToolUseInput) (types.PostToolUseOutput, error)
}

type NotificationHandler interface {
	HandleNotification(input types.NotificationInput) (types.NotificationOutput, error)
}

type UserPromptSubmitHandler interface {
	HandleUserPromptSubmit(input types.UserPromptSubmitInput) (types.UserPromptSubmitOutput, error)
}

type StopHandler interface {
	HandleStop(input types.StopInput) (types.StopOutput, error)
}

type SubagentStopHandler interface {
	HandleSubagentStop(input types.SubagentStopInput) (types.SubagentStopOutput, error)
}

type PreCompactHandler interface {
	HandlePreCompact(input types.PreCompactInput) (types.PreCompactOutput, error)
}

type SessionStartHandler interface {
	HandleSessionStart(input types.SessionStartInput) (types.SessionStartOutput, error)
}

type MultiHandler struct {
	PreToolUseHandler     PreToolUseHandler
	PostToolUseHandler    PostToolUseHandler
	NotificationHandler   NotificationHandler
	UserPromptSubmitHandler UserPromptSubmitHandler
	StopHandler           StopHandler
	SubagentStopHandler   SubagentStopHandler
	PreCompactHandler     PreCompactHandler
	SessionStartHandler   SessionStartHandler
}

func (m *MultiHandler) HandleEvent(input types.HookInput, eventName types.EventName) (types.HookOutput, error) {
	switch eventName {
	case types.EventPreToolUse:
		if m.PreToolUseHandler != nil {
			if preInput, ok := input.(types.PreToolUseInput); ok {
				return m.PreToolUseHandler.HandlePreToolUse(preInput)
			}
		}
	case types.EventPostToolUse:
		if m.PostToolUseHandler != nil {
			if postInput, ok := input.(types.PostToolUseInput); ok {
				return m.PostToolUseHandler.HandlePostToolUse(postInput)
			}
		}
	case types.EventNotification:
		if m.NotificationHandler != nil {
			if notifInput, ok := input.(types.NotificationInput); ok {
				return m.NotificationHandler.HandleNotification(notifInput)
			}
		}
	case types.EventUserPromptSubmit:
		if m.UserPromptSubmitHandler != nil {
			if promptInput, ok := input.(types.UserPromptSubmitInput); ok {
				return m.UserPromptSubmitHandler.HandleUserPromptSubmit(promptInput)
			}
		}
	case types.EventStop:
		if m.StopHandler != nil {
			if stopInput, ok := input.(types.StopInput); ok {
				return m.StopHandler.HandleStop(stopInput)
			}
		}
	case types.EventSubagentStop:
		if m.SubagentStopHandler != nil {
			if subStopInput, ok := input.(types.SubagentStopInput); ok {
				return m.SubagentStopHandler.HandleSubagentStop(subStopInput)
			}
		}
	case types.EventPreCompact:
		if m.PreCompactHandler != nil {
			if compactInput, ok := input.(types.PreCompactInput); ok {
				return m.PreCompactHandler.HandlePreCompact(compactInput)
			}
		}
	case types.EventSessionStart:
		if m.SessionStartHandler != nil {
			if sessionInput, ok := input.(types.SessionStartInput); ok {
				return m.SessionStartHandler.HandleSessionStart(sessionInput)
			}
		}
	}
	
	return types.Success(), nil
}

type FuncHandler func(input types.HookInput, eventName types.EventName) (types.HookOutput, error)

func (f FuncHandler) HandleEvent(input types.HookInput, eventName types.EventName) (types.HookOutput, error) {
	return f(input, eventName)
}