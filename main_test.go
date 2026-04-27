package main

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/docopt/docopt-go"
)

func TestServerFailBindInvalidPort(t *testing.T) {
	os.Args = []string{"scraper", "server", "-p", "-1", "-d", "."}
	exit = func(i int) {
		if i == 0 {
			t.Error("exit status must not be 0")
		}
	}
	main()
}

func TestValidateConfig(t *testing.T) {
	os.Args = []string{"scraper", "validate", "-c", "../test_assets/scraper-config.json"}
	exit = func(i int) {
		if i == 0 {
			t.Error("exit status must not be 0")
		}
	}
	main()
}

func TestValidateConfigInvalid(t *testing.T) {
	os.Args = []string{"scraper", "validate", "-c", "../test_assets/invalid-config.json"}
	exit = func(i int) {
		if i == 0 {
			t.Error("exit status must not be 0")
		}
	}
	main()
}

func TestValidateConfigNotExists(t *testing.T) {
	os.Args = []string{"scraper", "validate", "-c", "../test_assets/not-exists.json"}
	exit = func(i int) {
		if i == 0 {
			t.Error("exit status must not be 0")
		}
	}
	main()
}

func TestBasicInvalidOutput(t *testing.T) {
	os.Args = []string{"scraper", "-c", "/dev/null", "-i", "/dev/null", "-o", "/"}
	exit = func(i int) {
		if i == 0 {
			t.Error("exit status must not be 0")
		}
	}
	main()
}

func TestBasicInvalidInput(t *testing.T) {
	os.Args = []string{"scraper", "-c", "/dev/null", "-i", ":", "-o", "/dev/null"}
	exit = func(i int) {
		if i == 0 {
			t.Error("exit status must not be 0")
		}
	}
	main()
}

func TestBasicInvalidConfigFile(t *testing.T) {
	os.Args = []string{"scraper", "-c", ":", "-i", "/dev/null", "-o", "/dev/null"}
	exit = func(i int) {
		if i == 0 {
			t.Error("exit status must not be 0")
		}
	}
	main()
}

func TestBasicWritesOutputFile(t *testing.T) {
	oldArgs := os.Args
	oldExit := exit
	defer func() {
		os.Args = oldArgs
		exit = oldExit
	}()

	outputFilepath := filepath.Join(t.TempDir(), "scraper-output.json")
	os.Args = []string{
		"scraper",
		"-c", "test_assets/scraper-config.json",
		"-i", "test_assets/ok.html",
		"-o", outputFilepath,
	}
	exit = func(i int) {
		t.Fatalf("unexpected exit status: %d", i)
	}

	main()

	output, err := os.ReadFile(outputFilepath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(output), "test passed") {
		t.Errorf("output does not contain scraped title: %s", output)
	}
}

func TestServerUsesConfDirOption(t *testing.T) {
	oldArgs := os.Args
	oldExit := exit
	oldCreateServer := createServer
	defer func() {
		os.Args = oldArgs
		exit = oldExit
		createServer = oldCreateServer
	}()

	confDir := t.TempDir()
	stopErr := errors.New("stop before listen")
	os.Args = []string{"scraper", "server", "-d", confDir}
	exit = func(i int) {}
	createServer = func(bindHost string, bindPort int, configDir string) (*http.Server, error) {
		if configDir != confDir {
			t.Errorf("configDir = %q, want %q", configDir, confDir)
		}
		return nil, stopErr
	}

	main()
}

func TestMainCommand(t *testing.T) {
	opts, err := docopt.ParseArgs(usage, []string{}, "")
	if err != nil {
		// t.Errorf(err.Error())
		t.Errorf("%s", err)
	}
	configFilePath, err := opts.String("--config")
	if err != nil {
		t.Errorf("%s", err)
	}
	if configFilePath != "scraper-config.json" {
		t.Errorf("%s", configFilePath)
	}

	inputFilePath, err := opts.String("--input")
	if err == nil {
		t.Errorf("%s", err)
	}
	if inputFilePath != "" {
		t.Errorf("%s", inputFilePath)
	}

	outputFilePath, err := opts.String("--output")
	if err == nil {
		t.Errorf("%s", err)
	}
	if outputFilePath != "" {
		t.Errorf("%s", outputFilePath)
	}
}

func TestValidateSubCommand(t *testing.T) {
	opts, err := docopt.ParseArgs(usage, []string{"validate"}, "")
	if err != nil {
		t.Fail()
	}
	if validate, err := opts.Bool("validate"); err != nil || !validate {
		t.Fail()
	}
}

func TestServerSubCommand(t *testing.T) {
	opts, err := docopt.ParseArgs(usage, []string{"server"}, "")
	if err != nil {
		t.Fail()
	}

	if server, err := opts.Bool("server"); err != nil || !server {
		t.Fail()
	}

	host, err := opts.String("--host")
	if err != nil {
		t.Errorf("%s", err)
	}
	if host != "127.0.0.1" {
		t.Errorf("%s != 127.0.0.1", host)
	}

	port, err := opts.Int("--port")
	if err != nil {
		t.Errorf("%s", err)
	}
	if port != 8080 {
		t.Errorf("%d != 8080", port)
	}

	confDir, err := opts.String("--conf-dir")
	if err != nil {
		t.Errorf("%s", err)
	}
	if confDir != "." {
		t.Errorf("%s != .", host)
	}
}
