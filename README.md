# natsk - NATS Client Wrapper

📨 Custom Konsultin NATS client wrapper with auto-reconnect support.

## Installation

```bash
go get github.com/konsultin/natsk
```

## Quick Start

```go
import "github.com/konsultin/natsk"

// Connect to NATS
client, err := natsk.New("nats://localhost:4222")
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// Subscribe to a subject
client.Subscribe("orders.created", func(msg *nats.Msg) {
    fmt.Printf("Received: %s\n", string(msg.Data))
})

// Publish a message
client.Publish("orders.created", []byte(`{"id": "123"}`))

// Queue Subscribe (load-balanced)
client.QueueSubscribe("orders.created", "order-workers", func(msg *nats.Msg) {
    fmt.Printf("Worker received: %s\n", string(msg.Data))
})
```

## Features

- **Auto-reconnect** with configurable retry attempts
- **Queue subscriptions** for load balancing
- **Graceful shutdown** with drain support
- **Simple API** - just `New()`, `Subscribe()`, `Publish()`, `Close()`

## License

MIT License - see [LICENSE](LICENSE)
