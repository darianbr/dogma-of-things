# DoT (Dogma of Things) — Implementation Plan

> **Date:** 2026-06-27
> **Status:** Planning
> **Stack:** Go (Dogma engine) + TypeScript/React (UI) + SQLite (projections)

---

## Overview

**DoT (Dogma of Things)** — an IoT smart home management system built on the Dogma framework (event sourcing, CQRS, process managers) with a TypeScript/React dashboard. The system manages devices, rooms, scenes, automations, and voice commands — integrating with Home Assistant, MQTT, LM Studio (LLM), and external services.

---

## Phase 1: Project Scaffolding & Foundation

### 1.1 Go Backend Structure
- [ ] Initialize Go module (`github.com/darian/dot`)
- [ ] Create directory structure:
  ```
  cmd/server/           — main entry point
  domain/               — aggregates + processes
  messages/
    commands/           — all command types
    events/             — all event types
    internal/           — shared types (device types, failure reasons)
  integrations/         — external system handlers
  ui/
    api/                — HTTP API handlers (JSON endpoints)
    projections/        — SQL projections + schema
    websocket/          — real-time event streaming
  ```
- [ ] Create `app.go` — Dogma application definition wiring all handlers
- [ ] Create `cmd/server/main.go` — engine bootstrap, HTTP server, WebSocket upgrade

### 1.2 TypeScript/React Frontend Structure
- [ ] Initialize Vite + React + TypeScript project
- [ ] Install dependencies: `@tanstack/react-query`, `zod`, `lucide-react`, `tailwindcss`
- [ ] Create directory structure:
  ```
  src/
    components/
      dashboard/
      devices/
      scenes/
      automations/
      voice/
      common/
    hooks/
    api/
    types/
    pages/
    App.tsx
    main.tsx
  ```
- [ ] Set up API client with WebSocket support
- [ ] Set up TypeScript type definitions mirroring Go message types

### 1.3 Shared Type Contract
- [ ] Define all command and event TypeScript interfaces
- [ ] Define projection query response types
- [ ] Create API client methods for each command dispatch

---

## Phase 2: Core Aggregates & Messages

### 2.1 Device Aggregate
- [ ] `messages/internal/types.go` — DeviceType, SubType, Protocol, DeviceStatus enums
- [ ] `messages/commands/device.go` — RegisterDevice, UpdateDeviceConfig, SetDeviceState, ReportDeviceState, RetireDevice
- [ ] `messages/events/device.go` — DeviceRegistered, DeviceConfigUpdated, DeviceStateChanged, DeviceStateReported, DeviceRetired, DeviceOffline, DeviceBackOnline
- [ ] `domain/device.go` — Aggregate implementation with state transitions and validation
- [ ] `domain/device_test.go` — Unit tests for all command scenarios
- [ ] TypeScript: `src/types/devices.ts` — Mirror types

### 2.2 Room Aggregate
- [ ] `messages/commands/room.go` — CreateRoom, UpdateRoom, AssignDeviceToRoom, RemoveDeviceFromRoom, DeleteRoom
- [ ] `messages/events/room.go` — RoomCreated, RoomUpdated, DeviceAssigned, DeviceRemoved, RoomDeleted
- [ ] `domain/room.go` — Aggregate implementation
- [ ] `domain/room_test.go` — Unit tests
- [ ] TypeScript: `src/types/rooms.ts`

### 2.3 Scene Aggregate
- [ ] `messages/commands/scene.go` — CreateScene, UpdateScene, ActivateScene, DeactivateScene, DeleteScene
- [ ] `messages/events/scene.go` — SceneCreated, SceneUpdated, SceneActivated, SceneDeactivated, SceneDeleted
- [ ] `domain/scene.go` — Aggregate implementation
- [ ] `domain/scene_test.go` — Unit tests
- [ ] TypeScript: `src/types/scenes.ts`

### 2.4 AutomationRule Aggregate
- [ ] `messages/commands/automation.go` — CreateAutomation, UpdateAutomation, EnableAutomation, DisableAutomation, DeleteAutomation
- [ ] `messages/events/automation.go` — AutomationCreated, AutomationUpdated, AutomationEnabled, AutomationDisabled, AutomationDeleted, AutomationTriggered
- [ ] `domain/automation.go` — Aggregate with Trigger, Condition, Action types
- [ ] `domain/automation_test.go` — Unit tests
- [ ] TypeScript: `src/types/automations.ts`

