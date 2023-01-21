package event

import "prc_hub_back/domain/model/user"

func CompleteEvent(id string, requestUser user.User) (Event, error) {
	completed := true
	return UpdateEvent(id, UpdateEventParam{Completed: &completed}, requestUser)
}
