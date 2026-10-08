package httpserver

import (
	"log/slog"
	"net"
	"time"
)

type Option func(*Server)

func Port(port string) Option {
	return func(s *Server) { s.server.Addr = net.JoinHostPort("", port) }
}

func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) { s.server.ReadTimeout = timeout }
}

func ReadHeaderTimeout(timeout time.Duration) Option {
	return func(s *Server) { s.server.ReadHeaderTimeout = timeout }
}

func WriteTimeout(timeout time.Duration) Option {
	return func(s *Server) { s.server.WriteTimeout = timeout }
}

func IdleTimeout(timeout time.Duration) Option {
	return func(s *Server) { s.server.IdleTimeout = timeout }
}

func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) { s.shutdownTimeout = timeout }
}

func Logger(log *slog.Logger) Option {
	return func(s *Server) { s.server.ErrorLog = slog.NewLogLogger(log.Handler(), slog.LevelError) }
}
