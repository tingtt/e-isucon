package user

import (
	"prc_hub_back/domain/model/mysql"
	"strings"
)

type GetUserListQueryParam struct {
	Name                *string `query:"name"`
	NameContain         *string `query:"name_contain"`
	PostEventAvailabled *bool   `json:"post_event_availabled"`
	Manage              *bool   `json:"manage"`
	Admin               *bool   `json:"admin"`
}

func GetList(q GetUserListQueryParam) ([]User, error) {
	// MySQLサーバーに接続
	db, err := mysql.Open()
	if err != nil {
		return nil, err
	}
	// return時にMySQLサーバーとの接続を閉じる
	defer db.Close()

	// クエリを作成
	query :=
		`SELECT
			id, name, email, password,
			post_event_availabled, manage, admin,
			twitter_id, github_username, star_count
		FROM users u
		WHERE`
	queryParams := []interface{}{}
	if q.PostEventAvailabled != nil {
		// 権限で絞り込み
		query += " post_event_availabled = ? AND"
		queryParams = append(queryParams, *q.PostEventAvailabled)
	}
	if q.Manage != nil {
		// 権限で絞り込み
		query += " manage = ? AND"
		queryParams = append(queryParams, *q.Manage)
	}
	if q.Admin != nil {
		// 権限で絞り込み
		query += " admin = ? AND"
		queryParams = append(queryParams, *q.Admin)
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
	query = strings.TrimSuffix(query, "	WHERE")
	query = strings.TrimSuffix(query, " AND")

	// `users`テーブルからを取得
	r, err := db.Query(query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	// 取得したテーブルを１行ずつ処理
	// 配列`users`に代入する
	var users []User
	for r.Next() {
		// 変数に割り当て
		u := User{}
		err = r.Scan(
			&u.Id, &u.Name, &u.Email, &u.Password, &u.PostEventAvailabled,
			&u.Manage, &u.Admin, &u.TwitterId, &u.GithubUsername, &u.StarCount,
		)
		if err != nil {
			return nil, err
		}

		// 配列に追加
		users = append(users, u)
	}

	return users, nil
}
