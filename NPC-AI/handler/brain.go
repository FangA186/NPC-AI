package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"npc-brain/ai"
	"npc-brain/models"
	"npc-brain/npc"
	"npc-brain/prompt"
)

// BrainHandler 处理来自 Unity 的 NPC 思考请求
func BrainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("有无调用思考?")
	// CORS 支持
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		return
	}

	// 解析 Unity 发来的感知数据
	var sense models.SenseData
	if err := json.NewDecoder(r.Body).Decode(&sense); err != nil {
		http.Error(w, "无法理解 Unity 的数据", http.StatusBadRequest)
		return
	}
	// fmt.Printf("\n--- [%s] 收到的真实环境情报 ---\n%s\n", sense.NPCName, sense.Context)
	// 根据 NPC 名字查找配置
	npcCfg, ok := npc.GetNPC(sense.NPCName)
	if !ok {
		errMsg := fmt.Sprintf("未找到 NPC 配置: %s", sense.NPCName)
		fmt.Println("⚠️", errMsg)
		http.Error(w, errMsg, http.StatusNotFound)
		return
	}

	fmt.Printf("🧠 [%s] 正在思考...\n", npcCfg.Name)

	// 构建 prompt 并调用 AI
	builtPrompt := prompt.BuildPrompt(npcCfg, &sense)
	response, err := ai.Think(npcCfg.Model, builtPrompt)
	if err != nil {
		fmt.Printf("❌ AI 思考出错: %v\n", err)
		http.Error(w, "AI 大脑宕机", http.StatusInternalServerError)
		return
	}

	// 返回 AI 响应
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(response))
}
