package event

import "strings"

type GetDocumentQueryParam struct {
	EventId     *int64  `query:"event_id"`
	Name        *string `query:"name"`
	NameContain *string `query:"name_contain"`
}

func GetDocumentList(q GetDocumentQueryParam) ([]EventDocument, error) {
	// MySQLサーバーに接続
	db, err := OpenMysql()
	if err != nil {
		return nil, err
	}
	// return時にMySQLサーバーとの接続を閉じる
	defer db.Close()

	// クエリを作成
	query := "SELECT * FROM documents WHERE"
	queryParams := []interface{}{}
	if q.EventId != nil {
		// イベントIDで絞り込み
		query += " event_id = ? AND"
		queryParams = append(queryParams, *q.EventId)
	}
	if q.Name != nil {
		// ドキュメント名の一致で絞り込み
		query += " name = ? AND"
		queryParams = append(queryParams, *q.Name)
	}
	if q.NameContain != nil {
		// ドキュメント名に文字列が含まれるかで絞り込み
		query += " name LIKE ?"
		queryParams = append(queryParams, "%"+*q.NameContain+"%")
	}
	// 不要な末尾の句を切り取り
	query = strings.TrimSuffix(query, " WHERE")
	query = strings.TrimSuffix(query, " AND")

	// `documents`テーブルからを取得し、変数`documents`に代入する
	var documents []EventDocument
	r, err := db.Query(query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	// 1行ずつ読込
	for r.Next() {
		// カラム読み込み用変数
		var (
			id      int64
			eventId int64
			name    string
			url     string
		)
		err = r.Scan(&id, &eventId, &name, &url)
		if err != nil {
			return nil, err
		}
		documents = append(documents, EventDocument{id, eventId, name, url})
	}

	return documents, nil
}
