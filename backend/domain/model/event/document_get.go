package event

import "prc_hub_back/domain/model/user"

func GetDocument(id int64, requestUser user.User) (EventDocument, error) {
	// Get document

	// MySQLサーバーに接続
	db, err := OpenMysql()
	if err != nil {
		return EventDocument{}, err
	}
	// return時にMySQLサーバーとの接続を閉じる
	defer db.Close()

	// `documents`テーブルから`id`が一致する行を取得し、変数`ed`に代入する
	var ed EventDocument
	r, err := db.Query("SELECT * FROM documents WHERE id = ?", id)
	if err != nil {
		return EventDocument{}, err
	}
	defer r.Close()
	if !r.Next() {
		// 1行もレコードが無い場合
		// not found
		return EventDocument{}, ErrEventDocumentNotFound
	}
	err = r.Scan(&ed.Id, &ed.EventId, &ed.Name, &ed.Url)
	if err != nil {
		return EventDocument{}, err
	}

	// Get event
	e, err := GetEvent(ed.EventId, GetEventQueryParam{}, requestUser)
	if err != nil {
		return EventDocument{}, err
	}

	// 権限の検証
	if !requestUser.Admin && !requestUser.Manage &&
		!e.Published && e.UserId != requestUser.Id {
		// `User`が`Admin`・`Manage`のいずれでもなく
		// `Published`でない 且つ 自分のものでない`Event`は取得不可
		return EventDocument{}, ErrEventDocumentNotFound
	}

	return ed, nil
}
