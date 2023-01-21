package event

import (
	"errors"
	"prc_hub_back/domain/model/logger"
	"prc_hub_back/domain/model/user"

	"github.com/google/uuid"
)

// Errors
var (
	ErrCannotCreateEventDocument = errors.New("sorry, you cannot create document to this event")
)

type CreateEventDocumentParam struct {
	EventId string `json:"event_id"`
	Name    string `json:"name"`
	Url     string `json:"url"`
}

func (p CreateEventDocumentParam) validate(requestUser user.User) error {
	// フィールドの検証
	err := validateDocumentName(p.Name)
	if err != nil {
		return err
	}
	err = validateUrl(p.Url)
	if err != nil {
		return err
	}
	err = validateEventId(p.EventId, requestUser)
	if err != nil {
		return err
	}

	// 権限の検証
	if !requestUser.Admin && !requestUser.Manage {
		// Eventを取得
		e, err := GetEvent(p.EventId, GetEventQueryParam{}, requestUser)
		if err != nil {
			return err
		}

		if requestUser.Id != e.UserId {
			// `Admin`・`Manage`のいずれでもなく`Event.UserId`が自分ではない場合は追加不可
			return ErrCannotCreateEventDocument
		}
	}

	return nil
}

func CreateEventDocument(p CreateEventDocumentParam, requestUser user.User) (EventDocument, error) {
	err := p.validate(requestUser)
	if err != nil {
		return EventDocument{}, err
	}

	id := uuid.New().String()
	e := EventDocument{
		Id:      id,
		EventId: p.EventId,
		Name:    p.Name,
		Url:     p.Url,
	}

	go func() {
		// MySQLサーバーに接続
		db, err := OpenMysql()
		if err != nil {
			logger.Logger().Fatalf("Failed:\n\terr: %v", err)
			return
		}
		// return時にMySQLサーバーとの接続を閉じる
		defer db.Close()

		// `documents`テーブルに追加
		_, err = db.Exec(
			`INSERT INTO documents (id, event_id, name, url) VALUES (?, ?, ?, ?)`,
			id, p.EventId, p.Name, p.Url,
		)
		if err != nil {
			logger.Logger().Fatalf("Failed:\n\terr: %v", err)
			return
		}
	}()

	return e, nil
}
