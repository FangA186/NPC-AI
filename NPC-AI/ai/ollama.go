package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"npc-brain/models"
)

const ollamaURL = "http://localhost:11434/api/generate"

// Think 调用 Ollama API，返回 AI 的原始 JSON 响应字符串
func Think(model string, prompt string) (string, error) {
	reqBody := models.OllamaRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
		Format: "json",
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	fmt.Println("\n--- 收到 Unity 情报，思考中 ---")

	resp, err := http.Post(ollamaURL, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", fmt.Errorf("Ollama API 调用失败: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取 Ollama 响应失败: %w", err)
	}

	var ollamaResp models.OllamaResponse
	if err := json.Unmarshal(bodyBytes, &ollamaResp); err != nil {
		return "", fmt.Errorf("解析 Ollama 响应失败: %w", err)
	}

	fmt.Printf("AI 决定: %s\n", ollamaResp.Response)
	return ollamaResp.Response, nil
}