### 2.5 User Aggregate
- [ ] `messages/commands/user.go` — RegisterUser, UpdateUserPreferences, SetVoiceProfile
- [ ] `messages/events/user.go` — UserRegistered, PreferencesUpdated, VoiceProfileSet
- [ ] `domain/user.go` — Aggregate implementation
- [ ] `domain/user_test.go` — Unit tests
- [ ] TypeScript: `src/types/users.ts`

---

## Phase 3: Process Managers

### 3.1 SceneActivationProcess
- [ ] `messages/commands/scene_activation.go` — Internal commands for device state application
- [ ] `messages/events/scene_activation.go` — SceneActivationStarted, DeviceStateApplied, SceneActivationCompleted, SceneActivationPartiallyFailed
- [ ] `domain/scene_activation.go` — Process: iterate devices, execute SetDeviceState, handle timeouts with retries (3x), record partial failures
- [ ] `domain/scene_activation_test.go` — Test: success, partial failure, timeout, retry

### 3.2 AutomationEvaluationProcess
- [ ] `messages/events/automation_evaluation.go` — AutomationTriggered, ConditionsEvaluated, ActionsExecuted, CooldownStarted
- [ ] `domain/automation_evaluation.go` — Process: watch for trigger events, evaluate conditions, check cooldown, execute actions, start cooldown deadline
- [ ] `domain/automation_evaluation_test.go` — Test: match/no-match, cooldown skip, multi-action

### 3.3 VoiceCommandProcess
- [ ] `messages/commands/voice.go` — VoiceCommandReceived, VoiceClarificationResponse
- [ ] `messages/events/voice.go` — IntentParsed, ClarificationRequested, CommandExecutedFromVoice, VoiceResponseDelivered, VoiceCommandFailed
- [ ] `messages/events/llm.go` — IntentParsed (LLM response), ResponseGenerated
- [ ] `domain/voice_command.go` — Process: receive text → LLM parse → confidence check → execute or clarify → TTS response
- [ ] `domain/voice_command_test.go` — Test: clear intent, low confidence, ambiguous device, LLM failure

### 3.4 DeviceProvisioningProcess
- [ ] `messages/events/provisioning.go` — ProvisioningStarted, DeviceDiscovered, CapabilitiesRead, RoomAssigned, TestCommandSent, ProvisioningCompleted, ProvisioningFailed
- [ ] `domain/provisioning.go` — Process: discover → read capabilities → assign room → verify connectivity → complete or fail
- [ ] `domain/provisioning_test.go` — Test: success path, discovery failure, verification retry

---

## Phase 4: Integrations

### 4.1 HomeAssistantIntegration
- [ ] `messages/commands/homeassistant.go` — ExecuteHAService, QueryHAState, CallHAScript
- [ ] `messages/events/homeassistant.go` — HAServiceExecuted, HAServiceFailed, HAStateQueried
- [ ] `integrations/homeassistant.go` — HTTP client to HA REST API, service execution, state queries
- [ ] `integrations/homeassistant_test.go` — Mocked HTTP tests

### 4.2 MQTTIntegration
- [ ] `messages/commands/mqtt.go` — PublishMQTTMessage, SubscribeToTopic
- [ ] `messages/events/mqtt.go` — MQTTMessagePublished, MQTTMessageReceived
- [ ] `integrations/mqtt.go` — MQTT client wrapper, publish/subscribe handlers
- [ ] `integrations/mqtt_test.go` — Mocked broker tests

### 4.3 LLMIntegration
- [ ] `messages/commands/llm.go` — ParseIntent, GenerateResponse, AskClarification
- [ ] `messages/events/llm.go` — IntentParsed, ResponseGenerated, LLMRequestFailed
- [ ] `integrations/llm.go` — OpenAI-compatible API client (LM Studio), prompt templating with device registry injection
- [ ] `integrations/llm_test.go` — Mocked LLM response tests

