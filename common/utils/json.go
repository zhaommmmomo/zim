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
