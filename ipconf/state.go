package ipconf

import (
	"encoding/json"
	"math/rand"
)

type State struct {
	MaxCpuUse       float64 `json:"max_cpu_use"`
	CpuUse          float64 `json:"cpu_use"`
	MaxMemUse       float64 `json:"max_mem_use"`
	MemUse          float64 `json:"mem_use"`
	MaxConnectCount float64 `json:"max_connect_count"`
	ConnectCount    float64 `json:"connect_count"`
	MaxBandwidth    float64 `json:"max_bandwidth"`
	MessageBytes    float64 `json:"message_bytes"`
}

func convertState(m *map[string]interface{}) (*State, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	state := &State{}
	err = json.Unmarshal(bytes, state)
	return state, err
}

// todo: 基于state信息计算对应的分数
func (s *State) calculateScore() int16 {
	return int16(rand.Intn(101))
}
