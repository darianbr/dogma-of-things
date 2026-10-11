## 9. Implementation Roadmap

### Phase 1: Foundation (Week 1–2)
- [ ] Install Home Assistant (HA OS on Raspberry Pi / x86 / Docker)
- [ ] Configure Eufy integration and verify all devices are discovered
- [ ] Install and configure LM Studio
- [ ] Test LM Studio API connectivity from HA
- [ ] Set up basic HA dashboards for manual device control

### Phase 2: Voice Pipeline (Week 3–4)
- [ ] Install openWakeWord addon
- [ ] Train "Hey Dibbles" wake word model (Colab notebook)
- [ ] Deploy model to HA
- [ ] Configure Assist pipeline with:
  - openWakeWord (wake word)
  - Whisper (speech-to-text)
  - Piper (text-to-speech)
- [ ] Test basic voice commands with rule-based intent processing

### Phase 3: LLM Integration (Week 5–6)
- [ ] Configure LM Studio as OpenAI-compatible endpoint in HA
- [ ] Create AI Personality for Assist with LM Studio backend
- [ ] **Build device alias mapping** — inventory all devices, document how you refer to them in conversation
- [ ] **Build intent mapping** — list all verb variants for each action
- [ ] Write system prompt with structured mappings (see Appendix B)
- [ ] Set up dynamic device registry injection from HA
- [ ] Test natural language commands against the mapping
- [ ] Implement function calling for structured device control
- [ ] Add command validation and safety checks
- [ ] Iterate on mappings based on real usage (add missing aliases)

### Phase 4: Hardware & Polish (Week 7–8)
- [ ] Purchase and set up ATOM Echo devices (or alternative hardware)
- [ ] Configure Google Minis as TTS output speakers via Cast
- [ ] Set up room-aware audio routing
- [ ] Create automations and scenes
- [ ] Test end-to-end voice control flow
- [ ] Iterate on prompts and wake word tuning

### Phase 5: Advanced Features (Ongoing)
- [ ] Multi-turn conversations with the LLM
- [ ] Context-aware responses (time of day, presence, etc.)
- [ ] Custom sentences for complex automations
- [ ] Dashboard for monitoring system health
- [ ] Backup and recovery procedures

---
