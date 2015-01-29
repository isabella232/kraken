package main

import (
	"testing"

	"github.com/xoom/jira"
)

func TestFindMapping(t *testing.T) {
	mappings := map[int]jira.Mapping{
		1: jira.Mapping{ProjectID: 2, ComponentID: 4, VersionID: 8},
		2: jira.Mapping{ProjectID: 4, ComponentID: 16, VersionID: 32},
	}

	mapping, present := findMapping(mappings, "2", "4", "8")
	if !present {
		t.Fatalf("Want true\n")
	}
	if mapping.ProjectID != 2 {
		t.Fatalf("Want 2 but got %d\n", mapping.ProjectID)
	}
	if mapping.ComponentID != 4 {
		t.Fatalf("Want 4 but got %d\n", mapping.ComponentID)
	}
	if mapping.VersionID != 8 {
		t.Fatalf("Want 8 but got %d\n", mapping.VersionID)
	}

	if _, present := findMapping(mappings, "2", "4", "2"); present {
		t.Fatalf("Want false\n")
	}
}
