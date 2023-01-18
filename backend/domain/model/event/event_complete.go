package event

import "prc_hub_back/domain/model/user"

func CompleteEvent(id int64, requestUser user.User) (Event, error) {
	completed := true
	return UpdateEvent(id, UpdateEventParam{Completed: &completed}, requestUser)
}
