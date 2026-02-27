package models

// SenseData 是 Unity 发送过来的感知数据
type SenseData struct {
	NPCName       string  `json:"npc_name"`
	Context       string  `json:"context"`
	PlayerMessage string  `json:"player_message"`
	MapMinX       float32 `json:"map_min_x"`
	MapMaxX       float32 `json:"map_max_x"`
	MapMinY       float32 `json:"map_min_y"`
	MapMaxY       float32 `json:"map_max_y"`
}

// OllamaRequest 是发给 Ollama API 的请求
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	Format string `json:"format"`
}

// OllamaResponse 是 Ollama API 返回的响应
type OllamaResponse struct {
	Response string `json:"response"`
}

// NPCConfig 是 NPC 的配置信息，从 JSON 文件加载
type NPCConfig struct {
	Name        string   `json:"name"`        // NPC 名字
	Personality string   `json:"personality"` // 性格描述
	Model       string   `json:"model"`       // 使用的 AI 模型
	Behaviors   []string `json:"behaviors"`   // 行为准则列表
}
