package main

import (
	"errors"
	"testing"

	"github.com/xoom/jira"
)

func TestGetOrCreateVersionMapHit(t *testing.T) {
	versions := map[string]jira.Version{"v1": jira.Version{Name: "v1"}, "v2": jira.Version{Name: "v2"}}
	v, err := getOrCreateVersion("1", "v1", versions, core{})
	if err != nil {
		t.Fatalf("Unexpected error: %v\n", err)
	}
	if v.Name != "v1" {
		t.Fatalf("Want v1 but got: %v\n", v.Name)
	}
}

func TestGetOrCreateVersionMapMiss(t *testing.T) {
	client := core{version: jira.Version{Name: "v1"}}
	v, err := getOrCreateVersion("1", "v1", map[string]jira.Version{}, client)
	if err != nil {
		t.Fatalf("Unexpected error: %v\n", err)
	}
	if v.Name != "v1" {
		t.Fatalf("Want v1 but got: %v\n", v.Name)
	}
}

func TestGetOrCreateVersionMapMissAndError(t *testing.T) {
	client := core{err: errors.New("Boom")}
	_, err := getOrCreateVersion("1", "v1", map[string]jira.Version{}, client)
	if err == nil {
		t.Fatalf("Expected an error\n")
	}
}
