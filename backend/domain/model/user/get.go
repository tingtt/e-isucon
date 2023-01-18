package user

func Get(id int64) (User, error) {
	// MySQLサーバーに接続
	db, err := OpenMysql()
	if err != nil {
		return User{}, err
	}
	// return時にMySQLサーバーとの接続を閉じる
	defer db.Close()

	// `users`テーブルから`id`が一致する行を取得し、変数`e`に代入する
	r, err := db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return User{}, err
	}
	defer r.Close()
	if !r.Next() {
		// 1行もレコードが無い場合
		// not found
		err = ErrUserNotFound
		return User{}, err
	}

	// 一時変数に割り当て
	var (
		id2                 int64
		name                string
		email               string
		password            string
		postEventAvailabled bool
		manage              bool
		admin               bool
		twitterId           *string
		githubUsername      *string
	)
	err = r.Scan(
		&id2, &name, &email, &password, &postEventAvailabled,
		&manage, &admin, &twitterId, &githubUsername,
	)
	if err != nil {
		return User{}, err
	}

	// スター数を取得
	var count uint64 = 0
	r2, err := db.Query("SELECT COUNT(*) FROM user_stars WHERE target_user_id = ?", id)
	if err != nil {
		return User{}, err
	}
	defer r2.Close()
	if !r2.Next() {
		return User{}, ErrConflictUserStars
	}
	err = r2.Scan(&count)
	if err != nil {
		return User{}, err
	}

	u := User{
		Id:                  id,
		Name:                name,
		Email:               email,
		Password:            password,
		StarCount:           count,
		PostEventAvailabled: postEventAvailabled,
		Manage:              manage,
		Admin:               admin,
		TwitterId:           twitterId,
		GithubUsername:      githubUsername,
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
		"SELECT * FROM users WHERE email = ?",
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

	// 一時変数に割り当て
	var (
		id                  int64
		name                string
		email2              string
		password            string
		postEventAvailabled bool
		manage              bool
		admin               bool
		twitterId           *string
		githubUsername      *string
	)
	err = r.Scan(
		&id, &name, &email2, &password, &postEventAvailabled,
		&manage, &admin, &twitterId, &githubUsername,
	)
	if err != nil {
		return User{}, err
	}

	// スター数を取得
	var count uint64 = 0
	r2, err := db.Query("SELECT COUNT(*) FROM user_stars WHERE target_user_id = ?", id)
	if err != nil {
		return User{}, err
	}
	defer r2.Close()
	if !r2.Next() {
		return User{}, ErrConflictUserStars
	}
	err = r2.Scan(&count)
	if err != nil {
		return User{}, err
	}

	u := User{
		Id:                  id,
		Name:                name,
		Email:               email,
		Password:            password,
		StarCount:           count,
		PostEventAvailabled: postEventAvailabled,
		Manage:              manage,
		Admin:               admin,
		TwitterId:           twitterId,
		GithubUsername:      githubUsername,
	}
	return u, nil
}
