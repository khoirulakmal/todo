package main

import (
	"net/http"
	"testing"

	"todo.khoirulakmal.dev/internal/assert"
)

func TestPing(t *testing.T) {
	// Initialize apps struct for logger dependency needed by middleware
	apps := newApps(t)

	// Start tls server for testing
	server := startTestServer(t, apps.routes())

	// Get url
	status, _, body := server.testGet(t, "/ping")

	// Assert result status
	assert.Equal(t, status, http.StatusOK)

	// Assert result body
	assert.Equal(t, body, "OK")

}

func TestGetList(t *testing.T) {
	app := newApps(t)

	// handler := app.session.LoadAndSave(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	expected := 1
	// 	app.session.Put(r.Context(), "dataID", 1)
	// 	actual := app.session.GetInt(r.Context(), "dataID")
	// 	assert.Equal(t, expected, actual)
	// }))

	testServer := startTestServer(t, app.session.LoadAndSave(app.routes()))
	defer testServer.Close()

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/todo/context", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Context-Name", "dataID")
	req.Header.Set("Context-Content", "1")
	_, _, _ = testServer.contextGet(t, req)

	req, err = http.NewRequest(http.MethodGet, testServer.URL+"/todo/created", nil)
	if err != nil {
		t.Fatal(err)
	}

	_, _, body := testServer.contextGet(t, req)

	assert.StringsContain(t, body, "Content mock for test")

}
