package event

import (
	"prc_hub_back/domain/model/user"
	"time"
)

type GetEventQueryParam struct {
	Embed *[]string `query:"embed"`
}

func GetEvent(id int64, q GetEventQueryParam, requestUser user.User) (EventEmbed, error) {
	// Get event
	// MySQLサーバーに接続
	db, err := OpenMysql()
	if err != nil {
		return EventEmbed{}, err
	}
	// return時にMySQLサーバーとの接続を閉じる
	defer db.Close()

	embedUser := false
	embedDocuments := false
	if q.Embed != nil {
		for _, e := range *q.Embed {
			if e == "user" {
				embedUser = true
			}
			if e == "documents" {
				embedDocuments = true
			}
		}
	}

	// `Event`を取得
	r1, err := db.Query("SELECT * FROM events WHERE id = ?", id)
	if err != nil {
		return EventEmbed{}, err
	}
	defer r1.Close()
	if !r1.Next() {
		// 1行もレコードが無い場合
		// not found
		return EventEmbed{}, ErrEventNotFound
	}
	// 一時変数に割当
	var (
		id2         int64
		name        string
		description *string
		location    *string
		published   bool
		completed   bool
		userId      int64
	)
	err = r1.Scan(&id2, &name, &description, &location, &published, &completed, &userId)
	if err != nil {
		return EventEmbed{}, err
	}

	// 返り値用変数
	event := EventEmbed{
		Event: Event{
			Id:          id,
			Name:        name,
			Description: description,
			Location:    location,
			Datetimes:   []EventDatetime{},
			Published:   published,
			Completed:   completed,
			UserId:      userId,
		},
	}

	// `EventDatetime`を取得
	r2, err := db.Query("SELECT * FROM event_datetimes WHERE event_id = ?", id)
	if err != nil {
		return EventEmbed{}, err
	}
	defer r2.Close()
	for r2.Next() {
		var (
			eId   string
			start *time.Time
			end   *time.Time
		)
		err = r2.Scan(&eId, &start, &end)
		if err != nil {
			return EventEmbed{}, err
		}
		// 配列に追加
		event.Event.Datetimes = append(event.Event.Datetimes, EventDatetime{*start, *end})
	}

	if embedUser {
		// `User`を取得
		u, err := user.Get(event.UserId)
		if err != nil {
			return EventEmbed{}, err
		}
		// 変数に追加
		event.User = &u
	}

	if embedDocuments {
		// `Documents`を取得
		ed, err := GetDocumentList(GetDocumentQueryParam{EventId: &id})
		if err != nil {
			return EventEmbed{}, err
		}
		event.Documents = &ed
	}

	return event, nil
}
