package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/kitsuyui/scraper/scraper"
)

type ServerContext struct {
	ConfigDirectory string
}

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
	}
}

func (s *ServerContext) confFilePathFromRequest(r *http.Request) string {
	// To avoid directory traversal
	resolvedPath := filepath.Join(s.ConfigDirectory, filepath.FromSlash(path.Clean("/"+r.URL.Path)))
	return resolvedPath
}

func errorStatus(w http.ResponseWriter, err error) {
	if os.IsNotExist(err) {
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
	confFile, err := os.Open(s.confFilePathFromRequest(r))
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
	confFile, err := os.Open(s.confFilePathFromRequest(r))
	if err != nil {
		errorStatus(w, err)
		return
	}
	defer confFile.Close()
	scraper.ScrapeByConfFile(confFile, r.Body, w)
}

func (s *ServerContext) handlerPUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	confFile, err := os.Create(s.confFilePathFromRequest(r))
	if err != nil {
		errorStatus(w, err)
		return
	}
	defer confFile.Close()
	_, err = io.Copy(confFile, r.Body)
	if err != nil {
		errorStatus(w, err)
		return
	}
}

func (s *ServerContext) handlerDELETE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := os.Remove(s.confFilePathFromRequest(r))
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
	server := &http.Server{Addr: bindAddr, Handler: http.HandlerFunc(sc.handler)}
	return server, nil
}
