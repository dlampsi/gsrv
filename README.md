# gsrv

Go module with HTTP and GRPC servers. Servers gracefully shut down when root context is cancelled.

## Usage

Basic usage:
```go
import "github.com/dlampsi/gsrv"

func main() {
    srv, err := gsrv.New("0.0.0.0:8080")
    if err != nil {
        // Process the error 
    }

    handler := CreateYourCustomHttpHandler()

    if err := srv.ServeHTTP(ctx, handler); if err != nil {
        // Process the error 
    }
}
```

With custom shutdown timeout:
```go
import (
    "time"

    "github.com/dlampsi/gsrv"
)

func main() {
    shutdownTimeout := 15 * time.Second
    srv, err := gsrv.New("0.0.0.0:8080", gsrv.WithTimeout(shutdownTimeout))
    if err != nil {
        // Process the error 
    }

    handler := CreateYourCustomHttpHandler()

    if err := srv.ServeHTTP(ctx, handler); if err != nil {
        // Process the error 
    }
}
```
