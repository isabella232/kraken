package main

import "testing"

func TestNextVersion(t *testing.T) {
	for _, v := range []string{"1.0", "1.0-SNAPSHOT"} {
		next := nextVersion(v)
		if next != "1.0" {
			t.Fatalf("Want 1.0 but got %s\n", next)
		}
	}
}
