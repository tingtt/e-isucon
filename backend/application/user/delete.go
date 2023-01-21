package user

import "prc_hub_back/domain/model/user"

func Delete(id string, requestUserId string) error {
	// リクエスト元のユーザーを取得
	u, err := Get(id)
	if err != nil {
		return err
	}

	return user.DeleteUesr(id, u)
}
