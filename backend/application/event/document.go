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

func CreateDocument(p CreateEventDocumentParam, requestUserId string) (_ event.EventDocument, err error) {
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

func GetDocument(id string, requestUserId string) (_ event.EventDocument, err error) {
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

func GetDocumentList(q GetDocumentQueryParam, requestUserId string) ([]event.EventDocument, error) {
	return event.GetDocumentList(
		event.GetDocumentQueryParam{
			EventId:     q.EventId,
			Name:        q.Name,
			NameContain: q.NameContain,
		},
	)
}

func UpdateDocument(id string, p UpdateEventDocumentParam, requestUserId string) (event.EventDocument, error) {
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

func DeleteDocument(id string, requestUserId string) error {
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
