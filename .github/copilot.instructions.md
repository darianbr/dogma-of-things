# DoT (Dogma of Things) — Project Instructions

> These instructions apply to **every AI interaction** in this workspace. Always read the research and plan before making changes or answering questions.

---

## Project Identity

- **Name:** DoT (Dogma of Things)
- **Tagline:** An IoT smart home management system built on the Dogma event-sourcing framework
- **Concept:** "Dogma of Things" — applying Dogma's message-driven, event-sourced architecture to IoT device management
- **Stack:** Go (Dogma engine) + TypeScript/React (Vite) + SQLite (projections)

---

## Always Read Before Responding

### Research Documents (Ground Truth)
The `/context/research/` directory contains all architectural decisions, device inventory, platform choices, and technical rationale. **Always reference these before making architectural or technology decisions.**

| File | Covers |
|---|---|
| `context/research/RESEARCH-01-Project-Overview.md` | Goals, core requirements, design principles (privacy-first, local LLM, custom wake word) |
| `context/research/RESEARCH-02-Device-Inventory-&-Compatibility.md` | Eufy devices, Google Home Mini limitations, integration status |
| `context/research/RESEARCH-03-Platform-Decision-Home-Assistant.md` | Why Home Assistant, Assist pipeline, LM Studio integration point |
| `context/research/RESEARCH-04-LM-Studio-Integration.md` | LLM integration approaches, model recommendations |
| `context/research/RESEARCH-05-LLM-Prompt-Strategy-Device-&-Intent-Mapping.md` | Frontloaded mapping, entity resolution, ambiguity handling |
| `context/research/RESEARCH-06-Voice-Activation-&-Custom-Wake-Words.md` | "Hey Dibbles" wake word, openWakeWord, hardware options |
| `context/research/RESEARCH-07-Google-Home-Mini-Repurposing.md` | Google Mini as TTS output only |
| `context/research/RESEARCH-08-Recommended-Architecture.md` | Full system diagram, data flow for voice commands |
| `context/research/RESEARCH-09-Implementation-Roadmap.md` | Five-phase implementation plan |
| `context/research/RESEARCH-10-Risks-&-Open-Questions.md` | Technical risks, mitigations, open questions |
| `context/research/RESEARCH-11-References.md` | External documentation links |
| `context/research/RESEARCH-12-Quick-Decision-Summary.md` | One-page decision table (platform, LLM, wake word, STT, TTS, model) |
| `context/research/RESEARCH-13-Sample-System-Prompt-for-LLM.md` | LLM system prompt template with device registry injection |

### Implementation Plan
- **`context/PLAN.md`** — The 8-phase implementation plan with checkbox-tracked tasks. Always check this for current phase, task status, and dependency order before implementing features.

---

## Architecture Principles

1. **Event Sourcing First** — Every state change is an event. The event log is the source of truth.
2. **CQRS** — Commands mutate state (write side), projections serve queries (read side).
3. **Process Managers for Workflows** — Multi-step operations (scene activation, voice commands, provisioning) use Dogma processes with deadlines and retries.
4. **Integrations for External Systems** — All external calls (HA API, MQTT, LLM, weather) go through integration handlers that record events.
5. **Privacy-First** — Local LLM (LM Studio), local STT (Whisper), local TTS (Piper). No cloud dependency for core functionality.
6. **Home Assistant as Device Hub** — HA handles device discovery and protocol translation; DoT handles business logic and orchestration.

---

## Domain Model

### Aggregates
- **Device** — Physical IoT device (lights, switches, sensors, thermostats, speakers, cameras, vacuums)
- **Room** — Physical location grouping devices
- **Scene** — Named multi-device state configuration
- **AutomationRule** — If-then rules with triggers, conditions, and actions
- **User** — User profiles with preferences and voice profiles

### Processes
- **SceneActivationProcess** — Orchestrates multi-device scene activation with retries and partial failure
- **AutomationEvaluationProcess** — Evaluates automation rules against events with cooldown management
- **VoiceCommandProcess** — Voice pipeline: raw text → LLM intent parsing → device resolution → execution → TTS response
- **DeviceProvisioningProcess** — Device onboarding: discovery → capabilities → room assignment → verification

