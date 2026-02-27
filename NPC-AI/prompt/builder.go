package prompt

import (
	"fmt"
	"npc-brain/models"
	"strings"
)

// BuildPrompt 根据 NPC 配置和感知数据动态生成 prompt
func BuildPrompt(npc *models.NPCConfig, sense *models.SenseData) string {
	// 构建行为准则（由你的 json 配置文件动态传入）
	var behaviorsText strings.Builder
	for i, b := range npc.Behaviors {
		behaviorsText.WriteString(fmt.Sprintf("%d. %s\n", i+1, b))
	}

	// 动态处理玩家的话
	playerMsg := sense.PlayerMessage
	if playerMsg == "" {
		playerMsg = "无"
	}

	return fmt.Sprintf(`你是一个2D开放世界游戏里的NPC，名字叫【%s】。
性格设定：%s

【你当前接收到的环境情报】
%s

【玩家对你说的话】
"%s"

【地图边界限制（极其重要，违背将导致游戏崩溃！）】
你当前所在的地图是一个封闭的房间。
X 坐标范围：%.1f 到 %.1f
Y 坐标范围：%.1f 到 %.1f
你的物理躯体绝对不能超出这个范围！

【你的动态行为准则】
%s
【action 字段与移动规则（极其重要，必须严格遵守！）】
1. 你的移动完全由你输出的 target_x 和 target_y 决定。
2. 移动 (move)：如果你想漫游瞎逛，或者想靠近某个物品/NPC，action 填 "move"。你必须自己决定一个地图范围内的新坐标，并填入 target_x 和 target_y。
3. 待机 (idle)：如果你正在和玩家对话，或者明确想留在原地，action 填 "idle"。此时你必须将 target_x 和 target_y 填为你当前的坐标。

【NPC 互动与群聊规则（极其重要！）】
1. 你的环境情报中，不仅有物体的位置，还有其他 NPC 此时此刻【头顶气泡正在说的话】。
2. 如果你在环境情报里看到其他 NPC 刚刚说了某些话，你【必须】仔细倾听，并在你的 json 'say' 字段中，根据你的性格，给出符合你行为准则的直接接话、反驳或回应！你们是在进行一场真实的群聊社交！不要自说自话！

请严格输出以下 JSON 格式（不要输出任何其他内容）：
{
  "npc_name": "%s",
  "player_input": "原封不动重复玩家说的话，没有则填无",
  "thought": "指明你的推理过程，包括你看到了谁在说话，你打算怎么回应他",
  "action": "move 或者 idle",
  "target_x": 浮点数 (想去的目标X坐标，如果不动则填当前X坐标),
  "target_y": 浮点数 (想去的目标Y坐标，如果不动则填当前Y坐标),
  "say": "你想说的话，如果有其他 NPC 在说话，必须在这里直接回应他们！"
}`,
		npc.Name,
		npc.Personality,
		sense.Context,
		playerMsg, // 使用处理过的 playerMsg
		sense.MapMinX, sense.MapMaxX,
		sense.MapMinY, sense.MapMaxY,
		behaviorsText.String(),
		npc.Name,
	)

}
