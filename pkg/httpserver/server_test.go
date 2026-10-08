package httpserver

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestNewDoesNotStartAndStartReportsBindFailure(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := listener.Close(); err != nil {
			t.Error(err)
		}
	}()
	_, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	srv := New(http.NotFoundHandler(), Port(port))
	select {
	case err := <-srv.Notify():
		t.Fatalf("constructor started a server: %v", err)
	default:
	}
	if err := srv.Start(); err == nil {
		t.Fatal("expected occupied port error")
	}
}

func TestShutdownDrainsActiveRequests(t *testing.T) {
	started, finish := make(chan struct{}), make(chan struct{})
	srv := New(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		close(started)
		<-finish
		if _, err := w.Write([]byte("ok")); err != nil {
			t.Error(err)
		}
	}), ShutdownTimeout(2*time.Second))
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	if err := srv.serve(listener); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := srv.Close(); err != nil {
			t.Error(err)
		}
	}()

	response := make(chan error, 1)
	go func() {
		client := &http.Client{Timeout: 3 * time.Second}
		resp, err := client.Get("http://" + listener.Addr().String())
		if err == nil {
			body, readErr := io.ReadAll(resp.Body)
			err = errors.Join(readErr, resp.Body.Close())
			if string(body) != "ok" {
				err = errors.Join(err, errors.New("response was interrupted"))
			}
		}
		response <- err
	}()
	select {
	case <-started:
	case <-time.After(3 * time.Second):
		t.Fatal("request did not reach handler")
	}
	shutdown := make(chan error, 1)
	go func() { shutdown <- srv.Shutdown() }()
	select {
	case err := <-srv.Notify():
		if !errors.Is(err, http.ErrServerClosed) {
			t.Fatalf("serve error: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("listener did not close")
	}
	select {
	case err := <-shutdown:
		t.Fatalf("shutdown returned before active request finished: %v", err)
	default:
	}
	close(finish)
	if err := <-response; err != nil {
		t.Fatal(err)
	}
	if err := <-shutdown; err != nil {
		t.Fatal(err)
	}
}

func TestShutdownTimeoutCanBeFollowedByClose(t *testing.T) {
	started := make(chan struct{})
	srv := New(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		close(started)
		<-r.Context().Done()
	}), ShutdownTimeout(50*time.Millisecond))
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	if err := srv.serve(listener); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := srv.Close(); err != nil {
			t.Error(err)
		}
	}()
	done := make(chan struct{})
	go func() {
		defer close(done)
		client := &http.Client{Timeout: 2 * time.Second}
		resp, err := client.Get("http://" + listener.Addr().String())
		if err == nil {
			if err := resp.Body.Close(); err != nil {
				t.Error(err)
			}
		}
	}()
	select {
	case <-started:
	case <-time.After(3 * time.Second):
		t.Fatal("request did not start")
	}
	if err := srv.Shutdown(); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("shutdown error = %v", err)
	}
	if err := srv.Close(); err != nil {
		t.Fatal(err)
	}
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("forced close did not release request")
	}
}
