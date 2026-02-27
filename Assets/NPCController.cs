using UnityEngine;
using UnityEngine.UI;
using UnityEngine.Networking;
using System.Collections;

[System.Serializable]
public class SenseData
{
    public string npc_name;
    public string context;
    public string player_message;
    public float map_min_x;
    public float map_max_x;
    public float map_min_y;
    public float map_max_y;
}

[System.Serializable]
public class AIResponse
{
    public string npc_name;
    public string player_input;
    public string thought;
    public string action;
    public float target_x;
    public float target_y;
    public string say;
}

public class NPCController : MonoBehaviour
{
    public string npcName = "神秘法师";
    public Transform player;
    public Text chatBubbleText;
    public float speed = 1.0f;
    public float thinkInterval = 8f;

    [Header("感知与地图限制")]
    public Vector2 visionSize = new Vector2(10f, 6f);
    public Vector2 mapMinBounds = new Vector2(-5f, -3f); // 地图左下角限制
    public Vector2 mapMaxBounds = new Vector2(5f, 3f);   // 地图右上角限制
    public float scanRadius = 3.0f; // 【新增】雷达扫描半径

    [Header("玩家交互 UI")]
    public GameObject playerUICanvas;
    public InputField chatInput;

    private Vector2 targetPosition;
    private float timer = 0f;
    private bool isThinking = false;

    void Start()
    {
        targetPosition = transform.position;
        if (playerUICanvas != null) playerUICanvas.SetActive(false);
        StartCoroutine(ThinkRoutine(""));
    }

    void Update()
    {
        if (playerUICanvas != null && playerUICanvas.activeSelf)
        {
            if (Input.GetKeyDown(KeyCode.Return) || Input.GetKeyDown(KeyCode.KeypadEnter))
            {
                SendPlayerMessage();
            }
            return;
        }

        transform.position = Vector2.MoveTowards(transform.position, targetPosition, speed * Time.deltaTime);

        // ================= 紧急修复：双重边界限制 =================
        // 强行把目前的坐标约束在地图内，防止物理引擎或浮点错误把它挤出去
        float currentX = Mathf.Clamp(transform.position.x, mapMinBounds.x, mapMaxBounds.x);
        float currentY = Mathf.Clamp(transform.position.y, mapMinBounds.y, mapMaxBounds.y);
        transform.position = new Vector2(currentX, currentY);
        // ==========================================================

        timer += Time.deltaTime;
        if (timer >= thinkInterval && !isThinking)
        {
            timer = 0f;
            StartCoroutine(ThinkRoutine(""));
        }

        if (player != null && Vector2.Distance(transform.position, player.position) < 2.0f)
        {
            if (Input.GetKeyDown(KeyCode.E))
            {
                playerUICanvas.SetActive(true);
                chatInput.text = "";
                chatInput.ActivateInputField();
            }
        }
    }

    public void SendPlayerMessage()
    {
        string msg = chatInput.text;
        if (!string.IsNullOrEmpty(msg) && !isThinking)
        {
            playerUICanvas.SetActive(false);
            StartCoroutine(ThinkRoutine(msg));
        }
    }

