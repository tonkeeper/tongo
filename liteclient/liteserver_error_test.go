package liteclient

import (
	"context"
	"errors"
	"testing"
)

func TestLiteserverError_Error(t *testing.T) {
	inner := errors.New("boom")
	e := &LiteServerError{Address: "1.2.3.4:5050", Err: inner}
	if got, want := e.Error(), "liteserver 1.2.3.4:5050: boom"; got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

func TestLiteserverError_Unwrap(t *testing.T) {
	inner := errors.New("boom")
	e := &LiteServerError{Address: "1.2.3.4:5050", Err: inner}
	if got := errors.Unwrap(e); got != inner {
		t.Fatalf("Unwrap() = %v, want %v", got, inner)
	}
}

func TestLiteserverError_ErrorsIs_Transport(t *testing.T) {
	e := &LiteServerError{Address: "h", Err: context.DeadlineExceeded}
	if !errors.Is(e, context.DeadlineExceeded) {
		t.Fatalf("errors.Is should walk through LiteserverError to inner sentinel")
	}
}

func TestLiteserverError_ErrorsAs_Inner(t *testing.T) {
	protocol := LiteServerErrorC{Code: 42, Message: "boom"}
	wrapped := &LiteServerError{Address: "h", Err: protocol}

	var lerr *LiteServerError
	if !errors.As(wrapped, &lerr) {
		t.Fatalf("errors.As did not find *LiteserverError")
	}
	if lerr.Address != "h" {
		t.Fatalf("Address = %q, want %q", lerr.Address, "h")
	}

	var perr LiteServerErrorC
	if !errors.As(wrapped, &perr) {
		t.Fatalf("errors.As did not find LiteServerErrorC through Unwrap")
	}
	if perr.Code != 42 || perr.Message != "boom" {
		t.Fatalf("LiteServerErrorC fields lost: %+v", perr)
	}
}

func TestClient_wrapErr(t *testing.T) {
	c := &Client{
		connections: []*Connection{{host: "1.2.3.4:5050"}},
	}

	if got := c.wrapErr(nil); got != nil {
		t.Fatalf("wrapErr(nil) = %v, want nil", got)
	}

	// Transport-style error.
	transport := context.DeadlineExceeded
	wrapped := c.wrapErr(transport)
	var lerr *LiteServerError
	if !errors.As(wrapped, &lerr) {
		t.Fatalf("transport: errors.As did not find *LiteserverError")
	}
	if lerr.Address != "1.2.3.4:5050" {
		t.Fatalf("transport: Address = %q", lerr.Address)
	}
	if lerr.Err != transport {
		t.Fatalf("transport: Err = %v, want %v", lerr.Err, transport)
	}
	if !errors.Is(wrapped, context.DeadlineExceeded) {
		t.Fatalf("transport: errors.Is(context.DeadlineExceeded) lost")
	}

	// Protocol-style error: wrapErr should wrap LiteServerErrorC the same way;
	// callers reach the inner via errors.As.
	protocol := LiteServerErrorC{Code: 651, Message: "not found"}
	wrapped = c.wrapErr(protocol)
	if !errors.As(wrapped, &lerr) || lerr.Address != "1.2.3.4:5050" {
		t.Fatalf("protocol: wrapped error missing address: %+v", lerr)
	}
	var perr LiteServerErrorC
	if !errors.As(wrapped, &perr) || perr.Code != 651 {
		t.Fatalf("protocol: inner LiteServerErrorC not reachable: %+v", perr)
	}
}
