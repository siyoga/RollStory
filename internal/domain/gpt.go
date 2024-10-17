package domain

import "github.com/sashabaranov/go-openai"

type (
	ReturnOrder string
)

const (
	Asc  ReturnOrder = "asc"
	Desc ReturnOrder = "desc"
)

var (
	PendingStatuses = []openai.RunStatus{openai.RunStatusQueued, openai.RunStatusInProgress}
	FailedStatuses  = []openai.RunStatus{
		openai.RunStatusFailed,
		openai.RunStatusRequiresAction,
		openai.RunStatusIncomplete,
		openai.RunStatusCancelled,
		openai.RunStatusCancelling,
		openai.RunStatusExpired,
	}
	CompletedStatuses = []openai.RunStatus{openai.RunStatusCompleted}
)
