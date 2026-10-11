## 4. LM Studio Integration

### 4.1 How LM Studio Fits In

LM Studio runs a local LLM server with an OpenAI-compatible API:

```
Endpoint: http://localhost:1235/v1/chat/completions
```

This means any tool that supports OpenAI API format works with LM Studio out of the box.

### 4.2 Integration Approaches

#### Option A: Home Assistant Assist with LLM Personality (Recommended)

Home Assistant supports creating an "AI Personality" for Assist that uses an OpenAI-compatible backend:

1. Configure LM Studio as an OpenAI-compatible endpoint in HA
2. Create an Assist pipeline with LM Studio as the intent processor
3. The LLM receives transcribed speech and returns device commands

**Pros:**
- Fully integrated into HA's voice stack
- Works with wake words, TTS, and the full Assist pipeline
- No custom code needed

**Cons:**
- LLM handles intent processing, which may be less deterministic than rule-based approaches
- Requires careful prompt engineering for reliable device control

#### Option B: Custom Service Layer

Build a separate Python/Node.js service that:
1. Receives voice commands (from HA or directly)
2. Sends them to LM Studio for processing
3. Calls HA's REST API or WebSocket API to execute commands

**Pros:**
- Full control over the conversation flow
- Can implement complex multi-step reasoning
- Easier to add custom validation and safety checks

**Cons:**
- More development effort
- Maintains a separate service alongside HA

#### Option C: HA Blueprints + LLM Function Calling

Use LM Studio's function calling capabilities (if your model supports it) to generate structured commands:

```json
{
  "tool_calls": [{
    "function": {
      "name": "turn_on_light",
      "arguments": {"entity_id": "light.living_room", "brightness": 80}
    }
  }]
}
```

**Pros:**
- Structured output is reliable and safe
- Easy to validate before execution
- Best of both worlds (natural language + deterministic actions)

**Cons:**
- Requires a model with good function calling support (Llama 3.1, Qwen 2.5, etc.)
- Some smaller models may not support this well

### 4.3 Recommended Model Choices for LM Studio

| Model | Size | Function Calling | Smart Home Suitability |
|-------|------|-----------------|----------------------|
| **Qwen 2.5** | 7B–32B | ✅ Good | Excellent — strong instruction following |
| **Llama 3.1** | 8B–70B | ✅ Good | Excellent — well-tested function calling |
| **Mistral Small** | 22B | ✅ Good | Good — efficient and capable |
| **Phi 3** | 3.8B | ⚠️ Basic | Good for constrained hardware |
| **Gemma 2** | 9B–27B | ⚠️ Basic | Good general purpose |

**Recommendation:** Start with **Qwen 2.5 7B** or **Llama 3.1 8B** for a good balance of capability and resource usage. If you have a GPU with 16GB+ VRAM, consider the 14B–32B variants.

---
