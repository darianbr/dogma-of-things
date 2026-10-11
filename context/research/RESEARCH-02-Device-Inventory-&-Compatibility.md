## 2. Device Inventory & Compatibility

### 2.1 Eufy Devices

**Integration Status:** ✅ Supported (Legacy integration)

Home Assistant has a built-in `eufy` integration for EufyHome devices. Key facts:

| Aspect | Detail |
|--------|--------|
| **Integration type** | Community-maintained, Legacy quality scale |
| **Supported device types** | Light, Switch |
| **IoT class** | Local Polling |
| **Setup** | Auto-discovery via `configuration.yaml` with EufyHome credentials |
| **Alternative** | Static device configuration with IP + access tokens |

**Basic configuration:**
```yaml
eufy:
  username: YOUR_EUFY_EMAIL
  password: YOUR_EUFY_PASSWORD
```

**Limitations to be aware of:**
- Marked as a **Legacy integration** — basic functionality works but may lack advanced features
- Supports Lights and Switches primarily
- Some newer Eufy devices (security cameras, vacuums, doorbells) may need separate integrations (e.g., `eufy_security` for cameras/doorbells)
- Uses the EufyHome cloud API for authentication, then communicates locally where possible

**Recommendation:** Test all your Eufy devices after initial HA setup to confirm discovery. For any missed devices, investigate the `eufy_security` HACS integration.

---

### 2.2 Google Home Mini Speakers

**Integration Status:** ⚠️ Partially supported / Repurposing required

| Aspect | Detail |
|--------|--------|
| **Native HA integration** | Google Assistant integration exists (controls HA devices *via* Google) |
| **Local control** | Supports "Local Fulfillment" for faster response |
| **Custom wake word** | ❌ **NOT POSSIBLE** — Google locks wake words to "Ok Google" / "Hey Google" |
| **Repurposing potential** | ⚠️ Limited — firmware is locked down |

#### Can you change the Google Mini wake word?

**Short answer: No.**

Google Home Mini devices have the wake word detection **hardcoded in firmware**. You cannot:
- Change "Ok Google" to a custom phrase
- Replace the wake word model
- Flash custom firmware (no community firmware exists for Google Mini)
- Disable Google Assistant and replace it with another voice assistant

The wake word processing happens on-device in a dedicated low-power processor that is not accessible or reprogrammable.

#### Can Google Home Mini respond to Home Assistant commands?

**Yes, but with caveats:**

1. **Google Assistant Integration** — You can link Home Assistant to Google Assistant, so "Ok Google, turn on the living room light" works. However:
   - This routes through Google's cloud (privacy concern)
   - You still need to say "Ok Google"
   - Google processes the voice, not your local LLM

2. **Local Fulfillment** — HA supports local command execution via mDNS, reducing latency, but the voice recognition still happens on Google's servers.

**Bottom line:** Google Home Minis are **not suitable** as voice input devices for a custom wake word + local LLM system. They can still be used as **output devices** (speakers for TTS responses) via the Cast/Chromecast integration in HA.

---
