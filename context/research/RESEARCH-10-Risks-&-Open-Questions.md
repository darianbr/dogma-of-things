## 10. Risks & Open Questions

### Technical Risks

| Risk | Impact | Mitigation |
|------|--------|-----------|
| Eufy legacy integration lacks device support | Medium | Test early; use `eufy_security` HACS integration as fallback |
| LLM hallucinates device names or actions | **Low** | Frontloaded device mapping + dynamic registry; model selects from provided list |
| Intent misinterpretation | **Low** | Explicit intent mapping covers common verb variants |
| Wake word false positives/negatives | Medium | Iterate on "Hey Dibbles" model; adjust sensitivity thresholds |
| Whisper STT accuracy in noisy environments | Medium | Use better microphone hardware; consider Faster-Whisper large model |
| LM Studio latency on CPU-only hardware | Medium | Use a model sized for your hardware; consider GPU acceleration |

### Open Questions

1. **Which Eufy devices exactly?** — List all devices to confirm integration support
2. **What hardware will run HA + LM Studio?** — CPU/GPU specs determine model choice
3. **How many rooms need voice coverage?** — Determines number of ATOM Echo devices needed
4. **Is Home Assistant Cloud acceptable for 30-day trial?** — Simplifies initial setup significantly
5. **What's the budget for additional hardware?** — ATOM Echo ($13 each) vs. Voice PE ($130)

### Privacy Considerations

- ✅ All voice processing stays local (no cloud STT needed)
- ✅ LLM runs locally via LM Studio (no data leaves your network)
- ✅ Home Assistant is open-source and auditable
- ⚠️ Google Home Mini fallback routes through Google cloud
- ⚠️ Eufy integration may require initial cloud authentication

---
