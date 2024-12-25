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
```
### `/rolldice/{player}`
```bash
curl "http://localhost:8080/rolldice/?player=Alice"
```
### `/logs`
```bash
curl http://localhost:8080/logs
```
### `/metrics/reset`
```bash
curl http://localhost:8080/metrics/reset
```
## Installation and Running
### Clone Repository
```bash
git clone https://github.com/OTEL_IN_GO.git
cd OTEL_IN_GO
```
### Install Dependencies
```bash
go mod tidy
```
### Run the Server
```bash
go run .
```
### Access Server
http://localhost:8080/rolldice

### Observability Exporters
- Traces: `stdouttrace.New`
- Metrics: `stdoutmetric.New`
- Logs: `stdoutlog.New`
