package user

type User struct {
	Id                  int64   `json:"id"`
	Name                string  `json:"name"`
	Email               string  `json:"-"`
	Password            string  `json:"-"`
	PostEventAvailabled bool    `json:"post_event_availabled"`
	Manage              bool    `json:"manage"`
	Admin               bool    `json:"admin"`
	TwitterId           *string `json:"twitter_id,omitempty"`
	GithubUsername      *string `json:"github_username,omitempty"`
	StarCount           uint64  `json:"star_count"`
}

type UserWithToken struct {
	User
	Token string `json:"token"`
}
