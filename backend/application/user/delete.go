package user

import "prc_hub_back/domain/model/user"

func Delete(id int64, requestUserId int64) error {
	// リクエスト元のユーザーを取得
	u, err := Get(id)
	if err != nil {
		return err
	}

	return user.DeleteUesr(id, u)
}
