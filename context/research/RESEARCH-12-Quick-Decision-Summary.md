## Appendix A: Quick Decision Summary

| Decision | Choice | Rationale |
|----------|--------|-----------|
| **Platform** | Home Assistant | Best Eufy support, local-first, custom wake words |
| **LLM Backend** | LM Studio | OpenAI-compatible, local, flexible model choice |
| **Prompt Strategy** | Frontloaded mappings | Eliminates entity hallucination; structured alias + intent tables |
| **Device Resolution** | Hybrid (static aliases + dynamic registry) | Aliases capture natural speech; registry stays current |
| **Wake Word** | "Hey Dibbles" via openWakeWord | Fully customizable, works with HA Assist |
| **Voice Hardware** | ATOM Echo / HA Voice PE | Google Mini can't be repurposed for wake words |
| **Google Mini Role** | TTS output speakers only | Can't change wake word, but good for audio output |
| **STT Engine** | Whisper (local) | Privacy-first, good accuracy |
| **TTS Engine** | Piper (local) | Natural-sounding, fast, local |
| **Model** | Qwen 2.5 7B or Llama 3.1 8B | Good function calling, runs on modest hardware |

---
