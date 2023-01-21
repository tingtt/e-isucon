package event

import (
	"prc_hub_back/domain/model/user"
	"time"
)

type Event struct {
	Id          string          `json:"id"`
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	Location    *string         `json:"location,omitempty"`
	Datetimes   []EventDatetime `json:"datetimes"`
	Published   bool            `json:"published"`
	Completed   bool            `json:"completed"`
	UserId      string          `json:"user_id"`
}

type EventDatetime struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end" dh:"end"`
}

type EventDocument struct {
	EventId string `json:"event_id"`
	Id      string `json:"id"`
	Name    string `json:"name"`
	Url     string `json:"url"`
}

type EventEmbed struct {
	Event
	User      *user.User       `json:"user,omitempty"`
	Documents *[]EventDocument `json:"documents,omitempty"`
}
