## 5. LLM Prompt Strategy: Device & Intent Mapping

### 5.1 The Problem

The single biggest failure mode for LLM-based smart home control is **entity resolution** — the LLM guessing wrong device names, inventing entity IDs, or misinterpreting what the user means. Relying on the LLM to "figure out" that "turn on" and "switch on" mean the same thing, or that "floor lamp" and "study lamp" both refer to `light.study_lamp_02`, is asking for intermittent failures.

**Solution:** Frontload the LLM with explicit mappings so it never has to guess.

### 5.2 Two-Layer Mapping Approach

#### Layer 1: Intent Normalization (Static)

Map natural language verbs to canonical actions. This is baked into the system prompt and rarely changes:

```
INTENT MAPPING:
  turn_on → [turn on, switch on, activate, light up, power on, come on, boot up]
  turn_off → [turn off, switch off, deactivate, shut off, power off, kill, stop]
  set_brightness → [dim, brighten, set brightness to, make it brighter/dimmer]
  toggle → [toggle, flip, switch]
  query_status → [is it on, what's the status, check, tell me about]
```

#### Layer 2: Device Aliases (Static + Dynamic)

Map natural language references to actual HA entity IDs. This is the critical layer:

```
DEVICE MAPPING:
  light.study_lamp_02 → aliases: [floor lamp, study lamp, study floor lamp, lamp in study,
                                   the lamp, that lamp (when context=study)]
  light.kitchen_overhead → aliases: [kitchen light, kitchen overhead, kitchen ceiling,
                                      the kitchen, overhead in kitchen]
  switch.garage_door → aliases: [garage door, garage, open the garage, close the garage]
```

**Key insight:** Aliases should capture how you *actually* refer to devices in conversation — not how they're labelled in HA.

### 5.3 Why Structured Format Matters

Structured/tabular mappings are referenced more reliably by LLMs than prose descriptions. The model can scan a table faster than parse a paragraph. The format above (entity → alias list) is compact and unambiguous.

### 5.4 Static vs Dynamic Registry

| Approach | Pros | Cons |
|----------|------|------|
| **Static in system prompt** | Simple, always available, no API call needed | Token cost grows with device count; requires manual update |
| **Dynamic from HA registry** | Always up-to-date; auto-discovers new devices | Extra API call per request; slightly more complex |
| **Hybrid (recommended)** | Best of both: aliases static, entity list dynamic | Slightly more prompt engineering |

**Recommendation:** Use the hybrid approach. Keep intent mappings and aliases static in the system prompt. Inject the actual entity ID list dynamically per request from HA's device registry. This way:
- Adding a device doesn't require prompt surgery
- The LLM picks from what actually exists, not what it memorized
- Token cost stays predictable (aliases are compact)

### 5.5 Ambiguity Resolution

The LLM must have explicit rules for when a request could apply to multiple devices:

```
RULES:
  - If the user's request maps to exactly one device, execute it.
  - If the request could apply to multiple devices, ask for clarification.
  - If the device is not in the mapping, tell the user it's not recognized.
  - Never guess when uncertain.
  - Use conversational context: "the lamp" in a follow-up refers to the last-mentioned lamp.
```

### 5.6 Context-Aware Aliases

Some aliases are only valid in certain contexts. Handle this with conditional aliases:

```
  light.study_lamp_02 → aliases: [
    "floor lamp",           # unambiguous
    "study lamp",           # unambiguous
    "the lamp",             # context-dependent: only if user is in study or last topic was study
    "that lamp",            # context-dependent: same as above
  ]
```

The system prompt should instruct the LLM to use room context (from the user's location or conversation history) to resolve ambiguous aliases.

### 5.7 Impact on Risk Profile

This approach directly mitigates the highest-impact risk:

| Risk | Before | After |
|------|--------|-------|
| LLM hallucinates entity IDs | High — model invents plausible-sounding IDs | Low — model selects from provided list |
| Intent misinterpretation | Medium — model guesses from context | Low — explicit mapping covers common variants |
| Adding new devices | Manual prompt update needed | Dynamic registry auto-includes; only aliases need adding |

---
