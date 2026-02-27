package main

import (
	"fmt"
	"net/http"
	"npc-brain/handler"
	"npc-brain/npc"
)

func main() {
	// 启动时加载所有 NPC 配置
	fmt.Println("🔄 正在加载 NPC 配置...")
	if err := npc.LoadAllNPCs("config/npcs"); err != nil {
		fmt.Printf("❌ 加载 NPC 配置失败: %v\n", err)
		return
	}

	// 注册路由
	http.HandleFunc("/think", handler.BrainHandler)

	fmt.Println("\n🧠 NPC 大脑已就绪！正在监听 8080 端口...")
	http.ListenAndServe(":8080", nil)
}
