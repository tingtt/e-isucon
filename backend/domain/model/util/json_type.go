package util

import (
	"encoding/json"
)

type NullableJSONString struct {
	Value **string
}

func (n *NullableJSONString) UnmarshalJSON(data []byte) error {
	// jsonにキーが存在する場合にこの関数が呼び出される
	var valueP *string = nil
	if string(data) == "null" {
		// jsonにキーが存在し、値がnull
		n.Value = &valueP
		return nil
	}

	var tmp string
	tmpP := &tmp
	if err := json.Unmarshal(data, &tmp); err != nil {
		// typeエラー
		return err
	}
	// valid value
	n.Value = &tmpP
	return nil
}

func (n NullableJSONString) KeyExists() bool {
	return n.Value != nil
}

func (n NullableJSONString) IsNull() bool {
	return n.KeyExists() && *n.Value == nil
}

type NullableJSONBool struct {
	Value **bool
}

func (n *NullableJSONBool) UnmarshalJSON(data []byte) error {
	// jsonにキーが存在する場合にこの関数が呼び出される
	var valueP *bool = nil
	if string(data) == "null" {
		// jsonにキーが存在し、値がnull
		n.Value = &valueP
		return nil
	}

	var tmp bool
	tmpP := &tmp
	if err := json.Unmarshal(data, &tmp); err != nil {
		// typeエラー
		return err
	}
	// valid value
	n.Value = &tmpP
	return nil
}

func (n NullableJSONBool) KeyExists() bool {
	return n.Value != nil
}

func (n NullableJSONBool) IsNull() bool {
	return n.KeyExists() && *n.Value == nil
}
