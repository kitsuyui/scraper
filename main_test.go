package main

import (
	"bytes"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/docopt/docopt-go"
)

type commandResult struct {
	exited   bool
	exitCode int
	stdout   string
	stderr   string
}

func runCLI(t *testing.T, args []string, setup func()) commandResult {
	t.Helper()

	oldArgs := os.Args
	oldExit := exit
	oldCreateServer := createServer
	oldStandardOutput := standardOutput
	oldErrorOutput := errorOutput

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	result := commandResult{exitCode: -1}

	os.Args = append([]string(nil), args...)
	standardOutput = &stdout
	errorOutput = &stderr
	exit = func(i int) {
		result.exited = true
		result.exitCode = i
	}
	if setup != nil {
		setup()
	}

	t.Cleanup(func() {
		os.Args = oldArgs
		exit = oldExit
		createServer = oldCreateServer
		standardOutput = oldStandardOutput
		errorOutput = oldErrorOutput
	})

	main()
	result.stdout = stdout.String()
	result.stderr = stderr.String()
	return result
}

func assertExitedNonZero(t *testing.T, result commandResult) {
	t.Helper()

	if !result.exited {
		t.Fatal("expected command to exit")
	}
	if result.exitCode == 0 {
		t.Fatal("exit status must not be 0")
	}
	if result.stdout != "" {
		t.Errorf("stdout = %q, want empty", result.stdout)
	}
	if result.stderr == "" {
		t.Error("stderr is empty")
	}
}

func TestServerFailBindInvalidPort(t *testing.T) {
	result := runCLI(t, []string{"scraper", "server", "-p", "-1", "-d", "."}, nil)
	assertExitedNonZero(t, result)
}

func TestValidateConfig(t *testing.T) {
	result := runCLI(t, []string{"scraper", "validate", "-c", "test_assets/scraper-config.json"}, nil)
	if result.exited {
		t.Fatalf("unexpected exit status: %d", result.exitCode)
	}
	if result.stdout != "" {
		t.Errorf("stdout = %q, want empty", result.stdout)
	}
	if result.stderr != "" {
		t.Errorf("stderr = %q, want empty", result.stderr)
	}
}

func TestValidateConfigInvalid(t *testing.T) {
	result := runCLI(t, []string{"scraper", "validate", "-c", "test_assets/invalid-config.json"}, nil)
	assertExitedNonZero(t, result)
}

func TestValidateConfigNotExists(t *testing.T) {
	result := runCLI(t, []string{"scraper", "validate", "-c", "test_assets/not-exists.json"}, nil)
	if !result.exited {
		t.Fatal("expected command to exit")
	}
	if result.exitCode != exitConfigFile {
		t.Fatalf("exit status = %d, want %d", result.exitCode, exitConfigFile)
	}
	if result.stdout != "" {
		t.Errorf("stdout = %q, want empty", result.stdout)
	}
	if result.stderr == "" {
		t.Error("stderr is empty")
	}
}

func TestBasicInvalidOutput(t *testing.T) {
	result := runCLI(t, []string{"scraper", "-c", "/dev/null", "-i", "/dev/null", "-o", "/"}, nil)
	assertExitedNonZero(t, result)
}

func TestBasicInvalidInput(t *testing.T) {
	result := runCLI(t, []string{"scraper", "-c", "/dev/null", "-i", ":", "-o", "/dev/null"}, nil)
	assertExitedNonZero(t, result)
}

func TestBasicInvalidConfigFile(t *testing.T) {
	result := runCLI(t, []string{"scraper", "-c", ":", "-i", "/dev/null", "-o", "/dev/null"}, nil)
	assertExitedNonZero(t, result)
}

func TestBasicWritesOutputFile(t *testing.T) {
	outputFilepath := filepath.Join(t.TempDir(), "scraper-output.json")
	result := runCLI(t, []string{
		"scraper",
		"-c", "test_assets/scraper-config.json",
		"-i", "test_assets/ok.html",
		"-o", outputFilepath,
	}, nil)
	if result.exited {
		t.Fatalf("unexpected exit status: %d", result.exitCode)
	}
	if result.stdout != "" {
		t.Errorf("stdout = %q, want empty", result.stdout)
	}
	if result.stderr != "" {
		t.Errorf("stderr = %q, want empty", result.stderr)
	}

	output, err := os.ReadFile(outputFilepath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(output), "test passed") {
		t.Errorf("output does not contain scraped title: %s", output)
	}
}

