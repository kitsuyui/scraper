package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/kitsuyui/scraper/scraper"
)

type responseRecorder struct {
	http.ResponseWriter
	status int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.status = code
	rr.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &responseRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, rec.status, time.Since(start))
	})
}

const maxBodyBytes = 10 * 1024 * 1024 // 10 MB

type ServerContext struct {
	ConfigDirectory string
}

const allowedMethods = "GET, POST, PUT, DELETE"

func (s *ServerContext) setConfigDirectory(configDir string) error {
	absPath, err := filepath.Abs(configDir)
	if err != nil {
		return err
	}
	fi, err := os.Stat(absPath)
	if err != nil {
		return err
	}
	if !fi.Mode().IsDir() {
		return fmt.Errorf("%s is not a directory", absPath)
	}
	s.ConfigDirectory = absPath
	return nil
}

func (s *ServerContext) handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.handlerGET(w, r)
	case "POST":
		s.handlerPOST(w, r)
	case "PUT":
		s.handlerPUT(w, r)
	case "DELETE":
		s.handlerDELETE(w, r)
	default:
		w.Header().Set("Allow", allowedMethods)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (s *ServerContext) confFilePathFromRequest(r *http.Request) (string, error) {
	// To avoid directory traversal
	resolvedPath := filepath.Join(s.ConfigDirectory, filepath.FromSlash(path.Clean("/"+r.URL.Path)))
	if filepath.Ext(resolvedPath) != ".json" {
		return "", fmt.Errorf("only .json files are accessible")
	}
	return resolvedPath, nil
}

type errorResponse struct {
	status int
	body   string
	cause  string
}

func classifyError(err error) errorResponse {
	var maxBytesErr *http.MaxBytesError
	if errors.As(err, &maxBytesErr) {
		return errorResponse{
			status: http.StatusRequestEntityTooLarge,
			body:   http.StatusText(http.StatusRequestEntityTooLarge),
			cause:  "max_bytes",
		}
	} else if os.IsNotExist(err) {
		return errorResponse{
			status: http.StatusNotFound,
			body:   `{"Error": "The file does not exists"}`,
			cause:  "not_found",
		}
	} else if os.IsPermission(err) {
		return errorResponse{
			status: http.StatusForbidden,
			body:   `{"Error": "Forbidden"}`,
			cause:  "permission_denied",
		}
	}
	return errorResponse{
		status: http.StatusBadRequest,
		body:   `{"Error": "Something Wrong. Bad Request"}`,
		cause:  "bad_request",
	}
}

func writeErrorStatus(w http.ResponseWriter, r *http.Request, err error) {
	resp := classifyError(err)
	log.Printf("request failure method=%s path=%s status=%d cause=%s err=%v", r.Method, r.URL.Path, resp.status, resp.cause, err)
	if resp.cause == "max_bytes" {
		http.Error(w, resp.body, resp.status)
		return
	}
	w.WriteHeader(resp.status)
	w.Write([]byte(resp.body))
}

func (s *ServerContext) handlerGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	confPath, err := s.confFilePathFromRequest(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	confFile, err := os.Open(confPath)
	if err != nil {
		writeErrorStatus(w, r, err)
		return
	}
	defer confFile.Close()
	_, err = io.Copy(w, confFile)
	if err != nil {
		writeErrorStatus(w, r, err)
		return
	}
}

func (s *ServerContext) handlerPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	confPath, err := s.confFilePathFromRequest(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	confFile, err := os.Open(confPath)
	if err != nil {
		writeErrorStatus(w, r, err)
		return
	}
	defer confFile.Close()
	if err := scraper.ScrapeByConfFile(confFile, r.Body, w); err != nil {
		writeErrorStatus(w, r, err)
	}
}

func (s *ServerContext) handlerPUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	targetPath, err := s.confFilePathFromRequest(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	tmpFile, err := os.CreateTemp(filepath.Dir(targetPath), ".put-tmp-*")
	if err != nil {
		writeErrorStatus(w, r, err)
		return
	}
	tmpPath := tmpFile.Name()

	_, copyErr := io.Copy(tmpFile, r.Body)
	closeErr := tmpFile.Close()
	if copyErr != nil || closeErr != nil {
		os.Remove(tmpPath)
		if copyErr != nil {
			writeErrorStatus(w, r, copyErr)
		} else {
			writeErrorStatus(w, r, closeErr)
		}
		return
	}

	if err := os.Rename(tmpPath, targetPath); err != nil {
		os.Remove(tmpPath)
		writeErrorStatus(w, r, err)
		return
	}
}

func (s *ServerContext) handlerDELETE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	confPath, err := s.confFilePathFromRequest(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = os.Remove(confPath)
	if err != nil {
		writeErrorStatus(w, r, err)
		return
	}
}

func CreateServer(bindHost string, bindPort int, configDir string) (*http.Server, error) {
	sc := ServerContext{}
	err := sc.setConfigDirectory(configDir)
	if err != nil {
		return nil, err
	}
	bindAddr := fmt.Sprintf("%s:%d", bindHost, bindPort)
	handler := loggingMiddleware(http.MaxBytesHandler(http.HandlerFunc(sc.handler), maxBodyBytes))
	server := &http.Server{
		Addr:         bindAddr,
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	return server, nil
}
