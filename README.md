# General Purpose Signal Generator (GPSG)

GPSG stands for general purpose signal generator.

- GPSG helps in real-time software waveform broadcasting.
- It can create waveforms of arbitrary shapes, and periodicity
- One can use it to create signals for simulators for
  - Electronic Circuits, for hardware analogue and digitization,
  - Data flow control, Flynn Taxonomy
  - Quantum Computing - probability simulation
  - Computer Networks - package transmission, noise, etc
  - Noise simulation
  - Even music production, Synths can be triggered using this


<br>
<br>

# How to Start GPSG

## 1. Start infrastructure services

From the project root directory:

```bash
sudo docker compose up -d
```

This starts:
- RabbitMQ (message broker)
- Redis / Valkey (state backend)

You can verify containers are running with:

```bash
docker ps
```

## 2. Start the GPSG server

From the same root directory:

```bash
go run main.go
```

If everything is working, you should see logs similar to:

```
Connected to Valkey
Successfully connected to RabbitMQ and declared queue: gpsg_queue
Starting GPSG server on :8080
```

At this point, GPSG is running and waiting for pattern requests.

## Starting a Signal Pattern

GPSG exposes an HTTP API to start signal patterns.

### Example: Start a pattern

```bash
curl -X POST http://localhost:8080/start-pattern \
  -H "Content-Type: application/json" \
  -d '{"waveform":"sine","time_script":[200,300,600]}'
```

### What this means

- A new signal pattern is created
- The pattern emits a signal event:
  - after 200 ms
  - then after 300 ms
  - then after 600 ms
  - and repeats this sequence indefinitely
- Each event is broadcast via RabbitMQ
- The response will be a pattern ID

## Checking Running Patterns

To see all active patterns:

```bash
curl http://localhost:8080/status
```

This returns:
- Server status
- Current server time
- Active pattern IDs
- Number of running patterns

## Stopping a Pattern

To stop a running pattern, use the pattern ID returned earlier:

```bash
curl "http://localhost:8080/stop-pattern?id=<PATTERN_ID>"
```

The pattern will stop immediately.
