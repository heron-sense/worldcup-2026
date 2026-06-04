package model

import "encoding/json"

// SlotState 单个 slot 的展示状态。
type SlotState struct {
	Visible string `json:"visible"` // Y 或 N
	Value   string `json:"value"`
}

// SceneSettings 序列化在 game_profile.scene_settings 中的 JSON 结构。
type SceneSettings struct {
	Credit int                  `json:"credit"`
	Slots  map[string]SlotState `json:"slots"`
}

// DefaultSceneSettings 新档案的默认场景配置。
func DefaultSceneSettings() SceneSettings {
	return SceneSettings{
		Credit: 5,
		Slots: map[string]SlotState{
			"slot1": {Visible: "N", Value: "1USD"},
			"slot2": {Visible: "N", Value: "2USD"},
			"slot3": {Visible: "N", Value: "5USD"},
		},
	}
}

// ParseSceneSettings 从 DB 文本解析；空串则返回默认值。
func ParseSceneSettings(raw string) (SceneSettings, error) {
	if raw == "" {
		return DefaultSceneSettings(), nil
	}
	var s SceneSettings
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		return SceneSettings{}, err
	}
	if s.Slots == nil {
		s.Slots = map[string]SlotState{}
	}
	return s, nil
}

func (s SceneSettings) JSON() (string, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// LayoutResponse 转为 design 文档中的 layout 数组格式。
func (s SceneSettings) LayoutResponse() []map[string]SlotState {
	out := make([]map[string]SlotState, 0, len(s.Slots))
	for name, slot := range s.Slots {
		out = append(out, map[string]SlotState{name: slot})
	}
	return out
}
