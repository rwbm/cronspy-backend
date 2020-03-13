package server_test

import (
	"cronspy/backend/pkg/util/server"
	"testing"
)

func TestNew(t *testing.T) {
	e := server.New()
	if e == nil {
		t.Errorf("Server should not be nil")
	}
}
