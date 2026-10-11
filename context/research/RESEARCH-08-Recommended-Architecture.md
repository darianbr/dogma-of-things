## 8. Recommended Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                      SMART HOME ARCHITECTURE                    │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────────┐  │
│  │  Voice Input  │    │   Eufy Hub   │    │  Other Devices   │  │
│  │              │    │              │    │                  │  │
│  │ • ATOM Echo  │    │ • Eufy Lights│    │ • Sensors        │  │
│  │ • HA Voice PE│    │ • Eufy Switch│    │ • Thermostats    │  │
│  │ • Android App│    │ • Cameras    │    │ • etc.           │  │
│  │  (wake word) │    │              │    │                  │  │
│  └──────┬───────┘    └──────┬───────┘    └────────┬─────────┘  │
│         │                   │                      │            │
│         ▼                   ▼                      ▼            │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │              HOME ASSISTANT (Core Hub)                   │  │
│  │                                                          │  │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐  │  │
│  │  │ openWakeWord │  │   Whisper   │  │   Piper TTS     │  │  │
│  │  │  (wake word) │  │  (STT)      │  │   (voice out)   │  │  │
│  │  └──────┬──────┘  └──────┬──────┘  └────────▲────────┘  │  │
│  │         │                │                   │            │  │
│  │         ▼                ▼                   │            │  │
│  │  ┌─────────────────────────────────┐         │            │  │
│  │  │    Assist Pipeline (Intent)     │         │            │  │
│  │  │                                 │         │            │  │
│  │  │  Text → LLM → Action → TTS     │         │            │  │
│  │  └──────────────┬──────────────────┘         │            │  │
│  │                 │                             │            │  │
│  └─────────────────┼─────────────────────────────┘            │
│                    │                                          │
│                    ▼                                          │
│  ┌──────────────────────────┐                                 │
│  │     LM STUDIO (LLM)      │                                 │
│  │                          │                                 │
│  │  Model: Qwen 2.5 7B     │                                 │
│  │  API: :1235/v1/chat      │                                 │
│  │  Role: Intent parsing,   │                                 │
│  │         conversation,    │                                 │
│  │         function calling │                                 │
│  └──────────────────────────┘                                 │
│                                                               │
│  ┌──────────────────────────┐                                 │
│  │  Google Home Minis       │                                 │
│  │  (Output Only)           │                                 │
│  │  • TTS playback via Cast │                                 │
│  │  • Music playback        │                                 │
│  │  • Google Assistant      │                                 │
│  │    (fallback only)       │                                 │
│  └──────────────────────────┘                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Data Flow for a Voice Command

```
1. User says: "Hey Dibbles, turn on the living room lights"
2. ATOM Echo detects voice → streams audio to HA
3. openWakeWord detects "Hey Dibbles" → triggers Assist
4. Whisper transcribes: "turn on the living room lights"
5. HA sends text to LM Studio via OpenAI API
6. LLM returns structured action: {service: "light.turn_on", entity_id: "light.living_room"}
7. HA validates and executes the action
8. Piper generates TTS: "Okay, I've turned on the living room lights"
9. Audio is cast to the Google Mini in the living room
```

---
