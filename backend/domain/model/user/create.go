package user

import (
	"prc_hub_back/domain/model/jwt"

	"golang.org/x/crypto/bcrypt"
)

type CreateUserParam struct {
	Name           string  `json:"name"`
	Email          string  `json:"email"`
	Password       string  `json:"password"`
	TwitterId      *string `json:"twitter_id,omitempty"`
	GithubUsername *string `json:"github_username,omitempty"`
}

func (p CreateUserParam) validate() error {
	err := validateName(p.Name)
	if err != nil {
		return err
	}
	err = validateEmail(p.Email)
	if err != nil {
		return err
	}
	err = validatePassword(p.Password)
	if err != nil {
		return err
	}
	return nil
}

func CreateUser(p CreateUserParam) (UserWithToken, error) {
	// バリデーション
	err := p.validate()
	if err != nil {
		return UserWithToken{}, err
	}

	// パスワードをハッシュ化
	hashed, err := bcrypt.GenerateFromPassword([]byte(p.Password), 10)
	if err != nil {
		return UserWithToken{}, err
	}

	// ""(空文字)を`null`に置き換え
	if p.TwitterId != nil && *p.TwitterId == "" {
		p.TwitterId = nil
	}
	if p.GithubUsername != nil && *p.GithubUsername == "" {
		p.GithubUsername = nil
	}

	// リポジトリに追加
	// MySQLサーバーに接続
	d, err := OpenMysql()
	if err != nil {
		return UserWithToken{}, err
	}
	// return時にMySQLサーバーとの接続を閉じる
	defer d.Close()

	// `users`テーブルに追加
	r, err := d.Exec(
		`INSERT INTO users (name, email, password, post_event_availabled, manage, admin, twitter_id, github_username) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		p.Name, p.Email, string(hashed), false, false, false, p.TwitterId, p.GithubUsername,
	)
	if err != nil {
		return UserWithToken{}, err
	}
	id, err := r.LastInsertId()
	if err != nil {
		return UserWithToken{}, err
	}
	u := User{
		Id:                  id,
		Name:                p.Name,
		Email:               p.Email,
		Password:            string(hashed),
		StarCount:           0,
		PostEventAvailabled: false,
		Manage:              false,
		Admin:               false,
		TwitterId:           p.TwitterId,
		GithubUsername:      p.GithubUsername,
	}

	// jwtを生成
	uwt := UserWithToken{User: u}
	uwt.Token, err = jwt.GenerateToken(jwt.GenerateTokenParam{Id: u.Id, Email: u.Email, Admin: u.Admin})
	if err != nil {
		return UserWithToken{}, err
	}

	return uwt, nil
}
