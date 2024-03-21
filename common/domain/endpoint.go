package domain

type Endpoint struct {
	Name     string                 `json:"name"`
	Ip       string                 `json:"ip"`
	Port     int16                  `json:"port"`
	MetaData map[string]interface{} `json:"meta"`
}