    IEnumerator ThinkRoutine(string playerMsg)
    {
        isThinking = true;

        bool canSeePlayer = false;
        float dist = 0f;

        if (player != null)
        {
            float distanceX = Mathf.Abs(player.position.x - transform.position.x);
            float distanceY = Mathf.Abs(player.position.y - transform.position.y);
            if (distanceX <= visionSize.x / 2f && distanceY <= visionSize.y / 2f)
            {
                canSeePlayer = true;
                dist = Vector2.Distance(transform.position, player.position);
            }
        }

        string envContext = $"我当前在坐标 ({transform.position.x:F1}, {transform.position.y:F1})。\n";
        if (canSeePlayer) envContext += $"看到玩家进入了视线！距离我 {dist:F1} 米。\n";
        else envContext += "视线内没有玩家。\n";

        // ================= 核心升级：动态物理雷达扫描 =================
        Collider2D[] colliders = Physics2D.OverlapCircleAll(transform.position, scanRadius);
        string itemsFound = "";
        foreach (Collider2D col in colliders)
        {
            // 如果扫到的是其他 NPC 或可交互物品
            if (col.gameObject != this.gameObject && col.transform != player)
            {
                if (col.CompareTag("NPC"))
                {
                    NPCController otherNPC = col.GetComponent<NPCController>();
                    string otherName = otherNPC != null ? otherNPC.npcName : col.name;
                    string whatTheyAreSaying = "";
                    
                    if (otherNPC != null && otherNPC.chatBubbleText != null && !string.IsNullOrEmpty(otherNPC.chatBubbleText.text))
                    {
                        whatTheyAreSaying = $"对方头上正冒出泡泡说：“{otherNPC.chatBubbleText.text}”；";
                    }
                    
                    itemsFound += $"【NPC】[{otherName}] 在坐标({col.transform.position.x:F1}, {col.transform.position.y:F1})。{whatTheyAreSaying}";
                }
                else if (col.CompareTag("Interactable"))
                {
                    itemsFound += $"物品 [{col.name}] 在坐标({col.transform.position.x:F1}, {col.transform.position.y:F1})；";
                }
            }
        }

        if (itemsFound != "") envContext += $"我周围 {scanRadius} 米内发现了：" + itemsFound;
        else envContext += $"我周围 {scanRadius} 米内没有任何特殊物品或其他角色。";
        // ==========================================================

        SenseData sense = new SenseData
        {
            npc_name = npcName,
            context = envContext,
            player_message = playerMsg,
            map_min_x = mapMinBounds.x,
            map_max_x = mapMaxBounds.x,
            map_min_y = mapMinBounds.y,
            map_max_y = mapMaxBounds.y
        };

        string jsonData = JsonUtility.ToJson(sense);

        using (UnityWebRequest request = new UnityWebRequest("http://localhost:8080/think", "POST"))
        {
            byte[] bodyRaw = System.Text.Encoding.UTF8.GetBytes(jsonData);
            request.uploadHandler = new UploadHandlerRaw(bodyRaw);
            request.downloadHandler = new DownloadHandlerBuffer();
            request.SetRequestHeader("Content-Type", "application/json");

            yield return request.SendWebRequest();

            if (request.result == UnityWebRequest.Result.Success)
            {
                AIResponse aiDecision = JsonUtility.FromJson<AIResponse>(request.downloadHandler.text);
                if (aiDecision != null)
                {
                    // 打印 AI 的思考过程和决定的动作
                    Debug.Log($"【{npcName}内心】: {aiDecision.thought} | 决定动作: {aiDecision.action}");
                    
                    if (!string.IsNullOrEmpty(aiDecision.say)) Speak(aiDecision.say);

                    // ================= 核心：完全由 AI 控制坐标 =================
                    // 算出 AI 给的新坐标距离当前坐标有多远
                    float distanceToNewTarget = Vector2.Distance(
                        new Vector2(aiDecision.target_x, aiDecision.target_y), 
                        transform.position
                    );

                    if (distanceToNewTarget > 0.1f)
                    {
                        // 如果坐标变化明显，说明 AI 想移动（不管是瞎逛还是有目的地）
                        float safeX = Mathf.Clamp(aiDecision.target_x, mapMinBounds.x, mapMaxBounds.x);
                        float safeY = Mathf.Clamp(aiDecision.target_y, mapMinBounds.y, mapMaxBounds.y);
                        targetPosition = new Vector2(safeX, safeY);
                        Debug.Log($"[大脑决策-移动] 目标点: {targetPosition}");
                    }
                    else
                    {
                        // 如果坐标几乎没变，说明 AI 选择待在原地
                        targetPosition = transform.position;
                        Debug.Log("[大脑决策-待机] 动作锁定为 Idle，保持当前位置");
                    }
                    // ==========================================================
                }
            }
            else
            {
                Debug.LogError($"【{npcName}】请求后端失败: {request.error}");
            }
        }
        
        // 无论成功还是失败，都必须把 isThinking 改回 false，否则它就永远卡死了！
        isThinking = false;
    }

    public void Speak(string content)
    {
        if (chatBubbleText != null) chatBubbleText.text = content;
    }

    private void OnDrawGizmosSelected()
    {
        // 画出视野框 (红)
        Gizmos.color = new Color(1, 0, 0, 0.3f);
        Gizmos.DrawCube(transform.position, new Vector3(visionSize.x, visionSize.y, 1f));
        Gizmos.color = Color.red;
        Gizmos.DrawWireCube(transform.position, new Vector3(visionSize.x, visionSize.y, 1f));

        // 画出地图边界 (黄)
        Gizmos.color = Color.yellow;
        Vector3 mapCenter = new Vector3((mapMinBounds.x + mapMaxBounds.x) / 2f, (mapMinBounds.y + mapMaxBounds.y) / 2f, 0);
        Vector3 mapSize = new Vector3(mapMaxBounds.x - mapMinBounds.x, mapMaxBounds.y - mapMinBounds.y, 1f);
        Gizmos.DrawWireCube(mapCenter, mapSize);

        // 【新增】：画出 3 米物理雷达扫描圈（绿色圆圈）
        Gizmos.color = Color.green;
        Gizmos.DrawWireSphere(transform.position, scanRadius);
    }
}