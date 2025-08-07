package types

type EventName string

const (
	EventPreToolUse        EventName = "PreToolUse"
	EventPostToolUse       EventName = "PostToolUse"
	EventNotification      EventName = "Notification"
	EventUserPromptSubmit  EventName = "UserPromptSubmit"
	EventStop              EventName = "Stop"
	EventSubagentStop      EventName = "SubagentStop"
	EventPreCompact        EventName = "PreCompact"
	EventSessionStart      EventName = "SessionStart"
)

func (e EventName) String() string {
	return string(e)
}

func (e EventName) IsValid() bool {
	switch e {
	case EventPreToolUse, EventPostToolUse, EventNotification, EventUserPromptSubmit,
		EventStop, EventSubagentStop, EventPreCompact, EventSessionStart:
		return true
	default:
		return false
	}
}