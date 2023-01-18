package event

import (
	"context"
	"database/sql"
	"errors"
	"prc_hub_back/domain/model/user"
	"prc_hub_back/domain/model/util"
	"strings"
)

// Errors
var (
	ErrCannotUpdateEvent = errors.New("sorry, you cannot update this event")
)

type UpdateEventParam struct {
	Name        *string                     `json:"name,omitempty"`
	Description util.NullableJSONString     `json:"description,omitempty"`
	Location    util.NullableJSONString     `json:"location,omitempty"`
	Datetimes   *[]CreateEventDatetimeParam `json:"datetimes,omitempty"`
	Published   *bool                       `json:"published,omitempty"`
	Completed   *bool                       `json:"completed,omitempty"`
}

func (p UpdateEventParam) validate(id int64, requestUser user.User) error {
	/**
	 * フィールドの検証
	**/
	if p.Name == nil &&
		p.Description.KeyExists() &&
		p.Location.KeyExists() &&
		p.Datetimes == nil &&
		p.Published == nil &&
		p.Completed == nil {
		return ErrNoUpdates
	}
	// `Name`
	if p.Name != nil {
		err := validateTitle(*p.Name)
		if err != nil {
			return err
		}
	}
	// `Datetimes`
	if p.Datetimes != nil {
		if len(*p.Datetimes) == 0 {
			return ErrValidateEventDatetimesCannotBeEmpty
		}
		for _, d := range *p.Datetimes {
			err := d.validate()
			if err != nil {
				return err
			}
		}
	}

	// 権限の検証
	if !requestUser.Admin && !requestUser.Manage {
		// Eventを取得
		e, err := GetEvent(id, GetEventQueryParam{}, requestUser)
		if err != nil {
			return err
		}

		if requestUser.Id != e.UserId {
			// `Admin`・`Manage`のいずれでもなく`Event.UserId`が自分ではない場合は変更不可
			return ErrCannotUpdateEvent
		}
	}

	return nil
}

func UpdateEvent(id int64, p UpdateEventParam, requestUser user.User) (Event, error) {
	// バリデーション
	err := p.validate(id, requestUser)
	if err != nil {
		return Event{}, err
	}

	// MySQLサーバーに接続
	db, err := OpenMysql()
	if err != nil {
		return Event{}, err
	}
	// return時にMySQLサーバーとの接続を閉じる
	defer db.Close()

	// トランザクション開始
	tx, err := db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return Event{}, err
	}
	defer func() {
		// return時にトランザクションの後処理
		//* 17行目の`defer`より先に実行される
		if err != nil {
			// 失敗時はロールバック
			tx.Rollback()
		} else {
			// 成功時はコミット
			tx.Commit()
		}
	}()

	// `events`テーブル用のクエリを作成
	query1 := "UPDATE events SET"
	queryParams1 := []interface{}{}
	if p.Name != nil {
		// `name`を変更
		query1 += " name = ?,"
		queryParams1 = append(queryParams1, *p.Name)
	}
	if p.Description.KeyExists() {
		// `description`を変更
		query1 += " description = ?,"
		if p.Description.IsNull() {
			queryParams1 = append(queryParams1, nil)
		} else {
			queryParams1 = append(queryParams1, *p.Description.Value)
		}
	}
	if p.Location.KeyExists() {
		// `location`を変更
		query1 += " location = ?,"
		if p.Location.IsNull() {
			queryParams1 = append(queryParams1, nil)
		} else {
			queryParams1 = append(queryParams1, *p.Location.Value)
		}
	}
	if p.Published != nil {
		// `published`を変更
		query1 += " published = ?,"
		queryParams1 = append(queryParams1, *p.Published)
	}
	if p.Completed != nil {
		// `completed`を変更
		query1 += " completed = ?"
		queryParams1 = append(queryParams1, *p.Completed)
	}
	// 更新するフィールドがあるか確認
	if strings.HasSuffix(query1, "SET") {
		// 更新するフィールドが無いため中断
		err = ErrNoUpdates
		return Event{}, err
	}
	// 不要な末尾の句を切り取り
	query1 = strings.TrimSuffix(query1, ",")

	// `events`テーブルの`id`が一致する行を更新
	r2, err := tx.Exec(query1+" WHERE id = ?", append(queryParams1, id))
	if err != nil {
		return Event{}, err
	}
	i, err := r2.RowsAffected()
	if err != nil {
		return Event{}, err
	}
	if i != 1 {
		// 変更された行数が1ではない場合
		// `id`に一致する`event`が存在しない
		return Event{}, ErrEventNotFound
	}

	if p.Datetimes != nil {
		// `event_datetimes`テーブルの更新

		// 既存のデータを削除
		_, err = tx.Exec(
			"DELETE FROM event_datetimes WHERE event_id = ?",
			id,
		)
		if err != nil {
			return Event{}, err
		}

		// 新規データの追加
		// `event_datetimes`テーブルに追加
		for _, dt := range *p.Datetimes {
			_, err = tx.Exec(
				"INSERT INTO event_datetimes (event_id, start, end) VALUES (?, ?, ?)",
				id, dt.Start, dt.End,
			)
			if err != nil {
				return Event{}, err
			}
		}
	}

	// 更新後のデータを取得
	ee, err := GetEvent(id, GetEventQueryParam{}, requestUser)
	if err != nil {
		return Event{}, err
	}

	return ee.Event, nil
}
