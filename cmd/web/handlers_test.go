package main

import (
	"net/http"
	"testing"

	"todo.khoirulakmal.dev/internal/assert"
)

func TestPing(t *testing.T) {
	// Initialize apps struct for logger dependency needed by middleware
	apps := newApps()

	// Start tls server for testing
	server := startTestServer(apps.routes())

	// Get url
	status, body := testGet(t, server, "/ping")

	// Assert result status
	assert.Equal(t, status, http.StatusOK)

	// Assert result body
	assert.Equal(t, body, "OK")

}
