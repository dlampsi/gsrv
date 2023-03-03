package gsrv

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"
)

const defaultTimeout = 5 * time.Second

type Server struct {
	address  string
	listener net.Listener
	timeout  time.Duration
	logger   Logger
}

// Initializes new server and tries to create net listener on provided address if not specified.
func New(address string, opts ...Option) (*Server, error) {
	s := &Server{
		address: address,
		timeout: defaultTimeout,
		logger:  newnopLogger(),
	}
	for _, opt := range opts {
		opt(s)
	}
	if s.listener == nil {
		listener, err := net.Listen("tcp", address)
		if err != nil {
			return nil, fmt.Errorf("Failed to create listener on %s: %w", address, err)
		}
		s.listener = listener
	}
	return s, nil
}

type Option func(*Server)

func WithListener(l net.Listener) Option {
	return func(s *Server) {
		s.listener = l
	}
}

func WithLogger(l Logger) Option {
	return func(s *Server) {
		s.logger = l
	}
}

func WithTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.timeout = t
	}
}

type Logger interface {
	Info(args ...interface{})
	Debug(args ...interface{})
	Infof(template string, args ...interface{})
	Debugf(template string, args ...interface{})
}

type nopLogger struct{}

// Creates a new dummy non-operate logger that never prints any messages.
func newnopLogger() *nopLogger {
	return &nopLogger{}
}

func (l *nopLogger) Info(args ...interface{})                    {}
func (l *nopLogger) Debug(args ...interface{})                   {}
func (l *nopLogger) Infof(template string, args ...interface{})  {}
func (l *nopLogger) Debugf(template string, args ...interface{}) {}

// Serves HTTP server with a graceful shutdown on provided context done and shutdown timeout.
func (s *Server) ServeHTTP(ctx context.Context, handler http.Handler) error {
	s.logger.Info("Listening HTTP on ", s.address)

	srv := &http.Server{
		Handler: handler,
	}

	errCh := make(chan error, 1)
	go func() {
		<-ctx.Done()

		s.logger.Debug("Server context closed")
		shutdownCtx, done := context.WithTimeout(context.Background(), s.timeout)
		defer done()

		s.logger.Debug("Shutting down server")
		if err := srv.Shutdown(shutdownCtx); err != nil {
			select {
			case errCh <- err:
			default:
			}
		}
	}()

	// This will block until the provided context is closed.
	if err := srv.Serve(s.listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("Failed to serve: %w", err)
	}

	s.logger.Debugf("Server stopped")

	select {
	case err := <-errCh:
		return fmt.Errorf("Failed to shutdown: %w", err)
	default:
		return nil
	}
}

// Serves gRPC server with a graceful shutdown on provided context done.
func (s *Server) ServeGRPC(ctx context.Context, srv *grpc.Server) error {
	s.logger.Info("Listening GRPC on ", s.address)

	errCh := make(chan error, 1)
	go func() {
		<-ctx.Done()
		s.logger.Debug("Server context closed")
		s.logger.Debug("Shutting down server")
		srv.GracefulStop()
	}()

	if err := srv.Serve(s.listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return fmt.Errorf("Failed to serve: %w", err)
	}

	s.logger.Debugf("Server stopped")

	// Return any errors that happened during shutdown.
	select {
	case err := <-errCh:
		return fmt.Errorf("Failed to shutdown: %w", err)
	default:
		return nil
	}
}
