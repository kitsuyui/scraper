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
	CLIRoot.SetArgs([]string{"validate", "-c", "../test_assets/invalid-config.json"})
}