### 4.4 WeatherIntegration
- [ ] `messages/commands/weather.go` — FetchWeather
- [ ] `messages/events/weather.go` — WeatherUpdated
- [ ] `integrations/weather.go` — Weather API client
- [ ] `integrations/weather_test.go` — Mocked API tests

---

## Phase 5: Projections (Read Models)

### 5.1 SQL Schema
- [ ] `ui/projections/schema.sql` — All projection tables:
  - `device_states` — current state of all devices
  - `room_dashboard` — devices grouped by room
  - `automation_log` — automation execution history
  - `voice_command_log` — voice command history
  - `energy_usage` — power consumption tracking
  - `device_health` — connectivity and uptime

### 5.2 Projection Handlers
- [ ] `ui/projections/device_state.go` — Handle device events → update device_states table
- [ ] `ui/projections/room_dashboard.go` — Handle room + device events → update room_dashboard
- [ ] `ui/projections/automation_log.go` — Handle automation events → append to automation_log
- [ ] `ui/projections/voice_command_log.go` — Handle voice events → append to voice_command_log
- [ ] `ui/projections/device_health.go` — Handle online/offline events → update device_health
- [ ] `ui/projections/energy_usage.go` — Handle power reports → aggregate daily kWh
- [ ] `ui/projections/db.go` — SQLite initialization, schema migration
- [ ] Tests for each projection handler

---

## Phase 6: API Layer (Go HTTP + WebSocket)

### 6.1 REST API Endpoints
- [ ] `ui/api/handlers.go` — Command dispatch handler (generic command executor)
- [ ] `ui/api/devices.go` — GET /api/devices, GET /api/devices/:id, POST /api/devices/:id/command
- [ ] `ui/api/rooms.go` — GET /api/rooms, POST /api/rooms, PUT /api/rooms/:id
- [ ] `ui/api/scenes.go` — GET /api/scenes, POST /api/scenes, POST /api/scenes/:id/activate
- [ ] `ui/api/automations.go` — GET /api/automations, POST /api/automations, PUT /api/automations/:id
- [ ] `ui/api/voice.go` — POST /api/voice/command, GET /api/voice/log
- [ ] `ui/api/projections.go` — GET /api/dashboard, GET /api/energy, GET /api/health

### 6.2 WebSocket Event Stream
- [ ] `ui/websocket/hub.go` — Connection hub, broadcast Dogma facts to subscribers
- [ ] `ui/websocket/handler.go` — WebSocket upgrade, message routing
- [ ] Real-time event streaming for live dashboard updates

---

## Phase 7: React Dashboard UI

### 7.1 Infrastructure
- [ ] `src/api/client.ts` — Fetch-based API client with TypeScript types
- [ ] `src/api/websocket.ts` — WebSocket connection manager with auto-reconnect
- [ ] `src/hooks/useDeviceState.ts` — React Query hook for device state with WebSocket updates
- [ ] `src/hooks/useCommand.ts` — Hook for dispatching commands
- [ ] `src/hooks/useEventStream.ts` — Hook for subscribing to live Dogma facts
- [ ] `src/hooks/useProjections.ts` — Hooks for querying read models

### 7.2 Dashboard Page
- [ ] `src/pages/Home.tsx` — Main dashboard layout
- [ ] `src/components/dashboard/Dashboard.tsx` — Overview with room cards, quick actions, recent activity
- [ ] `src/components/dashboard/RoomCard.tsx` — Room summary with device count and status
- [ ] `src/components/dashboard/DeviceCard.tsx` — Device state display + quick controls (toggle, slider)
- [ ] `src/components/dashboard/QuickActions.tsx` — Favorite scenes and automations

### 7.3 Devices Page
- [ ] `src/pages/Devices.tsx` — Device listing page
- [ ] `src/components/devices/DeviceList.tsx` — Filterable, sortable device grid
- [ ] `src/components/devices/DeviceDetail.tsx` — Full device controls, state history, config editing
- [ ] `src/components/devices/DeviceRegistry.tsx` — Manage aliases, room assignment, capabilities
- [ ] `src/components/devices/ProvisionWizard.tsx` — Step-by-step device onboarding wizard

### 7.4 Scenes Page
- [ ] `src/pages/Scenes.tsx` — Scene listing page
- [ ] `src/components/scenes/SceneList.tsx` — Scene cards with activate button
- [ ] `src/components/scenes/SceneBuilder.tsx` — Visual scene editor (add devices, set target states)
- [ ] `src/components/scenes/SceneDetail.tsx` — Scene details and activation history

