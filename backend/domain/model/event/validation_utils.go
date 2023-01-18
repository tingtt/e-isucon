package event

import (
	"errors"
	"prc_hub_back/domain/model/user"
	"time"
)

// Errors
var (
	ErrValidateEventTitleCannotBeEmpty         = errors.New("Event title cannot be empty")
	ErrValidateDocumentNameCannotBeEmpty       = errors.New("Event document name cannot be empty")
	ErrValidateUrlCannotBeEmpty                = errors.New("Event document url cannot be empty")
	ErrValidateEventDatetimesCannotBeEmpty     = errors.New("Event datetime cannot be empty")
	ErrValidateEventDatetimeStartMustBeforeEnd = errors.New("Event start datetime must be before end datetime")
)

func validateTitle(title string) error {
	// 空文字チェック
	if title == "" {
		return ErrValidateEventTitleCannotBeEmpty
	}
	return nil
}

func validateDocumentName(name string) error {
	// 空文字チェック
	if name == "" {
		return ErrValidateDocumentNameCannotBeEmpty
	}
	return nil
}

func validateUrl(url string) error {
	// 空文字チェック
	if url == "" {
		return ErrValidateUrlCannotBeEmpty
	}
	return nil
}

func validateEventId(id int64, requestUser user.User) error {
	// MySQLサーバーに接続
	db, err := OpenMysql()
	if err != nil {
		return err
	}
	// return時にMySQLサーバーとの接続を閉じる
	defer db.Close()

	// `documents`テーブルから`id`が一致する行を確認
	_, err = GetDocument(id, requestUser)
	if err != nil {
		return err
	}
	return nil
}

func validateEventDatetime(start time.Time, end time.Time) error {
	if !start.Before(end) {
		return ErrValidateEventDatetimeStartMustBeforeEnd
	}
	return nil
}
