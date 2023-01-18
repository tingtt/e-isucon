package event

import (
	"prc_hub_back/application/user"
	"prc_hub_back/domain/model/event"
)

type (
	CreateEventDocumentParam event.CreateEventDocumentParam
	UpdateEventDocumentParam event.UpdateEventDocumentParam
	GetDocumentQueryParam    event.GetDocumentQueryParam
)

func CreateDocument(p CreateEventDocumentParam, requestUserId int64) (_ event.EventDocument, err error) {
	// リクエスト元のユーザーを取得
	u, err := user.Get(requestUserId)
	if err != nil {
		return
	}

	return event.CreateEventDocument(
		event.CreateEventDocumentParam{
			EventId: p.EventId,
			Name:    p.Name,
			Url:     p.Url,
		},
		u,
	)
}

func GetDocument(id int64, requestUserId int64) (_ event.EventDocument, err error) {
	// リクエスト元のユーザーを取得
	u, err := user.Get(requestUserId)
	if err != nil {
		return
	}

	return event.GetDocument(
		id,
		u,
	)
}

func GetDocumentList(q GetDocumentQueryParam, requestUserId int64) ([]event.EventDocument, error) {
	return event.GetDocumentList(
		event.GetDocumentQueryParam{
			EventId:     q.EventId,
			Name:        q.Name,
			NameContain: q.NameContain,
		},
	)
}

func UpdateDocument(id int64, p UpdateEventDocumentParam, requestUserId int64) (event.EventDocument, error) {
	// リクエスト元のユーザーを取得
	u, err := user.Get(requestUserId)
	if err != nil {
		return event.EventDocument{}, err
	}

	return event.UpdateEventDocument(
		id,
		event.UpdateEventDocumentParam{
			Name: p.Name,
			Url:  p.Url,
		},
		u,
	)
}

func DeleteDocument(id int64, requestUserId int64) error {
	// リクエスト元のユーザーを取得
	u, err := user.Get(requestUserId)
	if err != nil {
		return err
	}

	return event.DeleteEventDocument(
		id,
		u,
	)
}
