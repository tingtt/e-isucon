package event

import (
	"prc_hub_back/domain/model/mysql"
	"prc_hub_back/domain/model/user"
	"time"
)

type GetEventQueryParam struct {
	Embed *[]string `query:"embed"`
}

func GetEvent(id string, q GetEventQueryParam, requestUser user.User) (EventEmbed, error) {
	// Get event
	// MySQLサーバーに接続
	d, err := mysql.Open()
	if err != nil {
		return EventEmbed{}, err
	}
	// return時にMySQLサーバーとの接続を閉じる
	defer d.Close()

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

	// クエリを作成
	query :=
		`WITH params AS ( SELECT ? as event_id )
		SELECT
			e.name, e.description, e.location, e.published, e.completed, e.user_id,
			null AS start, null AS end,
			null AS doc_id, null AS doc_name, null AS doc_url,`
	if embedUser {
		// `users`テーブル結合
		query +=
			` u.id, u.name, u.email, u.post_event_availabled, u.manage, u.admin, u.twitter_id, u.github_username`
	} else {
		query +=
			` null AS u_id, null AS u_name, null AS u_email, null AS post_event_availabled, null AS manage, null AS admin, null AS twitter_id, null AS github_username`
	}
	query +=
		` FROM events e`
	if embedUser {
		// `users`テーブル結合
		query +=
			` LEFT JOIN users u ON e.user_id = u.id`
	}
	query += `
		WHERE e.id IN (SELECT event_id FROM params)`
	// `event_datetimes`テーブルを結合
	query +=
		` UNION ALL
		SELECT
			null, null, null, null, null, null,
			dt.start, dt.end,
			null, null, null,
			null, null, null, null, null, null, null, null
		FROM event_datetimes dt
		WHERE dt.event_id IN (SELECT event_id FROM params)`
	if embedDocuments {
		// `documents`テーブルを結合
		query +=
			` UNION ALL
			SELECT
				null, null, null, null, null, null,
				null, null,
				doc.id, doc.name, doc.url,
				null, null, null, null, null, null, null, null
			FROM documents doc
			WHERE doc.event_id IN (SELECT event_id FROM params)`
	}
	// 順序を保証するためにUNION後にソート (nameがNULLではない行を最初に返す)
	query += " ORDER BY 1 IS NULL ASC"

	// クエリを実行
	r, err := d.Query(query, id)
	if err != nil {
		return EventEmbed{}, err
	}
	defer r.Close()

	// 読み込み用変数
	var (
		tmpEvent     *EventEmbed     = nil
		tmpDocuments []EventDocument = nil
	)
	// `id`に一致した`event`が読み込まれるまで仮のエラーを代入
	err = ErrEventNotFound
	// 1行ずつ読込
	for r.Next() {
		// カラム読み込み用変数
		var (
			eName        *string
			eDescription *string
			eLocation    *string
			ePublished   *bool
			eCompleted   *bool
			eUserId      *string

			eDtStart *time.Time
			eDtEnd   *time.Time

			eDocId   *string
			eDocName *string
			eDocUrl  *string

			uId                  *string
			uName                *string
			uEmail               *string
			uPostEventAvailabled *bool
			uManage              *bool
			uAdmin               *bool
			uTwitterId           *string
			uGithubId            *string
		)
		// 変数に読み込み
		err = r.Scan(
			&eName, &eDescription, &eLocation, &ePublished, &eCompleted, &eUserId,
			&eDtStart, &eDtEnd,
			&eDocId, &eDocName, &eDocUrl,
			&uId, &uName, &uEmail, &uPostEventAvailabled, &uManage, &uAdmin, &uTwitterId, &uGithubId,
		)
		if err != nil {
			return EventEmbed{}, err
		}
		// 読み込んだ内容によって読み込み用変数のそれぞれのフィールドに代入
		if tmpEvent == nil {
			if eName == nil || eUserId == nil {
				// `id`に一致する`event`が存在しない
				return EventEmbed{}, ErrEventNotFound
			}
			// Scanしたフィールドを代入
			tmpEvent = &EventEmbed{
				Event: Event{
					Id:          id,
					Name:        *eName,
					Description: eDescription,
					UserId:      *eUserId,
				},
			}
			if ePublished != nil {
				tmpEvent.Published = *ePublished
			}
			if eCompleted != nil {
				tmpEvent.Completed = *eCompleted
			}
			// `id`が一致した`event`が見つかったためエラーを解消
			err = nil

			if uId != nil && uName != nil && uEmail != nil && uPostEventAvailabled != nil && uManage != nil && uAdmin != nil {
				// `user`が取得された場合、Scanしたカラムの値を代入
				tmpEvent.User = &user.User{
					Id:                  *uId,
					Name:                *uName,
					Email:               *uEmail,
					PostEventAvailabled: *uPostEventAvailabled,
					Manage:              *uManage,
					Admin:               *uAdmin,
					TwitterId:           uTwitterId,
					GithubUsername:      uGithubId,
				}
			}
		}
		if tmpEvent != nil && eDtStart != nil && eDtEnd != nil {
			// `event_datetime`が取得された場合、Scanしたカラムの値を代入
			tmpEvent.Datetimes = append(
				tmpEvent.Datetimes,
				EventDatetime{
					Start: *eDtStart,
					End:   *eDtEnd,
				},
			)
		}
		if tmpEvent != nil && eDocId != nil && eDocName != nil && eDocUrl != nil {
			// `document`が取得された場合、Scanしたカラムの値を代入
			tmpDocuments = append(
				tmpDocuments,
				EventDocument{
					Id:      *eDocId,
					Name:    *eDocName,
					Url:     *eDocUrl,
					EventId: id,
				},
			)
		}
	}
	// 読み込み用変数を統合
	tmpEvent.Documents = &tmpDocuments

	return *tmpEvent, err
}