### 7.5 Automations Page
- [ ] `src/pages/Automations.tsx` — Automation listing page
- [ ] `src/components/automations/AutomationList.tsx` — Automation cards with enable/disable toggle
- [ ] `src/components/automations/AutomationBuilder.tsx` — Visual rule editor (trigger → conditions → actions)
- [ ] `src/components/automations/AutomationLog.tsx` — Execution history timeline

### 7.6 Voice Page
- [ ] `src/components/voice/VoiceLog.tsx` — Voice command history with intent display
- [ ] `src/components/voice/VoiceSettings.tsx` — Wake word config, TTS settings, output device selection

### 7.7 Common Components
- [ ] `src/components/common/EventTimeline.tsx` — Live event stream visualization (Dogma facts)
- [ ] `src/components/common/StateIndicator.tsx` — Online/offline/unknown status badge
- [ ] `src/components/common/ConfirmDialog.tsx` — Confirmation modal for destructive actions
- [ ] `src/components/common/LoadingSpinner.tsx` — Loading state component

### 7.8 Settings Page
- [ ] `src/pages/Settings.tsx` — User preferences, integrations config, system settings

---

## Phase 8: Polish & Testing

### 8.1 Integration Testing
- [ ] End-to-end test: device registration → room assignment → scene creation → scene activation
- [ ] End-to-end test: voice command → LLM parsing → device control → TTS response
- [ ] End-to-end test: automation creation → trigger event → condition evaluation → action execution
- [ ] Test device offline/online transitions and health tracking

### 8.2 UI Polish
- [ ] Responsive design for mobile/tablet
- [ ] Dark mode support
- [ ] Loading states and error boundaries
- [ ] Toast notifications for command results
- [ ] Keyboard shortcuts for common actions

### 8.3 Documentation
- [ ] README.md — Project overview, architecture diagram, setup instructions
- [ ] API documentation — Endpoint reference
- [ ] Domain documentation — Aggregate and process descriptions
- [ ] Sample configurations — Docker compose, environment variables

---

## Phase Order & Dependencies

```
Phase 1 (Scaffolding)
    ↓
Phase 2 (Aggregates) ──────────────────────┐
    ↓                                       ↓
Phase 3 (Processes) ───────────── Phase 5 (Projections)
    ↓                                   ↓
Phase 4 (Integrations)            Phase 6 (API Layer)
    └──────────┬─────────────────────┘
               ↓
          Phase 7 (React UI)
               ↓
          Phase 8 (Polish)
```

**Note:** Phases 4 (Integrations) and 5 (Projections) can be developed in parallel once Phase 2 aggregates and messages are defined. Phase 6 requires both 4 and 5. Phase 7 requires Phase 6.

---

## Key Design Decisions

| Decision | Choice | Rationale |
|---|---|---|
| **Backend language** | Go | Dogma is Go-native, excellent for event-driven systems |
| **Frontend framework** | React + Vite + TypeScript | Fast dev experience, strong typing, large ecosystem |
| **Styling** | Tailwind CSS | Utility-first, fast iteration, good component libraries |
| **Database** | SQLite (via projectionkit) | Zero-config, embedded, sufficient for single-home scale |
| **Real-time** | WebSocket | Low-latency event streaming for live dashboard |
| **API style** | REST + WebSocket | REST for commands/queries, WebSocket for live events |
| **LLM backend** | LM Studio | Local, OpenAI-compatible, privacy-first |
| **Device protocol** | Home Assistant API + MQTT | HA as hub, MQTT for direct device communication |

---

## Success Criteria

- [ ] Can register a device and see it appear in the React dashboard in real-time
- [ ] Can create a scene and activate it (multiple devices change state)
- [ ] Can create an automation rule and see it trigger on events
- [ ] Can send a voice command and see the full pipeline execute (LLM parse → action → TTS)
- [ ] Event timeline shows all Dogma facts in real-time
- [ ] Device offline/online transitions are tracked and displayed
- [ ] All aggregates have comprehensive unit tests
- [ ] All processes have scenario tests using the Dogma testkit
