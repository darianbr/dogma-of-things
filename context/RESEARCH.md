# DoT (Dogma of Things) — Research Document

> **Date:** 2026-06-20
> **Author:** Darian
> **Status:** Research & Planning

---

## Table of Contents

- [Project Overview](./research/RESEARCH-01-Project-Overview.md) — Goal, core requirements, and design principles (privacy-first, local LLM, custom wake word, Eufy devices).

- [Device Inventory & Compatibility](./research/RESEARCH-02-Device-Inventory-&-Compatibility.md) — Eufy integration status, supported device types, configuration, limitations, and Google Home Mini compatibility assessment.

- [Platform Decision: Home Assistant](./research/RESEARCH-03-Platform-Decision-Home-Assistant.md) — Why HA was chosen, the Assist voice stack, pluggable pipeline stages, and LM Studio integration point.

- [LM Studio Integration](./research/RESEARCH-04-LM-Studio-Integration.md) — Three integration approaches (Assist personality, custom service, function calling), model recommendations, and trade-offs.

- [LLM Prompt Strategy: Device & Intent Mapping](./research/RESEARCH-05-LLM-Prompt-Strategy-Device-&-Intent-Mapping.md) — Frontloaded mapping approach to eliminate entity hallucination: intent normalization, device aliases, ambiguity rules, context-aware resolution, and risk impact analysis.

- [Voice Activation & Custom Wake Words](./research/RESEARCH-06-Voice-Activation-&-Custom-Wake-Words.md) — "Hey Dibbles" wake word training via openWakeWord, Colab notebook workflow, hardware options (ATOM Echo, ESP32-S3-BOX-3, Voice PE), and phased rollout plan.

- [Google Home Mini Repurposing](./research/RESEARCH-07-Google-Home-Mini-Repurposing.md) — What Google Minis can still do (TTS output via Cast, fallback Google Assistant, music) and what they cannot (wake word change, custom firmware).

- [Recommended Architecture](./research/RESEARCH-08-Recommended-Architecture.md) — Full system diagram, data flow for a voice command (wake word → STT → LLM → action → TTS → speaker), and component roles.

- [Implementation Roadmap](./research/RESEARCH-09-Implementation-Roadmap.md) — Five-phase plan: foundation, voice pipeline, LLM integration (with mapping tasks), hardware, and advanced features.

- [Risks & Open Questions](./research/RESEARCH-10-Risks-&-Open-Questions.md) — Technical risks with mitigations, open questions, and privacy considerations.

- [References](./research/RESEARCH-11-References.md) — Home Assistant docs, openWakeWord, Piper, Whisper, LM Studio, Wyoming Protocol, and hardware purchase links.

- [Quick Decision Summary](./research/RESEARCH-12-Quick-Decision-Summary.md) — One-page decision table: platform, LLM backend, prompt strategy, device resolution, wake word, hardware, STT, TTS, model.

- [Sample System Prompt for LLM](./research/RESEARCH-13-Sample-System-Prompt-for-LLM.md) — Full system prompt template with structured intent mapping, device aliases, dynamic registry injection, rules, function schema, and worked examples.

---

*This document has been split into individual section files for easier reference.*
