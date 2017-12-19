package commands

import (
	"testing"
)

func TestValidateConfig(t *testing.T) {
	CLIRoot.SetArgs([]string{"validate", "-c", "../test_assets/scraper-config.json"})
	if err := CLIRoot.Execute(); err != nil {
		t.Error(err)
	}
}

func TestValidateConfigInvalid(t *testing.T) {
	exit = func(i int) {
		if i != 1 {
			t.Error("exit code Must be 1")
		}
	}
	CLIRoot.SetArgs([]string{"validate", "-c", "../test_assets/invalid-config.json"})
	CLIRoot.Execute()
}

func TestBasicInvalidOutput(t *testing.T) {
	exit = func(i int) {
		if i != 1 {
			t.Error("exit code Must be 1")
		}
	}
	CLIRoot.SetArgs([]string{"-c", "/dev/null", "-i", "/dev/null", "-o", ":"})
	CLIRoot.Execute()
}

func TestBasicInvalidInput(t *testing.T) {
	exit = func(i int) {
		if i != 1 {
			t.Error("exit code Must be 1")
		}
	}
	CLIRoot.SetArgs([]string{"-c", "/dev/null", "-i", ":", "-o", "/dev/null"})
	CLIRoot.Execute()
}

func TestBasicInvalidConfigFile(t *testing.T) {
	exit = func(i int) {
		if i != 1 {
			t.Error("exit code Must be 1")
		}
	}
	CLIRoot.SetArgs([]string{"-c", ":", "-i", "/dev/null", "-o", "/dev/null"})
	CLIRoot.Execute()
}
