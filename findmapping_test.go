package main

import (
	"testing"

	"github.com/xoom/jira"
)

func TestFindMapping(t *testing.T) {
	mappings := map[int]jira.Mapping{
		1: jira.Mapping{ProjectID: 2, ComponentID: 4, VersionName: "v1"},
		2: jira.Mapping{ProjectID: 4, ComponentID: 16, VersionName: "v2"},
	}

	mapping, present := findMapping(mappings, "2", "4", "v1")
	if !present {
		t.Fatalf("Want true\n")
	}
	if mapping.ProjectID != 2 {
		t.Fatalf("Want 2 but got %d\n", mapping.ProjectID)
	}
	if mapping.VersionName != "v1" {
		t.Fatalf("Want v1 but got %d\n", mapping.VersionName)
	}
	if mapping.ComponentID != 4 {
		t.Fatalf("Want 4 but got %d\n", mapping.ComponentID)
	}

	if _, present := findMapping(mappings, "2", "4", "v2"); present {
		t.Fatalf("Want false\n")
	}
}
