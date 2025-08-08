package handler

import (
	"github.com/HeroSizy/claude-code-hooks-go-sdk/types"
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

// HandlerAdapter wraps specific handler interfaces to implement the general Handler interface
type HandlerAdapter struct {
	PreToolUse       PreToolUseHandler
	PostToolUse      PostToolUseHandler
	Notification     NotificationHandler
	UserPromptSubmit UserPromptSubmitHandler
	Stop             StopHandler
	SubagentStop     SubagentStopHandler
	PreCompact       PreCompactHandler
	SessionStart     SessionStartHandler
}

func (h *HandlerAdapter) HandleEvent(input types.HookInput, eventName types.EventName) (types.HookOutput, error) {
	switch eventName {
	case types.EventPreToolUse:
		if h.PreToolUse != nil {
			if preInput, ok := input.(types.PreToolUseInput); ok {
				return h.PreToolUse.HandlePreToolUse(preInput)
			}
		}
	case types.EventPostToolUse:
		if h.PostToolUse != nil {
			if postInput, ok := input.(types.PostToolUseInput); ok {
				return h.PostToolUse.HandlePostToolUse(postInput)
			}
		}
	case types.EventNotification:
		if h.Notification != nil {
			if notifInput, ok := input.(types.NotificationInput); ok {
				return h.Notification.HandleNotification(notifInput)
			}
		}
	case types.EventUserPromptSubmit:
		if h.UserPromptSubmit != nil {
			if promptInput, ok := input.(types.UserPromptSubmitInput); ok {
				return h.UserPromptSubmit.HandleUserPromptSubmit(promptInput)
			}
		}
	case types.EventStop:
		if h.Stop != nil {
			if stopInput, ok := input.(types.StopInput); ok {
				return h.Stop.HandleStop(stopInput)
			}
		}
	case types.EventSubagentStop:
		if h.SubagentStop != nil {
			if subStopInput, ok := input.(types.SubagentStopInput); ok {
				return h.SubagentStop.HandleSubagentStop(subStopInput)
			}
		}
	case types.EventPreCompact:
		if h.PreCompact != nil {
			if compactInput, ok := input.(types.PreCompactInput); ok {
				return h.PreCompact.HandlePreCompact(compactInput)
			}
		}
	case types.EventSessionStart:
		if h.SessionStart != nil {
			if sessionInput, ok := input.(types.SessionStartInput); ok {
				return h.SessionStart.HandleSessionStart(sessionInput)
			}
		}
	}

	return types.Success(), nil
}

// Adapter functions to wrap specific handlers as general Handler
func AdaptPreToolUse(h PreToolUseHandler) Handler {
	return &HandlerAdapter{PreToolUse: h}
}

func AdaptPostToolUse(h PostToolUseHandler) Handler {
	return &HandlerAdapter{PostToolUse: h}
}

func AdaptNotification(h NotificationHandler) Handler {
	return &HandlerAdapter{Notification: h}
}

func AdaptUserPromptSubmit(h UserPromptSubmitHandler) Handler {
	return &HandlerAdapter{UserPromptSubmit: h}
}

func AdaptStop(h StopHandler) Handler {
	return &HandlerAdapter{Stop: h}
}

func AdaptSubagentStop(h SubagentStopHandler) Handler {
	return &HandlerAdapter{SubagentStop: h}
}

func AdaptPreCompact(h PreCompactHandler) Handler {
	return &HandlerAdapter{PreCompact: h}
}

func AdaptSessionStart(h SessionStartHandler) Handler {
	return &HandlerAdapter{SessionStart: h}
}
