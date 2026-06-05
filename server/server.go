package server

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/kitsuyui/scraper/scraper"
)

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

func errorStatus(w http.ResponseWriter, err error) {
	var maxBytesErr *http.MaxBytesError
	if errors.As(err, &maxBytesErr) {
		http.Error(w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
	} else if os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"Error": "The file does not exists"}`))
	} else if os.IsPermission(err) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"Error": "Forbidden"}`))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"Error": "Something Wrong. Bad Request"}`))
	}
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
		errorStatus(w, err)
		return
	}
	defer confFile.Close()
	_, err = io.Copy(w, confFile)
	if err != nil {
		errorStatus(w, err)
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
		errorStatus(w, err)
		return
	}
	defer confFile.Close()
	if err := scraper.ScrapeByConfFile(confFile, r.Body, w); err != nil {
		errorStatus(w, err)
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
		errorStatus(w, err)
		return
	}
	tmpPath := tmpFile.Name()

	_, copyErr := io.Copy(tmpFile, r.Body)
	closeErr := tmpFile.Close()
	if copyErr != nil || closeErr != nil {
		os.Remove(tmpPath)
		if copyErr != nil {
			errorStatus(w, copyErr)
		} else {
			errorStatus(w, closeErr)
		}
		return
	}

	if err := os.Rename(tmpPath, targetPath); err != nil {
		os.Remove(tmpPath)
		errorStatus(w, err)
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
		errorStatus(w, err)
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
	handler := http.MaxBytesHandler(http.HandlerFunc(sc.handler), maxBodyBytes)
	server := &http.Server{Addr: bindAddr, Handler: handler}
	return server, nil
}