func TestServerUsesConfDirOption(t *testing.T) {
	confDir := t.TempDir()
	stopErr := errors.New("stop before listen")
	result := runCLI(t, []string{"scraper", "server", "-d", confDir}, func() {
		createServer = func(bindHost string, bindPort int, configDir string) (*http.Server, error) {
			if configDir != confDir {
				t.Errorf("configDir = %q, want %q", configDir, confDir)
			}
			return nil, stopErr
		}
	})

	if !result.exited {
		t.Fatal("expected command to exit")
	}
	if result.exitCode != exitServer {
		t.Fatalf("exit status = %d, want %d", result.exitCode, exitServer)
	}
	if !strings.Contains(result.stderr, stopErr.Error()) {
		t.Errorf("stderr = %q, want to contain %q", result.stderr, stopErr.Error())
	}
}

func TestCLIErrorPathsUseDistinctExitCodes(t *testing.T) {
	cases := []struct {
		name           string
		args           func(t *testing.T) []string
		setup          func()
		wantExitCode   int
		stderrContains string
	}{
		{
			name: "validate config",
			args: func(t *testing.T) []string {
				return []string{"scraper", "validate", "-c", "test_assets/invalid-config.json"}
			},
			wantExitCode:   exitValidateConfig,
			stderrContains: "Error:",
		},
		{
			name: "server create",
			args: func(t *testing.T) []string {
				return []string{"scraper", "server", "-d", "."}
			},
			setup: func() {
				createServer = func(bindHost string, bindPort int, configDir string) (*http.Server, error) {
					return nil, errors.New("server unavailable")
				}
			},
			wantExitCode:   exitServer,
			stderrContains: "server unavailable",
		},
		{
			name: "input file",
			args: func(t *testing.T) []string {
				return []string{
					"scraper",
					"-c", "test_assets/config.json",
					"-i", "test_assets/not-exists.html",
					"-o", filepath.Join(t.TempDir(), "output.json"),
				}
			},
			wantExitCode: exitInputFile,
		},
		{
			name: "output file",
			args: func(t *testing.T) []string {
				return []string{
					"scraper",
					"-c", "test_assets/config.json",
					"-i", "test_assets/ok.html",
					"-o", t.TempDir(),
				}
			},
			wantExitCode: exitOutputFile,
		},
		{
			name: "config file",
			args: func(t *testing.T) []string {
				return []string{"scraper", "-c", "test_assets/not-exists.json"}
			},
			wantExitCode: exitConfigFile,
		},
		{
			name: "scrape",
			args: func(t *testing.T) []string {
				return []string{
					"scraper",
					"-c", "test_assets/config.json",
					"-i", "/",
					"-o", filepath.Join(t.TempDir(), "output.json"),
				}
			},
			wantExitCode: exitScrape,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			result := runCLI(t, tt.args(t), tt.setup)
			if !result.exited {
				t.Fatal("expected command to exit")
			}
			if result.exitCode != tt.wantExitCode {
				t.Fatalf("exit status = %d, want %d", result.exitCode, tt.wantExitCode)
			}
			if result.stdout != "" {
				t.Errorf("stdout = %q, want empty", result.stdout)
			}
			if result.stderr == "" {
				t.Fatal("stderr is empty")
			}
			if tt.stderrContains != "" && !strings.Contains(result.stderr, tt.stderrContains) {
				t.Errorf("stderr = %q, want to contain %q", result.stderr, tt.stderrContains)
			}
		})
	}
}

func TestOutputFilePreservedOnScrapeError(t *testing.T) {
	dir := t.TempDir()
	outputFilepath := filepath.Join(dir, "output.json")

	// Write initial content to the output file.
	initial := []byte(`{"previous": "content"}`)
	if err := os.WriteFile(outputFilepath, initial, 0o600); err != nil {
		t.Fatal(err)
	}

	// Trigger a scrape error by feeding a directory as the HTML input.
	result := runCLI(t, []string{
		"scraper",
		"-c", "test_assets/config.json",
		"-i", "/",
		"-o", outputFilepath,
	}, nil)

	if !result.exited || result.exitCode != exitScrape {
		t.Fatalf("expected exit with code %d, got exited=%v code=%d", exitScrape, result.exited, result.exitCode)
	}

	// The original file must not have been truncated.
	got, err := os.ReadFile(outputFilepath)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(initial) {
		t.Errorf("output file was modified on scrape error: got %q, want %q", got, initial)
	}
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
