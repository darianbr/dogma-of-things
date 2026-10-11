## 1. Project Overview

### Goal
Build a local smart home control service that:
- Uses natural language commands to control smart home devices
- Integrates with a local LLM running via LM Studio
- Leverages Home Assistant as the device management hub
- Supports custom wake word activation ("Hey Dibbles...")
- Works primarily with Eufy devices and Google Home Mini speakers

### Core Requirements
- **Privacy-first**: All processing stays local (no cloud dependencies)
- **Natural language**: Conversational device control via LLM
- **Custom wake word**: Personal activation phrase instead of "Ok Google"
- **Eufy support**: Primary device ecosystem
- **Google Mini reuse**: Repurpose existing hardware where possible

---
