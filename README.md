# NPC-AI: LLM-Driven 2D Open World NPCs

Welcome to **NPC-AI**! This project integrates a Go-based backend server with a Unity 2D environment to create intelligent, autonomous, and conversational NPCs powered by Large Language Models (LLMs).

## Core Features
1. **Fully Autonomous Movement**: NPCs move around the Unity map based *exclusively* on coordinates determined by the LLM (`target_x` and `target_y`), escaping the limitations of traditional random wander scripts.
2. **Context-Aware Perception**: NPCs possess a "radar" (Physics2D OverlapCircle) that allows them to "see" players, interactable items, and other NPCs within a set radius.
3. **Dynamic Social Interaction**: NPCs can "hear" the floating chat bubbles of other NPCs nearby. The AI is strictly prompted to engage in group conversations, allowing them to bicker, trade, and chat organically.
4. **JSON-Based Personalities**: Each NPC has a dedicated `.json` configuration file dictating their name, personality, underlying LLM model (e.g., `qwen2.5:7b`), and specific dynamic behavior rules.
5. **Modular Go Backend**: The backend is cleanly structured with separate packages for models, routing, prompt generation, AI API calls, and NPC configuration loading.

## Architecture
- **Frontend (Unity/C#)**: 
  - `NPCController.cs`: Handles boundary logic, calculates movement distances, scans the environment using 2D colliders, and parses AI decisions.
- **Backend (Go)**:
  - `main.go`: Entry point, starts the HTTP server on port `8080`.
  - `prompt/builder.go`: Injects dynamic sense data, neighboring conversations, and strict coordinate/action rules into the prompt.
  - `ai/ollama.go`: Communicates with the local Ollama service to stream logic.
  - `config/npcs/`: Stores the `.json` schemas for each individual character (e.g., `mage.json`, `merchant.json`).

## Requirements
- Unity (2022+ recommended)
- Go (1.20+)
- [Ollama](https://ollama.ai/) installed locally along with the `qwen2.5:7b` model.

## Setup & Running
1. Start the Ollama local service.
2. In the `NPC-AI/` directory, run the backend server:
   ```bash
   go run .\main.go
   ```
3. Open the `SampleScene` in Unity.
4. Ensure the `Player` tag is assigned to the player character, and the `NPC` tag is assigned to all NPC characters with a `CircleCollider2D`.
5. Press **Play** in the Unity Editor to watch the magic unfold!
