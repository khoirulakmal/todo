package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func newApps() *application {
	return &application{
		errorLog: log.New(os.Stdout, "", 0),
		infoLog:  log.New(os.Stdout, "", 0),
	}
}

func startTestServer(handler http.Handler) *httptest.Server {
	return httptest.NewTLSServer(handler)
}

func testGet(t *testing.T, server *httptest.Server, link string) (int, string) {
	result, err := server.Client().Get(server.URL + link)
	if err != nil {
		t.Fatal(err)
	}
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	return result.StatusCode, string(body)
}
