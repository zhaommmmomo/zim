package utils

import "encoding/json"

func Marshal(a any) string {
	data, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func UnMarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func MapToObj(m *map[string]interface{}, a any) error {
	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, a)
}
