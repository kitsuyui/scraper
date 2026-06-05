package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestLoggingMiddleware(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	server := httptest.NewServer(loggingMiddleware(inner))
	defer server.Close()

	res, err := http.Get(server.URL + "/test-path")
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()

	got := buf.String()
	if !strings.Contains(got, "GET") {
		t.Errorf("log line missing method: %q", got)
	}
	if !strings.Contains(got, "/test-path") {
		t.Errorf("log line missing path: %q", got)
	}
	if !strings.Contains(got, "404") {
		t.Errorf("log line missing status: %q", got)
	}
}

var testConfigDir string
var serverContext ServerContext

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "tmp")
	if err != nil {
		log.Fatal(err)
	}
	testConfigDir = dir
	serverContext = ServerContext{}
	err = serverContext.setConfigDirectory(dir)
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	retCode := m.Run()
	os.Exit(retCode)
}

func TestBasicCRUD(t *testing.T) {
	content := `[{"type": "css", "label": "title", "query": "title"}]`
	contentHTML := `<html><head><title>Yay</title></head><body><h1>dummy</h1></body></html>`
	configPath := "/something-config.json"
	server := httptest.NewServer(http.HandlerFunc(serverContext.handler))
	defer server.Close()
	status, body, err := putRequest(server, content, configPath)
	if err != nil {
		t.Errorf("this must not be error")
	}
	if body != "" {
		t.Errorf("this must be empty")
	}
	if status != http.StatusOK {
		t.Errorf("this must be 200")
	}
	status, body, err = getRequest(server, configPath)
	if err != nil {
		t.Errorf("must not be error")
	}
	if status != http.StatusOK {
		t.Errorf("must be 200")
	}
	if body != content {
		t.Errorf("The content will be same as last PUT request content.")
	}
	status, body, err = postRequest(server, contentHTML, configPath)
	if err != nil {
		t.Errorf("must not be error")
	}
	if status != http.StatusOK {
		t.Errorf("must be 200")
	}
	var result interface{}
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		t.Errorf("must not be error")
	}
	status, _, err = deleteRequest(server, configPath)
	if err != nil {
		t.Errorf("must not be error")
	}
	if status != http.StatusOK {
		t.Errorf("must be 200")
	}
}

func TestGetNotExists(t *testing.T) {
	configPath := "/not-exists.json"
	server := httptest.NewServer(http.HandlerFunc(serverContext.handler))
	defer server.Close()
	status, _, err := getRequest(server, configPath)
	if err != nil {
		t.Errorf("this must not be error")
	}
	if status != http.StatusNotFound {
		t.Errorf("this must be 404")
	}
}

func TestPostNotExists(t *testing.T) {
	configPath := "/not-exists.json"
	server := httptest.NewServer(http.HandlerFunc(serverContext.handler))
	defer server.Close()
	status, _, err := postRequest(server, "", configPath)
	if err != nil {
		t.Errorf("this must not be error")
	}
	if status != http.StatusNotFound {
		t.Errorf("this must be 404")
	}
}

func TestDeleteNotExists(t *testing.T) {
	configPath := "/not-exists.json"
	server := httptest.NewServer(http.HandlerFunc(serverContext.handler))
	defer server.Close()
	status, _, err := deleteRequest(server, configPath)
	if err != nil {
		t.Errorf("this must not be error")
	}
	if status != http.StatusNotFound {
		t.Errorf("this must be 404")
	}
}

func TestUnsupportedMethod(t *testing.T) {
	configPath := "/something-config.json"
	server := httptest.NewServer(http.HandlerFunc(serverContext.handler))
	defer server.Close()
	req, err := http.NewRequest(http.MethodPatch, server.URL+configPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("this must be 405")
	}
	if allow := res.Header.Get("Allow"); allow != "GET, POST, PUT, DELETE" {
		t.Errorf("Allow header must list supported methods")
	}
}

func TestBodySizeLimit(t *testing.T) {
	content := `[{"type": "css", "label": "title", "query": "title"}]`
	configPath := "/size-limit-test-config.json"

	const testLimit = 100
	limitedHandler := http.MaxBytesHandler(http.HandlerFunc(serverContext.handler), testLimit)
	server := httptest.NewServer(limitedHandler)
	defer server.Close()

	status, _, err := putRequest(server, content, configPath)
	if err != nil {
		t.Fatal(err)
	}
	if status != http.StatusOK {
		t.Errorf("PUT config: expected 200, got %d", status)
	}

	largeBody := strings.Repeat("x", testLimit+1)

	status, _, err = postRequest(server, largeBody, configPath)
	if err != nil {
		t.Fatal(err)
	}
	if status != http.StatusRequestEntityTooLarge {
		t.Errorf("POST over limit: expected 413, got %d", status)
	}

	status, _, err = putRequest(server, largeBody, configPath)
	if err != nil {
		t.Fatal(err)
	}
	if status != http.StatusRequestEntityTooLarge {
		t.Errorf("PUT over limit: expected 413, got %d", status)
	}
}

func postRequest(server *httptest.Server, postHTML string, postPath string) (int, string, error) {
	b := bytes.NewBufferString(postHTML)
	res, err := http.Post(server.URL+postPath, "application/octet-stream", b)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return res.StatusCode, string(body), err
}

func putRequest(server *httptest.Server, postJSON string, postPath string) (int, string, error) {
	b := bytes.NewBufferString(postJSON)
	client := &http.Client{}
	req, err := http.NewRequest("PUT", server.URL+postPath, b)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	return res.StatusCode, string(body), err
}

func getRequest(server *httptest.Server, getPath string) (int, string, error) {
	res, err := http.Get(server.URL + getPath)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return res.StatusCode, string(body), err
}

func deleteRequest(server *httptest.Server, deletePath string) (int, string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", server.URL+deletePath, nil)
	if err != nil {
		log.Fatal(err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	return res.StatusCode, string(body), err
}
