package user

func AddStar(userId uint64) (count uint64, err error) {
	_, err = Get(int64(userId))
	if err != nil {
		return 0, err
	}

	// MySQLサーバーに接続
	db, err := OpenMysql()
	if err != nil {
		return 0, err
	}
	// return時にMySQLサーバーとの接続を閉じる
	defer db.Close()

	// `user_stars`テーブルに追加
	_, err = db.Exec("INSERT INTO user_stars (target_user_id) VALUES (?)", userId)
	if err != nil {
		return 0, err
	}

	// スター数のカウントを取得
	r, err := db.Query("SELECT COUNT(*) FROM user_stars WHERE target_user_id = ?", userId)
	if err != nil {
		return 0, err
	}
	defer r.Close()

	if !r.Next() {
		return 0, ErrConflictUserStars
	}

	err = r.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, err
}
