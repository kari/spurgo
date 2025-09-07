package main

import (
	"strings"
	"testing"
)

func TestSample(t *testing.T) {
	s, err := Sample("test/data.txt", "")
	if err != nil {
		t.Errorf("got error: %s", err)
	}
	if s == "" {
		t.Errorf("sample returned empty, expected any line")
	}
}

func TestSampleSearch(t *testing.T) {
	s, err := Sample("test/data.txt", "akka")
	if err != nil {
		t.Errorf("got error: %s", err)
	}
	if !strings.Contains(strings.ToLower(s), "akka") {
		t.Errorf("sample does not contain substring")
	}

	_, err = Sample("test/data.txt", "Vantaa") // vantaasta ei ole vertauskuvia
	if err == nil {
		t.Errorf("expected error")
	}
	if err.Error() != "no matching lines found" {
		t.Errorf("expected 'no matching lines found', got %s", err)
	}
}

func TestInvalidFile(t *testing.T) {
	_, err := Sample("data/yolo.txt", "")
	if err == nil {
		t.Errorf("err should be non-null")
	}
}
