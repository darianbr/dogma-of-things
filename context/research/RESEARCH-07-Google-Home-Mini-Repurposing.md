## 7. Google Home Mini Repurposing

### What you CAN do with Google Home Minis

Even though they can't be repurposed for custom wake words, they still have value:

#### 1. TTS Output Speakers (Recommended)
Use them as speakers for Home Assistant's voice responses:
- HA has a **Cast integration** for Chromecast devices
- When Assist responds, route the TTS audio to the Google Mini in the current room
- This gives you room-aware voice responses

**Setup:**
```yaml
# In configuration.yaml
media_player:
  - platform: cast
    devices:
      - name: "Living Room Mini"
        host: 192.168.1.XX
```

#### 2. Google Assistant as a Fallback
Keep Google Assistant linked to HA for basic commands:
- "Ok Google, turn on the lights" still works
- Use as a backup when the LLM system is down
- Accept that this routes through Google's cloud

#### 3. Music & Media Playback
Continue using them for music streaming (Spotify, etc.) — this is unaffected by HA integration.

### What you CANNOT do

- ❌ Change the wake word
- ❌ Replace Google Assistant with a local voice assistant
- ❌ Flash custom firmware
- ❌ Access the microphone stream externally (for HA wake word detection)

---
