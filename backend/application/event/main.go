package event

import "prc_hub_back/domain/model/event"

func Init(user string, password string, db string) {
	event.InitRepository(user, password, db)
}
