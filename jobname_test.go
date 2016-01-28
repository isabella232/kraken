package main

import "testing"

func TestJobNameToComponent(t *testing.T) {
	j := "plat-a-bcd-release"
	component := componentNameFromJobname(j)
	if component != "a-bcd" {
		t.Fatalf("Want a-bcd but got %s\n", component)
	}
}
