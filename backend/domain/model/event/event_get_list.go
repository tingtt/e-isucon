package event

import (
	"errors"
	"fmt"
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
	d, err := OpenMysql()
	if err != nil {
		return nil, err
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

	// 絞り込み用の仮クエリ
	queryStrBase :=
		`SELECT
			id AS event_id
		FROM events
		WHERE`
	queryParams := []interface{}{}
	if q.Name != nil {
		// イベント名の一致で絞り込み
		queryStrBase += " name = ? AND"
		queryParams = append(queryParams, *q.Name)
	}
	if q.NameContain != nil {
		// イベント名に文字列が含まれるかで絞り込み
		queryStrBase += " name LIKE ? AND"
		queryParams = append(queryParams, "%"+*q.NameContain+"%")
	}
	if q.Location != nil {
		// `Location`の一致で絞り込み
		queryStrBase += " location = ? AND"
		queryParams = append(queryParams, *q.Location)
	}
	if q.LocationContain != nil {
		// `Location`に文字列が含まれるかで絞り込み
		queryStrBase += " location LIKE ? AND"
		queryParams = append(queryParams, "%"+*q.LocationContain+"%")
	}
	if q.Published != nil {
		// `Published`で絞り込み
		queryStrBase += " published = ?"
		queryParams = append(queryParams, *q.Published)
	}
	// 不要な末尾の句を切り取り
	queryStrBase = strings.TrimSuffix(queryStrBase, "WHERE")
	queryStrBase = strings.TrimSuffix(queryStrBase, "AND")

	// 本クエリ
	query :=
		`SELECT * FROM (
		WITH params AS ( ` + queryStrBase + ` )
		SELECT
			e.id, e.name, e.description, e.location, e.published, e.completed, e.user_id,
			null AS start, null AS end,
			null AS doc_id, null AS doc_name, null AS doc_url,`
	if embedUser {
		// `users`テーブル結合
		query +=
			` u.id AS u_id, u.name AS u_name, u.email AS u_email, u.post_event_availabled, u.manage, u.admin, u.twitter_id, u.github_username`
	} else {
		query +=
			` null AS u_id, null AS u_name, null AS u_email, null AS post_event_availabled, null AS manage, null AS admin, null AS twitter_id, null AS github_username`
	}
	query += `
		FROM events e`
	if embedUser {
		// `users`テーブル結合
		query += ` LEFT JOIN users u ON e.user_id = u.id`
	}
	query += `
		WHERE e.id IN (SELECT event_id FROM params)`
	// `event_datetimes`テーブルを結合
	query +=
		` UNION ALL
		SELECT
			dt.event_id, null, null, null, null, null, null,
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
				doc.event_id, null, null, null, null, null, null,
				null, null,
				doc.id, doc.name, doc.url,
				null, null, null, null, null, null, null, null
			FROM documents doc
			WHERE doc.event_id IN (SELECT event_id FROM params)`
	}
	// 順序を保証するためにUNION後にソート (event.idの昇順でソートした上でevent.nameがNULLではない行を最初に返す)
	query += ") AS e ORDER BY id, name IS NULL ASC"

	// クエリを実行
	r, err := d.Query(query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	// 読み込み用変数
	events := []EventEmbed{}
	var (
		tmpEvent     *EventEmbed      = nil
		tmpDocuments *[]EventDocument = &[]EventDocument{}
	)
	// 1行ずつ読込
	for r.Next() {
		// カラム読み込み用変数
		var (
			eId          string
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
			&eId, &eName, &eDescription, &eLocation, &ePublished, &eCompleted, &eUserId,
			&eDtStart, &eDtEnd,
			&eDocId, &eDocName, &eDocUrl,
			&uId, &uName, &uEmail, &uPostEventAvailabled, &uManage, &uAdmin, &uTwitterId, &uGithubId,
		)
		if err != nil {
			return nil, err
		}
		// 読み込んだ内容によって読み込み用変数のそれぞれのフィールドに代入
		if tmpEvent == nil || tmpEvent.Id != eId {
			if eName == nil || eUserId == nil {
				// 想定外の値なため処理を中断
				err = errors.New("invalid column set of row")
				fmt.Printf("err: %v\n", err)
				return nil, err
			}
			if tmpEvent != nil {
				// 読込中の`event`が存在する場合は、結果用配列に追加
				// 読み込み用変数を統合
				tmpEvent.Documents = tmpDocuments
				events = append(events, *tmpEvent)
				// 読み込み用変数をクリア
				tmpEvent = nil
				tmpDocuments = &[]EventDocument{}
			}
			// Scanしたフィールドを代入
			tmpEvent = &EventEmbed{
				Event: Event{
					Id:          eId,
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
			*tmpDocuments = append(
				*tmpDocuments,
				EventDocument{
					Id:      *eDocId,
					Name:    *eDocName,
					Url:     *eDocUrl,
					EventId: eId,
				},
			)
		}
	}
	// 最後に読み込んだ`event`が存在する場合は結果用配列に追加
	if tmpEvent != nil {
		tmpEvent.Documents = tmpDocuments
		events = append(events, *tmpEvent)
	}

	return events, err
}
