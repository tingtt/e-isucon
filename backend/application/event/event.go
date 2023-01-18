package event

import (
	"prc_hub_back/application/user"
	"prc_hub_back/domain/model/event"
	userDomain "prc_hub_back/domain/model/user"
)

type (
	CreateEventParam       = event.CreateEventParam
	UpdateEventParam       = event.UpdateEventParam
	GetEventListQueryParam = event.GetEventListQueryParam
	GetEventQueryParam     = event.GetEventQueryParam
)

func CreateEvent(p CreateEventParam, requestUserId int64) (event.Event, error) {
	// リクエスト元のユーザーを取得
	u, err := user.Get(requestUserId)
	if err != nil {
		return event.Event{}, err
	}

	return event.CreateEvent(p, u)
}

func GetEvent(id int64, q GetEventQueryParam, requestUserId *int64) (event.EventEmbed, error) {
	u := new(userDomain.User)

	if requestUserId != nil {
		// リクエスト元のユーザーを取得
		var u2 userDomain.User
		u2, err := user.Get(*requestUserId)
		if err != nil {
			return event.EventEmbed{}, err
		}
		u = &u2
	} else if requestUserId == nil {
		// リクエストユーザーが指定されていない場合は最小権限のユーザーを仮使用
		u = &userDomain.User{
			Id:                  0,
			PostEventAvailabled: false,
			Manage:              false,
			Admin:               false,
		}
	}

	return event.GetEvent(id, q, *u)
}

func GetEventList(q GetEventListQueryParam, requestUserId *int64) ([]event.EventEmbed, error) {
	u := new(userDomain.User)

	if requestUserId != nil {
		// リクエスト元のユーザーを取得
		var u2 userDomain.User
		u2, err := user.Get(*requestUserId)
		if err != nil {
			return nil, err
		}
		u = &u2
	} else if requestUserId == nil {
		// リクエストユーザーが指定されていない場合は最小権限のユーザーを仮使用
		u = &userDomain.User{
			Id:                  0,
			PostEventAvailabled: false,
			Manage:              false,
			Admin:               false,
		}
	}

	return event.GetEventList(q, *u)
}

func UpdateEvent(id int64, p UpdateEventParam, requestUserId int64) (event.Event, error) {
	// リクエスト元のユーザーを取得
	u, err := user.Get(requestUserId)
	if err != nil {
		return event.Event{}, err
	}

	return event.UpdateEvent(id, p, u)
}

func DeleteEvent(id int64, requestUserId int64) error {
	// リクエスト元のユーザーを取得
	u, err := user.Get(requestUserId)
	if err != nil {
		return err
	}

	return event.DeleteEvent(id, u)
}
