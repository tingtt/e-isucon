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
	_, err = db.Exec("UPDATE users SET star_count = star_count + 1 WHERE id = ?", userId)
	if err != nil {
		return 0, err
	}

	// スター数のカウントを取得
	r, err := db.Query("SELECT star_count FROM users WHERE id = ?", userId)
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
