package user

func Get(id string) (User, error) {
	// MySQLサーバーに接続
	db, err := OpenMysql()
	if err != nil {
		return User{}, err
	}
	// return時にMySQLサーバーとの接続を閉じる
	defer db.Close()

	// usersテーブルからidが一致する行を取得し、変数eに代入する
	r, err := db.Query(
		`SELECT
			id, name, email, password,
			post_event_availabled, manage, admin,
			twitter_id, github_username, star_count
		FROM users
		WHERE id = ?
		GROUP BY id`,
		id,
	)
	if err != nil {
		return User{}, err
	}
	defer r.Close()
	if !r.Next() {
		// 1行もレコードが無い場合
		// not found
		return User{}, ErrUserNotFound
	}

	// 変数に割り当て
	u := User{}
	err = r.Scan(
		&u.Id, &u.Name, &u.Email, &u.Password, &u.PostEventAvailabled,
		&u.Manage, &u.Admin, &u.TwitterId, &u.GithubUsername, &u.StarCount,
	)
	if err != nil {
		return User{}, err
	}

	return u, nil
}

func GetByEmail(email string) (User, error) {
	// MySQLサーバーに接続
	db, err := OpenMysql()
	if err != nil {
		return User{}, err
	}
	// return時にMySQLサーバーとの接続を閉じる
	defer db.Close()

	// `users`テーブルから`id`が一致する行を取得し、変数`e`に代入する
	r, err := db.Query(
		`SELECT id, password, admin FROM users WHERE email = ?`,
		email,
	)
	if err != nil {
		return User{}, err
	}
	defer r.Close()
	if !r.Next() {
		// 1行もレコードが無い場合
		// not found
		return User{}, ErrUserNotFound
	}

	// 変数に割り当て
	u := User{}
	err = r.Scan(&u.Id, &u.Password, &u.Admin)
	if err != nil {
		return User{}, err
	}

	u.Email = email
	return u, nil
}
