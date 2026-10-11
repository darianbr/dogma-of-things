## 3. Platform Decision: Home Assistant

### Why Home Assistant is the right choice

| Factor | Assessment |
|--------|------------|
| **Eufy support** | ✅ Built-in integration |
| **Local-first** | ✅ All processing can stay on-premise |
| **Voice assistant (Assist)** | ✅ Built-in with custom wake word support |
| **LLM integration** | ✅ OpenAI-compatible API support (LM Studio compatible) |
| **Extensibility** | ✅ HACS ecosystem, custom integrations, automations |
| **Community** | ✅ Large, active community |
| **Hardware flexibility** | ✅ Runs on Raspberry Pi, x86, Docker, etc. |

### Home Assistant Voice Stack (Assist)

Home Assistant has a complete voice assistant pipeline called **Assist**:

```
Microphone → Wake Word Detection → Speech-to-Text → Intent Processing → Action → Text-to-Speech → Speaker
```

Each stage is pluggable:

| Stage | Local Option | Cloud Option |
|-------|-------------|--------------|
| **Wake Word** | openWakeWord, microWakeWord | — |
| **Speech-to-Text** | Whisper (local), Faster-Whisper | Google STT, Haas Assist |
| **Intent Processing** | Hass (rule-based), **LLM (LM Studio)** | OpenAI, Anthropic |
| **Text-to-Speech** | Piper (local, excellent quality) | Google TTS, ElevenLabs |

### Key Insight: You can plug LM Studio into the intent processing stage

Home Assistant supports OpenAI-compatible APIs for conversational Assist. Since LM Studio exposes an OpenAI-compatible endpoint (`http://localhost:1235/v1`), you can use it as the LLM backend for voice commands.

---
