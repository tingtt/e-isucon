package event

import "prc_hub_back/domain/model/event"

func Init(user string, password string, host string, port uint, db string) {
	event.InitRepository(user, password, host, port, db)
}
