# OpenTelemetry example in golang
[OpenTelemetry](https://opentelemetry.io/) is an observability framework – an API, SDK, and tools that are designed to aid in the generation and collection of application telemetry data such as metrics, logs, and traces.
## Features
- Dice Rolling API for anonymous or named players.
- Telemetry with OpenTelemetry for tracing, metrics, and logging.
- Additional routes for viewing logs and resetting metrics.
- Graceful shutdown handling for safe server termination.

## Routes
### `/rolldice/`
```bash
curl http://localhost:8080/rolldice/
