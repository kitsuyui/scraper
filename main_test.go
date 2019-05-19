package main

import (
	"os"
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
	os.Args = []string{"scraper", "-c", "/dev/null", "-i", "/dev/null", "-o", ":"}
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

func TestMainCommand(t *testing.T) {
	opts, err := docopt.ParseArgs(usage, []string{}, "")
	if err != nil {
		t.Errorf(err.Error())
	}
	configFilePath, err := opts.String("--config")
	if err != nil {
		t.Errorf(err.Error())
	}
	if configFilePath != "scraper-config.json" {
		t.Errorf(configFilePath)
	}

	inputFilePath, err := opts.String("--input")
	if err == nil {
		t.Errorf(err.Error())
	}
	if inputFilePath != "" {
		t.Errorf(inputFilePath)
	}

	outputFilePath, err := opts.String("--output")
	if err == nil {
		t.Errorf(err.Error())
	}
	if outputFilePath != "" {
		t.Errorf(outputFilePath)
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
		t.Errorf(err.Error())
	}
	if host != "127.0.0.1" {
		t.Errorf("%s != 127.0.0.1", host)
	}

	port, err := opts.Int("--port")
	if err != nil {
		t.Errorf(err.Error())
	}
	if port != 8080 {
		t.Errorf("%d != 8080", port)
	}

	confDir, err := opts.String("--conf-dir")
	if err != nil {
		t.Errorf(err.Error())
	}
	if confDir != "." {
		t.Errorf("%s != .", host)
	}
}
