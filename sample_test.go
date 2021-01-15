package main

import (
	"strings"
	"testing"
)

func TestSample(t *testing.T) {
	s := sample("", "")
	if s == "" {
		t.Errorf("sample returned empty")
	}
	t.Log(s)
}

func TestSampleSearch(t *testing.T) {
	f := "lehmä"
	s := sample("", f)
	if !strings.Contains(strings.ToLower(s), f) {
		t.Errorf("sample does not contain substring")
	}
	t.Log(s)
	s = sample("", "Vantaa") // vantaasta ei ole vertauskuvia
}

func TestInvalidFile(t *testing.T) {
	s := sample("data/yolo.txt", "")
	t.Log(s)
}
