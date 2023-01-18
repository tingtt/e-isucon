package util

import "github.com/google/uuid"

func UUID() string {
	// UUIDを生成
	id := uuid.New()
	for id.String() == "" {
		id = uuid.New()
	}
	return id.String()
}
