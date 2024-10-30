package domain

type UserInfo struct {
	ThreadId  string
	World     string
	Character string
	IsStarted bool
}

func (u UserInfo) IsEmpty() bool {
	if u.World == "" && u.Character == "" && u.ThreadId == "" {
		return true
	}

	return false
}
