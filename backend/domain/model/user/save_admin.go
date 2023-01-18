package user

import (
	"golang.org/x/crypto/bcrypt"
)

func SaveAdmin(email string, password string) error {
	tmpBool := true
	tmpName := "admin"
	u, err := GetList(GetUserListQueryParam{Admin: &tmpBool, Name: &tmpName})
	if err != nil {
		return err
	}

	if len(u) == 1 {
		// `Admin`の`User`が登録済
		var newEmail *string = nil
		var newPassword *string = nil
		if u[0].Email != email {
			// `Email`が不一致
			newEmail = &email
		}
		if verify, err := u[0].Verify(password); err != nil {
			return err
		} else if !verify {
			// `Password`が不一致
			// 新規パスワード
			newPassword = &password
		}

		// `User`更新
		_, err := Update(u[0].Id, UpdateUserParam{
			Email:    newEmail,
			Password: newPassword,
		}, u[0])
		if err != nil {
			return err
		}
	} else {
		// `Admin`の`User`が未登録

		// パスワードをハッシュ化
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
		if err != nil {
			return err
		}

		// リポジトリに追加
		// MySQLサーバーに接続
		d, err := OpenMysql()
		if err != nil {
			return err
		}
		// return時にMySQLサーバーとの接続を閉じる
		defer d.Close()

		// `users`テーブルに追加
		_, err = d.Exec(
			`INSERT INTO users (name, email, password, post_event_availabled, manage, admin, twitter_id, github_username) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			"admin", email, string(hashed), true, true, true, nil, nil,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
