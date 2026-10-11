## 6. Voice Activation & Custom Wake Words

### 6.1 The "Hey Dibbles" Question

**Can you create a custom wake word like "Hey Dibbles"?**

**Yes — but NOT on Google Home Mini devices.**

Here's the breakdown:

#### ❌ Google Home Mini
- Wake word is **firmware-locked** to "Ok Google" / "Hey Google"
- No way to change, replace, or customize
- No custom firmware community exists

#### ✅ Home Assistant (with openWakeWord)
Home Assistant has **first-class custom wake word support** via the openWakeWord project:

**How it works:**
1. Install the **openWakeWord** addon in HA
2. Use the [official wake word training Colab notebook](https://colab.research.google.com/drive/1q1oe2zOyZp7UsB3jJiQ1IFn8z5YfjwEb) to train a model for "Hey Dibbles"
3. The training process:
   - Uses Piper TTS to generate thousands of synthetic voice samples of your wake word
   - Augments with room acoustics, background noise, different speeds
   - Trains a TensorFlow model (~1 hour on free Colab)
   - Produces a `.tflite` model file
4. Drop the model into `/share/openwakeword/` in HA
5. Assign it to your voice assistant

**Training process summary:**
```
"Hey Dibbles" → Piper generates samples → Augmentation → ML Training → .tflite model → HA
```

**Constraints:**
- Currently **English only** (due to TTS model availability)
- Works best with 3–4 syllable phrases
- Should be uncommon enough to avoid false triggers
- "Hey Dibbles" is a good candidate (unique, 3 syllables)

### 6.2 Wake Word Hardware Options

Since Google Mini can't be repurposed for wake word detection, here are alternatives:

| Device | Cost | Wake Word Engine | Notes |
|--------|------|-----------------|-------|
| **M5Stack ATOM Echo** | ~$13 | openWakeWord (server-side) | Tiny, official HA tutorial exists |
| **ESP32-S3-BOX-3** | ~$30–40 | microWakeWord (on-device) | Has display, better mic array |
| **Home Assistant Voice PE** | ~$130 | microWakeWord (on-device) | Official HA hardware, best experience |
| **Android phone** | Free (if owned) | microWakeWord (on-device) | HA Companion App supports wake words |
| **Linux PC / HA Server** | Free (if owned) | openWakeWord (server-side) | Use built-in mic + Linux Voice Assistant |
| **Raspberry Pi + mic hat** | ~$20–30 | openWakeWord (server-side) | Flexible DIY option |

### 6.3 Recommended Wake Word Setup

For the best "Hey Dibbles" experience:

1. **Short-term (quick test):** Use your Android phone with HA Companion App + microWakeWord
   - Train "Hey Dibbles" model
   - Test the full pipeline before buying hardware

2. **Medium-term (per-room coverage):** Buy 2–3 ATOM Echo devices (~$40 total)
   - Place in key rooms
   - Each streams audio to HA for wake word detection

3. **Long-term (best experience):** Home Assistant Voice PE or ESP32-S3-BOX-3
   - On-device wake word detection (lower latency)
   - Better microphone arrays
   - Display for visual feedback

---
