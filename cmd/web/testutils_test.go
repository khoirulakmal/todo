package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form"
	"todo.khoirulakmal.dev/internal/models/mocks"
)

func newApps(t *testing.T) *application {
	// Parse template cache
	tmpl, err := parseTemplate()
	if err != nil {
		t.Fatal(err)
	}
	if len(tmpl) == 0 {
		t.Fatalf("Template not found!")
	}

	session := scs.New()
	session.Lifetime = 2 * time.Hour
	session.Cookie.Secure = true

	// Initialize form decoder
	formDecoder := form.NewDecoder()

	return &application{
		errorLog:      log.New(os.Stdout, "", 0),
		infoLog:       log.New(os.Stdout, "", 0),
		todos:         &mocks.TodoModel{},
		users:         &mocks.UserModel{},
		templateCache: tmpl,
		session:       session,
		formDecode:    formDecoder,
	}
}

// Define a custom testServer type which embeds a httptest.Server instance.
type testServer struct {
	*httptest.Server
}

func startTestServer(t *testing.T, handler http.Handler) *testServer {
	ts := httptest.NewTLSServer(handler)
	// Initialize a new cookie jar.
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add the cookie jar to the test server client. Any response cookies will
	// now be stored and sent with subsequent requests when using this client.
	ts.Client().Jar = jar

	// Disable redirect-following for the test server client by setting a custom
	// CheckRedirect function. This function will be called whenever a 3xx
	// response is received by the client, and by always returning a
	// http.ErrUseLastResponse error it forces the client to immediately return
	// the received response.
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func (ts *testServer) testGet(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}

func (ts *testServer) contextGet(t *testing.T, request *http.Request) (int, http.Header, string) {
	rs, err := ts.Client().Do(request)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}