### Integrations
- **HomeAssistantIntegration** — HA REST/WebSocket API for device control and state queries
- **MQTTIntegration** — MQTT broker for direct device communication
- **LLMIntegration** — LM Studio (OpenAI-compatible) for intent parsing and conversation
- **WeatherIntegration** — Weather API for automation conditions

### Projections (Read Models)
- `device_states` — Current state of all devices
- `room_dashboard` — Devices grouped by room
- `automation_log` — Automation execution history
- `voice_command_log` — Voice command history
- `energy_usage` — Power consumption tracking
- `device_health` — Device connectivity and uptime

---

## Coding Conventions

### Go (Backend) — Follow the [Dogma Example App](https://github.com/dogmatiq/example)

**Always reference the [Dogma example repository](https://github.com/dogmatiq/example) for patterns, conventions, and idiomatic style.** The example is the authoritative reference for how to structure Dogma applications.

#### Project Structure (mirrors the example)
```
app.go                  — Application struct + Configure (wires all handlers)
doc.go                  — Package-level documentation
cmd/server/main.go      — Engine bootstrap, HTTP server, fact observers
domain/                 — Aggregates and process managers
messages/
  commands/             — Command type definitions
  events/               — Event type definitions
  internal/             — Shared types (enums, constants used across messages)
integrations/           — External system integration handlers
ui/
  api/                  — HTTP API handlers
  projections/          — SQL projection handlers + schema
  websocket/            — Real-time event streaming
```

#### Aggregate Pattern (see `domain/account.go`, `domain/customer.go`)
- **Private aggregate root struct** (e.g., `type account struct`) with `dogma.NoSnapshotBehavior` embedded
- Implement `AggregateInstanceDescription() string` — returns a human-readable description for logging (empty string if no meaningful state yet)
- **Command handler methods** on the aggregate (e.g., `func (a *account) OpenAccount(s dogma.AggregateCommandScope[*account], m *commands.OpenAccount)`) — validate state, then `s.RecordEvent(...)` or `s.Log(...)` to reject
- **`ApplyEvent(m dogma.Event)` method** — rebuilds state from events (switch on event types, mutate fields)
- **Public handler struct** (e.g., `type AccountHandler struct{}`) with:
  - `Configure(c dogma.AggregateConfigurer)` — calls `c.Identity("name", "uuid")` and `c.Routes(...)` with `dogma.HandlesCommand[*T]()` and `dogma.RecordsEvent[*T]()`
  - `New() *aggregate` — returns a new zero-value aggregate instance
  - `RouteCommandToInstance(m dogma.Command) string` — extracts the aggregate ID from the command (panic with `dogma.UnexpectedMessage` for unknown types)
  - `HandleCommand(a *aggregate, s dogma.AggregateCommandScope[*aggregate], m dogma.Command)` — switch on command type, delegate to aggregate method

#### Process Pattern (see `domain/deposit.go`, `domain/transfer.go`)
- **Stateless processes:** embed `dogma.StatelessProcessBehavior` and `dogma.NoDeadlineMessagesBehavior[*dogma.StatelessProcessRoot]`
- **Stateful processes:** private struct with state fields, `MarshalBinary`/`UnmarshalBinary`, and `ProcessInstanceDescription(ended bool) string`
- Handler struct with:
  - `Configure(c dogma.ProcessConfigurer)` — `c.Identity()`, `c.Routes()` with `dogma.HandlesEvent[*T]()`, `dogma.ExecutesCommand[*T]()`
  - `RouteEventToInstance(ctx context.Context, m dogma.Event) (string, bool, error)` — returns process ID
  - `HandleEvent(ctx context.Context, root T, s dogma.ProcessEventScope, m dogma.Event)` — switch on event, execute commands or schedule deadlines
- Register deadlines in `init()` with `dogma.RegisterDeadline[*DeadlineType]("uuid")`

#### Integration Pattern (see `integrations/thirdpartybank.go`)
- Handler struct (e.g., `type ThirdPartyBankIntegrationHandler struct{}`)
- `Configure(c dogma.IntegrationConfigurer)` — `c.Identity()`, `c.Routes()` with `dogma.HandlesCommand[*T]()` and `dogma.RecordsEvent[*T]()`
- `HandleCommand(ctx context.Context, s dogma.IntegrationCommandScope, c dogma.Command)` — switch on command, make external call, `s.RecordEvent(...)` for success or failure
- Add `doc.go` in the package with a description

#### Projection Pattern (see `ui/projections/customer.go`, `ui/projections/ledger.go`)
- Handler struct with `sqlprojection.NoCompactBehavior` embedded
- `Configure(c dogma.ProjectionConfigurer)` — `c.Identity()`, `c.Routes()` with `dogma.HandlesEvent[*T]()`
- `HandleEvent(ctx context.Context, tx *sql.Tx, s dogma.ProjectionEventScope, m dogma.Event)` — switch on event, execute SQL within the transaction
- `Reset(ctx context.Context, tx *sql.Tx, _ dogma.ProjectionResetScope)` — DELETE from all managed tables
- Schema in `*.sql` files embedded via `//go:embed *.sql`
- `db.go` with `NewDB()`, `MustNewDB()`, and `CreateSchema()` using embedded SQL files

#### Message Pattern (see `messages/events/customer.go`, `messages/commands/transfer.go`)
- Struct with fields, registered in `init()` with `dogma.RegisterEvent[*T]("uuid")` or `dogma.RegisterCommand[*T]("uuid")`
- Implement `MessageDescription() string` — human-readable summary for logging
- Implement `Validate(dogma.EventValidationScope) error` — validate required fields
- Implement `MarshalBinary() ([]byte, error)` and `UnmarshalBinary([]byte) error` — use JSON
- Group related messages in the same file (e.g., all transfer commands in `commands/transfer.go`)

#### Application Wiring (see `app.go`)
- `App` struct embeds all handler instances + `*sql.DB`
- `Configure(c dogma.ApplicationConfigurer)` calls `c.Identity()` then `c.Routes()` with:
  - `dogma.ViaAggregate(handler)` for each aggregate
  - `dogma.ViaProcess(handler)` for each process
  - `dogma.ViaIntegration(handler)` for each integration
  - `dogma.ViaProjection(sqlprojection.New(db, sqlprojection.SQLiteDriver, &handler))` for each projection

#### Testing Pattern (see `ui/projections/customer_test.go`, `integrations/thirdpartybank_test.go`)
- Use `github.com/dogmatiq/testkit` with dot import: `. "github.com/dogmatiq/testkit"`
- External test package (`package projections_test`) importing the app
- `Begin(t, &example.App{ReadDB: db}).EnableHandlers("handler-name").Prepare(...)`
- `ExecuteCommand(&commands.SomeCommand{...})` to drive behavior
- Verify state by querying the SQLite database directly

#### Main Entry Point (see `cmd/bank/main.go`)
- Signal handling with `signal.NotifyContext`
- Create DB, create App, create engine with `engine.New(runtimeconfig.FromApplication(app))`
- Set up `fact.ObserverFunc` for logging, filtering on relevant fact types
- `engine.OperationOption` with `engine.EnableProjections(true)`

### TypeScript/React (Frontend)
- TypeScript strict mode enabled
- Component structure: one component per file, collocated styles
- Use React Query for server state, Zustand or context for client state
- API types mirror Go message types (keep them in sync)
- WebSocket for real-time event streaming (Dogma facts)
- Tailwind CSS for styling

---

## Key Context

- **Wake word:** "Hey Dibbles" (via openWakeWord)
- **LLM model:** Qwen 2.5 7B or Llama 3.1 8B (local via LM Studio)
- **STT:** Whisper (local)
- **TTS:** Piper (local)
- **Device protocols:** Eufy (via HA), MQTT, Zigbee, Z-Wave
- **Output speakers:** Google Home Minis (TTS via Cast only — cannot change wake word)
- **Database:** SQLite via projectionkit

---

## When Implementing

1. **Check `PLAN.md`** — What phase are we in? What's the next task?
2. **Check relevant research** — Does the research address this domain? What decisions were made?
3. **Follow Dogma patterns** — Aggregates handle commands and record events. Processes orchestrate workflows. Integrations handle external calls. Projections build read models.
4. **Test with Dogma testkit** — Write deterministic tests that verify event sequences.
5. **Keep the UI decoupled** — The React frontend talks to the Go API via REST + WebSocket. Never mix UI concerns with domain logic.
