package npc

import (
	"encoding/json"
	"fmt"
	"npc-brain/models"
	"os"
	"path/filepath"
	"strings"
)

// NPCRegistry 存储所有已加载的 NPC 配置
var NPCRegistry = make(map[string]*models.NPCConfig)

// LoadAllNPCs 从 config/npcs/ 目录加载所有 NPC 配置文件
func LoadAllNPCs(configDir string) error {
	files, err := os.ReadDir(configDir)
	if err != nil {
		return fmt.Errorf("无法读取 NPC 配置目录 %s: %w", configDir, err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		cfg, err := loadNPCFile(filepath.Join(configDir, file.Name()))
		if err != nil {
			fmt.Printf("⚠️ 加载 NPC 配置 %s 失败: %v\n", file.Name(), err)
			continue
		}

		NPCRegistry[cfg.Name] = cfg
		fmt.Printf("✅ 已加载 NPC: %s (模型: %s)\n", cfg.Name, cfg.Model)
	}

	if len(NPCRegistry) == 0 {
		return fmt.Errorf("没有找到任何有效的 NPC 配置文件")
	}

	return nil
}

// GetNPC 根据名字查找已加载的 NPC 配置
func GetNPC(name string) (*models.NPCConfig, bool) {
	cfg, ok := NPCRegistry[name]
	return cfg, ok
}

// loadNPCFile 读取并解析单个 NPC JSON 配置文件
func loadNPCFile(path string) (*models.NPCConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	var cfg models.NPCConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}

	if cfg.Name == "" {
		return nil, fmt.Errorf("NPC 配置缺少 name 字段")
	}

	return &cfg, nil
}
