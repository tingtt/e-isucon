package event

import (
	"prc_hub_back/domain/model/user"
	"strings"
	"time"
)

type GetEventListQueryParam struct {
	Published       *bool     `query:"published"`
	Name            *string   `query:"name"`
	NameContain     *string   `query:"name_contain"`
	Location        *string   `query:"location"`
	LocationContain *string   `query:"location_contain"`
	Embed           *[]string `query:"embed"`
}

func GetEventList(q GetEventListQueryParam, requestUser user.User) ([]EventEmbed, error) {
	// MySQLサーバーに接続
	db, err := OpenMysql()
	if err != nil {
		return nil, err
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

	// `Event`リストを取得

	// 取得用変数
	events := []EventEmbed{}

	// クエリを作成
	query := "SELECT * FROM events WHERE"
	queryParams := []interface{}{}
	if q.Name != nil {
		// イベント名の一致で絞り込み
		query += " name = ? AND"
		queryParams = append(queryParams, *q.Name)
	}
	if q.NameContain != nil {
		// イベント名に文字列が含まれるかで絞り込み
		query += " name LIKE ? AND"
		queryParams = append(queryParams, "%"+*q.NameContain+"%")
	}
	if q.Location != nil {
		// `Location`の一致で絞り込み
		query += " location = ? AND"
		queryParams = append(queryParams, *q.Location)
	}
	if q.LocationContain != nil {
		// `Location`に文字列が含まれるかで絞り込み
		query += " location LIKE ? AND"
		queryParams = append(queryParams, "%"+*q.LocationContain+"%")
	}
	if q.Published != nil {
		// `Published`で絞り込み
		query += " published = ?"
		queryParams = append(queryParams, *q.Published)
	}
	// 不要な末尾の句を切り取り
	query = strings.TrimSuffix(query, "WHERE")
	query = strings.TrimSuffix(query, "AND")

	// 実行
	r1, err := db.Query(query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer r1.Close()

	// １行ずつ処理
	for r1.Next() {
		// 一時変数に割当
		var (
			id          int64
			name        string
			description *string
			location    *string
			published   bool
			completed   bool
			userId      int64
		)
		err = r1.Scan(&id, &name, &description, &location, &published, &completed, &userId)
		if err != nil {
			return nil, err
		}
		// 配列追加用変数
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
			return nil, err
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
				return nil, err
			}
			// 配列に追加
			event.Event.Datetimes = append(event.Event.Datetimes, EventDatetime{*start, *end})
		}

		if embedUser {
			// `User`を取得
			u, err := user.Get(userId)
			if err != nil {
				return nil, err
			}
			// 変数に追加
			event.User = &u
		}

		if embedDocuments {
			// `Documents`を取得
			ed, err := GetDocumentList(GetDocumentQueryParam{EventId: &id})
			if err != nil {
				return nil, err
			}
			event.Documents = &ed
		}

		events = append(events, event)
	}

	return events, nil
}
