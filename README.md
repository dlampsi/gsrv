# gsrv

[![Go Reference](https://pkg.go.dev/badge/github.com/dlampsi/gsrv.svg)](https://pkg.go.dev/github.com/dlampsi/gsrv)

A useful wrapper for starting and maintaining HTTP and gRPC servers with graceful shutdown on the root context cancel.

## Usage

Basic usage:

```go
import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

    "github.com/dlampsi/gsrv"
)

var (
    shutdownTimeout := 15 * time.Second
)

func main() {
	// Root context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Listening for OS interupt signals
	terminateCh := make(chan os.Signal, 1)
	signal.Notify(terminateCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-terminateCh:
			cancel()
		case <-ctx.Done():
		}
	}()

	// GSRV server
	srv, err := gsrv.New(
		"0.0.0.0:8080", 
		gsrv.WithTimeout(shutdownTimeout),
	)
	if err != nil {
		// Process the error 
	}

	/*
		Here comes your custom HTTP handler (router) 
		which can be any framework (like gin or gorilla)
	*/
	handler := CreateYourCustomHttpHandler()

	// This will block until the provided context is closed.
	err = srv.ServeHTTP(ctx, handler)
	cancel()
	if err != nil {
		// Process the error 
	}
}
```

### Custom logger

The module provides a [Logger](server.go#L63) interface that you can implement and use for your gsrv server.

For example, [zap logger](https://github.com/uber-go/zap) already satisfies this interface:

```go
logger, _ := zap.NewProduction()
defer logger.Sync()

srv, err := gsrv.New("0.0.0.0:8080", gsrv.WithLogger(logger))
if err != nil {
	// Process the error 
}
```