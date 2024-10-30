package models

import "github.com/siyoga/rollstory/internal/domain"

type (
	User struct {
		ThreadId  string `json:"thread_id"`
		IsStarted bool   `json:"is_started"`
		World     string `json:"world"`
		Character string `json:"character"`
	}
)

func (m User) ToDomain() domain.UserInfo {
	return domain.UserInfo{
		ThreadId:  m.ThreadId,
		World:     m.World,
		IsStarted: m.IsStarted,
		Character: m.Character,
	}
}

func (_ User) FromDomain(u domain.UserInfo) User {
	return User{
		ThreadId:  u.ThreadId,
		World:     u.World,
		IsStarted: u.IsStarted,
		Character: u.Character,
	}
}
