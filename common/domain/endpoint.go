package domain

type Endpoint struct {
	Ip    string `json:"ip"`
	Port  int16  `json:"port"`
	Score int16  `json:"-"`
}
