package domain

import "github.com/siyoga/rollstory/internal/models"

type Story struct {
	World     string
	Character string
}

func (_ Story) FromModel(m models.Story) Story {
	return Story{
		World:     m.World,
		Character: m.Character,
	}
}

func (s Story) ToModel() models.Story {
	return models.Story{
		World:     s.World,
		Character: s.Character,
	}
}

func (s Story) IsEmpty() bool {
	if s.World == "" && s.Character == "" {
		return true
	}

	return false
}
