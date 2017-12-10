package scraper

import (
	"bufio"
	"bytes"
	"os"
	"testing"
)

func TestScrapeByConfFileInvalid(t *testing.T) {
	testFilepath := "../test_assets/invalid-config.json"
	testHTMLFilepath := "../test_assets/ok.html"
	invalidConf, err := os.Open(testFilepath)
	if err != nil {
		t.Errorf("Must be opened %s", testFilepath)
	}
	input, err := os.Open(testHTMLFilepath)
	if err != nil {
		t.Errorf("Must be opened %s", testHTMLFilepath)
	}
	output := bufio.NewWriter(&bytes.Buffer{})
	err = ScrapeByConfFile(invalidConf, input, output)
	if err == nil {
		t.Errorf("Must be error on invalid config file: %s", testFilepath)
	}
}

func TestScrapeByConfFileInvalidHTML(t *testing.T) {
	testFilepath := "../test_assets/config.json"
	testHTMLFilepath := "../test_assets/broken.html"
	invalidConf, err := os.Open(testFilepath)
	if err != nil {
		t.Errorf("Must be opened %s", testFilepath)
	}
	input, err := os.Open(testHTMLFilepath)
	if err != nil {
		t.Errorf("Must be opened %s", testHTMLFilepath)
	}
	output := bufio.NewWriter(&bytes.Buffer{})
	err = ScrapeByConfFile(invalidConf, input, output)
	if err != nil {
		t.Errorf("Must not be error even if invalid HTML file: %s", testHTMLFilepath)
	}
}

func TestScrapeByConfFile(t *testing.T) {
	testFilepath := "../test_assets/config.json"
	testHTMLFilepath := "../test_assets/ok.html"
	invalidConf, err := os.Open(testFilepath)
	if err != nil {
		t.Errorf("Must be opened %s", testFilepath)
	}
	input, err := os.Open(testHTMLFilepath)
	if err != nil {
		t.Errorf("Must be opened %s", testHTMLFilepath)
	}
	output := bufio.NewWriter(&bytes.Buffer{})
	err = ScrapeByConfFile(invalidConf, input, output)
	if err != nil {
		t.Errorf("Must not be error on valid config file: %s", testFilepath)
	}
}

func TestInvalidConfig(t *testing.T) {
	testFilepath := "../test_assets/invalid-config.json"
	invalidConf, err := os.Open(testFilepath)
	if err != nil {
		t.Errorf("Must be opened %s", testFilepath)
	}
	err = ValidateConfigFile(invalidConf)
	if err == nil {
		t.Errorf("Must be error on invalid config file: %s", testFilepath)
	}
}

func TestBrokenConfig(t *testing.T) {
	testFilepath := "../test_assets/broken.json"
	invalidConf, err := os.Open(testFilepath)
	if err != nil {
		t.Errorf("Must be opened %s", testFilepath)
	}
	err = ValidateConfigFile(invalidConf)
	if err == nil {
		t.Errorf("Must be error on broken config file: %s", testFilepath)
	}
}

func TestValidConfig(t *testing.T) {
	testFilepath := "../test_assets/config.json"
	invalidConf, err := os.Open(testFilepath)
	if err != nil {
		t.Errorf("Must be opened %s", testFilepath)
	}
	err = ValidateConfigFile(invalidConf)
	if err != nil {
		t.Errorf("Must not be error on config file: %s", testFilepath)
	}
}
