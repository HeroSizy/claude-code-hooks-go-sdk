package types

type CompactTrigger string

const (
	CompactTriggerManual CompactTrigger = "manual"
	CompactTriggerAuto   CompactTrigger = "auto"
)

func (t CompactTrigger) String() string {
	return string(t)
}

func (t CompactTrigger) IsValid() bool {
	switch t {
	case CompactTriggerManual, CompactTriggerAuto:
		return true
	default:
		return false
	}
}

type SessionSource string

const (
	SessionSourceStartup SessionSource = "startup"
	SessionSourceResume  SessionSource = "resume"
	SessionSourceClear   SessionSource = "clear"
)

func (s SessionSource) String() string {
	return string(s)
}

func (s SessionSource) IsValid() bool {
	switch s {
	case SessionSourceStartup, SessionSourceResume, SessionSourceClear:
		return true
	default:
		return false
	}
}

type ToolName string

const (
	ToolTask      ToolName = "Task"
	ToolBash      ToolName = "Bash"
	ToolGlob      ToolName = "Glob"
	ToolGrep      ToolName = "Grep"
	ToolRead      ToolName = "Read"
	ToolEdit      ToolName = "Edit"
	ToolMultiEdit ToolName = "MultiEdit"
	ToolWrite     ToolName = "Write"
	ToolWebFetch  ToolName = "WebFetch"
	ToolWebSearch ToolName = "WebSearch"
)

func (t ToolName) String() string {
	return string(t)
}

func (t ToolName) IsValid() bool {
	switch t {
	case ToolTask, ToolBash, ToolGlob, ToolGrep, ToolRead, ToolEdit, ToolMultiEdit, ToolWrite, ToolWebFetch, ToolWebSearch:
		return true
	default:
		return false
	}
}